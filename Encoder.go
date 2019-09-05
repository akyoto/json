package json

// import (
// 	"fmt"
// 	"io"
// 	"reflect"
// 	"strconv"
// )

// var (
// 	bracketStart        = []byte{'{'}
// 	bracketEnd          = []byte{'}'}
// 	stringIndicator     = []byte{'"'}
// 	fieldValueSeparator = []byte{'"', ':'}
// 	comma               = []byte{','}
// )

// // Encode encodes the object to JSON
// // and writes the output to the given writer.
// func Encode(object interface{}, writer io.Writer) error {
// 	return writeValue(writer, reflect.ValueOf(object))
// }

// // writeValue writes the given value in JSON-encoded form to the writer.
// // nolint:errcheck
// func writeValue(writer io.Writer, v reflect.Value) error {
// 	stringWriter := writer.(io.StringWriter)
// 	t := v.Type()

// 	if v.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 		v = v.Elem()
// 	}

// 	switch v.Kind() {
// 	case reflect.String:
// 		writer.Write(stringIndicator)
// 		stringWriter.WriteString(v.String())
// 		writer.Write(stringIndicator)

// 	case reflect.Int:
// 		// TODO: Avoid allocation
// 		x := strconv.FormatInt(v.Int(), 10)
// 		stringWriter.WriteString(x)

// 	case reflect.Float64:
// 		// TODO: Avoid allocation
// 		x := strconv.FormatFloat(v.Float(), 'f', -1, 64)
// 		stringWriter.WriteString(x)

// 	case reflect.Struct:
// 		writer.Write(bracketStart)
// 		fieldCount := t.NumField()

// 		for i := 0; i < fieldCount; i++ {
// 			// Key
// 			tField := t.Field(i)
// 			writer.Write(stringIndicator)
// 			stringWriter.WriteString(tField.Name)
// 			writer.Write(fieldValueSeparator)

// 			// Value
// 			vField := v.Field(i)
// 			err := writeValue(writer, vField)

// 			if err != nil {
// 				return err
// 			}

// 			if i != fieldCount-1 {
// 				writer.Write(comma)
// 			}
// 		}

// 		writer.Write(bracketEnd)

// 	default:
// 		return fmt.Errorf("Can't encode field of kind '%v'", v.Kind())
// 	}

// 	return nil
// }
