package main

import (
	"fmt"
	"io/fs"

	"github.com/eiannone/keyboard"
)

func ExecuteComandsOnKeystroke(key keyboard.Key, startDir string, directory []fs.FileInfo, coordinates Coordinates, char rune) (string, []fs.FileInfo, Coordinates) {
	switch {
	case key == keyboard.KeyArrowLeft:
		startDir, directory, coordinates = ExitDirectory(startDir, coordinates)

	case key == keyboard.KeyArrowUp:
		if len(directory) < 1 {
			break
		}
		coordinates = GoUp(coordinates, directory)

	case key == keyboard.KeyArrowDown:

		if len(directory) < 1 {

			break
		}
		coordinates = GoDown(coordinates, directory)

	case key == keyboard.KeyArrowRight:
		startDir, coordinates, directory = EnterDirectory(startDir, coordinates, directory)

	case char == 'n':
		directory, coordinates = NewFolder(directory, coordinates, startDir)

	case char == 'm':
		directory, coordinates = NewFile(directory, coordinates, startDir)

	case char == 'r':
		directory, coordinates = Rename(directory, coordinates, startDir)

	case char == 'c':
		if len(SelectedForCopy) > 1 {
			SelectedForCopy = ""
		}
		SelectedForCut = UpdatePath(startDir, coordinates.selectedName)

	case char == 'v':
		if len(SelectedForCut) > 1 {
			SelectedForCut = ""
		}
		SelectedForCopy = UpdatePath(startDir, coordinates.selectedName)

	case char == 'p':
		if NameExists(startDir) {
			ClearConsole()
			fmt.Println("Name already exists")
			break
		}
		coordinates, directory = Paste(startDir, coordinates)

	case char == 'd':
		directory, coordinates = Delete(startDir, coordinates)

	case char == 'h':
		PrintHelp()
	}
	return startDir, directory, coordinates
}
