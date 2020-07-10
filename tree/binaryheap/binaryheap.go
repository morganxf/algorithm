package binaryheap

import (
	"fmt"
	"strings"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/morganxf/algorithm/util"
)

type Heap struct {
	list       *arraylist.List
	Comparator util.Comparator
}

func NewWith(comparator util.Comparator) *Heap {
	return &Heap{list: arraylist.New(), Comparator: comparator}
}

func NewWithIntComparator() *Heap {
	return &Heap{list: arraylist.New(), Comparator: util.IntComparator}
}

func NewWithStringComparator() *Heap {
	return &Heap{list: arraylist.New(), Comparator: util.StringComparator}
}

func (heap *Heap) Push(values ...interface{}) {
	if len(values) == 1 {
		heap.list.Add(values[0])
		heap.bubbleUp()
	} else {
		heap.list.Add(values...)
		// lastChildParentIndex 是最后一个孩子父亲的index
		// 以此为起点，遍历到根节点，调整树结构
		lastChildParentIndex := heap.list.Size()/2 - 1
		for i := lastChildParentIndex; i >= 0; i-- {
			heap.bubbleDownIndex(i)
		}
	}
}

func (heap *Heap) Pop() (interface{}, bool) {
	// 获取堆顶元素
	value, ok := heap.list.Get(0)
	if !ok {
		return nil, false
	}
	// 最后一个元素和堆顶元素交换，删除堆顶元素，从堆顶重新调整堆
	lastIndex := heap.list.Size() - 1
	heap.list.Swap(0, lastIndex)
	heap.list.Remove(lastIndex)
	heap.bubbleDown()
	return value, true
}

func (heap *Heap) Peek() (interface{}, bool) {
	return heap.list.Get(0)
}

func (heap *Heap) Empty() bool {
	return heap.list.Empty()
}

func (heap *Heap) Size() int {
	return heap.list.Size()
}

func (heap *Heap) Clear() {
	heap.list.Clear()
}

func (heap *Heap) Values() []interface{} {
	return heap.list.Values()
}

func (heap *Heap) String() string {
	str := "BinaryHeap\n"
	values := []string{}
	for _, value := range heap.list.Values() {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

// 二叉堆 以0为起始点
// leftChildIndex = 2*i+1, rightChildIndex = 2*i+2
// parentIndex = (childIndex-1)/2
func (heap *Heap) bubbleUp() {
	index := heap.list.Size() - 1
	// 遍历index节点到root节点之间的路径，次数为树的高度
	for parentIndex := (index - 1) >> 1; parentIndex >= 0; parentIndex = (index - 1) >> 1 {
		indexValue, _ := heap.list.Get(index)
		parentValue, _ := heap.list.Get(parentIndex)
		if heap.Comparator(parentValue, indexValue) <= 0 {
			// 小于父亲节点
			break
		}
		heap.list.Swap(index, parentIndex)
		index = parentIndex
	}
}

func (heap *Heap) bubbleDown() {
	heap.bubbleDownIndex(0)
}

func (heap *Heap) bubbleDownIndex(index int) {
	size := heap.list.Size()
	for leftIndex := index<<1 + 1; leftIndex < size; leftIndex = index<<1 + 1 {
		rightIndex := leftIndex + 1
		smallerIndex := leftIndex
		leftValue, _ := heap.list.Get(leftIndex)
		rightValue, _ := heap.list.Get(rightIndex)
		if rightIndex < size && heap.Comparator(leftValue, rightValue) > 0 {
			smallerIndex = rightIndex
		}
		indexValue, _ := heap.list.Get(index)
		smallerValue, _ := heap.list.Get(smallerIndex)
		// 向下迭代，直到父亲小于孩子，或者迭代到最后一个节点
		if heap.Comparator(indexValue, smallerValue) > 0 {
			heap.list.Swap(index, smallerIndex)
		} else {
			break
		}
		index = smallerIndex
	}
}

func (heap *Heap) withinRange(index int) bool {
	return index >= 0 && index < heap.Size()
}
