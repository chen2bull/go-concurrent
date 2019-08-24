package hash

import (
	"fmt"
	"unsafe"
)

type typeFuncs struct {
	hash   func(unsafe.Pointer, uintptr) uintptr
	equal  func(unsafe.Pointer, unsafe.Pointer) bool
	assert func(v interface{})
}

func assertInt32(v interface{}) {
	_, ok := v.(int32)
	if !ok {
		panic(fmt.Sprintf("key %v is not int32", v))
	}
}

func assertUint32(v interface{}) {
	_, ok := v.(uint32)
	if !ok {
		panic(fmt.Sprintf("key %v is not uint32", v))
	}
}

func assertInt64(v interface{}) {
	_, ok := v.(int64)
	if !ok {
		panic(fmt.Sprintf("key %v is not int64", v))
	}
}

func assertUint64(v interface{}) {
	_, ok := v.(uint64)
	if !ok {
		panic(fmt.Sprintf("key %v is not uint64", v))
	}
}

func assertString(v interface{}) {
	_, ok := v.(string)
	if !ok {
		panic(fmt.Sprintf("key %v is not string", v))
	}
}

func assertFloat32(v interface{}) {
	_, ok := v.(float32)
	if !ok {
		panic(fmt.Sprintf("key %v is not float32", v))
	}
}

func assertFloat64(v interface{}) {
	_, ok := v.(float64)
	if !ok {
		panic(fmt.Sprintf("key %v is not float64", v))
	}
}

func assertComplex64(v interface{}) {
	_, ok := v.(complex64)
	if !ok {
		panic(fmt.Sprintf("key %v is not complex64", v))
	}
}

func assertComplex128(v interface{}) {
	_, ok := v.(complex128)
	if !ok {
		panic(fmt.Sprintf("key %v is not complex128", v))
	}
}

var keyTypeMap [lenOfKeyTypeEnum] typeFuncs

func init() {
	keyTypeMap = [lenOfKeyTypeEnum]typeFuncs{
		KeyTypeInt32:   {algarray[alg_MEM32].hash, algarray[alg_MEM32].equal, assertInt32},
		KeyTypeUint32:   {algarray[alg_MEM32].hash, algarray[alg_MEM32].equal, assertUint32},
		KeyTypeInt64:   {algarray[alg_MEM64].hash, algarray[alg_MEM64].equal, assertInt64},
		KeyTypeUint64:   {algarray[alg_MEM64].hash, algarray[alg_MEM64].equal, assertUint64},
		KeyTypeString:  {algarray[alg_STRING].hash, algarray[alg_STRING].equal, assertString},
		KeyTypeFloat32: {algarray[alg_FLOAT32].hash, algarray[alg_FLOAT32].equal, assertFloat32},
		KeyTypeFloat64: {algarray[alg_FLOAT64].hash, algarray[alg_FLOAT64].equal, assertFloat64},
		KeyTypeCPLX64:  {algarray[alg_CPLX64].hash, algarray[alg_CPLX64].equal, assertComplex64},
		KeyTypeCPLX128: {algarray[alg_CPLX128].hash, algarray[alg_CPLX128].equal, assertComplex128},
	}
}
