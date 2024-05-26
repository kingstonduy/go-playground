package main

import (
	"fmt"
	"os"
)

func main() {
	var data string
	// data := []byte("This is some data to write to the file.")
	for i := 2; i <= 128; i++ {
		data = data + fmt.Sprintf("Field%d string `json:\"field%d:\"`\n", i, i)
	}
	err := os.WriteFile("myfile.txt", []byte(data), 0644) // 0644 is a common permission mode
	if err != nil {
		panic(err)
	}
}
