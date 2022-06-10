package scpalgo

import (
	"math"
	"math/rand"
	"scpmod/scpfunc"
	"scpmod/supmath"
	"time"
)

type NormType int64
type CollapseType int64

type CalcDistFun func([]float64, []float64) (float64, bool)
type CollapseFun func(solutions [][]float64, globalBest []float64,
	solFits []float64, globalFit float64)

const (
	NoneNorm NormType = iota
	MeanNorm
	MaxNorm
)

const (
	RandCollapse CollapseType = iota
	StandCollapse
)

type BHASolver struct {
	normalization NormType
	collapseType  CollapseType
	stagnMax      float64
}

func NewBHASolver(normalization NormType, collapse CollapseType, stagnMax float64) *BHASolver {
	return &BHASolver{normalization: normalization, collapseType: collapse, stagnMax: stagnMax}
}

func (bha *BHASolver) getCost(costs []float64) []float64 {
	if bha.normalization == NoneNorm {
		return costs
	}
	res := make([]float64, len(costs))
	copy(res, costs)
	if bha.normalization == MeanNorm {
		supmath.NormilizeFloatSliceMean(res)
	} else {
		supmath.NormilizeFloatSliceMax(res)
	}
	return res
}

func (bha *BHASolver) getCollapseFun() CollapseFun {
	if bha.collapseType == RandCollapse {
		return calcCollapseRand
	} else {
		return calcCollapseStand
	}
}

func (bha *BHASolver) Solve(popSize int, numIter int, costs []float64,
	repair *scpfunc.SolutionRepairer, binarizer *supmath.Binarizer) (float64, []float64) {

	rand.Seed(time.Now().UnixNano())
	costs_ := bha.getCost(costs)
	calcCollapse := bha.getCollapseFun()
	stagnMax := int(bha.stagnMax * float64(numIter))
	stagnCount := 0

	// generate population
	solutions := make([][]float64, popSize, popSize)
	for i := 0; i < popSize; i++ {
		solutions[i] = supmath.RandBinFloat(len(costs_))
	}
	solFits := make([]float64, popSize, popSize)

	//start global best
	globalBest := make([]float64, len(costs_), len(costs_))
	for i := range globalBest {
		globalBest[i] = 0.0
	}
	repair.RepairSolution(globalBest)
	globalFit, _ := scpfunc.CalcFitness(globalBest, costs_)

	for step := 0; step < numIter; step++ {

		//repair solutions
		for i := range solutions {
			repair.RepairSolution(solutions[i])
		}

		// calc fitness
		for i := 0; i < popSize; i++ {
			solFits[i], _ = scpfunc.CalcFitness(solutions[i], costs_)
		}

		//Update Black hole with the best star
		minInd, minFit := supmath.MinFloatSlice(solFits)
		if minFit < globalFit {
			globalFit, minFit = minFit, globalFit
			globalBest, solutions[minInd] = solutions[minInd], globalBest
			stagnCount = 0
		} else {
			stagnCount += 1
			if stagnCount > stagnMax {
				break
			}
		}

		// Collapse in Black Hole
		calcCollapse(solutions, globalBest, solFits, globalFit)

		//move toward Black hole
		for i := range solutions {
			calcRotation(solutions[i], globalBest)
		}

		// binarization
		for i := range solutions {
			binarizer.Binarization(solutions[i], globalBest)
		}

	}

	globalFit, _ = scpfunc.CalcFitness(globalBest, costs)
	return globalFit, globalBest
}

func calcRotation(solution []float64, globalBest []float64) {
	randVec := supmath.RandFloats(0.0, 1.0, len(solution))
	for i := range solution {
		solution[i] = solution[i] + randVec[i]*(globalBest[i]-solution[i])
	}
}

// calcCollapseRand - версия, где R > r - поглощение, r = [0,1]
func calcCollapseRand(solutions [][]float64, globalBest []float64,
	solFits []float64, globalFit float64) {

	count := countColumns(globalBest)
	R := globalFit / supmath.SumFloatSlice(solFits)
	randVec := supmath.RandFloats(0.0, 1.0, len(solutions))
	for i := range solutions {
		if R > randVec[i] {
			solutions[i] = supmath.RandBinLimit(len(solutions[0]), count) //getRand(len(solutions[i]), count) //supmath.RandBinFloat(len(solutions[i]))
		}
	}
}

func calcCollapseStand(solutions [][]float64, globalBest []float64,
	solFits []float64, globalFit float64) {
	count := countColumns(globalBest)
	R := globalFit / supmath.SumFloatSlice(solFits)

	for i := range solutions {
		if R > math.Abs(solFits[i]-globalFit) {
			solutions[i] = supmath.RandBinLimit(len(solutions[0]), count)
		}
	}
}

func countColumns(sol []float64) int {
	res := 0
	for i := range sol {
		if sol[i] > 0 {
			res += 1
		}
	}
	return res
}
