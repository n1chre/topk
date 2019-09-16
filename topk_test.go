package topk

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestComparablesSort(t *testing.T) {
	c := &comparables{
		values: []interface{}{5, 6, 2, 4, 2},
		isLess: IntComparator,
	}
	if l := c.Len(); l != 5 {
		t.Fatalf("expecting len to be 5, got %d", l)
	}
	if is := c.Less(0, 1); !is {
		t.Fatalf("expecting less(5,6) to be true, got false")
	}
	c.Swap(0, 1)
	if want := []interface{}{6, 5, 2, 4, 2}; !reflect.DeepEqual(want, c.values) {
		t.Fatalf("want %+v after swap, got %+v", want, c.values)
	}
	sort.Sort(c)
	if want := []interface{}{2, 2, 4, 5, 6}; !reflect.DeepEqual(want, c.values) {
		t.Fatalf("want %+v after sort, got %+v", want, c.values)
	}
}

func TestTopK(t *testing.T) {
	type step struct {
		x     interface{}
		state []interface{}
	}

	for i, c := range []struct {
		k     int
		steps []step
	}{
		{
			k: 1,
			steps: []step{
				{5, []interface{}{5}},
				{6, []interface{}{6}},
				{2, []interface{}{6}},
				{4, []interface{}{6}},
				{2, []interface{}{6}},
			},
		},
		{
			k: 2,
			steps: []step{
				{5, []interface{}{5}},
				{6, []interface{}{6, 5}},
				{2, []interface{}{6, 5}},
				{4, []interface{}{6, 5}},
				{2, []interface{}{6, 5}},
			},
		},
		{
			k: 3,
			steps: []step{
				{5, []interface{}{5}},
				{6, []interface{}{6, 5}},
				{2, []interface{}{6, 5, 2}},
				{4, []interface{}{6, 5, 4}},
				{2, []interface{}{6, 5, 4}},
			},
		},
		{
			k: 4,
			steps: []step{
				{5, []interface{}{5}},
				{6, []interface{}{6, 5}},
				{2, []interface{}{6, 5, 2}},
				{4, []interface{}{6, 5, 4, 2}},
				{2, []interface{}{6, 5, 4, 2}},
			},
		},
		{
			k: 5,
			steps: []step{
				{5, []interface{}{5}},
				{6, []interface{}{6, 5}},
				{2, []interface{}{6, 5, 2}},
				{4, []interface{}{6, 5, 4, 2}},
				{2, []interface{}{6, 5, 4, 2, 2}},
			},
		},
	} {
		t.Run(fmt.Sprintf("case-%d", i+1), func(t *testing.T) {
			topK := New(c.k, IntComparator)
			for _, s := range c.steps {
				Push(topK, s.x)
				if g := Get(topK); !reflect.DeepEqual(g, s.state) {
					t.Fatalf("state mismatch after pushing %v:\n\twant: %v\n\thave: %v", s.x, g, s.state)
				}
			}
		})
	}
}
