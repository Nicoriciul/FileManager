package main

import (
	"io/fs"
	"io/ioutil"

	tsize "github.com/kopoli/go-terminal-size"
)

var Directory = ReadFiles(".")
var Selected string
var StartPos = 0
var CurrentPosition = 0
var Size, _ = tsize.GetSize()
var WindowSize = Size.Height - 10
var EndPos = LastPos()
var SelectedForCopy string
var SelectedForCut string

func ResetCoordinates(directory []fs.FileInfo) {
	StartPos = 0
	CurrentPosition = 0
	EndPos = LastPos()
	Selected = directory[CurrentPosition].Name()
}

func LastPos() int {
	return Size.Height - 10 + StartPos
}

func ReadFiles(path string) []fs.FileInfo {
	files, _ := ioutil.ReadDir(path)
	return files
}
