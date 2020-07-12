package btree

type Iterator struct {
	tree     *Tree
	node     *Node
	entry    *Entry
	position position
}

type position byte

const (
	begin, between, end position = 0, 1, 2
)

func (t *Tree) Iterator() Iterator {
	return Iterator{tree: t, node: nil, position: begin}
}

func (it *Iterator) Next() bool {
	if it.position == end {
		goto end
	}
	if it.position == begin {
		// most left leaf node
		left := it.tree.Left()
		// 根据插入和删除的rebalance规则，孩子的顺序是从左到右。如果left==nil，则说明没有后续的node
		if left == nil {
			goto end
		}
		it.node = left
		it.entry = left.Entries[0]
		goto between
	}
	// if it.position == between
	{
		e, _ := it.tree.search(it.node, it.entry.Key)
		// 如果当前node存在Children 且当前node的Children没有迭代完。即 e+1 < len(it.node.Children), 更新node为孩子 开始新node的迭代
		if e+1 < len(it.node.Children) {
			// e+1为下一个node, 更新node
			it.node = it.node.Children[e+1]
			// 迭代是按顺序迭代，所以要一直向下找到most left leaf child
			// 更新node and entry为most left leaf child and it's entry
			for len(it.node.Children) > 0 {
				it.node = it.node.Children[0]
			}
			it.entry = it.node.Entries[0]
			goto between
		}
		// 当前node 迭代entry。e是当前entry的index，e+1是下一个entry
		// if e+1 == len(it.node.Entries), 当前node的entry已经迭代完成。next需要迭代当前node的parent
		if e+1 < len(it.node.Entries) {
			it.entry = it.node.Entries[e+1]
			goto between
		}
	}
	// 当前node的entry已经迭代完成，next迭代当前node的parent对应的entry
	for it.node.Parent != nil {
		// 当前node更新为node.Parent
		it.node = it.node.Parent
		// 找到parent对应entry的index
		e, _ := it.tree.search(it.node, it.entry.Key)
		// if e < len(it.node.Entries) 更新entry
		// if e == len(it.node.Entries) 当前node已经迭代结束，需要继续向上找parent
		if e < len(it.node.Entries) {
			// 更新entry为parent对应的entry
			it.entry = it.node.Entries[e]
			goto between
		}
	}

end:
	it.End()
	return false

between:
	it.position = between
	return true
}

func (it *Iterator) Prev() bool {
	if it.position == begin {
		goto begin
	}
	if it.position == end {
		// most right leaf node
		right := it.tree.Right()
		if right == nil {
			goto begin
		}
		it.node = right
		it.entry = right.Entries[len(right.Entries)-1]
		goto between
	}
	// if it.position == between
	{
		e, _ := it.tree.search(it.node, it.entry.Key)
		if e < len(it.node.Children) {
			it.node = it.node.Children[e]
			// update node to the most right leaf node
			for len(it.node.Children) > 0 {
				it.node = it.node.Children[len(it.node.Children)-1]
			}
			it.entry = it.node.Entries[len(it.node.Entries)-1]
			goto between
		}
		// 在当前node迭代entry, e-1为前一个entry
		if e-1 >= 0 {
			it.entry = it.node.Entries[e-1]
			goto between
		}
	}
	// if e-1 < 0, 即e==0. 当前node已经迭代完，需要迭代parent
	for it.node.Parent != nil {
		it.node = it.node.Parent
		e, _ := it.tree.search(it.node, it.entry.Key)
		// if e == 0, 当前node已经迭代完成。继续向上找parent
		if e-1 >= 0 {
			it.entry = it.node.Entries[e-1]
			goto between
		}
	}

begin:
	it.Begin()
	return false

between:
	it.position = between
	return true
}

func (it *Iterator) Key() interface{} {
	return it.entry.Key
}

func (it *Iterator) Value() interface{} {
	return it.entry.Value
}

func (it *Iterator) Begin() {
	it.node = nil
	it.position = begin
	it.entry = nil
}

func (it *Iterator) End() {
	it.node = nil
	it.position = end
	it.entry = nil
}

func (it *Iterator) First() bool {
	it.Begin()
	return it.Next()
}

func (it *Iterator) Last() bool {
	it.End()
	return it.Prev()
}
