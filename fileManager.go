package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
	tsize "github.com/kopoli/go-terminal-size"
	cp "github.com/otiai10/copy"
)

type Coordinates struct {
	windowFirstElemIndex int
	selectedIndex        int
	windowLastElemIndex  int
	selectedName         string
}

var SelectedForCopy string
var SelectedForCut string
var Size, _ = tsize.GetSize()
var WindowSize = Size.Height - 10

func main() {

	RunProgram(".")

}

func Printdirectory(dir []fs.FileInfo) {
	for i := 0; i < len(dir); i++ {
		fmt.Println(dir[i].Name())
	}
}

func RunProgram(startDir string) {
	directory := ReadFiles(startDir)
	coordinates := ResetCoordinates(directory)
	PrintDir(directory, coordinates)

	for char, key, _ := keyboard.GetSingleKey(); char != 'q'; {
		startDir, directory, coordinates = ExecuteComandsOnKeystroke(key, startDir, directory, coordinates, char)
		char, key, _ = keyboard.GetSingleKey()
	}
}

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

func PrintHelp() {
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	ClearConsole()
	commands := []string{"COPY = V", "CUT = C", "PASTE = P", "RENAME = R", "NEW FILE = M", "NEW FOLDER = N", "DELETE = D", "HELP = H", "QUIT = Q"}
	for _, current := range commands {
		fmt.Println(string(colorRed), current)
	}
	fmt.Print(colorReset)
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

func GoDownWhenListFits(directory []fs.FileInfo, coordinates Coordinates) Coordinates {
	if coordinates.selectedIndex < len(directory) {
		coordinates.selectedName = directory[coordinates.selectedIndex].Name()
	} else {
		coordinates = ResetCoordinates(directory)
	}
	PrintDir(directory, coordinates)
	return coordinates
}

func GoDownWhenListDoesntFit(directory []fs.FileInfo, coordinates Coordinates) Coordinates {
	if coordinates.selectedIndex <= LastPos(coordinates.windowFirstElemIndex) {
		coordinates.selectedName = directory[coordinates.selectedIndex].Name()
		PrintDir(directory, coordinates)
		return coordinates
	}
	coordinates = ScrollDown(directory, coordinates)
	return coordinates
}

func ScrollDown(directory []fs.FileInfo, coordinates Coordinates) Coordinates {
	if coordinates.selectedIndex > LastPos(coordinates.windowFirstElemIndex) && LastPos(coordinates.windowFirstElemIndex)+1 <= len(directory)-1 {
		coordinates.windowFirstElemIndex++
		coordinates.windowLastElemIndex = LastPos(coordinates.windowFirstElemIndex)
		coordinates.selectedName = directory[coordinates.selectedIndex].Name()
		PrintDir(directory, coordinates)
	} else {
		coordinates = ResetCoordinates(directory)
		PrintDir(directory, coordinates)
	}
	return coordinates
}

func GoUpWhenListFits(directory []fs.FileInfo, coordinates Coordinates) Coordinates {
	if coordinates.selectedIndex >= 0 {
		coordinates.selectedName = directory[coordinates.selectedIndex].Name()
	} else if len(directory) > 0 {
		coordinates.selectedIndex = len(directory) - 1
		coordinates.selectedName = directory[coordinates.selectedIndex].Name()
	}
	PrintDir(directory, coordinates)
	return coordinates
}

func GoUpWhenListDoesntFit(directory []fs.FileInfo, coordinates Coordinates) Coordinates {

	if coordinates.selectedIndex >= coordinates.windowFirstElemIndex {
		coordinates.selectedName = directory[coordinates.selectedIndex].Name()
		coordinates.windowLastElemIndex = LastPos(coordinates.windowFirstElemIndex)
		PrintDir(directory, coordinates)
		return coordinates
	}
	coordinates = ScrollUp(directory, coordinates)
	return coordinates

}

func ScrollUp(directory []fs.FileInfo, coordinates Coordinates) Coordinates {
	if coordinates.selectedIndex < coordinates.windowLastElemIndex && coordinates.windowFirstElemIndex-1 >= 0 {
		coordinates.windowFirstElemIndex--
	} else {
		coordinates.windowFirstElemIndex = len(directory) - WindowSize - 1
		coordinates.selectedIndex = len(directory) - 1
	}
	coordinates.windowLastElemIndex = LastPos(coordinates.windowFirstElemIndex)
	coordinates.selectedName = directory[coordinates.selectedIndex].Name()
	PrintDir(directory, coordinates)
	return coordinates
}

func UpdatePath(currentPath string, Selected string) string {
	return currentPath + `/` + Selected
}

func ClearConsole() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func UpdateCoordinatesBeforePrinting(directory []fs.FileInfo, coordinates Coordinates) Coordinates {
	indexOfSelected := coordinates.selectedIndex
	directoryLength := len(directory)
	if IsInsideCurrentWindow(indexOfSelected, coordinates.windowFirstElemIndex, coordinates.windowLastElemIndex) {
		return coordinates
	}
	if IsInsideFirstWindow(indexOfSelected) {
		coordinates.windowFirstElemIndex = 0
		coordinates.windowLastElemIndex = LastPos(0)
	}
	if IsInsideLastWindow(directoryLength, indexOfSelected) {
		coordinates.windowLastElemIndex = directoryLength - 1
		coordinates.windowFirstElemIndex = coordinates.windowLastElemIndex - WindowSize
	} else {
		halfWindow := WindowSize / 2
		coordinates.windowFirstElemIndex = coordinates.selectedIndex - halfWindow
		coordinates.windowLastElemIndex = coordinates.selectedIndex + halfWindow
	}

	return coordinates
}

func IsInsideLastWindow(directoryLength int, selectedIndex int) bool {
	lastWindowElemIndex := directoryLength - 1
	firstWindowElemIndex := lastWindowElemIndex - WindowSize
	return IsInsideCurrentWindow(selectedIndex, firstWindowElemIndex, lastWindowElemIndex)
}

func IsInsideFirstWindow(selectedIndex int) bool {
	return IsInsideCurrentWindow(selectedIndex, 0, LastPos(0))
}

func IsInsideCurrentWindow(selectedIndex int, start int, end int) bool {
	return selectedIndex >= start &&
		selectedIndex <= end
}

func GetGap(directory []fs.FileInfo, height int) int {

	if len(directory)+10 > height {

		return 6
	}
	return height - len(directory) - 3
}

func PrintOnLine(name string, line int) {
	for i := 0; i < line; i++ {
		fmt.Println("")
	}
	fmt.Print(name)
}

func WillFit(directory []fs.FileInfo) bool {
	return len(directory)-1 < WindowSize
}

func IndexOf(files []fs.FileInfo, Selected string) int {
	for i, current := range files {
		if current.Name() == Selected {
			return i
		}
	}
	return -1
}

func ReadText() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func PrintDir(files []fs.FileInfo, coordinates Coordinates) {
	ClearConsole()
	width := 70
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	for i := coordinates.windowFirstElemIndex; i <= coordinates.windowLastElemIndex && i < len(files); i++ {

		dirName := files[i].Name()
		dirSize := strconv.Itoa(int(files[i].Size()))

		if coordinates.selectedName == dirName {
			fmt.Print(string(colorGreen), GetFullText(dirName, dirSize, width))
			fmt.Println(string(colorReset))
		} else {
			fmt.Println(GetFullText(dirName, dirSize, width))
		}
	}
}

func GetFullText(input string, fileSize string, maxLength int) string {

	if len(input)+len(fileSize)+5 >= maxLength {
		extraText := len(input) + len(fileSize) - maxLength
		input = input[0:len(input)-extraText-8] + "..."
	}
	availableSpace := maxLength - len(input) - len(fileSize)
	return input + strings.Repeat(" ", availableSpace) + fileSize + " KB"
}

func ResetCoordinates(directory []fs.FileInfo) Coordinates {
	selectedName := ""
	if len(directory) > 0 {
		selectedName = directory[0].Name()
	}
	return Coordinates{0, 0, LastPos(0), selectedName}
}

func LastPos(windowFirstElemIndex int) int {
	return WindowSize + windowFirstElemIndex
}

func ReadFiles(path string) []fs.FileInfo {
	files, _ := ioutil.ReadDir(path)
	return files
}
