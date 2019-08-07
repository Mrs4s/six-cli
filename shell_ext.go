package main

import (
	"github.com/Mrs4s/readline"
	"github.com/Mrs4s/six-cli/models"
	"github.com/mattn/go-runewidth"
	"math"
	"strings"
	"syscall"
)

//参考 https://github.com/acarl005/textcol/blob/master/textcol.go 修复重置
func PrintColumns(strs []string, margin int) {
	maxLength := 0
	marginStr := strings.Repeat(" ", margin)
	var lengths []int
	for _, str := range strs {
		length := runewidth.StringWidth(str)
		maxLength = models.Max(maxLength, length)
		lengths = append(lengths, length)
	}
	fd, _ := syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	width, _, _ := readline.GetSize(int(fd))
	width = int(float32(width) * 1.1)
	numCols, numRows := calculateTableSize(width, margin, maxLength, len(strs))
	if numCols == 1 {
		for _, str := range strs {
			shell.Println(str)
		}
		return
	}
	for i := 0; i < numCols*numRows; i++ {
		x, y := rowIndexToTableCoords(i, numCols)
		j := tableCoordsToColIndex(x, y, numRows)
		strLen := 0
		str := ""
		if j < len(lengths) {
			strLen = lengths[j]
			str = (strs)[j]
		}
		numSpacesRequired := maxLength - strLen
		spaceStr := strings.Repeat(" ", numSpacesRequired)
		shell.Print(str)
		if x+1 == numCols {
			shell.Print("\n")
		} else {
			shell.Print(spaceStr)
			shell.Print(marginStr)
		}
	}
}

func PrintTables(table [][]string, margin int) {
	var maxLens []int
	for i, col := range table {
		for j, row := range col {
			if i == 0 {
				maxLens = append(maxLens, runewidth.StringWidth(row))
				continue
			}
			length := runewidth.StringWidth(row)
			if maxLens[j] < length {
				maxLens[j] = length
			}
		}
	}
	for _, col := range table {
		for i, row := range col {
			shell.Print(row)
			if i != len(col)-1 {
				shell.Print(strings.Repeat(" ", maxLens[i]-runewidth.StringWidth(row)+margin))
			}
		}
		shell.Println()
	}
}

func calculateTableSize(width, margin, maxLength, numCells int) (int, int) {
	numCols := (width + margin) / (maxLength + margin)
	if numCols == 0 {
		numCols = 1
	}
	numRows := int(math.Ceil(float64(numCells) / float64(numCols)))
	return numCols, numRows
}

func rowIndexToTableCoords(i, numCols int) (int, int) {
	x := i % numCols
	y := i / numCols
	return x, y
}

func tableCoordsToColIndex(x, y, numRows int) int {
	return y + numRows*x
}
