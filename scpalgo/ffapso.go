package scpalgo

import (
	"math/rand"
	"scpmod/scpfunc"
	"scpmod/supmath"
	"time"
)

type PSOSolver struct {
	*FFASolver
}

func NewPSOSolver(alpha []float64, _, betta, bettapow float64, move MoveType, changeBest ChangeBestFf) *PSOSolver {
	res := &PSOSolver{FFASolver: NewFFASolver(alpha, 0.0, betta, bettapow, move, changeBest)}
	res.SetAlpha(alpha)
	return res
}

func (ffa *PSOSolver) Solve(popSize int, numIter int, costs []float64,
	repair *scpfunc.SolutionRepairer, binarizer *supmath.Binarizer) (float64, []float64) {

	rand.Seed(time.Now().UnixNano())
	// generate population
	solutions := make([][]float64, popSize, popSize)
	for i := 0; i < popSize; i++ {
		solutions[i] = supmath.RandBinFloat(len(costs))
	}

	solFits := make([]float64, popSize, popSize)
	//repair solutions
	for i := range solutions {
		repair.RepairSolution(solutions[i])
	}
	// calc fitness
	for i := 0; i < popSize; i++ {
		solFits[i], _ = scpfunc.CalcFitness(solutions[i], costs)
	}

	//start global best
	globalBest := make([]float64, len(costs), len(costs))
	for i := range globalBest {
		globalBest[i] = 0.0
	}
	repair.RepairSolution(globalBest)
	globalFit, _ := scpfunc.CalcFitness(globalBest, costs)

	var alpha float64

	for step := 0; step < numIter; step++ {

		alpha = ffa.calcAlpha(step)
		//Update Global best
		minInd, minFit := supmath.MinFloatSlice(solFits)
		if minFit < globalFit {
			globalFit = minFit
			copy(globalBest, solutions[minInd])
		}

		for i := 0; i < popSize; i++ {
			if solFits[i] > globalFit {
				ffa.moveOperator(solutions[i], globalBest, 0, alpha, globalBest)
				binarizer.Binarization(solutions[i], globalBest)
				repair.RepairSolution(solutions[i])
				solFits[i], _ = scpfunc.CalcFitness(solutions[i], costs)
			}

		}

	}
	return globalFit, globalBest
}
