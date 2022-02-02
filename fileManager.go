package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
	tsize "github.com/kopoli/go-terminal-size"
	cp "github.com/otiai10/copy"
)

func main() {

	RunProgram(`D:\`)
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
					fmt.Println("path is ", startDir)
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
				selectedIndex := IndexOf(mainDir, selected)
				if selectedIndex != -1 && mainDir[selectedIndex].IsDir() {

					fmt.Println("This : ", mainDir[selectedIndex].Name(), "is file??", mainDir[selectedIndex].IsDir())
					fmt.Print("\033[H\033[2J")
					startDir = startDir + `\` + selected
					selected = ""
					position = 0
					startPos = 0
					fmt.Println("empty")
				} else {
					fmt.Print("\033[H\033[2J")
					if filepath.Base(startDir) != selected {
						startDir = startDir + `\` + selected
					}
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
			if len(selectedForCopy) > 1 {
				selectedForCopy = ""
			}
			selectedForCut = startDir + `\` + selected
		}

		if char == 'v' { //copy
			if len(selectedForCut) > 1 {
				selectedForCut = ""
			}
			selectedForCopy = startDir + `\` + selected
		}

		if char == 'p' { //paste

			toCut, _ := os.Stat(selectedForCut)
			toCopy, _ := os.Stat(selectedForCopy)
			if len(selectedForCut) > 0 && selectedForCut != startDir+`\`+filepath.Base(selectedForCut) { //cut-paste file/folder
				selected = filepath.Base(selectedForCut)

				if toCut.IsDir() { //folder

				} else {

					fmt.Println(selected)
					err := os.Rename(selectedForCut, startDir+`\`+selected)
					if err != nil {
						log.Fatal(err)
					}

					mainDir = ReadFiles(startDir)
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
				}
			}

			if len(selectedForCopy) > 0 && selectedForCopy != startDir+`\`+filepath.Base(selectedForCopy) { //copy-paste file/folder
				selected = filepath.Base(selectedForCopy)
				if toCopy.IsDir() { //folder
					cp.Copy(startDir, startDir+`\`+selected)

				} else { //file
					data, _ := ioutil.ReadFile(selectedForCopy)
					ioutil.WriteFile(startDir+`\`+selected, data, 0644)
				}
				mainDir = ReadFiles(startDir)
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
			}

		}

		if char == 'd' {
		_:
			os.RemoveAll(startDir + `\` + selected)
			ReadFiles(startDir)
			selected = mainDir[position].Name()
			PrintDir(mainDir, selected, startPos, lastElemPos)
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
	return paths + `\` + selectedFolder
}

func ReadFiles(path string) []fs.FileInfo {
	files, _ := ioutil.ReadDir(path)
	return files
}
