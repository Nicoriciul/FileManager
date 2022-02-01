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

	RunProgram(`D:\Games\`)

}

func RunProgram(startDir string) {
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
		_, key, _ := keyboard.GetSingleKey()
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
			currentDirAsStrings := ReadFiles(startDir)
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
				selectedIndex := IndexOf(currentDirAsStrings, selected)
				if selectedIndex != -1 && currentDir[selectedIndex].IsDir() {
					startDir = startDir + `\` + selected
					selected = ""
					position = 0
					startPos = 0
					fmt.Print("\033[H\033[2J")
					fmt.Println("empty")
				}
				if selectedIndex != -1 && !currentDir[selectedIndex].IsDir() {
					fmt.Print("\033[H\033[2J")
					fmt.Println("Can't open")
				}
			}
		}

		if key == keyboard.KeyCtrlN {
			err := os.Mkdir(`D:\Games\`+ReadText(), 0755)
			if err != nil {
				log.Fatal(err)
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

func MakeNewFolder(destination string) {
	myfile, e := os.Create(destination)
	if e != nil {
		log.Fatal(e)
	}
	myfile.Close()
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
