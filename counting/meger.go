package counting

type Merger struct {
	half  []*Merger
	layer []*Balancer
	size  int
}

func NewMerger(size int) *Merger {
	if !isPowerOfTwo(size) {
		panic("size is not power of 2")
	}
	halfSize := size / 2
	layer := make([]*Balancer, halfSize)
	for i := 0; i < halfSize; i++ {
		layer[i] = NewBlancer()
	}
	if size > 2 {
		// 根据图12-12 左边的逻辑图,Merger[2k] 由两个Merger[k] 后加一个balancer组成
		half := []*Merger{NewMerger(halfSize), NewMerger(halfSize)}
		return &Merger{half: half, layer: layer}
	}
	return &Merger{layer: layer}
}

func (m *Merger) Traverse(input int) int {
	output := 0
	if m.size > 2 {
		output = m.half[input%2].Traverse(input / 2)
	}
	return output + m.layer[output].Traverse()
}
