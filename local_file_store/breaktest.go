package day01

import(
	"fmt"
	"math/rand"
	"time"
)

func main(){
	rand.Seed(time.Now().UnixNano())
	var count int = 0
	for{
		count++
		var num0 int = rand.Intn(100)
		if num0 == 99 {
			break
		}
	}
	fmt.Printf("生成99这个数, 一共运行了%d次;\n", count)
}
