package main

import (
	"bufio"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
)

type Data struct {
	path      string
	directory []fs.FileInfo
	list      *List
}

func NewComponents() *Data {
	directory := ReadFiles(".")
	list := NewList(GetNames(directory))
	return &Data{
		path:      ".",
		directory: directory,
		list:      list,
	}
}

func (d *Data) UpdateData(path string) {
	d.directory = ReadFiles(path)
	d.list.elements = GetNames(d.directory)
}

func main() {
	data := NewComponents()
	data.list.InitialPrint()
	RunProgram(data)
}

func RunProgram(data *Data) {
	for char, key, _ := keyboard.GetSingleKey(); char != 'q'; {
		ExecuteComandsOnKeystroke(key, char, data)
		char, key, _ = keyboard.GetSingleKey()
	}
}

func UpdatePath(currentPath string, Selected string) string {
	return currentPath + `/` + Selected
}

func ReadText() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func ReadFiles(path string) []fs.FileInfo {
	files, _ := ioutil.ReadDir(path)
	return files
}

func InitialRead(path string) []fs.FileInfo {

	return ReadFiles(path)
}

func GetNames(directory []fs.FileInfo) []string {
	var names []string
	width := 70
	for i := 0; i < len(directory); i++ {
		dirName := directory[i].Name()
		dirSize := strconv.Itoa(int(directory[i].Size()))
		names = append(names, GetFullText(dirName, dirSize, width))
	}
	return names
}

func GetFullText(input string, fileSize string, maxLength int) string {

	if len(input)+len(fileSize)+5 >= maxLength {
		extraText := len(input) + len(fileSize) - maxLength
		input = input[0:len(input)-extraText-8] + "..."
	}
	availableSpace := maxLength - len(input) - len(fileSize)
	return input + strings.Repeat(" ", availableSpace) + fileSize + " KB"
}
