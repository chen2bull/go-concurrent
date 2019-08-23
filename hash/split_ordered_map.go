package hash

const (
	KeyTypeInt32 = iota
	KeyTypeInt64
	KeyTypeString
	KeyTypeFloat32
	KeyTypeFloat64
	KeyTypeCPLX64
	KeyTypeCPLX128
	lenOfKeyTypeEnum
)

type SplitOrderedMap struct {
	t *typeFuncs

}


