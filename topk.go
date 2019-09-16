package topk

import (
	"container/heap"
	"sort"
)

// Interface type describes the requirements for a type using the routines in this package
// Note that Push in this interface (heap.Interface) is for package topk's implementation to call.
// To add things from the heap, use topk.Push
// Useful for large streams of data when you only care about the top K elements, where K is small. Then it only
// keeps the top K elements in memory and performs Push in O(log k) time.
type Interface interface {
	// Interface is a heap with a size, so it needs to implement the heap.Interface
	heap.Interface
	// IsLess is used to determine whether a new value will get pushed to the structure or not
	// IsLess(interface{}, interface{}) should be consistent with Less(int, int) from heap.Interface
	IsLess(a, b interface{}) bool
	// Peek should return the minimum element. Same as Pop(), but without removing it
	Peek() interface{}
	// Get must return top K values, in any order
	// it mustn't modify (push/pop) the Interface, and the complexity should be O(1)
	// returned slice won't be modified, so one can/should return a "live view" of the data
	// if we pushed L values, then len(result) is for
	//	- L <  K: L
	//	- L >= K: K
	Get() []interface{}
	// K returns the size
	K() int
}

type comparables struct {
	values []interface{}
	isLess Less
}

// Len is a part of the sort.Interface type
func (c *comparables) Len() int {
	return len(c.values)
}

// Less is a part of the sort.Interface type
func (c *comparables) Less(i, j int) bool {
	return c.IsLess(c.values[i], c.values[j])
}

// Swap is a part of the sort.Interface type
func (c *comparables) Swap(i, j int) {
	c.values[i], c.values[j] = c.values[j], c.values[i]
}

// IsLess is part of the Interface type
func (c *comparables) IsLess(a, b interface{}) bool {
	return c.isLess(a, b)
}

// topK implements {heap,sort,}.Interface
type topK struct {
	k int
	*comparables
}

// Push is a part of the heap.Interface type
func (t *topK) Push(x interface{}) {
	t.values = append(t.values, x)
}

// Pop is a part of the heap.Interface type, panics if topK is empty
func (t *topK) Pop() interface{} {
	n := len(t.values)
	x := t.values[n-1]
	t.values = t.values[:n-1]
	return x
}

// Peek is a part of the Interface type, panics if topK is empty
func (t *topK) Peek() interface{} {
	return t.values[0]
}

// Get is a part of the Interface type
func (t *topK) Get() []interface{} {
	return t.values
}

// K is a part of the Interface type
func (t *topK) K() int {
	return t.k
}

// TOP LEVEL FUNCTIONS

// Less is used to determine whether the first argument is "less than" the second argument
type Less func(a, b interface{}) bool

// New returns an empty Interface with size k and ordering determined by the comparator
func New(k int, l Less) Interface {
	return &topK{k, &comparables{make([]interface{}, 0, k), l}}
}

// Push pushes the element x onto the topK. The complexity is O(log k)
func Push(t Interface, x interface{}) {
	if l := t.Len(); l < t.K() {
		heap.Push(t, x)
	} else if l > 0 && t.IsLess(t.Peek(), x) {
		heap.Pop(t)     // remove the smallest element
		heap.Push(t, x) // push the new one
	}
}

// Get returns the top K values from the given interface. The complexity is O(k * log k)
func Get(t Interface) []interface{} {
	vals := t.Get()
	valsCopy := make([]interface{}, len(vals))
	copy(valsCopy, vals)

	cmp := &comparables{valsCopy, t.IsLess}
	sort.Sort(sort.Reverse(cmp))

	return cmp.values
}

// IntComparator compares ints
func IntComparator(a, b interface{}) bool { return a.(int) < b.(int) }

// Float64Comparator compares float64s
func Float64Comparator(a, b interface{}) bool { return a.(float64) < b.(float64) }

// StringComparator compares strings
func StringComparator(a, b interface{}) bool { return a.(string) < b.(string) }
