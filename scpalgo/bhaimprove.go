package scpalgo

import (
	"math/rand"
	"scpmod/scpfunc"
	"scpmod/supmath"
	"time"
)

type ImproveBHASolver struct {
	*BHASolver
	LocalSeacrhMax int
}

func NewImproveBHASolver(normalization NormType, collapse CollapseType, stagnMax float64) *ImproveBHASolver {
	return &ImproveBHASolver{&BHASolver{normalization: normalization, collapseType: collapse, stagnMax: stagnMax}, 30}
}

func (bha *ImproveBHASolver) Solve(popSize int, numIter int, costs []float64,
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

	// new part. Memory best solution
	globalMem := make([]float64, len(costs), len(costs))
	var globalFitMem float64 = 0.0
	globalFitMem += 0.0
	flagLocalSearch := false
	localCount := 0

	for step := 0; step < numIter; step++ {

		if stagnCount > stagnMax {
			flagLocalSearch = true
			stagnCount = 0
			localCount = 0
			copy(globalMem, globalBest)
			globalFitMem = globalFit
		}

		if flagLocalSearch {
			for i := range solutions {
				repair.RemoveRedudancy(solutions[i])
				localCount += 1
			}
		} else {
			for i := range solutions {
				repair.RepairSolution(solutions[i])
			}
		}

		// end local search
		if localCount > bha.LocalSeacrhMax {
			stagnCount = 0
			flagLocalSearch = false
			for i := range solutions {
				repair.RepairSolution(solutions[i])
			}
			repair.RepairSolution(globalBest)
			globalFit, _ = scpfunc.CalcFitness(globalBest, costs_)
			if globalFit > globalFitMem {
				copy(globalBest, globalMem)
				globalFit = globalFitMem
			}
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
		}

		// Collapse in Black Hole
		calcCollapse(solutions, globalBest, solFits, globalFit)

		//move toward Black hole
		for i := range solutions {
			calcRotation(solutions[i], globalBest)
		}

		// binarization
		for i := range solutions {
			binarizer.Binarization(solutions[i], []interface{}{globalBest})
		}

	}

	// new part
	for i := range solutions {
		repair.RepairSolution(solutions[i])
	}
	repair.RepairSolution(globalBest)
	globalFit, _ = scpfunc.CalcFitness(globalBest, costs_)
	if globalFit > globalFitMem {
		copy(globalBest, globalMem)
		globalFit = globalFitMem
	}

	globalFit, _ = scpfunc.CalcFitness(globalBest, costs)
	return globalFit, globalBest
}
