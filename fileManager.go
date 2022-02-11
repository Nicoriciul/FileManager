package main

import (
	"bufio"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	tsize "github.com/kopoli/go-terminal-size"
)

type Coordinates struct {
	windowFirstElemIndex int
	selectedIndex        int
	windowLastElemIndex  int
	selectedName         string
}

var Size, _ = tsize.GetSize()
var WindowSize = Size.Height - 10

func main() {
	RunProgram(".")
}

func RunProgram(startDir string) {
	directory, coordinates := InitialRead(startDir)
	for char, key, _ := keyboard.GetSingleKey(); char != 'q'; {
		startDir, directory, coordinates = ExecuteComandsOnKeystroke(key, startDir, directory, coordinates, char)
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
