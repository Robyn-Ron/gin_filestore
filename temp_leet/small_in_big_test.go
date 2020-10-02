package temp_leet

import (
	"fmt"
	"math"
	"testing"
)

func TestModifyMapValue(t *testing.T) {
	m := make(map[int][]string)
	m[1] = []string{"hello"}
	m[1] = append(m[1], "hhh")
	fmt.Println(m)
}


func TestShortSeq(t *testing.T) {
	shortestSeq([]int{7,5,9,0,2,1,3,5,7,9,1,1,5,8,8,9,7}, []int{1,5,9})
}

func shortestSeq(big []int, small []int) []int {
	m := MySmall(small)
	left, right := -1, -1
	for i := 0; i < len(big); i++ {
		if m.inSmall(big[i]) {
			left, right = i, i
			break
		}
	}
	if left == -1 || right == -1 {
		return []int{}
	}
	min := int(math.MaxInt32)
	minLeft := left
	counter := InitCounter(small) //用来记录窗口中small中各个数的数量
	for i := right; i < len(big); i++ {
		if ! m.inSmall(big[i]){
			continue
		}
		counter.add(i,big[i],1)
		right = i
		if counter.check() {
			if right - left + 1 < min {
				min = right - left + 1
				minLeft = left
			}
			counter.add(i,big[left], -1)
			left = counter.lefts[0]
		}
	}
	if min == int(math.MaxInt32){
		return []int{}
	}
	return []int{minLeft, minLeft+min-1}
}

type MySmall []int
func (this MySmall)inSmall(num int) bool{
	for _, number := range this {
		if number == num {
			return true
		}
	}
	return false
}

type Counter struct{
	small []int
	count []int
	lefts []int
}

func InitCounter(small []int)  Counter{
	return Counter{
		small: small,
		count: make([]int, len(small)),
	}
}

func (this *Counter)check() bool{
	for _, num := range this.count {
		if num < 1 {
			return false
		}
	}
	return true
}

func (this *Counter) add(left, number, cnt int) {
	if cnt > 0 {
		this.lefts = append(this.lefts, left)
	}
	for i:=0; i < len(this.small); i++{
		if number == this.small[i] {
			this.count[i] += cnt
			return
		}
	}
}
