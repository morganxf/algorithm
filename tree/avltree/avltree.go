package avltree

import (
	"fmt"

	"github.com/morganxf/algorithm/util"
)

type Tree struct {
	Root       *Node
	Comparator util.Comparator
	size       int
}

type Node struct {
	Key      interface{}
	Value    interface{}
	Parent   *Node
	Children [2]*Node // Left - Right
	balance  int8
}

func NewWith(comparator util.Comparator) *Tree {
	return &Tree{Comparator: comparator}
}

func NewWithIntComparator() *Tree {
	return &Tree{Comparator: util.IntComparator}
}

func NewWithStringComparator() *Tree {
	return &Tree{Comparator: util.StringComparator}
}

func (t *Tree) Put(key interface{}, value interface{}) {
	t.put(key, value, nil, &t.Root)
}

func (t *Tree) Get(key interface{}) (interface{}, bool) {
	cur := t.Root
	for cur != nil {
		cmp := t.Comparator(key, cur.Key)
		switch {
		case cmp == 0:
			return cur.Value, true
		case cmp < 0:
			cur = cur.Children[0]
		case cmp > 0:
			cur = cur.Children[1]
		}
	}
	return nil, false
}

func (t *Tree) Remove(key interface{}) {
	t.remove(key, &t.Root)
}

func (t *Tree) Left() *Node {
	return t.bottom(0)
}

func (t *Tree) Right() *Node {
	return t.bottom(1)
}

func (t *Tree) Empty() bool {
	return t.size == 0
}

func (t *Tree) Size() int {
	return t.size
}

func (t *Tree) Keys() []interface{} {
	it := t.Iterator()
	keys := make([]interface{}, 0, t.size)
	for it.Next() {
		keys = append(keys, it.Key())
	}
	return keys
}

func (t *Tree) Values() []interface{} {
	it := t.Iterator()
	values := make([]interface{}, 0, t.size)
	for it.Next() {
		values = append(values, it.Value())
	}
	return values
}

func (t *Tree) Clear() {
	t.Root = nil
	t.size = 0
}

func (t *Tree) Floor(key interface{}) (*Node, bool) {
	var floor *Node
	var found = false
	cur := t.Root
	for cur != nil {
		cmp := t.Comparator(key, cur.Key)
		switch {
		case cmp < 0:
			cur = cur.Children[0]
		case cmp > 0:
			// floor是比key小的node
			floor, found = cur, true
			// 继续 是为了找到最大的比key小的node
			cur = cur.Children[1]
		case cmp == 0:
			return cur, true
		}
	}
	if found {
		return floor, found
	}
	return nil, false
}

func (t *Tree) Ceiling(key interface{}) (*Node, bool) {
	var ceiling *Node
	var found = false
	cur := t.Root
	for cur != nil {
		cmp := t.Comparator(key, cur.Key)
		switch {
		case cmp < 0:
			ceiling, found = cur, true
			cur = cur.Children[0]
		case cmp > 0:
			cur = cur.Children[1]
		case cmp == 0:
			return cur, true
		}
	}
	if found {
		return ceiling, true
	}
	return nil, false
}

func (t *Tree) String() string {
	str := "Tree\n"
	if !t.Empty() {
		output(t.Root, "", true, &str)
	}
	return str
}

// false: already balanced
func (t *Tree) put(key interface{}, value interface{}, parent *Node, target **Node) bool {
	cur := *target
	if cur == nil {
		t.size++
		*target = &Node{Key: key, Value: value, Parent: parent}
		return true
	}

	cmp := t.Comparator(key, cur.Key)
	if cmp == 0 {
		cur.Key = key
		cur.Value = value
		return false
	}

	var newTarget **Node
	if cmp < 0 {
		newTarget = &cur.Children[0]
	} else {
		newTarget = &cur.Children[1]
	}
	imbalanced := t.put(key, value, cur, newTarget)
	if imbalanced {
		// *target == newTarget.Parent
		// 增加node可能导致祖先的不平衡，重新平衡最小不平衡数
		return putRebalance(int8(cmp), target)
	}
	return false
}

func (t *Tree) remove(key interface{}, target **Node) bool {
	cur := *target
	if cur == nil {
		return false
	}

	cmp := t.Comparator(key, cur.Key)
	// cur是要被remove的node
	if cmp == 0 {
		t.size--
		if cur.Children[1] == nil {
			if cur.Children[0] != nil {
				// 更新左子节点的父节点为当前节点的父节点
				cur.Children[0].Parent = cur.Parent
			}
			// 更新当前节点为当前节点的左子节点，即更新当前节点的父节点的子节点为当前节点的左子节点
			*target = cur.Children[0]
			return true
		}
		// 使用右子树的最小node值替换当前node值
		return removeMin(&cur.Children[1], &cur.Key, &cur.Value)
	}

	var newTarget **Node
	if cmp < 0 {
		newTarget = &cur.Children[0]
	} else {
		newTarget = &cur.Children[1]
	}
	fix := t.remove(key, newTarget)
	if fix {
		return removeFix(int8(-cmp), target)
	}
	return false
}

