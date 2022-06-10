package scpexpt

import (
	"encoding/csv"
	"fmt"
	"io"
	"scpmod/scpalgo"
	"scpmod/scpfunc"
	"scpmod/supmath"
	"time"

	"github.com/olekukonko/tablewriter"
)

type Colors int64
type OutputStyle int64

const (
	TableStyle OutputStyle = iota
	CsvStyle
)

const (
	Default Colors = iota
	Green
	Red
)

type Parser func(filename string) (table [][2]int, costs []float64,
	alpha, betta map[int][]int, problem error)

type ExptParams struct {
	PopSize    int
	NumIter    int
	CountExpts int
	Binarizer  *supmath.Binarizer
}

func NewExptParamsBr(popSize, numIter, countExpts int, transfer, discrete string) *ExptParams {
	return &ExptParams{PopSize: popSize,
		NumIter:    numIter,
		CountExpts: countExpts,
		Binarizer:  supmath.NewBinarizer(supmath.GetTransferByStr(transfer), supmath.GetDiscreteByStr(discrete))}
}

func NewExptParams(popSize, numIter, countExpts int, bin *supmath.Binarizer) *ExptParams {
	return &ExptParams{PopSize: popSize, NumIter: numIter, CountExpts: countExpts, Binarizer: bin}
}

type ScpExptMaker struct {
	resultsHeader []string
	colors        map[Colors]string
}

func NewScpExptMaker() *ScpExptMaker {
	return &ScpExptMaker{resultsHeader: []string{"File", "Min cost", "Mean cost", "Max cost", "Mean size", "Mean time", "Success"},
		colors: map[Colors]string{Default: "\033[0m", Green: "\033[32m", Red: "\033[31m"}}
}

func (expt *ScpExptMaker) Save2File(writer io.Writer, data [][]string, headers []string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(headers)
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func (expt *ScpExptMaker) Save2Csv(writer *csv.Writer, data [][]string, headers []string) {
	err := writer.Write(headers)
	if err != nil {
		fmt.Println("Problem with saving in csv file, ", err.Error())
		return
	}
	writer.WriteAll(data)
}

func (expt *ScpExptMaker) TestSetInstanceFiles(filenames []string, params *ExptParams,
	solver scpalgo.ScpSolver, parser Parser) ([][]string, []string) {

	instances := make([]*scpfunc.SolutionRepairer, len(filenames))
	for i := range filenames {
		_, costs, alpha, betta, err := parser(filenames[i])

		if err != nil {
			fmt.Println("Problem in instance: ", filenames[i], err.Error())
			panic(err)
		}
		instances[i] = scpfunc.NewSolutionRepairer(alpha, betta, costs)
	}

	return expt.TestSetInstance(filenames, instances, params, solver)
}

func (expt *ScpExptMaker) TestSetInstance(filenames []string, instances []*scpfunc.SolutionRepairer, params *ExptParams,
	solver scpalgo.ScpSolver) ([][]string, []string) {

	data := make([][]string, len(instances))
	for i := range data {
		data[i] = make([]string, len(expt.resultsHeader))
	}

	for index := range instances {
		costsSlice, sizesSlice, timesSlice, ok := TestOneInstance(filenames[index], instances[index], params, solver)

		data[index][0] = string(filenames[index])
		data[index][1] = fmt.Sprint(supmath.MinFloat64(costsSlice))
		data[index][2] = fmt.Sprint(supmath.MeanFloat64(costsSlice))
		data[index][3] = fmt.Sprint(supmath.MaxFloat64(costsSlice))
		data[index][4] = fmt.Sprint(supmath.MeanFloat64(sizesSlice))
		data[index][5] = fmt.Sprintf("%v", time.Duration(supmath.MeanInt64(timesSlice)))
		if ok == true {
			data[index][6] = "OK"
		} else {
			data[index][6] = "FAIl"
		}
	}
	return data, expt.resultsHeader
}

func TestOneInstance(instance string, repairer *scpfunc.SolutionRepairer, params *ExptParams,
	solver scpalgo.ScpSolver) ([]float64, []float64, []int64, bool) {
	// expt info
	costsSlice := make([]float64, params.CountExpts, params.CountExpts)
	sizesSlice := make([]float64, params.CountExpts, params.CountExpts)
	timeSlice := make([]int64, params.CountExpts, params.CountExpts)
	okNormal := true

	//fmt.Println(params.CountExpts)
	for i := 0; i < params.CountExpts; i++ {
		startTime := time.Now()

		// call function
		value, optimum := solver.Solve(params.PopSize, params.NumIter, repairer.GetCosts(), repairer, params.Binarizer)
		elapsedTime := time.Since(startTime)
		costsSlice[i] = value
		timeSlice[i] = int64(elapsedTime)

		// calc num of columns
		sizesSlice[i] = 0
		for _, e := range optimum {
			if e > 0 {
				sizesSlice[i] += 1.0
			}
		}

		//check solution
		ok, okWhat := repairer.CheckSolution(optimum)
		if !ok {
			fmt.Println("-----Problem with file", instance, okWhat)
			okNormal = false
		}
	}

	return costsSlice, sizesSlice, timeSlice, okNormal

}
