package temp_leet

import (
	"container/heap"
	"fmt"
	"testing"
)

func TestMaxHeapInSmallestK(t *testing.T) {
	m := &MaxHeap{}
	heap.Init(m)
	heap.Push(m,3)
	heap.Push(m,5)
	heap.Push(m,1)
	heap.Push(m,2)
	heap.Push(m,1)
	fmt.Println((*m)[0], (*m)[1], (*m)[2], (*m)[3], (*m)[4])
	fmt.Println("heap = ", m)
	fmt.Println(heap.Pop(m))
	fmt.Println(heap.Pop(m))
	fmt.Println(heap.Pop(m))
	fmt.Println(heap.Pop(m))
	fmt.Println(heap.Pop(m))
}

func smallestK(arr []int, k int) []int {
	if len(arr) == 0 || k == 0 {
		return []int{}
	}
	m := MaxHeap{}
	for _, number := range arr {
		if m.Len() < k {
			m.Push(number)
		} else {
			top := m[0]
			if number >= top {
				continue
			} else {
				m.Pop()
				m.Push(number)
			}
		}
	}
	return []int(m)
}

type MaxHeap []int

func (this MaxHeap) Len() int {
	return len(this)
}

func (this MaxHeap) Swap(i, j int) {
	(this)[i], (this)[j] = (this)[j], (this)[i]
}

func (this MaxHeap)Less(i, j int)bool {
	return (this)[i] > (this)[j]
}

func (this *MaxHeap)Push(x interface{}) {
	v := x.(int)
	(*this) = append(*this, v)
}

func (this *MaxHeap)Pop() interface{}{
	old := *this
	res := (old)[len(*this)-1]
	*this = (old)[:len(*this)-1]
	return res
}
