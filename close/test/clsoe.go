package main

import (
	"fmt"
	"github.com/helays/utils/close"
	"os"
)

func main() {
	file, err := os.Open("clsoe.go")
	defer close.Close(file)
	fmt.Println(err, "文件")
}
