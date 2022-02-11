package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func InitialRead(startDir string) ([]fs.FileInfo, Coordinates) {
	directory := ReadFiles(startDir)
	coordinates := ResetCoordinates(directory)
	PrintDir(directory, coordinates)
	return directory, coordinates
}

func ClearConsole() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
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
