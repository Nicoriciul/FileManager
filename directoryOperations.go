package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/eiannone/keyboard"
	cp "github.com/otiai10/copy"
)

var SelectedForCopy string
var SelectedForCut string

func ExitDirectory(data *Data) {

	data.path = strings.TrimRight(filepath.Dir(data.path), `/`)
	data.list.ResetIndexes()
	data.UpdateData(data.path)
	data.list.Print()
}

func EnterDirectory(data *Data) {

	forward := ReadFiles(UpdatePath(data.path, data.directory[data.list.selectedIndex].Name()))
	if len(forward) > 0 {
		data.path = UpdatePath(data.path, data.directory[data.list.selectedIndex].Name())
		data.list.ResetIndexes()
		data.UpdateData(data.path)

	} else if data.directory[data.list.selectedIndex].IsDir() {
		data.path = UpdatePath(data.path, data.directory[data.list.selectedIndex].Name())
		data.UpdateData(data.path)
		ClearConsole()
		fmt.Println("empty")
		return
	} else {

		return
	}
	data.list.Print()
}

func IndexOf(files []fs.FileInfo, Selected string) int {
	for i, current := range files {
		if current.Name() == Selected {
			return i
		}
	}
	return -1
}

func NewFolder(data *Data) {
	newName := NameFileOrFolder(data, "New Folder")
	err := os.Mkdir(UpdatePath(data.path, newName), 0755)
	if err != nil {
		log.Fatal(err)
	}
	data.UpdateData(data.path)
	data.list.selectedIndex = IndexOf(data.directory, newName)
	data.list.UpdateCoordinates()
	data.list.Print()
}

func NewFile(data *Data) {
	newName := NameFileOrFolder(data, "New File")
	file, e := os.Create(UpdatePath(data.path, newName))
	if e != nil {
		log.Fatal(e)
	}
	file.Close()
	data.UpdateData(data.path)
	data.list.selectedIndex = IndexOf(data.directory, newName)
	data.list.UpdateCoordinates()
	data.list.Print()
}

func Rename(data *Data) {
	selected := data.directory[data.list.selectedIndex].Name()
	newName := NameFileOrFolder(data, selected)
	originalPath := UpdatePath(data.path, selected)
	newPath := UpdatePath(data.path, newName)
	e := os.Rename(originalPath, newPath)
	if e != nil {
		log.Fatal(e)
	}
	data.UpdateData(data.path)
	data.list.selectedIndex = IndexOf(data.directory, newName)
	data.list.UpdateCoordinates()
	data.list.Print()
}

func NameFileOrFolder(data *Data, newName string) string {
	PrintOnLine(newName, data.list.GetGap())
	_, key, _ := keyboard.GetSingleKey()
	if key == keyboard.KeyBackspace2 || key == keyboard.KeyBackspace {
		data.list.Print()
		PrintOnLine("", data.list.GetGap())
		newName = ReadText()
	} else if key == keyboard.KeyEnter {
		newName = data.directory[data.list.selectedIndex].Name()
	}
	return newName
}

func Paste(data *Data) {
	var copiedFolder string
	if len(SelectedForCut) > 0 {
		copiedFolder = filepath.Base(SelectedForCut)
		CutPaste(data.path)
	} else if len(SelectedForCopy) > 0 {
		copiedFolder = filepath.Base(SelectedForCopy)
		CopyPaste(data.path)
	} else {
		PrintOnLine("No File or Folder Selected", data.list.GetGap())
		return
	}
	data.UpdateData(data.path)
	data.list.selectedIndex = IndexOf(data.directory, copiedFolder)
	data.list.UpdateCoordinates()
	data.list.Print()
}

func CopyPaste(path string) {
	folderName := filepath.Base(SelectedForCopy)
	cp.Copy(SelectedForCopy, UpdatePath(path, folderName))
}

func CutPaste(path string) {
	folderName := filepath.Base(SelectedForCut)
	cp.Copy(SelectedForCut, UpdatePath(path, folderName))
	os.RemoveAll(SelectedForCut)
}

func Delete(data *Data) {
	os.RemoveAll(UpdatePath(data.path, data.directory[data.list.selectedIndex].Name()))
	data.UpdateData(data.path)
	data.list.ResetIndexes()
	if len(data.directory) < 1 {
		ClearConsole()
		fmt.Println("empty")
		return
	}
	data.list.Print()
}

func NameExists(startDir string) bool {
	return SelectedForCut == UpdatePath(startDir, filepath.Base(SelectedForCut)) ||
		SelectedForCopy == UpdatePath(startDir, filepath.Base(SelectedForCopy))
}
