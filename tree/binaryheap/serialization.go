package binaryheap

func (heap *Heap) ToJSON() ([]byte, error) {
	return heap.list.ToJSON()
}

func (heap *Heap) FromJSON(data []byte) error {
	return heap.list.FromJSON(data)
}
