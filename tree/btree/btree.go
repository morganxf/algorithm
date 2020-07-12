package btree

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/morganxf/algorithm/util"
)

type Tree struct {
	Root       *Node
	Comparator util.Comparator
	size       int
	m          int
}

func NewWith(order int, comparator util.Comparator) *Tree {
	if order < 3 {
		panic("Invalid order, should be at least 3")
	}
	return &Tree{m: order, Comparator: comparator}
}

func NewWithIntComparator(order int) *Tree {
	return NewWith(order, util.IntComparator)
}

func NewWithStringComparator(order int) *Tree {
	return NewWith(order, util.StringComparator)
}

type Node struct {
	Parent   *Node
	Entries  []*Entry
	Children []*Node
}

func (n *Node) height() int {
	height := 0
	for ; n != nil; n = n.Children[0] {
		height++
		if len(n.Children) == 0 {
			break
		}
	}
	return height
}

type Entry struct {
	Key   interface{}
	Value interface{}
}

func (entry *Entry) String() string {
	return fmt.Sprintf("%v", entry.Key)
}

func (t *Tree) Put(key interface{}, value interface{}) {
	entry := &Entry{Key: key, Value: value}
	if t.Root == nil {
		t.Root = &Node{Entries: []*Entry{entry}, Children: []*Node{}}
		t.size++
		return
	}
	if t.insert(t.Root, entry) {
		t.size++
	}
}

func (t *Tree) Remove(key interface{}) {
	node, index, found := t.searchRecursively(t.Root, key)
	if found {
		t.delete(node, index)
		t.size--
	}
}

func (t *Tree) Get(key interface{}) (value interface{}, found bool) {
	node, index, found := t.searchRecursively(t.Root, key)
	if found {
		return node.Entries[index].Value, true
	}
	return nil, false
}

func (t *Tree) Empty() bool {
	return t.size == 0
}

func (t *Tree) Size() int {
	return t.size
}

