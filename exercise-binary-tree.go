package main

import "golang.org/x/tour/tree"
import "fmt"

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	var walk func(_t *tree.Tree)
	walk = func (_t *tree.Tree) {
		if _t == nil {
			return
		}
		walk(_t.Left)
		ch <- _t.Value
		walk(_t.Right)
	}
	walk(t)
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for i := range ch1 {
		j, ok := <-ch2
		if !ok || i != j {
			return false
		}
	}
	_, ok := <-ch2
	return !ok
}

func main() {
	fmt.Println(Same(tree.New(2), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
