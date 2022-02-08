package main

import (
	"bufio"
	"fmt"
	"io/fs"
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

func main() {

	RunProgram(".")

}

func PrintDirectory(dir []fs.FileInfo) {
	for i := 0; i < len(dir); i++ {
		fmt.Println(dir[i].Name())
	}
}

func RunProgram(startDir string) {
	var SelectedForCopy string
	var SelectedForCut string
	ResetCoordinates(Directory)
	PrintDir(Directory, Selected, 0, LastPos())

	for {
		switch char, key, _ := keyboard.GetSingleKey(); {
		case key == keyboard.KeyArrowLeft: // --------------LEFT------------
			back := ReadFiles(filepath.Dir(startDir))
			startDir = strings.TrimRight(filepath.Dir(startDir), `/`)
			Directory = back
			ResetCoordinates(Directory)
			PrintDir(Directory, Selected, StartPos, LastPos())

		case key == keyboard.KeyArrowUp: // --------------UP------------
			if len(Directory) < 1 {
				break
			}
			willFit := len(Directory)-1 < WindowSize
			CurrentPosition--
			if willFit {
				GoUpWhenListFits()

			} else {
				GoUpWhenListDoesntFit()
			}

		case key == keyboard.KeyArrowDown: // --------------DOWN------------

			if len(Directory) < 1 {

				break
			}
			willFit := len(Directory)-1 < WindowSize
			CurrentPosition++
			if willFit {
				GoDownWhenListFits()
			} else {
				if CurrentPosition < len(Directory) {
					if CurrentPosition <= LastPos() {
						Selected = Directory[CurrentPosition].Name()
						PrintDir(Directory, Selected, StartPos, LastPos())
					}
					if CurrentPosition > LastPos() {
						StartPos++
						Selected = Directory[CurrentPosition].Name()
						PrintDir(Directory, Selected, StartPos, LastPos())
					}
				} else {
					StartPos = 0
					CurrentPosition = 0
					Selected = Directory[CurrentPosition].Name()
					PrintDir(Directory, Selected, StartPos, WindowSize)

				}
			}

		case key == keyboard.KeyArrowRight: // --------------RIGHT------------
			forward := ReadFiles(goForward(startDir, Selected))
			SelectedIndex := IndexOf(Directory, Selected)
			if len(forward) > 0 {
				Directory = forward
				startDir = UpdatePath(startDir, Selected)
				ResetCoordinates(Directory)
				PrintDir(Directory, Selected, 0, WindowSize)

			} else {
				if SelectedIndex != -1 && Directory[SelectedIndex].IsDir() {
					startDir = UpdatePath(startDir, Selected)
					Directory = forward
					CurrentPosition = 0
					StartPos = 0
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
			Selected = folderName
			Directory = ReadFiles(startDir)
			CurrentPosition, Selected, StartPos = Rename(Directory, startDir, CurrentPosition, Selected, StartPos, WindowSize, Size)

		case char == 'm': //new file

			fileName := "New File"

			file, e := os.Create(UpdatePath(startDir, fileName))
			if e != nil {
				log.Fatal(e)
			}
			file.Close()
			CurrentPosition = IndexOf(Directory, Selected)
			Selected = fileName
			Directory = ReadFiles(startDir)
			CurrentPosition, Selected, StartPos = Rename(Directory, startDir, CurrentPosition, Selected, StartPos, WindowSize, Size)

		case char == 'r': //rename
			CurrentPosition = IndexOf(Directory, Selected)
			CurrentPosition, Selected, StartPos = Rename(Directory, startDir, CurrentPosition, Selected, StartPos, WindowSize, Size)

		case char == 'c': //cut
			if len(SelectedForCopy) > 1 {
				SelectedForCopy = ""
			}
			SelectedForCut = UpdatePath(startDir, Selected)

		case char == 'v': //copy
			if len(SelectedForCut) > 1 {
				SelectedForCut = ""
			}
			SelectedForCopy = UpdatePath(startDir, Selected)

		case char == 'p': //paste

			if len(SelectedForCut) > 0 && SelectedForCut != UpdatePath(startDir, filepath.Base(SelectedForCut)) { //cut-paste file/folder
				Selected = filepath.Base(SelectedForCut)
				cp.Copy(SelectedForCut, UpdatePath(startDir, Selected))
				os.RemoveAll(SelectedForCut)
				Directory = ReadFiles(startDir)
				StartPos, CurrentPosition = Print(Directory, Selected, WindowSize, StartPos, CurrentPosition)

			}
			if len(SelectedForCopy) > 0 && SelectedForCopy != UpdatePath(startDir, filepath.Base(SelectedForCopy)) { //copy-paste file/folder
				Selected := filepath.Base(SelectedForCopy)
				cp.Copy(SelectedForCopy, UpdatePath(startDir, Selected))
				Directory = ReadFiles(startDir)
				StartPos, CurrentPosition = Print(Directory, Selected, WindowSize, StartPos, CurrentPosition)
			}

		case char == 'd':
			os.RemoveAll(UpdatePath(startDir, Selected))
			Directory = ReadFiles(startDir)
			if len(Directory) > 0 {
				Selected = Directory[CurrentPosition].Name()
				PrintDir(Directory, Selected, StartPos, WindowSize)
			} else {
				Selected = ""
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

func GoDownWhenListFits() {
	if CurrentPosition < len(Directory) {
		Selected = Directory[CurrentPosition].Name()
	} else {
		CurrentPosition = 0
		Selected = Directory[CurrentPosition].Name()
	}
	PrintDir(Directory, Selected, 0, len(Directory)-1)
}

func GoDownWhenListDoesntFit() {

	if CurrentPosition <= LastPos() {
		Selected = Directory[CurrentPosition].Name()
		PrintDir(Directory, Selected, StartPos, LastPos())
		return
	}
	ScrollDown()
}

func ScrollDown() {
	if CurrentPosition > LastPos() && LastPos()+1 <= len(Directory)-1 {
		StartPos++
		Selected = Directory[CurrentPosition].Name()
		PrintDir(Directory, Selected, StartPos, LastPos())
	} else {

		StartPos = 0
		CurrentPosition = 0
		Selected = Directory[CurrentPosition].Name()
		PrintDir(Directory, Selected, StartPos, WindowSize)
	}
}

func GoUpWhenListFits() {
	if CurrentPosition >= 0 {
		Selected = Directory[CurrentPosition].Name()
	} else if len(Directory) > 0 {
		CurrentPosition = len(Directory) - 1
		Selected = Directory[CurrentPosition].Name()
	}
	PrintDir(Directory, Selected, 0, len(Directory)-1)
}

func GoUpWhenListDoesntFit() {

	if CurrentPosition >= StartPos {
		Selected = Directory[CurrentPosition].Name()
		PrintDir(Directory, Selected, StartPos, LastPos())
		return
	}
	ScrollUp()
}

func ScrollUp() {
	if CurrentPosition < StartPos && StartPos-1 >= 0 {
		StartPos--
		Selected = Directory[CurrentPosition].Name()
		PrintDir(Directory, Selected, StartPos, LastPos())
	} else {
		StartPos = len(Directory) - WindowSize - 1
		CurrentPosition = len(Directory) - 1
		Selected = Directory[CurrentPosition].Name()
		PrintDir(Directory, Selected, StartPos, LastPos())
	}
}

func UpdatePath(currentPath string, Selected string) string {
	return currentPath + `/` + Selected
}
func ClearConsole() {

	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
func Rename(Directory []fs.FileInfo, startDir string, CurrentPosition int, Selected string, StartPos int, WindowSize int, size tsize.Size) (int, string, int) {
	PrintDir(Directory, Selected, StartPos, LastPos())
	PrintOnLine(Selected, GetGap(Directory, size.Height))
	var newName string
	_, key, _ := keyboard.GetSingleKey()
	if key == keyboard.KeyBackspace2 || key == keyboard.KeyBackspace {

		PrintDir(Directory, Selected, CurrentPosition, LastPos())
		PrintOnLine("", GetGap(Directory, size.Height))
		newName = ReadText()

	} else if key == keyboard.KeyEnter {
		newName = Selected
	}

	originalPath := UpdatePath(startDir, Selected)
	newPath := UpdatePath(startDir, newName)

	e := os.Rename(originalPath, newPath)
	if e != nil {
		log.Fatal(e)
	}
	Directory = ReadFiles(startDir)
	Selected = newName
	CurrentPosition = IndexOf(Directory, Selected)
	StartPos = CurrentPosition
	PrintDir(Directory, Selected, CurrentPosition, WindowSize)

	return CurrentPosition, Selected, StartPos
}

func GetGap(Directory []fs.FileInfo, height int) int {

	if len(Directory)+10 > height {

		return 6
	}
	return height - len(Directory) - 3
}

func PrintOnLine(name string, line int) {
	for i := 0; i < line; i++ {
		fmt.Println("")
	}
	fmt.Print(name)
}
func Print(Directory []fs.FileInfo, Selected string, WindowSizeint, StartPos int, CurrentPosition int) (int, int) {
	indexOfSel := IndexOf(Directory, Selected)
	if len(Directory)-1 > WindowSize {
		if indexOfSel >= StartPos && indexOfSel <= WindowSize {
			CurrentPosition = indexOfSel
			PrintDir(Directory, Selected, StartPos, WindowSize)

		} else {
			StartPos = indexOfSel
			CurrentPosition = indexOfSel
			PrintDir(Directory, Selected, StartPos, WindowSize)

		}
	}
	PrintDir(Directory, Selected, StartPos, WindowSize)
	return StartPos, CurrentPosition
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

func PrintDir(files []fs.FileInfo, Selected string, start int, end int) {
	ClearConsole()
	width := 70
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	for i := start; i <= end && i < len(files); i++ {

		dirName := files[i].Name()
		dirSize := strconv.Itoa(int(files[i].Size()))

		if Selected == dirName {
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
