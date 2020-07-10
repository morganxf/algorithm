package binaryheap

type Iterator struct {
	heap  *Heap
	index int
}

func (heap *Heap) Iterator() Iterator {
	return Iterator{heap: heap, index: -1}
}

func (it *Iterator) Next() bool {
	if it.index < it.heap.Size() {
		it.index++
	}
	return it.heap.withinRange(it.index)
}

func (it *Iterator) Prev() bool {
	if it.index >= 0 {
		it.index--
	}
	return it.heap.withinRange(it.index)
}

func (it *Iterator) Value() interface{} {
	v, _ := it.heap.list.Get(it.index)
	return v
}

func (it *Iterator) Index() int {
	return it.index
}

func (it *Iterator) Begin() {
	it.index = -1
}

func (it *Iterator) End() {
	it.index = it.heap.Size()
}

func (it *Iterator) First() bool {
	it.Begin()
	return it.Next()
}

func (it *Iterator) Last() bool {
	it.End()
	return it.Prev()
}
