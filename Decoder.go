package json

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/akyoto/stringutils/unsafe"
)

const (
	readBufferSize   = 4096
	stringBufferSize = 4096
)

var decoderPool = sync.Pool{
	New: func() interface{} {
		return &decoder{
			buffer:  make([]byte, readBufferSize),
			strings: make([]byte, stringBufferSize),
			types:   make(map[reflect.Type]fieldIndexMap),
		}
	},
}

type fieldIndexMap = map[string]int

type decoder struct {
	reader        io.Reader
	buffer        []byte
	strings       []byte
	stringsLength int
	types         map[reflect.Type]fieldIndexMap
}

// NewDecoder creates a new JSON decoder.
func NewDecoder(reader io.Reader) *decoder {
	poolObj := decoderPool.Get()
	decoder := poolObj.(*decoder)
	decoder.reader = reader
	return decoder
}

// Decode deserializes the JSON data into the given object.
func (decoder *decoder) Decode(object interface{}) error {
	v := reflect.ValueOf(object)
	v = v.Elem()
	t := v.Type()
	captureStart := -1
	fieldIndices := decoder.fieldIndices(t)
	fieldIndex := 0
	fieldExists := false

	for {
		n, err := decoder.reader.Read(decoder.buffer)

		for i := 0; i < n; i++ {
			c := decoder.buffer[i]

			switch c {
			case '"':
				if captureStart >= 0 {
					captured := decoder.buffer[captureStart:i]

					if fieldExists {
						length := len(captured)

						if decoder.stringsLength+length > len(decoder.strings) {
							newBufferLength := stringBufferSize

							if newBufferLength < length {
								newBufferLength = length
							}

							decoder.strings = make([]byte, newBufferLength)
							decoder.stringsLength = 0
						}

						destination := decoder.strings[decoder.stringsLength : decoder.stringsLength+length]
						copy(destination, captured)
						decoder.stringsLength += length
						v.Field(fieldIndex).SetString(unsafe.BytesToString(destination))
						fieldExists = false
					} else {
						fieldIndex, fieldExists = fieldIndices[string(captured)]

						if !fieldExists {
							return fmt.Errorf("Field does not exist: %s", string(captured))
						}
					}

					captureStart = -1
				} else {
					captureStart = i + 1
				}
			}
		}

		if err == io.EOF {
			return nil
		}
	}
}

// fieldIndices returns a map of field names mapped to their index.
func (decoder *decoder) fieldIndices(structType reflect.Type) fieldIndexMap {
	fieldsObj, exists := decoder.types[structType]

	if exists {
		return fieldsObj
	}

	fields := fieldIndexMap{}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		jsonName := field.Tag.Get("json")

		if jsonName != "" {
			comma := strings.Index(jsonName, ",")

			if comma != -1 {
				jsonName = jsonName[:comma]
			}

			fields[jsonName] = i
		} else {
			fields[field.Name] = i
		}
	}

	decoder.types[structType] = fields
	return fields
}

// Close frees up resources and returns the decoder to the pool.
func (decoder *decoder) Close() {
	decoder.reader = nil
	decoderPool.Put(decoder)
}
