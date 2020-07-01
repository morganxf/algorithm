package container

// Container is base interface that all data structures implement.
type Container interface {
	Empty() bool
	Size() int
	Clear()
	Values() []interface{}
}
