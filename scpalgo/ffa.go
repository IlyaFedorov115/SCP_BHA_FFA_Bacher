package scpalgo

import (
	"fmt"
	"math"
	"math/rand"
	"scpmod/scpfunc"
	"scpmod/supmath"
	"time"
)

type MoveType int64
type ChangeBestFf int64

const (
	StandardMove MoveType = iota
	BestFFMove
)

const (
	NoChange ChangeBestFf = iota
	ChangeRand
	ChangeMutation
	ChangeBest
	ChangeZero
)

type FFASolver struct {
	paramAlpha     []float64
	paramGamma     float64
	paramBetta     float64
	paramBettaPow  float64
	moveType       MoveType
	changeBestType ChangeBestFf
	CalcDist       CalcDistFun
	MoveOne        bool
}

func (ffa *FFASolver) SetAlpha(alpha []float64) {
	if len(alpha) == 1 {
		ffa.paramAlpha = alpha
	} else if len(alpha) == 2 {
		ffa.paramAlpha = alpha
		if ffa.paramAlpha[0] < ffa.paramAlpha[1] {
			ffa.paramAlpha[0], ffa.paramAlpha[1] = ffa.paramAlpha[1], ffa.paramAlpha[0]
		}
	} else {
		panic(fmt.Sprintf("Invalid len of param alpha, must be 1 or 2, but given: %v", len(alpha)))
	}
}

func NewFFASolver(alpha []float64, gamma, betta, bettapow float64, move MoveType, changeBest ChangeBestFf) *FFASolver {
	res := &FFASolver{
		paramAlpha:     alpha,
		paramGamma:     gamma,
		paramBetta:     betta,
		paramBettaPow:  bettapow,
		moveType:       move,
		changeBestType: changeBest,
		CalcDist:       supmath.CalcEuclidDist,
		MoveOne:        false,
	}
	res.SetAlpha(alpha)
	return res
}

func NewFFASolverStand() *FFASolver {
	return &FFASolver{
		paramAlpha:     []float64{1.0},
		paramGamma:     1.0,
		paramBetta:     1.0,
		paramBettaPow:  2,
		moveType:       StandardMove,
		changeBestType: NoChange,
	}
}

func (ffa *FFASolver) Solve(popSize int, numIter int, costs []float64,
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
				ffa.changeFireFly(solutions[i], globalBest, binarizer, repair)
				solFits[i], _ = scpfunc.CalcFitness(solutions[i], costs)
			}
		}

	}
	return globalFit, globalBest
}

func (ffa *FFASolver) calcAlpha(step int) float64 {

	if len(ffa.paramAlpha) == 1 {
		return ffa.paramAlpha[0]
	} else {
		return ffa.paramAlpha[1] + (ffa.paramAlpha[0]-ffa.paramAlpha[1])*math.Exp(float64(-step))
	}

}

func (ffa *FFASolver) calcBetta(dist float64) float64 {
	return ffa.paramBetta * math.Exp(-ffa.paramGamma*math.Pow(dist, ffa.paramBettaPow))
}

func (ffa *FFASolver) moveOperator(x1, x2 []float64, dist, alpha float64, best []float64) (ok bool) {
	if len(x1) != len(x2) {
		ok = false
		return
	}
	betta := ffa.calcBetta(dist)
	randVec := supmath.RandFloats(-1.0, 1.0, len(x1))

	if ffa.moveType == StandardMove {
		for i := range x1 {
			x1[i] = x1[i] + betta*(x2[i]-x1[i]) + alpha*randVec[i]
		}
	} else if ffa.moveType == BestFFMove {
		for i := range x1 {
			x1[i] = x1[i] + betta*(x2[i]-x1[i]) + alpha*(randVec[i])*(x1[i]-best[i])
		}
	}
	return true
}

func (ffa *FFASolver) changeFireFly(sol, global []float64,
	binarizer *supmath.Binarizer, repair *scpfunc.SolutionRepairer) {
	if ffa.changeBestType == NoChange {
		return
	} else if ffa.changeBestType == ChangeBest {
		dist, _ := ffa.CalcDist(sol, global)
		ffa.moveOperator(sol, global, dist, ffa.calcAlpha(0), global)
	} else if ffa.changeBestType == ChangeRand {
		copy(sol, supmath.RandBinLimit(len(sol), countColumns(global)))
	} else if ffa.changeBestType == ChangeZero {
		zero := make([]float64, len(sol), len(sol))
		dist, _ := ffa.CalcDist(sol, zero)
		ffa.moveOperator(sol, zero, dist, ffa.calcAlpha(0), global)
	} else {
		supmath.LimitMutation(sol, 1, len(sol))
	}

	binarizer.Binarization(sol, []interface{}{global})
	repair.RepairSolution(sol)
}
