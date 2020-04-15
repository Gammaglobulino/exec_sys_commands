package find_largest_files_on_dir

import (
	"../find_largest_files_on_dir"
	"container/list"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestGetDir(t *testing.T) {
	path := "c:\\go"
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	assert.Nil(t, err)
	assert.NotEmpty(t, dirFiles)

}
func TestAddNodeToList(t *testing.T) {
	path := "c:\\go"
	dirFiles, err := ioutil.ReadDir(path)
	fileList := list.List{}

	if err != nil {
		log.Fatal(err)
	}
	for _, file := range dirFiles {
		if file.IsDir() {
			fullPath := filepath.Join(path, file.Name())
			fmt.Printf("[dir %s]: %s\n", fullPath, file.Name())
		} else {
			fmt.Println("[file]", file.Name())
			fullPath := filepath.Join(path, file.Name())
			find_largest_files_on_dir.InsertSortedNodeInfo(&fileList, find_largest_files_on_dir.FileNode{
				FullPath: fullPath,
				Info:     file,
			})

		}
	}
	assert.NotEmpty(t, fileList)
	for element := fileList.Front(); element != nil; element = element.Next() {

		fileName := element.Value.(find_largest_files_on_dir.FileNode).Info.Name()
		fileSize := element.Value.(find_largest_files_on_dir.FileNode).Info.Size()

		fmt.Printf("file: %s size:%v \n", fileName, fileSize)
	}
}

func TestGetFilesInDirBySizeFirtDirLEvel(t *testing.T) {
	path := "c:\\go"
	fileList := list.List{}
	find_largest_files_on_dir.GetFilesInDirBySize(&fileList, path)
	assert.NotEmpty(t, fileList)
	for element := fileList.Front(); element != nil; element = element.Next() {

		fileName := element.Value.(find_largest_files_on_dir.FileNode).Info.Name()
		filePath := element.Value.(find_largest_files_on_dir.FileNode).FullPath
		fileSize := element.Value.(find_largest_files_on_dir.FileNode).Info.Size()

		fmt.Printf("file: %s\\%s size:%v \n", filePath, fileName, fileSize)
	}
}

func TestAddNodeWithLambdasSortingByFileSize(t *testing.T) {
	path := "c:\\go"
	dirFiles, err := ioutil.ReadDir(path)
	fileList := list.List{}

	sortingFunction := func(a, b find_largest_files_on_dir.FileNode) bool {
		if a.Info.Size() < b.Info.Size() {
			return true
		}
		return false
	}

	if err != nil {
		log.Fatal(err)
	}
	for _, file := range dirFiles {
		if file.IsDir() {
			fullPath := filepath.Join(path, file.Name())
			fmt.Printf("[dir %s]: %s\n", fullPath, file.Name())
		} else {
			fmt.Println("[file]", file.Name())
			fullPath := filepath.Join(path, file.Name())
			find_largest_files_on_dir.InsertSortedNodeWithLambdas(&fileList, find_largest_files_on_dir.FileNode{
				FullPath: fullPath,
				Info:     file,
			}, sortingFunction)

		}
	}
	assert.NotEmpty(t, fileList)
	for element := fileList.Front(); element != nil; element = element.Next() {

		fileName := element.Value.(find_largest_files_on_dir.FileNode).Info.Name()
		fileSize := element.Value.(find_largest_files_on_dir.FileNode).Info.Size()

		fmt.Printf("file: %s size:%v \n", fileName, fileSize)
	}
}

func TestAddNodeWithLambdasSortingByModifiedDate(t *testing.T) {
	path := "c:\\go"
	dirFiles, err := ioutil.ReadDir(path)
	fileList := list.List{}

	sortingFunction := func(a, b find_largest_files_on_dir.FileNode) bool {
		if a.Info.ModTime().Before(b.Info.ModTime()) {
			return true
		}
		return false
	}

	if err != nil {
		log.Fatal(err)
	}
	for _, file := range dirFiles {
		if file.IsDir() {
			fullPath := filepath.Join(path, file.Name())
			fmt.Printf("[dir %s]: %s\n", fullPath, file.Name())
		} else {
			fmt.Println("[file]", file.Name())
			fullPath := filepath.Join(path, file.Name())
			find_largest_files_on_dir.InsertSortedNodeWithLambdas(&fileList, find_largest_files_on_dir.FileNode{
				FullPath: fullPath,
				Info:     file,
			}, sortingFunction)

		}
	}
	assert.NotEmpty(t, fileList)
	for element := fileList.Front(); element != nil; element = element.Next() {

		fileName := element.Value.(find_largest_files_on_dir.FileNode).Info.Name()
		fileSize := element.Value.(find_largest_files_on_dir.FileNode).Info.Size()
		fileDate := element.Value.(find_largest_files_on_dir.FileNode).Info.ModTime()

		fmt.Printf("file: %s size:%v last modified:%s\n", fileName, fileSize, fileDate)
	}
}
