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
	fieldIndices := decoder.fieldIndexMap(v.Type())
	stringStart := -1

	var (
		i             int
		c             byte
		fieldIndex    int
		fieldExists   bool
		currentNumber int64
		inNumber      bool
		isFloat       bool
	)

	for {
		n, err := decoder.reader.Read(decoder.buffer)

		for i, c = range decoder.buffer[:n] {
			// String capture
			if stringStart >= 0 {
				if c == '"' {
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
				}

				continue
			}

			// Number capture
			if inNumber {
				if c >= '0' && c <= '9' {
					currentNumber = (currentNumber * 10) + (int64(c) - '0')
					continue
				}

				if c == '.' {
					isFloat = true
					continue
				}

				if c == ',' || c == '}' {
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
					continue
				}

				continue
			}

			switch c {
			case '"':
				stringStart = i + 1

			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				currentNumber = int64(c) - '0'
				inNumber = true
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