func (t *Tree) bottom(d int) *Node {
	if t.Root == nil {
		return nil
	}
	var cur = t.Root
	for ; cur.Children[d] != nil; cur = cur.Children[d] {
	}
	return cur
}

func removeMin(target **Node, minKey *interface{}, minValue *interface{}) bool {
	cur := *target
	// 找到本子树中最小的node
	if cur.Children[0] == nil {
		// 使用找到的node的value替换minKey and minValue. 采用值拷贝的方法？不能直接替换node么？值拷贝最简单
		*minKey = cur.Key
		*minValue = cur.Value
		// remove cur node
		if cur.Children[1] != nil {
			cur.Children[1].Parent = cur.Parent
		}
		*target = cur.Children[1]
		return true
	}
	return removeMin(&cur.Children[0], minKey, minValue)
}

func putRebalance(c int8, root **Node) bool {
	cur := *root
	if cur.balance == 0 {
		cur.balance = c
		return true
	}
	if cur.balance == -c {
		cur.balance = 0
		return false
	}
	if cur.Children[(c+1)/2].balance == c {
		cur = singleRotate(c, cur)
	} else {
		cur = doubleRotate(c, cur)
	}
	*root = cur
	return false
}

func removeFix(c int8, root **Node) bool {
	cur := *root
	if cur.balance == 0 {
		cur.balance = c
		return false
	}

	if cur.balance == -c {
		cur.balance = 0
		return true
	}

	d := (c + 1) / 2
	if cur.Children[d].balance == 0 {
		cur = rotate(c, cur)
		cur.balance = -c
		*root = cur
		return false
	}

	if cur.Children[d].balance == c {
		cur = singleRotate(c, cur)
	} else {
		cur = doubleRotate(c, cur)
	}
	*root = cur
	return true
}

func singleRotate(c int8, root *Node) *Node {
	root.balance = 0
	// 旋转 得到新的平衡树的根
	root = rotate(c, root)
	root.balance = 0
	return root
}

func doubleRotate(c int8, root *Node) *Node {
	d := (c + 1) / 2
	child := root.Children[d]
	// 第一次rotate
	root.Children[d] = rotate(-c, root.Children[d])
	// 第二次rotate
	newRoot := rotate(c, root)
	switch {
	case newRoot.balance == c:
		root.balance = -c
		child.balance = 0
	case newRoot.balance == -c:
		root.balance = 0
		child.balance = c
	default:
		root.balance = 0
		child.balance = 0
	}
	newRoot.balance = 0
	return newRoot
}

// 旋转
// d == 0, LL, c == -1
// d == 1, RR, c == +1
func rotate(c int8, root *Node) *Node {
	d := (c + 1) / 2
	// 旋转后 新的根节点
	newRoot := root.Children[d]
	// 构建新的子树
	root.Children[d] = newRoot.Children[d^1]
	if root.Children[d] != nil {
		root.Children[d].Parent = root
	}
	// 构建新的树
	newRoot.Children[d^1] = root
	newRoot.Parent = root.Parent
	root.Parent = newRoot
	return newRoot
}

func (n *Node) Left() *Node {
	return n.Children[0]
}

func (n *Node) Right() *Node {
	return n.Children[1]
}

func (n *Node) Next() *Node {
	return n.walk(1)
}

func (n *Node) Prev() *Node {
	return n.walk(0)
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.Key)
}

// 如果d==1, 则代表寻找第一个比之大的节点。为当前节点右子孩子中最小的节点或者祖先节点的第一个比之大的右孩子
// 如果d==0, 则代表寻找第一个比之小的节点。为当前节点左子孩子中最大的节点或者祖先节点的第一个比之小的左孩子
func (n *Node) walk(d int) *Node {
	if n == nil {
		return nil
	}
	cur := n
	if cur.Children[d] == nil {
		parent := cur.Parent
		for parent != nil && parent.Children[d] == cur {
			cur = parent
			parent = parent.Parent
		}
		return parent
	}
	child := cur.Children[d]
	for child.Children[d^1] != nil {
		child = child.Children[d^1]
	}
	return child
}

// 格式化的后置遍历
func output(n *Node, prefix string, isTail bool, str *string) {
	// 递归Tree高度，prefix为每一层Tree的前缀
	if n.Children[1] != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(n.Children[1], newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += n.String() + "\n"
	if n.Children[0] != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(n.Children[0], newPrefix, true, str)
	}
}

// 格式化的后置遍历，与output等价的实现
func output2(n *Node, prefix string, isTail bool, str *string) {
	//// 后置遍历
	//if n == nil {
	//   return
	//}
	//output(n.Right(), "", false, str)
	//*str += n.String() + "\n"
	//output(n.Left(), "", false, str)
	if n == nil {
		return
	}
	newPrefix := prefix
	if isTail {
		newPrefix += "│   "
	} else {
		newPrefix += "    "
	}
	output(n.Right(), newPrefix, false, str)
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += n.String() + "\n"
	newPrefix = prefix
	if isTail {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}
	output(n.Left(), newPrefix, true, str)
}
