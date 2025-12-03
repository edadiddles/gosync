package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	NO_OP int = iota
	CREATED
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

			dest_name := "~~~~~"
			if dest_idx < len(dest_tree.children) {
				dest_name = dest_tree.children[dest_idx].name
			}

			fmt.Println(src_name, "<------>", dest_name)
			switch strings.Compare(src_name, dest_name) {
			case -1:
				new_file := src_tree.children[i]
				new_file.modified = CREATED
				tempSlice := append(dest_tree.children[:dest_idx], new_file)
				if dest_idx+1 >= len(dest_tree.children) {
					dest_tree.children = tempSlice
				} else {
					dest_tree.children = append(tempSlice, dest_tree.children[dest_idx+1:]...)
				}
			case 1:
				if dest_idx >= len(dest_tree.children) {
					fmt.Println("breaking early")
					break loop
				}

				dest_tree.children[dest_idx].modified = DELETED
				dest_idx++
			case 0:
				compare_ftree(src_tree.children[i], dest_tree.children[dest_idx])
				break loop
			}
		}

		dest_idx++
	}
}

func generate_checksum(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", filepath)
		return []byte{}, err
	}
	defer file.Close()

	h := sha256.New()

	_, err = io.Copy(h, file)
	if err != nil {
		fmt.Println("Error on io copy")
		return []byte{}, err
	}

	return h.Sum(nil), nil

}
