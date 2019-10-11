package json

import (
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/akyoto/stringutils/unsafe"
)

const (
	readBufferSize = 4096
)

var decoderPool = sync.Pool{
	New: func() interface{} {
		return &decoder{
			buffer:       make([]byte, readBufferSize),
			stringsSlice: make([]string, 0, 32),
			types:        make(map[reflect.Type]fieldIndexMap),
		}
	},
}

type decoderState struct {
	value       reflect.Value
	keys        fieldIndexMap
	field       reflect.Value
	fieldExists bool
}

type decoder struct {
	// Initialized once
	reader       io.Reader
	buffer       []byte
	stringsSlice []string
	types        map[reflect.Type]fieldIndexMap
	states       [16]decoderState

	// Initialized on every Decode call
	stackDepth    int
	state         *decoderState
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

// reset resets the iterator state.
func (decoder *decoder) reset(object interface{}) {
	decoder.stackDepth = -1
	decoder.push(reflect.ValueOf(object))

	decoder.stringStart = -1
	decoder.numbersStart = -1
	decoder.commaPosition = -1
	decoder.divideFloatBy = 1
	decoder.arrayIndex = -1
}

// push creates a new element on the stack.
func (decoder *decoder) push(value reflect.Value) {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	decoder.stackDepth++
	decoder.state = &decoder.states[decoder.stackDepth]
	decoder.state.value = value
	decoder.state.keys = decoder.fieldIndexMap(value.Type())
	decoder.state.field = reflect.Value{}
	decoder.state.fieldExists = false
}

// pop removes the last element on the stack.
func (decoder *decoder) pop() {
	decoder.stackDepth--
	decoder.state = &decoder.states[decoder.stackDepth]
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

				if decoder.state.fieldExists {
					length := len(captured)
					tmp := make([]byte, length)
					copy(tmp, captured)

					if decoder.arrayIndex >= 0 {
						decoder.stringsSlice = append(decoder.stringsSlice, unsafe.BytesToString(tmp))
						decoder.arrayIndex++
					} else {
						decoder.state.field.SetString(unsafe.BytesToString(tmp))
						decoder.state.fieldExists = false
					}
				} else {
					var fieldIndex int
					fieldIndex, decoder.state.fieldExists = decoder.state.keys[string(captured)]

					if !decoder.state.fieldExists {
						return fmt.Errorf("Field does not exist: %s", string(captured))
					}

					decoder.state.field = decoder.state.value.Field(fieldIndex)
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
						decoder.state.field.SetFloat(float64(decoder.currentNumber) / float64(decoder.divideFloatBy))
					} else {
						decoder.state.field.SetInt(decoder.currentNumber)
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

			switch c {
			case '"':
				decoder.stringStart = i + 1

			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				decoder.currentNumber = int64(c) - '0'
				decoder.numbersStart = i

			case '[':
				decoder.arrayIndex = 0

			case ']':
				decoder.arrayIndex = -1
				decoder.state.fieldExists = false

				if len(decoder.stringsSlice) > 0 {
					tmp := make([]string, len(decoder.stringsSlice))
					copy(tmp, decoder.stringsSlice)
					decoder.state.field.Set(reflect.ValueOf(tmp))
					decoder.stringsSlice = decoder.stringsSlice[:0]
				}

			case '{':
				if !decoder.state.fieldExists {
					continue
				}

				switch decoder.state.field.Kind() {
				case reflect.Map:

				}

			case '}':
				if !decoder.state.fieldExists {
					continue
				}

				decoder.pop()
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
