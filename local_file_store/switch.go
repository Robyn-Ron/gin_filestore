package day01

import(
	"fmt"
)

func numToDate(number int) string {
	switch number {
	case 1:
	 	return "monday"
	case 2:
		return "tuesday"
	case 3:
		return "wednesday"
	case 4:
		return "thursday"
	case 5:
		return "friday"
	case 6:
		return "saturday"
	case 7:
		return "sunday"
	}
	return "please input a regular number"
}

func main() {
	var date int
	for true {
		fmt.Scanf("%d", &date)
		switch numToDate(date) {
		case "monday", "tuesday", "wednesday", "thursday", "friday":
			fmt.Println("今天上班, 别想着出去浪了......")
			fmt.Println("除非你财务自由了!")
		case "saturday", "sunday":
			fmt.Println("今天还是可以出去浪的, 不过你也可以选择加班 (*╹▽╹*)")
			fallthrough
		default:
			fmt.Println(numToDate(date))
		}
	}

}
