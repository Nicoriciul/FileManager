package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

func ExecuteComandsOnKeystroke(key keyboard.Key, char rune, data *Data) {
	switch {

	case key == keyboard.KeyArrowDown || key == keyboard.KeyArrowUp:
		data.list.UpDown(key)
	case key == keyboard.KeyArrowLeft:
		ExitDirectory(data)
	case key == keyboard.KeyArrowRight:
		if len(data.directory) < 1 {
			return
		}
		EnterDirectory(data)
	case char == 'n':
		NewFolder(data)

	case char == 'm':
		NewFile(data)

	case char == 'r':
		Rename(data)

	case char == 'c':
		if len(SelectedForCopy) > 1 {
			SelectedForCopy = ""
		}
		if len(data.directory) > 0 {
			SelectedForCut = UpdatePath(data.path, data.directory[data.list.selectedIndex].Name())
		}

	case char == 'v':
		if len(SelectedForCut) > 1 {
			SelectedForCut = ""
		}
		if len(data.directory) > 0 {
			SelectedForCopy = UpdatePath(data.path, data.directory[data.list.selectedIndex].Name())
		}
	case char == 'p':
		if NameExists(data.path) {
			ClearConsole()
			fmt.Println("Name already exists")
			break
		}
		Paste(data)

	case char == 'd':
		Delete(data)
		// case char == 'h':
		// 	PrintHelp()
	}
}
