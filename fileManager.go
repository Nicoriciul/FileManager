package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/eiannone/keyboard"
)

const lastElemPos = 24

func main() {
	RunProgram(`D:\Games`)

}

func RunProgram(startDir string) {
	mainDir, _ := ReadFiles(startDir)
	selected := mainDir[0]
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
			back, _ := ReadFiles(goBack(startDir))
			startDir = filepath.Dir(startDir)
			mainDir = back
			selected = mainDir[0]
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
					selected = mainDir[position]
				} else {
					position = arrayLength - 1
					selected = mainDir[position]
				}
				PrintDir(mainDir, selected, 0, arrayLength-1)

			} else {
				if position >= 0 {
					if position >= startPos {
						selected = mainDir[position]
						PrintDir(mainDir, selected, startPos, startPos+lastElemPos)
					}
					if position < startPos {
						startPos--
						selected = mainDir[position]
						PrintDir(mainDir, selected, startPos, startPos+lastElemPos)
					}
				} else {
					startPos = arrayLength - lastElemPos - 1
					position = arrayLength - 1
					selected = mainDir[position]
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
					selected = mainDir[position]
				} else {
					position = 0
					selected = mainDir[position]
				}
				PrintDir(mainDir, selected, 0, arrayLength-1)

			} else {
				if position < arrayLength {
					if position <= startPos+lastElemPos {
						selected = mainDir[position]
						PrintDir(mainDir, selected, startPos, startPos+lastElemPos)
					}
					if position > startPos+lastElemPos {
						startPos++
						selected = mainDir[position]
						PrintDir(mainDir, selected, startPos, startPos+lastElemPos)
					}
				} else {
					startPos = 0
					position = 0
					selected = mainDir[position]
					PrintDir(mainDir, selected, startPos, lastElemPos)

				}
			}
		}

		if key == keyboard.KeyArrowRight {
			currentDir := returnDirectory(startDir)
			currentDirAsStrings, _ := ReadFiles(startDir)
			forward, _ := ReadFiles(goForward(startDir, selected))
			mainDir = forward
			if len(mainDir) > 0 {
				startDir = startDir + `\` + selected
				selected = mainDir[0]
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
			if err := os.Mkdir(ReadText(), os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func IndexOf(files []string, selected string) int {
	for i, current := range files {
		if current == selected {
			return i
		}
	}
	return -1
}

func ReadText() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

func MakeNewFolder(destination string) {
	myfile, e := os.Create(destination)
	if e != nil {
		log.Fatal(e)
	}
	myfile.Close()
}

func PrintDir(fileNames []string, selected string, start int, end int) {
	fmt.Print("\033[H\033[2J")
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	for i := start; i <= end && i < len(fileNames); i++ {

		dirName := fileNames[i]
		if selected == dirName {
			fmt.Print(string(colorGreen), dirName)
			fmt.Println(string(colorReset))

		} else {
			fmt.Println(fileNames[i])
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

func ReadFiles(path string) ([]string, []string) {
	var updatedPaths []string
	var fileNames []string
	files, _ := ioutil.ReadDir(path)
	for i := 0; i < len(files); i++ {
		updatedPaths = append(updatedPaths, path+`\`+files[i].Name())
		fileNames = append(fileNames, files[i].Name())
	}
	return fileNames, updatedPaths
}
