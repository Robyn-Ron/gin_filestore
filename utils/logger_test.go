package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	fmt.Println(log.Prefix())
	fmt.Println(log.Flags())
	log.Println("haha")
	//log.Panicln("haha")
	//log.Fatalln("haha")
}
