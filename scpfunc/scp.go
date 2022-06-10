package scpfunc

import (
	"fmt"
	"math"
)

type SolutionRepairer struct {
	alpha, betta map[int][]int
	costs        []float64
	S            []int
	W            map[int]int
	U            []int
}

func (repairer *SolutionRepairer) GetCosts() []float64 {
	return repairer.costs
}

func NewSolutionRepairer(alpha, betta map[int][]int, costs []float64) *SolutionRepairer {
	S := make([]int, len(costs))
	W := make(map[int]int)
	U := make([]int, 0, len(alpha))
	return &SolutionRepairer{alpha, betta, costs, S, W, U}
}

func (repairer *SolutionRepairer) CheckSolution(solution []float64) (ok bool, response string) {
	S := make([]int, 0, cap(solution))
	for ind, elem := range solution {
		if elem > 0.0 {
			S = append(S, ind+1)
		}
	}

	checked := make(map[int]bool)

	for _, elem := range S {
		for _, row := range repairer.betta[elem] {
			checked[row] = true
		}
	}

	if len(repairer.alpha) == len(checked) {
		ok = true
	} else {
		ok = false
	}
	response = fmt.Sprintf("Allright: %v Count rows: %v, was covered: %v", ok, len(repairer.alpha), len(checked))
	return
}

func (repairer *SolutionRepairer) RepairSolution(solution []float64) {
	S := make(map[int]bool)
	for ind, elem := range solution {
		if elem > 0.0 {
			S[ind+1] = true
		}
	}

	repairer.W = make(map[int]int)
	for key := range repairer.alpha {
		repairer.W[key] = MapIntSliceIntersection(S, repairer.alpha[key])
	}

	//repairer.U = repairer.U[:0]
	for key, value := range repairer.W {
		if value == 0 {
			repairer.U = append(repairer.U, key)
		}
	}

	// repair uncovered
	for {
		if len(repairer.U) == 0 {
			break
		}

		var currMinCol int
		var currMin float64 = math.MaxFloat64
		for _, col := range repairer.alpha[repairer.U[0]] {
			tmp := repairer.costs[col-1] / float64(IntIntersection(repairer.U, repairer.betta[col]))
			if tmp < currMin {
				currMin = tmp
				currMinCol = col
			}
		}
		//repairer.S = append(repairer.S, currMinCol)
		solution[currMinCol-1] = 1.0
		for _, elem := range repairer.betta[currMinCol] {
			repairer.W[elem] += 1
			ind := findElementInt(repairer.U, elem)
			if ind != -1 {
				repairer.U = append(repairer.U[:ind], repairer.U[ind+1:]...)
			}
			if len(repairer.U) == 0 {
				break
			}
		}

	}

	//delete redundancy
	for i := len(solution) - 1; i >= 0; i-- {
		if solution[i] < 1.0 {
			continue
		}
		flag := true
		for _, elem := range repairer.betta[i+1] {
			if repairer.W[elem] < 2 {
				flag = false
				break
			}
		}

		if flag {
			solution[i] = 0.0

			for _, elem := range repairer.betta[i+1] {
				repairer.W[elem] -= 1
			}
		}
	}
}

//RemoveRedudancy method
// for remove redudancy
func (repairer *SolutionRepairer) RemoveRedudancy(solution []float64) {
	S := make(map[int]bool)
	for ind, elem := range solution {
		if elem > 0.0 {
			S[ind+1] = true
		}
	}

	repairer.W = make(map[int]int)
	for key := range repairer.alpha {
		repairer.W[key] = MapIntSliceIntersection(S, repairer.alpha[key])
	}

	//delete redundancy
	for i := len(solution) - 1; i >= 0; i-- {
		if solution[i] < 1.0 {
			continue
		}
		flag := true
		for _, elem := range repairer.betta[i+1] {
			if repairer.W[elem] < 2 {
				flag = false
				break
			}
		}

		if flag {
			solution[i] = 0.0

			for _, elem := range repairer.betta[i+1] {
				repairer.W[elem] -= 1
			}
		}
	}

}

func IntIntersection(a []int, b []int) (count int) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				count++
			}
		}
	}
	return
}

func MapIntSliceIntersection(a map[int]bool, b []int) (c int) {
	for _, item := range b {
		if _, ok := a[item]; ok {
			c++
		}
	}
	return
}

func CalcFitness(xVec []float64, costs []float64) (res float64, ok bool) {
	if len(xVec) != len(costs) {
		ok = false
		return
	}
	ok = true

	for i := range xVec {
		res += xVec[i] * costs[i]
	}
	return
}

func findElementInt(x []int, elem int) int {
	for i := range x {
		if x[i] == elem {
			return i
		}
	}
	return -1
}

func GetLimitCost(costs []float64, fit float64) int {
	index := len(costs) - 1
	for i := range costs {
		if costs[i] > fit {
			return i
		}
	}
	return index
}
