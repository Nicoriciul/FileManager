package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/eiannone/keyboard"
	tsize "github.com/kopoli/go-terminal-size"
)

type List struct {
	elements      []string
	winFirstElem  int
	selectedIndex int
}

func NewList(strings []string) *List {
	return &List{
		elements:      strings,
		winFirstElem:  0,
		selectedIndex: 0,
	}
}

var Size, _ = tsize.GetSize()
var WindowSize = Size.Height - 10

func (l *List) UpDown(key keyboard.Key) {
	if len(l.elements) < 1 {
		ClearConsole()
		fmt.Println("empty")
		return
	}
	switch {

	case key == keyboard.KeyArrowUp:

		l.GoUp()

	case key == keyboard.KeyArrowDown:
		l.GoDown()
	}
	l.Print()
}

func (l *List) InitialPrint() {
	l.Print()
}

func LastElem(winFirstElem int) int {
	return WindowSize + winFirstElem
}

func (l *List) Print() {
	ClearConsole()
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	for i := l.winFirstElem; i <= LastElem(l.winFirstElem) && i < len(l.elements); i++ {

		name := l.elements[i]
		if l.selectedIndex == i {
			fmt.Print(string(colorGreen), name)
			fmt.Println(string(colorReset))
		} else {
			fmt.Println(name)
		}
	}
}

func (l *List) GoUp() {
	l.selectedIndex--

	if l.WillFit() {
		l.GoUpWhenListFits()

	} else {
		l.GoUpWhenListDoesntFit()
	}
}

func (l *List) GoDown() {
	l.selectedIndex++
	if l.WillFit() {
		l.GoDownWhenListFits()
	} else {
		l.GoDownWhenListDoesntFit()
	}
}

func (l *List) GoDownWhenListFits() {
	if l.selectedIndex >= len(l.elements) {
		l.ResetIndexes()
	}
}

func (l *List) GoDownWhenListDoesntFit() {
	if l.selectedIndex > LastElem(l.winFirstElem) {
		l.ScrollDown()
		return
	} else if l.selectedIndex >= len(l.elements)-1 {
		l.ResetIndexes()
		return
	}
}

func (l *List) ScrollDown() {
	if l.selectedIndex > LastElem(l.winFirstElem) && LastElem(l.winFirstElem)+1 <= len(l.elements)-1 {
		l.winFirstElem++
	} else {
		l.ResetIndexes()
	}
}

func (l *List) GoUpWhenListFits() {
	if l.selectedIndex < 0 {
		l.selectedIndex = len(l.elements) - 1
	}
}

func (l *List) GoUpWhenListDoesntFit() {

	if l.selectedIndex < l.winFirstElem {

		l.ScrollUp()
	}

}

func (l *List) ScrollUp() {
	if l.selectedIndex < LastElem(l.winFirstElem) && l.winFirstElem-1 >= 0 {
		l.winFirstElem--
	} else {
		l.winFirstElem = len(l.elements) - WindowSize - 1
		l.selectedIndex = len(l.elements) - 1
	}
}

func (l *List) WillFit() bool {
	return len(l.elements)-1 < WindowSize
}

func (l *List) ResetIndexes() {
	if len(l.elements) > 0 {
		l.selectedIndex = 0
		l.winFirstElem = 0
	}
}

func ClearConsole() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func (l *List) UpdateCoordinates() {
	indexOfSelected := l.selectedIndex
	directoryLength := len(l.elements)
	if IsInsideCurrentWindow(indexOfSelected, l.winFirstElem, LastElem(l.winFirstElem)) {
		return
	}
	if IsInsideFirstWindow(indexOfSelected) {
		l.winFirstElem = 0
	}
	if IsInsideLastWindow(directoryLength, indexOfSelected) {
		lastElement := directoryLength - 1
		l.winFirstElem = lastElement - WindowSize
	} else {
		halfWindow := WindowSize / 2
		l.winFirstElem = l.selectedIndex - halfWindow
	}
}

func IsInsideLastWindow(directoryLength int, selectedIndex int) bool {
	lastWindowElemIndex := directoryLength - 1
	firstWindowElemIndex := lastWindowElemIndex - WindowSize
	return IsInsideCurrentWindow(selectedIndex, firstWindowElemIndex, lastWindowElemIndex)
}

func IsInsideFirstWindow(selectedIndex int) bool {
	return IsInsideCurrentWindow(selectedIndex, 0, LastElem(0))
}

func IsInsideCurrentWindow(selectedIndex int, start int, end int) bool {
	return selectedIndex >= start &&
		selectedIndex <= end
}

func (l *List) GetGap() int {

	if len(l.elements)+10 > WindowSize {

		return 6
	}
	return WindowSize - len(l.elements) - 3
}

func PrintOnLine(name string, line int) {
	for i := 0; i < line; i++ {
		fmt.Println("")
	}
	fmt.Print(name)
}

func GetFullText(input string, fileSize string, maxLength int) string {

	if len(input)+len(fileSize)+5 >= maxLength {
		extraText := len(input) + len(fileSize) - maxLength
		input = input[0:len(input)-extraText-8] + "..."
	}
	availableSpace := maxLength - len(input) - len(fileSize)
	return input + strings.Repeat(" ", availableSpace) + fileSize + " KB"
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
