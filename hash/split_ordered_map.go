package hash

const (
	KeyTypeInt32 = iota
	KeyTypeUint32
	KeyTypeInt64
	KeyTypeUint64
	KeyTypeString
	KeyTypeFloat32
	KeyTypeFloat64
	KeyTypeCPLX64
	KeyTypeCPLX128
	lenOfKeyTypeEnum
)

type SplitOrderedMap struct {
	t *typeFuncs
	seed int32
	size int64
}


