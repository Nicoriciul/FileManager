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

func main() {

	RunProgram(".")
}

func RunProgram(startDir string) {
	var selectedForCopy string
	var selectedForCut string
	var size tsize.Size
	size, _ = tsize.GetSize()
	lastElemPos := size.Height - 10
	mainDir := ReadFiles(startDir)
	selected := mainDir[0].Name()
	startPos := 0
	if lastElemPos >= len(mainDir)-1 {
		PrintDir(mainDir, selected, 0, lastElemPos)
	} else {
		PrintDir(mainDir, selected, 0, lastElemPos)
	}
	position := 0
	for {
		char, key, _ := keyboard.GetSingleKey()
		if key == keyboard.KeyArrowLeft { // --------------LEFT------------
			back := ReadFiles(goBack(startDir))
			startDir = strings.TrimRight(filepath.Dir(startDir), `/`)
			mainDir = back
			fmt.Println("selecetd is  -----------------", selected)
			fmt.Println("startdir is ---------------", startDir)
			fmt.Println("main dir is -------------", mainDir[0].Name(), "   length  ", len(mainDir))
			selected = mainDir[0].Name()
			position = 0
			startPos = 0
			if lastElemPos >= len(mainDir)-1 {
				//PrintDir(mainDir, selected, 0, lastElemPos)
			} else {
				//PrintDir(mainDir, selected, 0, lastElemPos)
			}
		}

		if key == keyboard.KeyArrowUp { // --------------UP------------
			arrayLength := len(mainDir)
			willFit := arrayLength-1 < lastElemPos
			position--
			if willFit {
				if position >= 0 {
					selected = mainDir[position].Name()
				} else {
					position = arrayLength - 1
					selected = mainDir[position].Name()
				}
				PrintDir(mainDir, selected, 0, arrayLength-1)

			} else {
				if position >= 0 {
					if position >= startPos {
						selected = mainDir[position].Name()
						PrintDir(mainDir, selected, startPos, startPos+lastElemPos)
					}
					if position < startPos {
						startPos--
						selected = mainDir[position].Name()
						PrintDir(mainDir, selected, startPos, startPos+lastElemPos)
					}
				} else {
					startPos = arrayLength - lastElemPos - 1
					position = arrayLength - 1
					selected = mainDir[position].Name()
					PrintDir(mainDir, selected, startPos, startPos+lastElemPos)
				}
			}
		}

		if key == keyboard.KeyArrowDown { // --------------DOWN------------
			position++
			arrayLength := len(mainDir)
			if arrayLength < 1 {

				break
			}
			willFit := arrayLength-1 < lastElemPos
			if willFit {
				if position < arrayLength {
					selected = mainDir[position].Name()
				} else {
					position = 0
					selected = mainDir[position].Name()
				}
				PrintDir(mainDir, selected, 0, arrayLength-1)

			} else {
				if position < arrayLength {
					if position <= startPos+lastElemPos {
						selected = mainDir[position].Name()
						PrintDir(mainDir, selected, startPos, startPos+lastElemPos)
					}
					if position > startPos+lastElemPos {
						startPos++
						selected = mainDir[position].Name()
						PrintDir(mainDir, selected, startPos, startPos+lastElemPos)
					}
				} else {
					startPos = 0
					position = 0
					selected = mainDir[position].Name()
					PrintDir(mainDir, selected, startPos, lastElemPos)

				}
			}
		}
		if key == keyboard.KeyArrowRight { // --------------RIGHT------------
			forward := ReadFiles(goForward(startDir, selected))
			selectedIndex := IndexOf(mainDir, selected)
			if len(forward) > 0 {
				mainDir = forward
				startDir = UpdatePath(startDir, selected)
				selected = mainDir[0].Name()
				position = 0
				startPos = 0
				if lastElemPos >= len(mainDir)-1 {
					PrintDir(mainDir, selected, 0, lastElemPos)
				} else {
					PrintDir(mainDir, selected, 0, lastElemPos)
				}
			} else {
				if selectedIndex != -1 && mainDir[selectedIndex].IsDir() {
					startDir = UpdatePath(startDir, selected)
					selected = ""
					position = 0
					startPos = 0
					ClearConsole()
					fmt.Println("empty")
				}
			}

		}
		if char == 'q' {
			return
		}

		if char == 'n' { //new folder

			folderName := "New Folder"
			err := os.Mkdir(UpdatePath(startDir, folderName), 0755)
			if err != nil {
				log.Fatal(err)
			}
			selected = folderName
			mainDir = ReadFiles(startDir)
			position, selected, startPos = Rename(mainDir, startDir, position, selected, startPos, lastElemPos, size)
		}

		if char == 'm' { //new file

			fileName := "New File"

			file, e := os.Create(UpdatePath(startDir, fileName))
			if e != nil {
				log.Fatal(e)
			}
			file.Close()
			position = IndexOf(mainDir, selected)
			selected = fileName
			mainDir = ReadFiles(startDir)
			position, selected, startPos = Rename(mainDir, startDir, position, selected, startPos, lastElemPos, size)
		}

		if char == 'r' { //rename
			position = IndexOf(mainDir, selected)
			position, selected, startPos = Rename(mainDir, startDir, position, selected, startPos, lastElemPos, size)
		}

		if char == 'c' { //cut
			if len(selectedForCopy) > 1 {
				selectedForCopy = ""
			}
			selectedForCut = UpdatePath(startDir, selected)
		}

		if char == 'v' { //copy
			if len(selectedForCut) > 1 {
				selectedForCut = ""
			}
			selectedForCopy = UpdatePath(startDir, selected)
		}

		if char == 'p' { //paste

			if len(selectedForCut) > 0 && selectedForCut != UpdatePath(startDir, filepath.Base(selectedForCut)) { //cut-paste file/folder
				selected = filepath.Base(selectedForCut)
				cp.Copy(selectedForCut, UpdatePath(startDir, selected))
				os.RemoveAll(selectedForCut)
				mainDir = ReadFiles(startDir)
				startPos, position = Print(mainDir, selected, lastElemPos, startPos, position)

			}
			if len(selectedForCopy) > 0 && selectedForCopy != UpdatePath(startDir, filepath.Base(selectedForCopy)) { //copy-paste file/folder
				selected := filepath.Base(selectedForCopy)
				cp.Copy(selectedForCopy, UpdatePath(startDir, selected))
				mainDir = ReadFiles(startDir)
				startPos, position = Print(mainDir, selected, lastElemPos, startPos, position)

			}
		}

		if char == 'd' {
			os.RemoveAll(UpdatePath(startDir, selected))
			mainDir = ReadFiles(startDir)
			if len(mainDir) > 0 {
				selected = mainDir[position].Name()
				PrintDir(mainDir, selected, startPos, lastElemPos)
			} else {
				selected = ""
				ClearConsole()
				fmt.Println("empty")
			}

		}

		if char == 'h' {
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

func UpdatePath(currentPath string, selected string) string {
	return currentPath + `/` + selected
}
func ClearConsole() {

	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
func Rename(mainDir []fs.FileInfo, startDir string, position int, selected string, startPos int, lastElemPos int, size tsize.Size) (int, string, int) {
	PrintDir(mainDir, selected, startPos, startPos+lastElemPos)
	PrintOnLine(selected, GetGap(mainDir, size.Height))
	var newName string
	_, key, _ := keyboard.GetSingleKey()
	if key == keyboard.KeyBackspace2 || key == keyboard.KeyBackspace {

		PrintDir(mainDir, selected, position, position+lastElemPos)
		PrintOnLine("", GetGap(mainDir, size.Height))
		newName = ReadText()

	} else if key == keyboard.KeyEnter {
		newName = selected
	}

	originalPath := UpdatePath(startDir, selected)
	newPath := UpdatePath(startDir, newName)

	e := os.Rename(originalPath, newPath)
	if e != nil {
		log.Fatal(e)
	}
	mainDir = ReadFiles(startDir)
	selected = newName
	position = IndexOf(mainDir, selected)
	startPos = position
	PrintDir(mainDir, selected, position, position+lastElemPos)
	return position, selected, startPos
}

func GetGap(maindir []fs.FileInfo, height int) int {

	if len(maindir)+10 > height {

		return 6
	}
	return height - len(maindir) - 3
}

func PrintOnLine(name string, line int) {
	for i := 0; i < line; i++ {
		fmt.Println("")
	}
	fmt.Print(name)
}
func Print(mainDir []fs.FileInfo, selected string, lastElemPos int, startPos int, position int) (int, int) {
	indexOfSel := IndexOf(mainDir, selected)
	if len(mainDir)-1 > lastElemPos {
		if indexOfSel >= startPos && indexOfSel <= lastElemPos {
			position = indexOfSel
			PrintDir(mainDir, selected, startPos, lastElemPos)

		} else {
			startPos = indexOfSel
			position = indexOfSel
			PrintDir(mainDir, selected, startPos, lastElemPos)
		}
	}
	PrintDir(mainDir, selected, startPos, lastElemPos)
	return startPos, position
}

func IndexOf(files []fs.FileInfo, selected string) int {
	for i, current := range files {
		if current.Name() == selected {
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

func PrintDir(files []fs.FileInfo, selected string, start int, end int) {
	ClearConsole()
	width := 70
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	for i := start; i <= end && i < len(files); i++ {

		dirName := files[i].Name()
		dirSize := strconv.Itoa(int(files[i].Size()))

		if selected == dirName {
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

func goForward(paths string, selectedFolder string) string {
	return paths + `/` + selectedFolder
}

func ReadFiles(path string) []fs.FileInfo {
	files, _ := ioutil.ReadDir(path)
	return files
}
