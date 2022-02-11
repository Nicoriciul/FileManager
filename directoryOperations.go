package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/eiannone/keyboard"
	cp "github.com/otiai10/copy"
)

var SelectedForCopy string
var SelectedForCut string

func NewFolder(directory []fs.FileInfo, coordinates Coordinates, startDir string) ([]fs.FileInfo, Coordinates) {
	newName := NameFileOrFolder(directory, coordinates, "New Folder")
	err := os.Mkdir(UpdatePath(startDir, newName), 0755)
	if err != nil {
		log.Fatal(err)
	}
	directory = ReadFiles(startDir)
	coordinates.selectedIndex = IndexOf(directory, newName)
	coordinates.selectedName = newName
	coordinates = UpdateCoordinatesBeforePrinting(directory, coordinates)
	fmt.Println(coordinates)
	PrintDir(directory, coordinates)
	return directory, coordinates
}

func NewFile(directory []fs.FileInfo, coordinates Coordinates, startDir string) ([]fs.FileInfo, Coordinates) {
	newName := NameFileOrFolder(directory, coordinates, "New File")
	file, e := os.Create(UpdatePath(startDir, newName))
	if e != nil {
		log.Fatal(e)
	}
	file.Close()
	directory = ReadFiles(startDir)
	coordinates.selectedIndex = IndexOf(directory, newName)
	coordinates.selectedName = newName
	coordinates = UpdateCoordinatesBeforePrinting(directory, coordinates)
	PrintDir(directory, coordinates)
	return directory, coordinates
}

func Rename(directory []fs.FileInfo, coordinates Coordinates, startDir string) ([]fs.FileInfo, Coordinates) {
	newName := NameFileOrFolder(directory, coordinates, coordinates.selectedName)
	originalPath := UpdatePath(startDir, coordinates.selectedName)
	newPath := UpdatePath(startDir, newName)
	e := os.Rename(originalPath, newPath)
	if e != nil {
		log.Fatal(e)
	}
	directory = ReadFiles(startDir)
	coordinates.selectedName = newName
	coordinates.selectedIndex = IndexOf(directory, coordinates.selectedName)
	coordinates = UpdateCoordinatesBeforePrinting(directory, coordinates)
	PrintDir(directory, coordinates)
	return directory, coordinates
}

func Paste(startDir string, coordinates Coordinates) (Coordinates, []fs.FileInfo) {
	if len(SelectedForCut) > 0 {
		coordinates = CutPaste(coordinates, startDir)
	} else if len(SelectedForCopy) > 0 {
		coordinates = CopyPaste(coordinates, startDir)
	}
	directory := ReadFiles(startDir)
	coordinates.selectedIndex = IndexOf(directory, coordinates.selectedName)
	coordinates = UpdateCoordinatesBeforePrinting(directory, coordinates)
	PrintDir(directory, coordinates)
	return coordinates, directory
}

func CopyPaste(coordinates Coordinates, startDir string) Coordinates {
	coordinates.selectedName = filepath.Base(SelectedForCopy)
	cp.Copy(SelectedForCopy, UpdatePath(startDir, coordinates.selectedName))
	return coordinates
}

func CutPaste(coordinates Coordinates, startDir string) Coordinates {
	coordinates.selectedName = filepath.Base(SelectedForCut)
	cp.Copy(SelectedForCut, UpdatePath(startDir, coordinates.selectedName))
	os.RemoveAll(SelectedForCut)
	return coordinates
}

func NameExists(startDir string) bool {
	return SelectedForCut == UpdatePath(startDir, filepath.Base(SelectedForCut)) ||
		SelectedForCopy == UpdatePath(startDir, filepath.Base(SelectedForCopy))
}

func Delete(startDir string, coordinates Coordinates) ([]fs.FileInfo, Coordinates) {
	os.RemoveAll(UpdatePath(startDir, coordinates.selectedName))
	directory := ReadFiles(startDir)
	if len(directory) > 0 {
		coordinates = ResetCoordinates(directory)

	} else {
		coordinates.selectedName = ""
		ClearConsole()
		fmt.Println("empty")
	}
	return directory, coordinates
}

func NameFileOrFolder(directory []fs.FileInfo, coordinates Coordinates, newName string) string {
	PrintOnLine(newName, GetGap(directory, Size.Height))
	_, key, _ := keyboard.GetSingleKey()
	if key == keyboard.KeyBackspace2 || key == keyboard.KeyBackspace {
		PrintDir(directory, coordinates)
		PrintOnLine("", GetGap(directory, Size.Height))
		newName = ReadText()
	} else if key == keyboard.KeyEnter {
		newName = coordinates.selectedName
	}
	return newName
}
