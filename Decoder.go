package json

import (
	"fmt"
	"io"
	"reflect"
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
	stringStart := -1
	currentNumber := int64(0)
	inNumber := false
	isFloat := false
	fieldIndices := decoder.fieldIndexMap(t)
	fieldIndex := 0
	fieldExists := false

	for {
		n, err := decoder.reader.Read(decoder.buffer)

		for i := 0; i < n; i++ {
			c := decoder.buffer[i]

			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				currentNumber = (currentNumber * 10) + (int64(c) - '0')
				inNumber = true

			case '"':
				if stringStart >= 0 {
					captured := decoder.buffer[stringStart:i]

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

					stringStart = -1
				} else {
					stringStart = i + 1
				}

			case ',', '}':
				if inNumber {
					if isFloat {
						// TODO: ...
						number := 0.0
						v.Field(fieldIndex).SetFloat(number)
					} else {
						v.Field(fieldIndex).SetInt(currentNumber)
					}

					currentNumber = 0
					isFloat = false
					inNumber = false
					fieldExists = false
				}

			case '.':
				isFloat = true
			}
		}

		if err == io.EOF {
			return nil
		}
	}
}

// Close frees up resources and returns the decoder to the pool.
func (decoder *decoder) Close() {
	decoder.reader = nil
	decoderPool.Put(decoder)
}
