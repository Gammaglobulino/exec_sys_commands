package find_largest_files_on_dir

import (
	"container/list"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type FileNode struct {
	FullPath string
	Info     os.FileInfo
}

func InsertSortedNodeWithSortingLambda(fileList *list.List, fileNode FileNode, sorting func(a, b FileNode) bool) {
	if fileList.Len() == 0 {
		fileList.PushFront(fileNode)
		return
	}
	for element := fileList.Front(); element != nil; element = element.Next() {
		if sorting(element.Value.(FileNode), fileNode) {
			fileList.InsertBefore(fileNode, element)
			return
		}
	}
	fileList.PushBack(fileNode)
}

func InsertSortedNodeInfo(fileList *list.List, fileNode FileNode) {
	if fileList.Len() == 0 {
		fileList.PushFront(fileNode)
		return
	}
	for element := fileList.Front(); element != nil; element = element.Next() {
		if fileNode.Info.Size() > element.Value.(FileNode).Info.Size() {
			fileList.InsertBefore(fileNode, element)
			return
		}
	}
	fileList.PushBack(fileNode)
}

func GetFiles(fileList *list.List, path string) {
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range dirFiles {
		if file.IsDir() {
			fullPath := filepath.Join(path, file.Name())
			//fmt.Printf("[dir %s]: %s\n",fullPath,file.Name())
			GetFiles(fileList, fullPath)
		} else {
			//fmt.Println("[file]", file.Name())
			fullPath := filepath.Join(path, file.Name())
			InsertSortedNodeInfo(fileList, FileNode{
				FullPath: fullPath,
				Info:     file,
			})

		}
	}
}

func GetFilesWithLambdas(fileList *list.List, path string, sorting func(a, b FileNode) bool) {
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range dirFiles {
		if file.IsDir() {
			fullPath := filepath.Join(path, file.Name())
			//fmt.Printf("[dir %s]: %s\n",fullPath,file.Name())
			GetFiles(fileList, fullPath)
		} else {
			//fmt.Println("[file]", file.Name())
			fullPath := filepath.Join(path, file.Name())
			InsertSortedNodeWithSortingLambda(fileList, FileNode{
				FullPath: fullPath,
				Info:     file,
			}, sorting)

		}
	}
}
