package preprocess

import "fmt"

type ReduceTable struct {
	table        [][2]int
	costs        []float64
	alpha        map[int][]int
	betta        map[int][]int
	RemoveUnique bool
}

func (reducer *ReduceTable) GetAlpha() map[int][]int {
	return reducer.alpha
}

func (reducer *ReduceTable) GetBetta() map[int][]int {
	return reducer.betta
}

func NewReduceTable(tab [][2]int, cost []float64, alp map[int][]int, bet map[int][]int) *ReduceTable {
	return &ReduceTable{
		table:        tab,
		costs:        cost,
		alpha:        alp,
		betta:        bet,
		RemoveUnique: true,
	}
}

func removeElementByIndex(slice []int, index int) []int {
	if index == -1 {
		fmt.Println("Index", index, "Slice", slice)
		panic("error removeElementByIndex")
	}
	sliceLen := len(slice)
	sliceLastIndex := sliceLen - 1
	if index != sliceLastIndex {
		slice[index] = slice[sliceLastIndex]
	}

	return slice[:sliceLastIndex]
}

func findIndex(slice []int, element int) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == element {
			return i
		}
	}
	return -1
}

func checkSubset(a []int, b []int) bool {
	if len(a) > len(b) {
		return false
	}
	for i := range a {
		flag := false
		for j := range b {
			if a[i] == b[j] {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

func (reducer *ReduceTable) Reduce() []float64 {
	solution := make([]float64, len(reducer.costs), len(reducer.costs))
	columnsKeys := make([]int, len(reducer.betta))
	for i := range columnsKeys {
		columnsKeys[i] = i + 1
	}

	prevLenTable := len(reducer.betta)
	// main loop (while change table)
	for {

		for keyI := range reducer.betta {
			for keyJ := range reducer.betta {
				if keyI == keyJ {
					continue
				}

				if checkSubset(reducer.betta[keyI], reducer.betta[keyJ]) &&
					//reducer.costs[keyI-1]/float64(len(reducer.betta[keyI])) >= reducer.costs[keyJ-1]/float64(len(reducer.betta[keyJ])) {
					reducer.costs[keyI-1] >= reducer.costs[keyJ-1] {
					//delete from alpha
					for _, row := range reducer.betta[keyI] {
						reducer.alpha[row] = removeElementByIndex(reducer.alpha[row], findIndex(reducer.alpha[row], keyI))
					}
					//delete from betta
					delete(reducer.betta, keyI)
					break
				}

			}
		}

		//delete unique
		if reducer.RemoveUnique {
			for keyAlph := range reducer.alpha {
				if len(reducer.alpha[keyAlph]) > 1 {
					continue
				}
				uniqCol := reducer.alpha[keyAlph][0]
				solution[uniqCol-1] = 1.0
				// delet for all col their uniq row
				for _, row := range reducer.betta[uniqCol] {
					for _, col := range reducer.alpha[row] {
						if col == uniqCol {
							continue
						}
						reducer.betta[col] = removeElementByIndex(reducer.betta[col], findIndex(reducer.betta[col], row))
					}
					delete(reducer.alpha, row)
				}
				delete(reducer.betta, uniqCol)

			}
		}

		// end condition
		if prevLenTable == len(reducer.betta) {
			break
		} else {
			prevLenTable = len(reducer.betta)
		}
	}
	return solution
}

func (reducer *ReduceTable) Reduce1() []float64 {
	solution := make([]float64, len(reducer.costs), len(reducer.costs))
	columnsKeys := make([]int, len(reducer.betta))
	for i := range columnsKeys {
		columnsKeys[i] = i + 1
	}

	prevLenTable := len(reducer.betta)
	// main loop (while change table)
	for {

		// Column Domination Part
		for i := 0; i < len(columnsKeys); i++ {

			if _, ok := reducer.betta[columnsKeys[i]]; !ok {
				continue
			}

			for j := 0; j < i; j++ {

				if _, ok := reducer.betta[columnsKeys[j]]; !ok {
					continue
				}

				if len(reducer.betta[columnsKeys[i]]) <= len(reducer.betta[columnsKeys[j]]) {
					if checkSubset(reducer.betta[columnsKeys[i]], reducer.betta[columnsKeys[j]]) &&
						reducer.costs[columnsKeys[i]-1]/float64(len(reducer.betta[columnsKeys[i]])) >= reducer.costs[columnsKeys[j]-1]/float64(len(reducer.betta[columnsKeys[j]])) {

						//delete from alpha
						for _, row := range reducer.betta[columnsKeys[i]] {
							if len(reducer.alpha[row]) == 0 {
								fmt.Println("Alpha, row: ", row, "Log: ", columnsKeys[i], "Betta i: ", reducer.betta[columnsKeys[i]], "Betta j:", reducer.betta[columnsKeys[j]])
								fmt.Println("Cost i = ", columnsKeys[i], reducer.costs[columnsKeys[i]-1])
								fmt.Println("Cost j = ", columnsKeys[j], reducer.costs[columnsKeys[j]-1])
							}
							reducer.alpha[row] = removeElementByIndex(reducer.alpha[row], findIndex(reducer.alpha[row], columnsKeys[i]))
						}
						//delete from betta
						delete(reducer.betta, columnsKeys[i])
						columnsKeys = removeElementByIndex(columnsKeys, i)
						i -= 1
						break

					}

				} else {
					if checkSubset(reducer.betta[columnsKeys[j]], reducer.betta[columnsKeys[i]]) &&
						reducer.costs[columnsKeys[j]-1]/float64(len(reducer.betta[columnsKeys[j]])) >= reducer.costs[columnsKeys[i]-1]/float64(len(reducer.betta[columnsKeys[i]])) {

						//delete from alpha
						for _, row := range reducer.betta[columnsKeys[j]] {
							reducer.alpha[row] = removeElementByIndex(reducer.alpha[row], findIndex(reducer.alpha[row], columnsKeys[j]))
						}
						//delete from betta
						delete(reducer.betta, columnsKeys[j])
						columnsKeys = removeElementByIndex(columnsKeys, j)
						j -= 1
					}
				}

			}

		}

		if reducer.RemoveUnique {
			//delete unique
			for _, val := range reducer.alpha {
				if len(val) == 1 {
					solution[val[0]-1] = 1.0
					for _, e := range reducer.betta[val[0]] {
						for _, r := range reducer.alpha[e] {
							if r == val[0] {
								continue
							}
							reducer.betta[r] = removeElementByIndex(reducer.betta[r], findIndex(reducer.betta[r], e))
						}
						delete(reducer.alpha, e)
					}
					delete(reducer.betta, val[0])
				}
			}
		}

		// end condition
		if prevLenTable == len(reducer.betta) {
			break
		} else {
			prevLenTable = len(reducer.betta)
		}

	}

	return solution

}
