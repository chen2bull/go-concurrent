package hash

import (
	"fmt"
	"math/bits"
	"reflect"
	"unsafe"
)

type typeFuncs struct {
	hash                  func(unsafe.Pointer, uintptr) uintptr
	equal                 func(unsafe.Pointer, unsafe.Pointer) bool
	getInterfaceValueAddr func(v interface{}) unsafe.Pointer
}

func takeInt32Addr(v interface{}) unsafe.Pointer {
	param, ok := v.(int32)
	if !ok {
		panic(fmt.Sprintf("key %v's type is not int32, is %T", v, v))
	}
	return unsafe.Pointer(&param)
}

func takeUint32Addr(v interface{}) unsafe.Pointer {
	param, ok := v.(uint32)
	if !ok {
		panic(fmt.Sprintf("key %v's type is not uint32, is %T", v, v))
	}
	return unsafe.Pointer(&param)
}

func takeInt64Addr(v interface{})unsafe.Pointer  {
	param, ok := v.(int64)
	if !ok {
		panic(fmt.Sprintf("key %v's type is not int64, is %T", v, v))
	}
	return unsafe.Pointer(&param)
}

func takeUint64Addr(v interface{}) unsafe.Pointer {
	param, ok := v.(uint64)
	if !ok {
		panic(fmt.Sprintf("key %v's type is not uint64, is %T", v, v))
	}
	return unsafe.Pointer(&param)
}

func takeIntAddr(v interface{})unsafe.Pointer  {
	param, ok := v.(int)
	if !ok {
		panic(fmt.Sprintf("key %v's type is not int, is %T", v, v))
	}
	return unsafe.Pointer(&param)
}

func takeUintAddr(v interface{}) unsafe.Pointer {
	param, ok := v.(uint)
	if !ok {
		panic(fmt.Sprintf("key %v's type is not uint, is %T", v, v))
	}
	return unsafe.Pointer(&param)
}

func takeStringAddr(v interface{})unsafe.Pointer  {
	param, ok := v.(string)
	if !ok {
		panic(fmt.Sprintf("key %v's type is not string, is %T", v, v))
	}
	return unsafe.Pointer(&param)
}

func takeFloat32Addr(v interface{})unsafe.Pointer  {
	param, ok := v.(float32)
	if !ok {
		panic(fmt.Sprintf("key %v's type is not float32, is %T", v, v))
	}
	return unsafe.Pointer(&param)
}

func takeFloat64Addr(v interface{}) unsafe.Pointer {
	param, ok := v.(float64)
	if !ok {
		panic(fmt.Sprintf("key %v's type is not float64, is %T", v, v))
	}
	return unsafe.Pointer(&param)
}

func takeComplex64Addr(v interface{}) unsafe.Pointer {
	param, ok := v.(complex64)
	if !ok {
		panic(fmt.Sprintf("key %v's type is not complex64, is %T", v, v))
	}
	return unsafe.Pointer(&param)
}

func takeComplex128Addr(v interface{}) unsafe.Pointer {
	param, ok := v.(complex128)
	if !ok {
		panic(fmt.Sprintf("key %v's type is not complex128, is %T", v, v))
	}
	return unsafe.Pointer(&param)
}

const lenOfKeyTypeMap = reflect.UnsafePointer + 100 // UnsafePointer + 1就已经放得下，避免后续go库扩展了
var keyTypeMap [lenOfKeyTypeMap] typeFuncs
var isKeyTypeInit = false

func init() {
	keyTypeMap = [lenOfKeyTypeMap]typeFuncs{
		reflect.Int32:{algarray[alg_MEM32].hash, algarray[alg_MEM32].equal, takeInt32Addr},
		reflect.Uint32:{algarray[alg_MEM32].hash, algarray[alg_MEM32].equal, takeUint32Addr},
		reflect.Int64:{algarray[alg_MEM64].hash, algarray[alg_MEM64].equal, takeInt64Addr},
		reflect.Uint64:{algarray[alg_MEM64].hash, algarray[alg_MEM64].equal, takeUint64Addr},
		reflect.Int:{algarray[alg_MEM64].hash, algarray[alg_MEM64].equal, takeIntAddr},
		reflect.Uint:{algarray[alg_MEM64].hash, algarray[alg_MEM64].equal, takeUintAddr},
		reflect.String:{algarray[alg_STRING].hash, algarray[alg_STRING].equal, takeStringAddr},
		reflect.Float32:{algarray[alg_FLOAT32].hash, algarray[alg_FLOAT32].equal, takeFloat32Addr},
		reflect.Float64:{algarray[alg_FLOAT64].hash, algarray[alg_FLOAT64].equal, takeFloat64Addr},
		reflect.Complex64:{algarray[alg_CPLX64].hash, algarray[alg_CPLX64].equal, takeComplex64Addr},
		reflect.Complex128:{algarray[alg_CPLX128].hash, algarray[alg_CPLX128].equal, takeComplex128Addr},
	}
	if bits.UintSize == 32 {
		keyTypeMap[reflect.Int] = typeFuncs{algarray[alg_MEM32].hash, algarray[alg_MEM32].equal, takeIntAddr}
		keyTypeMap[reflect.Uint] = typeFuncs{algarray[alg_MEM32].hash, algarray[alg_MEM32].equal, takeUintAddr}
	}
	isKeyTypeInit = true
}
