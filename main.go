package main

import (
	"encoding/csv"
	_ "net/http/pprof"
	"os"
	"scpmod/parser_scp"
	"scpmod/scpalgo"
	"scpmod/scpexpt"
	"scpmod/supmath"
)

var FileNames4 = []string{"./OR/scp41.txt", "./OR/scp42.txt", "./OR/scp43.txt", "./OR/scp44.txt", "./OR/scp45.txt",
	"./OR/scp46.txt", "./OR/scp47.txt", "./OR/scp48.txt", "./OR/scp49.txt", "./OR/scp410.txt"}

func main() {
	// simple example to use Expt Package
	// FileNames4 - example, use your path
	paramsExpt := scpexpt.NewExptParams(20, 1000, 10,
		supmath.NewBinarizer(supmath.GetTransferByStr("s3"), supmath.ElitistDiscrete))
	//solver := scpalgo.NewFFASolver([]float64{0.1}, 0.0002, 1.0, 2, scpalgo.StandardMove, scpalgo.NoChange)
	solver := scpalgo.NewBHASolver(scpalgo.MaxNorm, scpalgo.StandCollapse, 1.0)
	expt := scpexpt.NewScpExptMaker()
	data1, headers1 := expt.TestSetInstanceFiles(FileNames4, paramsExpt, solver, parser_scp.ParseScp)
	expt.Save2File(os.Stdout, data1, headers1)

	file, _ := os.Create("bhaV1Stand.csv")
	defer file.Close()
	writer := csv.NewWriter(file)
	expt.Save2Csv(writer, data1, headers1)
}
