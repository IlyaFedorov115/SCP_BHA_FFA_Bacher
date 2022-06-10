package scpalgo

import (
	"scpmod/scpfunc"
	"scpmod/supmath"
)

type ScpSolver interface {
	Solve(popSize int, numIter int, costs []float64,
		repair *scpfunc.SolutionRepairer, binarizer *supmath.Binarizer) (float64, []float64)
}
