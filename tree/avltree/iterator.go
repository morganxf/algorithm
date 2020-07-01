package avltree

type Iterator struct {
	tree     *Tree
	node     *Node
	position position
}

type position byte

const (
	begin, between, end = 0, 1, 2
)

func (t *Tree) Iterator() *Iterator {
	return &Iterator{tree: t, node: nil, position: begin}
}

func (it *Iterator) Next() bool {
	switch it.position {
	case begin:
		it.node = it.tree.Left()
		it.position = between
	case between:
		it.node = it.node.Next()
	}
	if it.node == nil {
		it.position = end
		return false
	}
	return true
}

func (it *Iterator) Prev() bool {
	switch it.position {
	case end:
		it.node = it.tree.Right()
		it.position = between
	case between:
		it.node = it.node.Prev()
	}
	if it.node == nil {
		it.position = begin
		return false
	}
	return true
}

func (it *Iterator) Key() interface{} {
	if it.node == nil {
		return nil
	}
	return it.node.Key
}

func (it *Iterator) Value() interface{} {
	if it.node == nil {
		return nil
	}
	return it.node.Value
}

func (it *Iterator) Begin() {
	it.node = nil
	it.position = begin
}

func (it *Iterator) End() {
	it.node = nil
	it.position = end
}

func (it *Iterator) First() bool {
	it.Begin()
	return it.Next()
}

func (it *Iterator) Last() bool {
	it.End()
	return it.Prev()
}
