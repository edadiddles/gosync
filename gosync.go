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

	prog := os.Args[1]
	src := os.Args[2]
	dest := os.Args[3]

	//TODO: check which prog is passed so the correct actions can be taking
	fmt.Println("Prog:", prog)
	fmt.Println("src:", src, "--> dest:", dest)

	fs1, err := os.ReadDir(src)
	if err != nil {
		fmt.Println(err)
		return
	}

	fs2, err := os.ReadDir(dest)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := range len(fs1) {
		fmt.Println(fs1[i])
		if fs1[i].IsDir() {
			subdir, err := os.ReadDir(src + "/" + fs1[i].Name())
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(subdir)
		}
	}

	fmt.Println("-----")
	
	for i := range len(fs2) {
		fmt.Println(fs2[i])
		if fs2[i].IsDir() {
			subdir, err := os.ReadDir(dest + "/" + fs2[i].Name())
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(subdir)
		}
	}
}
