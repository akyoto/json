package json

import (
	"reflect"
	"strings"
)

// fieldIndexMap maps a field name to its position in the struct (index).
type fieldIndexMap = map[string]int

// fieldIndexMap returns a map of field names mapped to their index.
func (decoder *decoder) fieldIndexMap(structType reflect.Type) fieldIndexMap {
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
