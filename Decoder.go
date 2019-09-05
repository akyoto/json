package json

import (
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/akyoto/stringutils/unsafe"
)

const (
	readBufferSize      = 4096
	stringBufferSize    = 4096
	maxStringBufferSize = 16384
)

var decoderPool = sync.Pool{
	New: func() interface{} {
		return &decoder{
			buffer:       make([]byte, readBufferSize),
			strings:      make([]byte, stringBufferSize),
			stringsSlice: make([]string, 0, 32),
			types:        make(map[reflect.Type]fieldIndexMap),
		}
	},
}

type stack struct {
	value reflect.Value
	keys  fieldIndexMap
}

type decoder struct {
	// Initialized once
	reader        io.Reader
	buffer        []byte
	strings       []byte
	stringsLength int
	stringsSlice  []string
	types         map[reflect.Type]fieldIndexMap
	stack         [16]stack

	// Initialized on every Decode call
	stackDepth    int
	currentStack  *stack
	field         reflect.Value
	fieldKind     reflect.Kind
	fieldIndex    int
	fieldExists   bool
	currentNumber int64
	stringStart   int
	numbersStart  int
	commaPosition int
	divideFloatBy int
	arrayIndex    int
}

// NewDecoder creates a new JSON decoder.
func NewDecoder(reader io.Reader) *decoder {
	poolObj := decoderPool.Get()
	decoder := poolObj.(*decoder)
	decoder.reader = reader
	return decoder
}

// Reset resets the iterator state
func (decoder *decoder) Reset(object interface{}) {
	decoder.stackDepth = 0
	decoder.currentStack = &decoder.stack[0]
	decoder.currentStack.value = reflect.ValueOf(object).Elem()
	decoder.currentStack.keys = decoder.fieldIndexMap(decoder.currentStack.value.Type())
	decoder.field = reflect.Value{}
	decoder.fieldKind = reflect.Invalid
	decoder.stringStart = -1
	decoder.numbersStart = -1
	decoder.commaPosition = -1
	decoder.divideFloatBy = 1
	decoder.arrayIndex = -1
}

// Decode deserializes the JSON data into the given object.
func (decoder *decoder) Decode(object interface{}) error {
	decoder.Reset(object)

	var (
		i int
		c byte
	)

	for {
		n, err := decoder.reader.Read(decoder.buffer)

		for i = 0; i < n; i++ {
			c = decoder.buffer[i]

			// String capture
			if decoder.stringStart >= 0 {
				for c != '"' {
					i++

					if i >= n {
						goto end
					}

					c = decoder.buffer[i]
				}

				captured := decoder.buffer[decoder.stringStart:i]

				if decoder.fieldExists {
					length := len(captured)

					if decoder.stringsLength+length > len(decoder.strings) {
						newBufferLength := len(decoder.strings) * 2

						if newBufferLength > maxStringBufferSize {
							newBufferLength = maxStringBufferSize
						}

						if newBufferLength < length {
							newBufferLength = length
						}

						decoder.strings = make([]byte, newBufferLength)
						decoder.stringsLength = 0
					}

					destination := decoder.strings[decoder.stringsLength : decoder.stringsLength+length]
					copy(destination, captured)
					decoder.stringsLength += length

					if decoder.arrayIndex >= 0 {
						decoder.stringsSlice = append(decoder.stringsSlice, unsafe.BytesToString(destination))
						decoder.arrayIndex++
					} else {
						decoder.field.SetString(unsafe.BytesToString(destination))
						decoder.fieldExists = false
					}
				} else {
					decoder.fieldIndex, decoder.fieldExists = decoder.currentStack.keys[string(captured)]

					if !decoder.fieldExists {
						return fmt.Errorf("Field does not exist: %s", string(captured))
					}

					decoder.field = decoder.currentStack.value.Field(decoder.fieldIndex)
					decoder.fieldKind = decoder.field.Kind()
				}

				decoder.stringStart = -1
				continue
			}

			// Number capture
			if decoder.numbersStart >= 0 {
				for c >= '0' && c <= '9' {
					decoder.currentNumber = (decoder.currentNumber * 10) + (int64(c) - '0')

					if decoder.commaPosition >= 0 {
						decoder.divideFloatBy *= 10
					}

					i++

					if i >= n {
						goto end
					}

					c = decoder.buffer[i]
				}

				if c == '.' {
					decoder.commaPosition = i - decoder.numbersStart
					continue
				}

				if c == ',' || c == '}' {
					if decoder.commaPosition >= 0 {
						decoder.field.SetFloat(float64(decoder.currentNumber) / float64(decoder.divideFloatBy))
					} else {
						decoder.field.SetInt(decoder.currentNumber)
					}

					decoder.currentNumber = 0
					decoder.numbersStart = -1
					decoder.commaPosition = -1
					decoder.divideFloatBy = 1
					decoder.fieldExists = false
					continue
				}

				continue
			}

			switch c {
			case '"':
				decoder.stringStart = i + 1

			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				decoder.currentNumber = int64(c) - '0'
				decoder.numbersStart = i

			case '{':
				switch decoder.fieldKind {
				case reflect.Map:

				}

			case '[':
				decoder.arrayIndex = 0

			case ']':
				decoder.arrayIndex = -1
				decoder.fieldExists = false

				if len(decoder.stringsSlice) > 0 {
					tmp := make([]string, len(decoder.stringsSlice))
					copy(tmp, decoder.stringsSlice)
					decoder.field.Set(reflect.ValueOf(tmp))
					decoder.stringsSlice = decoder.stringsSlice[:0]
				}
			}
		}

	end:
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
