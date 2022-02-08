package main

import "io/fs"

var Directory = GetDirectory(".")
var Selected string
var StartPosition int
var CurrentPosition int
var EndPos int

func GetDirectory(path string) []fs.FileInfo {

	return ReadFiles(path)
}

func ResetCoordinates() {

}
