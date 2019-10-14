package json

import (
	"io"
	"reflect"
	"sync"

	"github.com/akyoto/mirror"
)

const readBufferSize = 4096

var decoderPool = sync.Pool{
	New: func() interface{} {
		return &decoder{
			buffer:       make([]byte, readBufferSize),
			stringsSlice: make([]string, 0, 32),
		}
	},
}

type decoderState struct {
	object      interface{}
	typ         mirror.Type
	kind        reflect.Kind
	fieldName   string
	fieldExists bool
}

func (state *decoderState) SetField(value interface{}) {
	state.typ.(mirror.StructType).SetFieldByJSONTag(state.object, state.fieldName, value)
}

type decoder struct {
	// Initialized once
	reader       io.Reader
	buffer       []byte
	stringsSlice []string
	states       [16]decoderState

	// Initialized on every Decode call
	stackDepth    int
	state         *decoderState
	currentNumber int64
	stringStart   int
	numbersStart  int
	commaPosition int
	divideFloatBy int
	isNegative    bool
	arrayIndex    int
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
	decoder.reset(object)

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

				switch decoder.state.kind {
				case reflect.Struct:
					if decoder.state.fieldExists {
						if decoder.arrayIndex >= 0 {
							decoder.stringsSlice = append(decoder.stringsSlice, string(captured))
							decoder.arrayIndex++
						} else {
							decoder.state.SetField(string(captured))
							decoder.state.fieldExists = false
						}
					} else {
						decoder.state.fieldName = string(captured)
						decoder.state.fieldExists = true
					}

				case reflect.Map:
					if decoder.state.fieldExists {
						// decoder.state.value.SetMapIndex(decoder.state.field, reflect.ValueOf(string(captured)))
						decoder.state.fieldExists = false
					} else {
						// decoder.state.field = reflect.ValueOf(string(captured))
						decoder.state.fieldExists = true
					}
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
					if decoder.isNegative {
						decoder.currentNumber = -decoder.currentNumber
						decoder.isNegative = false
					}

					if decoder.commaPosition >= 0 {
						result := float64(decoder.currentNumber) / float64(decoder.divideFloatBy)
						decoder.state.SetField(&result)
					} else {
						decoder.state.SetField(&decoder.currentNumber)
					}

					decoder.currentNumber = 0
					decoder.numbersStart = -1
					decoder.commaPosition = -1
					decoder.divideFloatBy = 1
					decoder.state.fieldExists = false
					continue
				}

				continue
			}

			// Based on the character we encounter, adjust the state
			switch c {
			case '"':
				decoder.stringStart = i + 1

			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				decoder.currentNumber = int64(c) - '0'
				decoder.numbersStart = i

			case '-':
				decoder.currentNumber = 0
				decoder.numbersStart = i + 1
				decoder.isNegative = true

			case '[':
				decoder.arrayIndex = 0

			case ']':
				decoder.arrayIndex = -1
				decoder.state.fieldExists = false

				if len(decoder.stringsSlice) > 0 {
					tmp := make([]string, len(decoder.stringsSlice))
					copy(tmp, decoder.stringsSlice)
					decoder.state.SetField(tmp)
					decoder.stringsSlice = decoder.stringsSlice[:0]
				}

			case '{':
				if !decoder.state.fieldExists {
					continue
				}

				// switch decoder.state.field.Kind() {
				// case reflect.Map:
				// 	object := reflect.MakeMap(decoder.state.field.Type())
				// 	decoder.state.SetField(object)
				// 	decoder.state.fieldExists = false
				// 	decoder.push(object)
				// }

			case '}':
				decoder.pop()

			case 't':
				i += len("rue")

				if decoder.state.fieldExists {
					decoder.state.SetField(true)
					decoder.state.fieldExists = false
				}

			case 'f':
				i += len("alse")

				if decoder.state.fieldExists {
					decoder.state.SetField(false)
					decoder.state.fieldExists = false
				}

			case 'n':
				i += len("ull")

				if decoder.state.fieldExists {
					// decoder.state.SetField(nil)
					decoder.state.fieldExists = false
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

// reset resets the iterator state.
func (decoder *decoder) reset(object interface{}) {
	decoder.stackDepth = -1
	decoder.push(object)

	decoder.stringStart = -1
	decoder.numbersStart = -1
	decoder.commaPosition = -1
	decoder.divideFloatBy = 1
	decoder.arrayIndex = -1
	decoder.isNegative = false
}

// push creates a new element on the stack.
func (decoder *decoder) push(object interface{}) {
	decoder.stackDepth++
	decoder.state = &decoder.states[decoder.stackDepth]
	decoder.state.object = object
	decoder.state.typ = mirror.TypeOf(object)

	if decoder.state.typ.Kind() == reflect.Ptr {
		decoder.state.typ = decoder.state.typ.(mirror.PointerType).Elem()
	}

	decoder.state.kind = decoder.state.typ.Kind()
	decoder.state.fieldExists = false
}

// pop removes the last element on the stack.
func (decoder *decoder) pop() {
	if decoder.stackDepth == 0 {
		return
	}

	decoder.stackDepth--
	decoder.state = &decoder.states[decoder.stackDepth]
}
