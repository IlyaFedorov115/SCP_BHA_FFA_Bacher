package parser_scp

import (
	"fmt"
	"io"
	"os"
)

func scanFromReader(reader io.Reader, template string,
	vals ...interface{}) (int, error) {
	return fmt.Fscanf(reader, template, vals...)
}
func scanSingle(reader io.Reader, val interface{}) (int,
	error) {
	return fmt.Fscan(reader, val)
}

func ParseScp(filename string) (table [][2]int, costs []float64,
	alpha, betta map[int][]int, problem error) {
	file, err := os.Open(filename)
	if err == nil {
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				problem = err
				fmt.Println("Problem with closing file ", file.Name())
			}
		}(file)

		var numRows, numColumns int
		var scanTemplate string = "%d %d"
		_, err = scanFromReader(file, scanTemplate, &numRows, &numColumns)
		if err != nil {
			problem = err
			fmt.Println("Problem with reading num of rows/columns!")
			return
		} else {
			//fmt.Println("Num rows", numRows, "Num columns", numColumns)
		}
		table = make([][2]int, 0, numColumns)
		costs = make([]float64, numColumns, numColumns)
		var tmpFloat float64

		scanTemplate = "%f"
		for i := 0; i < numColumns; i++ {
			_, err = scanSingle(file, &tmpFloat)
			if err != nil {
				problem = err
				return
			} else {
				costs[i] = tmpFloat
			}
		}

		var numCol int
		var tmpInt int

		alpha = make(map[int][]int, numRows)
		betta = make(map[int][]int, numColumns)
		for i := 0; i < numRows; i++ {
			_, problem = scanSingle(file, &numCol)
			if problem != nil {
				return
			} else {
				alpha[i] = make([]int, numCol)
				for j := 0; j < numCol; j++ {
					_, problem = scanSingle(file, &tmpInt)
					table = append(table, [2]int{i, tmpInt})
					alpha[i][j] = tmpInt
					betta[tmpInt] = append(betta[tmpInt], i)
				}
			}
		}
	}
	return
}

func ParseRail(filename string) (table [][2]int, costs []float64,
	alpha, betta map[int][]int, problem error) {

	file, err := os.Open(filename)
	if err == nil {
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				problem = err
				fmt.Println("Problem with closing file ", file.Name())
			}
		}(file)

		var numRows, numColumns int
		var scanTemplate string = "%d %d"
		_, err = scanFromReader(file, scanTemplate, &numRows, &numColumns)
		if err != nil {
			problem = err
			fmt.Println("Problem with reading num of rows/columns!")
			return
		} else {
			fmt.Println("Num rows", numRows, "Num columns", numColumns)
		}
		table = make([][2]int, 0, numColumns)
		costs = make([]float64, numColumns, numColumns)
		scanTemplate = "%f"

		alpha = make(map[int][]int, numRows)
		betta = make(map[int][]int, numColumns)
		var tmpFloat float64
		var numStrs int
		var tmpInt int

		for i := 0; i < numColumns; i++ {
			_, problem = scanSingle(file, &tmpFloat)
			if problem != nil {
				fmt.Println("Problem with scan float")
				return
			} else {
				costs[i] = tmpFloat
				_, problem = scanSingle(file, &numStrs)

				betta[i+1] = make([]int, numStrs)
				// add rows to column
				for j := 0; j < numStrs; j++ {
					_, problem = scanSingle(file, &tmpInt)
					table = append(table, [2]int{tmpInt - 1, i + 1})
					betta[i+1][j] = tmpInt - 1
					alpha[tmpInt-1] = append(alpha[tmpInt-1], i+1)
				}

			}
		}
	}
	return
}

func ParseAirline(filename string) (table [][2]int, costs []float64,
	alpha, betta map[int][]int, problem error) {
	file, err := os.Open(filename)
	if err == nil {
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				problem = err
				fmt.Println("Problem with closing file ", file.Name())
			}
		}(file)

		var numRows, numColumns int
		var scanTemplate string = "%d %d"
		_, err = scanFromReader(file, scanTemplate, &numColumns, &numRows)
		if err != nil {
			problem = err
			fmt.Println("Problem with reading num of rows/columns!")
			return
		} else {
			fmt.Println("Num rows", numRows, "Num columns", numColumns)
		}
		table = make([][2]int, 0, numColumns)
		costs = make([]float64, numColumns, numColumns)
		var tmpFloat float64

		//read costs
		scanTemplate = "%f"
		for i := 0; i < numColumns; i++ {
			_, err = scanSingle(file, &tmpFloat)
			if err != nil {
				problem = err
				return
			} else {
				costs[i] = tmpFloat
			}
		}

		alpha = make(map[int][]int, numRows)
		betta = make(map[int][]int, numColumns)
		var numStrs int
		var tmpInt int

		for i := 0; i < numColumns; i++ {
			_, problem = scanSingle(file, &numStrs)

			betta[i+1] = make([]int, numStrs)
			// add rows to column
			for j := 0; j < numStrs; j++ {
				_, problem = scanSingle(file, &tmpInt)
				table = append(table, [2]int{tmpInt - 1, i + 1})
				betta[i+1][j] = tmpInt - 1
				alpha[tmpInt-1] = append(alpha[tmpInt-1], i+1)
			}

		}
	}
	return
}
