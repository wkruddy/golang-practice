package main

import (
	"fmt"
	"log"
)

// TriangleRow maps to an array of values for a row in the triangle
type TriangleRow []int

// Triangle is a matrix of values
type Triangle []TriangleRow

func getCellValue(currentIndex int, previousRow TriangleRow) int {
	lookupLeft := currentIndex - 1
	lookupRight := currentIndex + 1

	cellSum := 0

	if lookupLeft >= 0 {
		cellSum += previousRow[lookupLeft]
	}

	if len(previousRow)-lookupRight >= 0 {
		cellSum += previousRow[lookupRight-1]
	}

	return cellSum
}

func buildRow(triangle Triangle, rowIndex int) TriangleRow {
	// Initialize it as the first row
	previousRow := TriangleRow{1}
	if rowIndex-1 > 0 {
		previousRow = triangle[rowIndex-1]
	}

	currentRow := TriangleRow{}

	for i := 0; i <= rowIndex; i++ {
		if rowIndex == 0 {
			currentRow = previousRow
			break
		}

		cellRowValue := getCellValue(i, previousRow)
		currentRow = append(currentRow, cellRowValue)
	}

	return currentRow
}

func buildTriangle(rowEnd int) Triangle {
	triangle := Triangle{}

	for i := 0; i <= rowEnd-1; i++ {
		currentRow := buildRow(triangle, i)
		triangle = append(triangle, currentRow)
	}

	return triangle
}

func showTriangle(triangle Triangle) {
	rowStr := "\n \n \n"
	for _, row := range triangle {
		rowStr = fmt.Sprintf(rowStr+"%v \n", row)
	}
	log.Println(rowStr)
}

func main() {

	triangleMatrix := buildTriangle(10)
	showTriangle(triangleMatrix)
}
