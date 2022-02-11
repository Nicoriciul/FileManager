package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func ExitDirectory(startDir string, coordinates Coordinates) (string, []fs.FileInfo, Coordinates) {
	back := ReadFiles(filepath.Dir(startDir))
	startDir = strings.TrimRight(filepath.Dir(startDir), `/`)
	directory := back
	coordinates = ResetCoordinates(directory)
	PrintDir(directory, coordinates)
	return startDir, directory, coordinates
}

func GoUp(coordinates Coordinates, directory []fs.FileInfo) Coordinates {
	coordinates.selectedIndex--
	if WillFit(directory) {
		coordinates = GoUpWhenListFits(directory, coordinates)

	} else {
		coordinates = GoUpWhenListDoesntFit(directory, coordinates)
	}
	return coordinates
}

func GoDown(coordinates Coordinates, directory []fs.FileInfo) Coordinates {
	coordinates.selectedIndex++
	if WillFit(directory) {
		coordinates = GoDownWhenListFits(directory, coordinates)
	} else {
		coordinates = GoDownWhenListDoesntFit(directory, coordinates)
	}
	return coordinates
}

func EnterDirectory(startDir string, coordinates Coordinates, directory []fs.FileInfo) (string, Coordinates, []fs.FileInfo) {
	forward := ReadFiles(UpdatePath(startDir, coordinates.selectedName))
	SelectedIndex := IndexOf(directory, coordinates.selectedName)
	if len(forward) > 0 {
		directory = forward
		startDir = UpdatePath(startDir, coordinates.selectedName)
		coordinates = ResetCoordinates(directory)
		PrintDir(directory, coordinates)

	} else if SelectedIndex != -1 && directory[SelectedIndex].IsDir() {
		startDir = UpdatePath(startDir, coordinates.selectedName)
		directory = forward
		coordinates = ResetCoordinates(directory)
		ClearConsole()
		fmt.Println("empty")
	}
	return startDir, coordinates, directory
}

func GoDownWhenListFits(directory []fs.FileInfo, coordinates Coordinates) Coordinates {
	if coordinates.selectedIndex < len(directory) {
		coordinates.selectedName = directory[coordinates.selectedIndex].Name()
	} else {
		coordinates = ResetCoordinates(directory)
	}
	PrintDir(directory, coordinates)
	return coordinates
}

func GoDownWhenListDoesntFit(directory []fs.FileInfo, coordinates Coordinates) Coordinates {
	if coordinates.selectedIndex <= LastPos(coordinates.windowFirstElemIndex) {
		coordinates.selectedName = directory[coordinates.selectedIndex].Name()
		PrintDir(directory, coordinates)
		return coordinates
	}
	coordinates = ScrollDown(directory, coordinates)
	return coordinates
}

func ScrollDown(directory []fs.FileInfo, coordinates Coordinates) Coordinates {
	if coordinates.selectedIndex > LastPos(coordinates.windowFirstElemIndex) && LastPos(coordinates.windowFirstElemIndex)+1 <= len(directory)-1 {
		coordinates.windowFirstElemIndex++
		coordinates.windowLastElemIndex = LastPos(coordinates.windowFirstElemIndex)
		coordinates.selectedName = directory[coordinates.selectedIndex].Name()
		PrintDir(directory, coordinates)
	} else {
		coordinates = ResetCoordinates(directory)
		PrintDir(directory, coordinates)
	}
	return coordinates
}

func GoUpWhenListFits(directory []fs.FileInfo, coordinates Coordinates) Coordinates {
	if coordinates.selectedIndex >= 0 {
		coordinates.selectedName = directory[coordinates.selectedIndex].Name()
	} else if len(directory) > 0 {
		coordinates.selectedIndex = len(directory) - 1
		coordinates.selectedName = directory[coordinates.selectedIndex].Name()
	}
	PrintDir(directory, coordinates)
	return coordinates
}

func GoUpWhenListDoesntFit(directory []fs.FileInfo, coordinates Coordinates) Coordinates {

	if coordinates.selectedIndex >= coordinates.windowFirstElemIndex {
		coordinates.selectedName = directory[coordinates.selectedIndex].Name()
		coordinates.windowLastElemIndex = LastPos(coordinates.windowFirstElemIndex)
		PrintDir(directory, coordinates)
		return coordinates
	}
	coordinates = ScrollUp(directory, coordinates)
	return coordinates

}

func ScrollUp(directory []fs.FileInfo, coordinates Coordinates) Coordinates {
	if coordinates.selectedIndex < coordinates.windowLastElemIndex && coordinates.windowFirstElemIndex-1 >= 0 {
		coordinates.windowFirstElemIndex--
	} else {
		coordinates.windowFirstElemIndex = len(directory) - WindowSize - 1
		coordinates.selectedIndex = len(directory) - 1
	}
	coordinates.windowLastElemIndex = LastPos(coordinates.windowFirstElemIndex)
	coordinates.selectedName = directory[coordinates.selectedIndex].Name()
	PrintDir(directory, coordinates)
	return coordinates
}

func IndexOf(files []fs.FileInfo, Selected string) int {
	for i, current := range files {
		if current.Name() == Selected {
			return i
		}
	}
	return -1
}
