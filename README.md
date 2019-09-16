# topk

TopK is a Go implementation of the Top K elements. 

## API

* `New(k int, l Less) Interface` - create a new Interface which supports `Push` and `Get`. Less determines the ordering
* `Push(t Interface, x interface{})` - push the new value onto the structure
* `Get(t Interface, x interface{})` - get the top K values

## Example usage

```golang
package main

import "github.com/n1chre/topk"

func main() {
    tk := topk.New(3, topk.IntComparator)
    for _, x := range []int{1, 5, 2, -1, 8, 9, -10} {
        topk.Push(tk, x)
    }
    // top3 = []interface{}{9, 8, 5}
    fmt.Println(topk.Get(tk))
}
```