func (t *Tree) Keys() []interface{} {
	keys := make([]interface{}, t.size)
	it := t.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

func (t *Tree) Values() []interface{} {
	values := make([]interface{}, t.size)
	it := t.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
}

func (t *Tree) Clear() {
	t.Root = nil
	t.size = 0
}

func (t *Tree) Height() int {
	return t.Root.height()
}

// 返回最左node
func (t *Tree) Left() *Node {
	return t.left(t.Root)
}

func (t *Tree) LeftKey() interface{} {
	if left := t.Left(); left != nil {
		return left.Entries[0].Key
	}
	return nil
}

func (t *Tree) LeftValue() interface{} {
	if left := t.Left(); left != nil {
		return left.Entries[0].Value
	}
	return nil
}

func (t *Tree) Right() *Node {
	return t.right(t.Root)
}

func (t *Tree) RightKey() interface{} {
	if right := t.Right(); right != nil {
		return right.Entries[len(right.Entries)-1].Key
	}
	return nil
}

func (t *Tree) RightValue() interface{} {
	if right := t.Right(); right != nil {
		return right.Entries[len(right.Entries)-1].Value
	}
	return nil
}

func (t *Tree) String() string {
	var buffer bytes.Buffer
	if _, err := buffer.WriteString("BTree\n"); err != nil {
	}
	if !t.Empty() {
		t.output(&buffer, t.Root, 0, true)
	}
	return buffer.String()
}

func (t *Tree) output(buffer *bytes.Buffer, node *Node, level int, isTail bool) {
	for e := 0; e < len(node.Entries)+1; e++ {
		if e < len(node.Children) {
			t.output(buffer, node.Children[e], level+1, true)
		}
		if e < len(node.Entries) {
			buffer.WriteString(strings.Repeat("    ", level))
			buffer.WriteString(fmt.Sprintf("%v", node.Entries[e].Key) + "\n")
		}
	}
}

func (t *Tree) left(n *Node) *Node {
	if t.Empty() {
		return nil
	}
	current := n
	for {
		if t.isLeaf(current) {
			return current
		}
		current = current.Children[0]
	}
}

func (t *Tree) right(n *Node) *Node {
	if t.Empty() {
		return nil
	}
	current := n
	for {
		if t.isLeaf(current) {
			return current
		}
		current = current.Children[len(current.Children)-1]
	}
}

func (t *Tree) maxChildren() int {
	return t.m
}

func (t *Tree) minChildren() int {
	// ceil(m/2)
	return (t.m + 1) / 2
}

// entriesNum + 1 = childrenNum
func (t *Tree) maxEntries() int {
	return t.maxChildren() - 1
}

func (t *Tree) minEntries() int {
	return t.minChildren() - 1
}

func (t *Tree) middle() int {
	return (t.m - 1) / 2
}

// insert只能在叶子节点插入
// insert是从上向下，直到叶子节点找到合适的位置，然后插入;
// split是从下向上，直到parent不再需要split
func (t *Tree) insert(node *Node, entry *Entry) (inserted bool) {
	if t.isLeaf(node) {
		return t.insertIntoLeaf(node, entry)
	}
	return t.insertIntoInternal(node, entry)
}

func (t *Tree) insertIntoLeaf(node *Node, entry *Entry) bool {
	insertPosition, found := t.search(node, entry.Key)
	// 找打同Key的元素
	if found {
		node.Entries[insertPosition] = entry
		return false
	}
	// insertPosition 整体后移. 采用链表是否会好一些 但是链表无法二分查找
	node.Entries = append(node.Entries, nil)
	copy(node.Entries[insertPosition+1:], node.Entries[insertPosition:])
	// 插入新entry
	node.Entries[insertPosition] = entry
	// 调整树节点的Entry数目
	t.split(node)
	return true
}

func (t *Tree) insertIntoInternal(node *Node, entry *Entry) (inserted bool) {
	insertPosition, found := t.search(node, entry.Key)
	if found {
		node.Entries[insertPosition] = entry
		return false
	}
	// 递归，直到叶子节点
	return t.insert(node.Children[insertPosition], entry)
}

// split是从下向上处理
func (t *Tree) split(node *Node) {
	if !t.shouldSplit(node) {
		return
	}
	if node == t.Root {
		t.splitRoot(node)
		return
	}
	t.splitNonRoot(node)
}

func (t *Tree) splitRoot(node *Node) {
	middle := t.middle()
	// split
	left := &Node{Entries: append([]*Entry(nil), t.Root.Entries[:middle]...)}
	right := &Node{Entries: append([]*Entry(nil), t.Root.Entries[middle+1:]...)}
	// 不是叶子节点， 孩子们需要split
	// 孩子们split 并分属新left right节点的孩子 顺序不变。left right内相对位置也没有变化
	if !t.isLeaf(t.Root) {
		left.Children = append([]*Node(nil), t.Root.Children[:middle+1]...)
		right.Children = append([]*Node(nil), t.Root.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}
	newRoot := &Node{
		Entries:  []*Entry{t.Root.Entries[middle]},
		Children: []*Node{left, right},
	}
	left.Parent = newRoot
	right.Parent = newRoot
	t.Root = newRoot
}

func (t *Tree) splitNonRoot(node *Node) {
	// 找到中间entry，取出，插入到父亲节点
	middle := t.middle()
	parent := node.Parent

	// 当前node的entries以middle为为中间 split左右node 作为父亲的新孩子
	left := &Node{Entries: append([]*Entry(nil), node.Entries[:middle]...), Parent: parent}
	right := &Node{Entries: append([]*Entry(nil), node.Entries[middle+1:]...), Parent: parent}
	// 如果不是叶子节点需要处理孩子
	if !t.isLeaf(node) {
		left.Children = append([]*Node(nil), node.Children[:middle+1]...)
		right.Children = append([]*Node(nil), node.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}
	// 计算middle entry在parent中得位置
	// 此时不会有重复的key，如果存在不会进入此逻辑
	insertPosition, _ := t.search(parent, node.Entries[middle].Key)
	// insertPosition整体后移，插入
	parent.Entries = append(parent.Entries, nil)
	copy(parent.Entries[insertPosition+1:], parent.Entries[insertPosition:])
	parent.Entries[insertPosition] = node.Entries[middle]
	// parent insertPosition原为node，现left替换原node
	parent.Children[insertPosition] = left
	// 由于原node 一份为二，parent.Children需要增加空间 且right的位置应该是insertPositon+1
	parent.Children = append(parent.Children, nil)
	copy(parent.Children[insertPosition+2:], parent.Children[insertPosition+1:])
	parent.Children[insertPosition+1] = right
	// 递归split
	// parent增加了一个node，需要更新当前的node为parent，看parent是否满足btree的条件，如不满足则split
	t.split(parent)
}

// 二分，搜索插入的位置
func (t *Tree) search(node *Node, key interface{}) (index int, found bool) {
	low, high := 0, len(node.Entries)-1
	var mid int
	// 因为childrenNum可以比entriesNum多一个，所以low<=high中的等号是很重要的
	// 当low==high时，如果key>node.Entries[low]时，return low==high+1
	// 当low==high时。如果key<node.Entries[low]时，return low==high
	for low <= high {
		// mid 偏向小的值 由于整数除法的原因
		mid = (low + high) / 2
		compare := t.Comparator(key, node.Entries[mid].Key)
		switch {
		case compare > 0:
			low = mid + 1
		case compare < 0:
			high = mid - 1
		case compare == 0:
			// 发现key已经存在
			return mid, true
		}
	}
	return low, false
}

func (t *Tree) isLeaf(node *Node) bool {
	return len(node.Children) == 0
}

func (t *Tree) shouldSplit(node *Node) bool {
	return len(node.Entries) > t.maxEntries()
}

func (t *Tree) searchRecursively(startNode *Node, key interface{}) (node *Node, index int, found bool) {
	if t.Empty() {
		return nil, -1, false
	}
	node = startNode
	for {
		index, found = t.search(node, key)
		if found {
			return node, index, true
		}
		if t.isLeaf(node) {
			return nil, -1, false
		}
		node = node.Children[index]
	}
}

func (t *Tree) delete(node *Node, index int) {
	// delete from a leaf node
	if t.isLeaf(node) {
		deletedKey := node.Entries[index].Key
		t.deleteEntry(node, index)
		t.rebalance(node, deletedKey)
		if len(t.Root.Entries) == 0 {
			t.Root = nil
		}
		return
	}
	// delete from an internal node
	leftLargestNode := t.right(node.Children[index])
	leftLargestEntryIndex := len(leftLargestNode.Entries) - 1
	// 用该entry左孩子的most right entry替换index所在entry, 相当于删除该index原先的值
	node.Entries[index] = leftLargestNode.Entries[leftLargestEntryIndex]
	// 删除替换的entry 且该entry是位于left node。即delete entry from leaf node
	deletedKey := leftLargestNode.Entries[leftLargestEntryIndex].Key
	t.deleteEntry(leftLargestNode, leftLargestEntryIndex)
	t.rebalance(leftLargestNode, deletedKey)
}

func (t *Tree) deleteEntry(node *Node, index int) {
	// delete by index
	copy(node.Entries[index:], node.Entries[index+1:])
	node.Entries[len(node.Entries)-1] = nil
	// update Entries
	node.Entries = node.Entries[:len(node.Entries)-1]
}

func (t *Tree) rebalance(node *Node, deletedKey interface{}) {
	// 不需要做rebalance
	if node == nil || len(node.Entries) >= t.minEntries() {
		return
	}
	// 不满足大于最小EntriesNum条件，需要rebalance

	// 首先尝试从左兄弟借entry
	leftSiblingNode, leftSiblingIndex := t.leftSibling(node, deletedKey)
	// 如果左兄弟的entry数目小于等于minEntry，不执行下面的逻辑
	// 从paren提取entry到当前node，从左兄弟提取最右entry，提升到parent
	if leftSiblingNode != nil && len(leftSiblingNode.Entries) > t.minEntries() {
		// parent提取entry，下方到当前节点，补充entry数目
		node.Entries = append([]*Entry{node.Parent.Entries[leftSiblingIndex]}, node.Entries...)
		// 左兄弟的最右entry提升到paren
		node.Parent.Entries[leftSiblingIndex] = leftSiblingNode.Entries[len(leftSiblingNode.Entries)-1]
		t.deleteEntry(leftSiblingNode, len(leftSiblingNode.Entries)-1)
		// 如果不是叶子，由于提出了最右entry，所以最右孩子需要移到当前node的最左孩子
		if !t.isLeaf(leftSiblingNode) {
			leftSiblingRightMostChild := leftSiblingNode.Children[len(leftSiblingNode.Children)-1]
			leftSiblingRightMostChild.Parent = node
			node.Children = append([]*Node{leftSiblingRightMostChild}, node.Children...)
			t.deleteChild(leftSiblingNode, len(leftSiblingNode.Children)-1)
		}
		// parent的entry数目没有变化，不用再递归rebalance parent
		return
	}

	// 尝试从右兄弟借entry。当右兄弟entry数据大于minEntries
	rightSiblingNode, rightSiblingIndex := t.rightSibling(node, deletedKey)
	if rightSiblingNode != nil && len(rightSiblingNode.Entries) > t.minEntries() {
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		node.Parent.Entries[rightSiblingIndex-1] = rightSiblingNode.Entries[0]
		t.deleteEntry(rightSiblingNode, 0)
		if !t.isLeaf(rightSiblingNode) {
			rightSiblingLeftMostChild := rightSiblingNode.Children[0]
			rightSiblingLeftMostChild.Parent = node
			node.Children = append(node.Children, rightSiblingLeftMostChild)
			t.deleteChild(rightSiblingNode, 0)
		}
		// parent的entry数目没有变化，不用再递归rebalance parent
		return
	}

	// merge两个节点会导致parent的entry数目发生变化，需要递归rebalance parent
	// merge两个节点。当左右兄弟的Entry数据都小于等于minEntries
	if rightSiblingNode != nil {
		// 合并两个node需要以对应parent的entry作为中间连接，即node.Parent.Entries[rightSiblingIndex-1]
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		// 合并右兄弟entries
		node.Entries = append(node.Entries, rightSiblingNode.Entries...)
		// ?
		deletedKey = node.Parent.Entries[rightSiblingIndex-1].Key
		t.deleteEntry(node.Parent, rightSiblingIndex-1)
		// 右兄弟的孩子append到当前节点
		t.appendChildren(node.Parent.Children[rightSiblingIndex], node)
		// 右兄弟从parent移除
		t.deleteChild(node.Parent, rightSiblingIndex)
	} else if leftSiblingNode != nil {
		// 与左兄弟合并到当前node
		entries := append([]*Entry(nil), leftSiblingNode.Entries...)
		entries = append(entries, node.Parent.Entries[leftSiblingIndex])
		node.Entries = append(entries, node.Entries...)
		// ?
		deletedKey = node.Parent.Entries[leftSiblingIndex].Key
		t.deleteEntry(node.Parent, leftSiblingIndex)
		t.prependChildren(node.Parent.Children[leftSiblingIndex], node)
		t.deleteChild(node.Parent, leftSiblingIndex)
	}
	// 如果左右兄弟都是nil，不做任何处理？

	// 如果parent已经是root 且root的entry数目变为空
	if node.Parent == t.Root && len(t.Root.Entries) == 0 {
		t.Root = node
		node.Parent = nil
		return
	}

	// 递归rebalance当前节点的parent
	t.rebalance(node.Parent, deletedKey)
}

// leftSibling node
func (t *Tree) leftSibling(node *Node, key interface{}) (*Node, int) {
	if node.Parent != nil {
		// 一定会发现，如果not found就不会进入这个逻辑
		index, _ := t.search(node.Parent, key)
		// left sibling node index
		index--
		// index是否合法
		if index >= 0 && index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}
	return nil, -1
}

func (t *Tree) rightSibling(node *Node, key interface{}) (*Node, int) {
	if node.Parent != nil {
		index, _ := t.search(node.Parent, key)
		index++
		if index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}
	return nil, -1
}

func (t *Tree) deleteChild(node *Node, index int) {
	if index >= len(node.Children) {
		return
	}
	copy(node.Children[index:], node.Children[index+1:])
	node.Children[len(node.Children)-1] = nil
	node.Children = node.Children[:len(node.Children)-1]
}

func (t *Tree) prependChildren(fromNode *Node, toNode *Node) {
	children := append([]*Node(nil), fromNode.Children...)
	toNode.Children = append(children, toNode.Children...)
	setParent(fromNode.Children, toNode)
}
func (t *Tree) appendChildren(fromNode *Node, toNode *Node) {
	toNode.Children = append(toNode.Children, fromNode.Children...)
	setParent(fromNode.Children, toNode)
}

func setParent(nodes []*Node, parent *Node) {
	for _, node := range nodes {
		node.Parent = parent
	}
}
