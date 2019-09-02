package json

import (
	"io"
	"reflect"
	"strings"
)

type fieldIndex = map[string]int

type decoder struct {
	reader io.Reader
	buffer []byte
	types  map[reflect.Type]fieldIndex
}

// NewDecoder creates a new JSON decoder.
func NewDecoder(reader io.Reader) *decoder {
	return &decoder{
		reader: reader,
		buffer: make([]byte, 4096),
		types:  map[reflect.Type]fieldIndex{},
	}
}

// Decode deserializes the JSON data into the given object.
func (decoder *decoder) Decode(object interface{}) error {
	v := reflect.ValueOf(object)
	v = v.Elem()
	t := v.Type()
	decoder.fieldIndex(t)

	for {
		n, err := decoder.reader.Read(decoder.buffer)

		for i := 0; i < n; i++ {
			c := decoder.buffer[i]

			switch c {
			case '"':
			case '{':
			case '}':
			}
		}

		if err == io.EOF {
			return nil
		}
	}
}

// fieldIndex returns a map of field names mapped to their index.
func (decoder *decoder) fieldIndex(structType reflect.Type) fieldIndex {
	fieldsObj, exists := decoder.types[structType]

	if exists {
		return fieldsObj
	}

	fields := fieldIndex{}

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
