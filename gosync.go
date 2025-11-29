package main

import (
	"fmt"
	"os"
	//"crypto/sha256"
	"strings"
)

const (
	CREATED int = iota
	UPDATED
	DELETED
)

type ftree struct {
	is_dir   bool
	name     string
	checksum string
	modified int
	children []*ftree
}

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

	srcTree := ftree{is_dir: true, name: src}
	build_ftree(&srcTree, src, fs1)

	destTree := ftree{is_dir: true, name: dest}
	build_ftree(&destTree, dest, fs2)


	compare_ftree(&srcTree, &destTree)
	
	walkFTree(&destTree)
}

func build_ftree(dirNode *ftree, pwd string, dirEntries []os.DirEntry) {
	for i := range len(dirEntries) {
		currFile := dirEntries[i]
		fNode := ftree{is_dir: currFile.IsDir(), name: currFile.Name()}
		dirNode.children = append(dirNode.children, &fNode)
		if fNode.is_dir {
			newDir := pwd + "/" + fNode.name
			dirList, err := os.ReadDir(newDir)
			if err != nil {
				fmt.Println(err)
				continue
			}
			build_ftree(&fNode, newDir, dirList)
		}
	}
}

func walkFTree(parent *ftree) {
	fmt.Println(parent)
	for i := range len(parent.children) {
		walkFTree(parent.children[i])
	}
}

func compare_ftree(src_tree *ftree, dest_tree *ftree) {
	dest_idx := 0

	if src_tree.checksum != dest_tree.checksum {
		dest_tree.modified = UPDATED
	}

	for i := range src_tree.children {

	loop:
		for {
			src_name := src_tree.children[i].name

			dest_name := ""
			if dest_idx >= len(dest_tree.children) {
				dest_name = dest_tree.children[dest_idx].name
			}

			switch strings.Compare(src_name, dest_name) {
			case 1:
				new_file := src_tree.children[i]
				fmt.Println(src_name, "created")
				new_file.modified = CREATED
				tempSlice := append(dest_tree.children[:dest_idx], new_file)
				dest_tree.children = append(tempSlice, dest_tree.children[dest_idx:]...)
				break loop
			case -1:
				fmt.Println(dest_name, "deleted")
				dest_tree.children[i].modified = DELETED
				dest_idx++
			case 0:
				fmt.Println(src_name, "matched")
				break loop
			}
		}

		fmt.Println("a;lkdfj")
		dest_idx++
	}
}
