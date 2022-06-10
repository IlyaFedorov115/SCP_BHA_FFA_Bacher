package scpalgo

import (
	"math/rand"
	"scpmod/scpfunc"
	"scpmod/supmath"
	"time"
)

type RSMSolver struct {
	*FFASolver
}

func NewRSMSolver(alpha []float64, _, betta, bettapow float64, move MoveType, changeBest ChangeBestFf) *RSMSolver {
	res := &RSMSolver{FFASolver: NewFFASolver(alpha, 0.0, betta, bettapow, move, changeBest)}
	res.SetAlpha(alpha)
	return res
}

func (ffa *RSMSolver) calcBetta(dist float64) float64 {
	return ffa.paramBetta
}

func (ffa *RSMSolver) moveOperator(x1, x2 []float64, dist, alpha float64, best []float64) (ok bool) {
	if len(x1) != len(x2) {
		ok = false
		return
	}
	//betta := ffa.calcBetta(dist)
	randVec := supmath.RandFloats(-1.0, 1.0, len(x1))

	if ffa.moveType == StandardMove {
		for i := range x1 {
			x1[i] = x1[i] + alpha*randVec[i]
		}
	} else if ffa.moveType == BestFFMove {
		for i := range x1 {
			x1[i] = x1[i] + alpha*(randVec[i])*(x1[i]-best[i])
		}
	}
	return true
}

func (ffa *RSMSolver) Solve(popSize int, numIter int, costs []float64,
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
			flagNoChange := true
			for j := 0; j < i; j++ {
				if solFits[i] > solFits[j] {
					dist, ok := ffa.CalcDist(solutions[i], solutions[j])
					if !ok {
						return 0, nil
					}
					ffa.moveOperator(solutions[i], solutions[j], dist, alpha, globalBest)
					binarizer.Binarization(solutions[i], globalBest)
					repair.RepairSolution(solutions[i])
					solFits[i], _ = scpfunc.CalcFitness(solutions[i], costs)
					flagNoChange = false
					if ffa.MoveOne {
						break
					}
				}
			}
			if flagNoChange {
				//ffa.changeFireFly(solutions[i], globalBest, binarizer, repair)
				//solFits[i], _ = scpfunc.CalcFitness(solutions[i], costs)
			}
		}

	}
	return globalFit, globalBest
}
