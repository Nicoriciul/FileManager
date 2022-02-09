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

	for {
		switch char, key, _ := keyboard.GetSingleKey(); {
		case key == keyboard.KeyArrowLeft: // --------------LEFT------------
			back := ReadFiles(filepath.Dir(startDir))
			startDir = strings.TrimRight(filepath.Dir(startDir), `/`)
			directory = back
			coordinates = ResetCoordinates(directory)
			PrintDir(directory, coordinates)

		case key == keyboard.KeyArrowUp: // --------------UP------------
			if len(directory) < 1 {
				break
			}
			willFit := len(directory)-1 < WindowSize
			coordinates.selectedIndex--
			if willFit {
				coordinates = GoUpWhenListFits(directory, coordinates)

			} else {
				coordinates = GoUpWhenListDoesntFit(directory, coordinates)
			}

		case key == keyboard.KeyArrowDown: // --------------DOWN------------

			if len(directory) < 1 {

				break
			}
			willFit := len(directory)-1 < WindowSize
			coordinates.selectedIndex++
			if willFit {
				coordinates = GoDownWhenListFits(directory, coordinates)
			} else {
				coordinates = GoDownWhenListDoesntFit(directory, coordinates)
			}

		case key == keyboard.KeyArrowRight: // --------------RIGHT------------
			forward := ReadFiles(goForward(startDir, coordinates.selectedName))
			SelectedIndex := IndexOf(directory, coordinates.selectedName)
			if len(forward) > 0 {
				directory = forward
				startDir = UpdatePath(startDir, coordinates.selectedName)
				coordinates = ResetCoordinates(directory)
				PrintDir(directory, coordinates)

			} else {
				if SelectedIndex != -1 && directory[SelectedIndex].IsDir() {
					startDir = UpdatePath(startDir, coordinates.selectedName)
					directory = forward
					coordinates.selectedIndex = 0
					coordinates.windowFirstElemIndex = 0
					coordinates.windowLastElemIndex = LastPos(0)
					ClearConsole()
					fmt.Println("empty")
				}
			}

		case char == 'q':
			return

		case char == 'n': //new folder

			folderName := "New Folder"
			err := os.Mkdir(UpdatePath(startDir, folderName), 0755)
			if err != nil {
				log.Fatal(err)
			}
			coordinates.selectedName = folderName
			directory = ReadFiles(startDir)
			coordinates = Rename(directory, startDir, coordinates)

		case char == 'm': //new file

			fileName := "New File"

			file, e := os.Create(UpdatePath(startDir, fileName))
			if e != nil {
				log.Fatal(e)
			}
			file.Close()
			coordinates.selectedIndex = IndexOf(directory, coordinates.selectedName)
			coordinates.selectedName = fileName
			directory = ReadFiles(startDir)
			coordinates = Rename(directory, startDir, coordinates)

		case char == 'r': //rename
			coordinates.selectedIndex = IndexOf(directory, coordinates.selectedName)
			coordinates = Rename(directory, startDir, coordinates)

		case char == 'c': //cut
			if len(SelectedForCopy) > 1 {
				SelectedForCopy = ""
			}
			SelectedForCut = UpdatePath(startDir, coordinates.selectedName)

		case char == 'v': //copy
			if len(SelectedForCut) > 1 {
				SelectedForCut = ""
			}
			SelectedForCopy = UpdatePath(startDir, coordinates.selectedName)

		case char == 'p': //paste
			coordinates.selectedName = filepath.Base(SelectedForCut)
			cp.Copy(SelectedForCut, UpdatePath(startDir, coordinates.selectedName))
			if len(SelectedForCut) > 0 && SelectedForCut != UpdatePath(startDir, filepath.Base(SelectedForCut)) { //cut-paste file/folder

				os.RemoveAll(SelectedForCut)
				directory = ReadFiles(startDir)
				coordinates = Print(directory, coordinates)

			}
			if len(SelectedForCopy) > 0 && SelectedForCopy != UpdatePath(startDir, filepath.Base(SelectedForCopy)) { //copy-paste file/folder
				Selected := filepath.Base(SelectedForCopy)
				cp.Copy(SelectedForCopy, UpdatePath(startDir, Selected))
				directory = ReadFiles(startDir)
				coordinates = Print(directory, coordinates)
			}

		case char == 'd':
			os.RemoveAll(UpdatePath(startDir, coordinates.selectedName))
			directory = ReadFiles(startDir)
			if len(directory) > 0 {
				coordinates.selectedName = directory[coordinates.selectedIndex].Name()
				PrintDir(directory, coordinates)
			} else {
				coordinates.selectedName = ""
				ClearConsole()
				fmt.Println("empty")
			}

		case char == 'h':
			colorReset := "\033[0m"
			colorRed := "\033[31m"
			ClearConsole()
			commands := []string{"COPY = V", "CUT = C", "PASTE = P", "RENAME = R", "NEW FILE = M", "NEW FOLDER = N", "DELETE = D", "HELP = H", "QUIT = Q"}
			for _, current := range commands {
				fmt.Println(string(colorRed), current)
			}
			fmt.Print(colorReset)
		}
	}
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
func Rename(directory []fs.FileInfo, startDir string, coordinates Coordinates) Coordinates {
	PrintDir(directory, coordinates)
	PrintOnLine(coordinates.selectedName, GetGap(directory, Size.Height))
	var newName string
	_, key, _ := keyboard.GetSingleKey()
	if key == keyboard.KeyBackspace2 || key == keyboard.KeyBackspace {

		PrintDir(directory, coordinates)
		PrintOnLine("", GetGap(directory, Size.Height))
		newName = ReadText()

	} else if key == keyboard.KeyEnter {
		newName = coordinates.selectedName
	}

	originalPath := UpdatePath(startDir, coordinates.selectedName)
	newPath := UpdatePath(startDir, newName)

	e := os.Rename(originalPath, newPath)
	if e != nil {
		log.Fatal(e)
	}
	directory = ReadFiles(startDir)
	coordinates.selectedName = newName
	coordinates.selectedIndex = IndexOf(directory, coordinates.selectedName)
	coordinates.windowFirstElemIndex = coordinates.selectedIndex
	PrintDir(directory, coordinates)

	return coordinates
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
func Print(directory []fs.FileInfo, coordinates Coordinates) Coordinates {
	indexOfSel := IndexOf(directory, coordinates.selectedName)
	if len(directory)-1 > WindowSize {
		if indexOfSel >= coordinates.windowFirstElemIndex && indexOfSel <= WindowSize {
			coordinates.selectedIndex = indexOfSel
			PrintDir(directory, coordinates)

		} else {
			coordinates.windowFirstElemIndex = indexOfSel
			coordinates.selectedIndex = indexOfSel
			PrintDir(directory, coordinates)

		}
	}
	PrintDir(directory, coordinates)
	return coordinates
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

func goBack(path string) string {
	return filepath.Dir(path)
}

func goForward(paths string, SelectedFolder string) string {
	return paths + `/` + SelectedFolder
}

func ResetCoordinates(directory []fs.FileInfo) Coordinates {
	return Coordinates{0, 0, LastPos(0), directory[0].Name()}
}

func LastPos(windowFirstElemIndex int) int {
	return Size.Height - 10 + windowFirstElemIndex
}

func ReadFiles(path string) []fs.FileInfo {
	files, _ := ioutil.ReadDir(path)
	return files
}
