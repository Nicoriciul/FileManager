package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/eiannone/keyboard"
)

const lastElemPos = 24

func main() {

	RunProgram(`D:\Games`)
	fmt.Println("saved?")

}

func RunProgram(startDir string) {
	var selectedForCopy string
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
		if key == keyboard.KeyArrowLeft {
			back := ReadFiles(goBack(startDir))
			startDir = strings.TrimRight(filepath.Dir(startDir), `\`)
			mainDir = back
			selected = mainDir[0].Name()
			position = 0
			startPos = 0
			if lastElemPos >= len(mainDir)-1 {
				PrintDir(mainDir, selected, 0, lastElemPos)
			} else {
				PrintDir(mainDir, selected, 0, lastElemPos)
			}
		}

		if key == keyboard.KeyArrowUp {
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

		if key == keyboard.KeyArrowDown {
			arrayLength := len(mainDir)
			willFit := arrayLength-1 < lastElemPos
			position++
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

		if key == keyboard.KeyArrowRight {
			currentDir := returnDirectory(startDir)
			forward := ReadFiles(goForward(startDir, selected))
			mainDir = forward
			if len(mainDir) > 0 {
				startDir = startDir + `\` + selected
				selected = mainDir[0].Name()
				position = 0
				startPos = 0
				if lastElemPos >= len(mainDir)-1 {
					PrintDir(mainDir, selected, 0, lastElemPos)
				} else {
					PrintDir(mainDir, selected, 0, lastElemPos)
				}
			} else {
				selectedIndex := IndexOf(currentDir, selected)
				if selectedIndex != -1 && currentDir[selectedIndex].IsDir() {
					fmt.Print("\033[H\033[2J")
					fmt.Println(startDir + `\` + selected)
					startDir = startDir + `\` + selected
					selected = ""
					position = 0
					startPos = 0
					fmt.Println("empty")
				}
				if selectedIndex != -1 && !currentDir[selectedIndex].IsDir() {
					fmt.Print("\033[H\033[2J")
					startDir = startDir + `\` + selected
					fmt.Println("Can't open")
				}
			}
		}

		if char == 'n' { //new folder

			folderName := ReadText()
			if len(folderName) < 1 {
				folderName = "New Folder"
			}
			err := os.Mkdir(startDir+`\`+folderName, 0755)
			if err != nil {
				log.Fatal(err)
			}
			selected = folderName
			mainDir = ReadFiles(startDir)
			PrintDir(mainDir, selected, position, lastElemPos)
		}

		if char == 'm' { //new file
			fileName := ReadText()
			if len(fileName) < 1 {
				fileName = "New File"
			}
			_, e := os.Create(startDir + `\` + fileName)
			if e != nil {
				log.Fatal(e)
			}
			selected = fileName
			mainDir = ReadFiles(startDir)
			PrintDir(mainDir, selected, position, lastElemPos)
		}

		if char == 'r' { //rename
			fmt.Print("\033[H\033[2J")
			fmt.Print("Rename ", selected, " to : ")
			name := ReadText()
			originalPath := startDir + `\` + selected
			newPath := startDir + `\` + name
			e := os.Rename(originalPath, newPath)
			if e != nil {
				log.Fatal(e)
			}
			mainDir = ReadFiles(startDir)
			selected = name
			PrintDir(mainDir, selected, position, lastElemPos)
		}

		if char == 'c' { //cut
			selectedForCopy = startDir + `\` + selected
		}

		if char == 'p' { //paste
			if len(selectedForCopy) > 0 && selectedForCopy != startDir {

				err := os.Rename(selectedForCopy, startDir)
				if err != nil {
					log.Fatal(err)
				}

				mainDir = ReadFiles(startDir)
				selected = filepath.Base(selectedForCopy)
				indexOfSel := IndexOf(mainDir, selected)
				ReadFiles(startDir)
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
			}
		}
	}
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
	fmt.Print("\033[H\033[2J")
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	for i := start; i <= end && i < len(files); i++ {

		dirName := files[i].Name()
		if selected == dirName {
			fmt.Print(string(colorGreen), dirName)
			fmt.Println(string(colorReset))

		} else {
			fmt.Println(files[i].Name())
		}
	}
}

func goBack(path string) string {
	return filepath.Dir(path)
}

func goForward(paths string, selectedFolder string) string {
	return paths + `\` + selectedFolder
}

func returnDirectory(path string) []fs.FileInfo {
	files, _ := ioutil.ReadDir(path)
	return files
}

func ReadFiles(path string) []fs.FileInfo {
	files, _ := ioutil.ReadDir(path)
	return files
}
