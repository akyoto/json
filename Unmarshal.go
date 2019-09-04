package json

import (
	"bytes"
	"sync"
)

var readerPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewReader(nil)
	},
}

// Unmarshal decodes the JSON into the given object.
func Unmarshal(data []byte, object interface{}) error {
	readerObject := readerPool.Get()
	reader := readerObject.(*bytes.Reader)
	reader.Reset(data)
	decoder := NewDecoder(reader)
	err := decoder.Decode(object)
	decoder.Close()
	readerPool.Put(reader)
	return err
}
