package main

import (
	"fmt"
	//"io/fs"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("At least three arguments expected")
		return
	}

	dir1 := os.Args[1]
	dir2 := os.Args[2]

	fmt.Println("Dir1:", dir1, "-- Dir2:", dir2)

	fs1, err := os.ReadDir(dir1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fs2, err := os.ReadDir(dir2)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fs1)
	fmt.Println(fs2)
}
