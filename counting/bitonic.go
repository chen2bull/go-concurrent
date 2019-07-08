package counting

type Bitonic struct {
	half  []*Bitonic
	layer *Merger // 图12-14 Bitonic[2k]由两个Bitonic[k]输出连到Merger[2k]组成
	size  int
}

func NewBitonic(size int) *Bitonic {
	if !isPowerOfTwo(size) {
		panic("size is not power of 2")
	}
	layer := NewMerger(size)
	halfSize := size / 2
	if size > 2 {
		half := []*Bitonic{NewBitonic(halfSize), NewBitonic(halfSize)}
		return &Bitonic{half: half, layer: layer}
	}
	return &Bitonic{layer: layer}
}

func (b *Bitonic) Traverse(input int) int {
	output := 0
	if b.size > 2 {
		output = b.half[input%2].Traverse(input / 2)
	}
	return output + b.layer.Traverse(output)
}
