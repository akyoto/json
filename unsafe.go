package json

import "unsafe"

//go:linkname mapassign reflect.mapassign
func mapassign(rtype unsafe.Pointer, m unsafe.Pointer, key, val unsafe.Pointer)

type eface struct {
	_    unsafe.Pointer
	data unsafe.Pointer
}

func unpackEFace(obj interface{}) *eface {
	return (*eface)(unsafe.Pointer(&obj))
}
