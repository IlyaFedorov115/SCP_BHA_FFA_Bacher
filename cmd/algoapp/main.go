package main

import (
	"encoding/csv"
	"fmt"
	"os"
	conf "scpmod/internal"
	"scpmod/parser_scp"
	"scpmod/scpalgo"
	"scpmod/scpexpt"
	"scpmod/supmath"
	"time"

	"path/filepath"

	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type algoParams struct {
	// ffa
	alpha    float64
	gamma    float64
	betta    float64
	move     string
	CalcDist string
	//bh
	stagn    float64
	collapse string
}

const (
	ffaDistEuclid = "euclid"
	ffaDistManhat = "manhat"

	ffaMoveStand = "stand"
	ffaMoveBest  = "best"

	bhCollapseStand = "stand"
	bhCollapseRand  = "rand"
)

const (
	ffaAlgo    = "ffa"
	bhAlgo     = "bh"
	ffaPsoAlgo = "ffapso"
	ffaRMSAlgo = "ffarms"

	typeSaveTable = "table"
	typeSaveCsv   = "csv"
)

func main() {
	var configData conf.ScpConfig

	var configName string
	var srcExpt string
	var save2Csv bool
	var saveTable bool

	var algoSettings algoParams
	var fileSave string
	var algoChoice string
	spin := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	// glag for file conf param
	pflag.StringVarP(&configName, "conf", "c", "../../configs/config.yaml", fmt.Sprintf("Name of config file for expt (without extension)"))

	// choose algorithm
	pflag.StringVarP(&algoChoice, "algo", "a", bhAlgo, fmt.Sprintf("Name of algo [%s|%s|%s|%s]", ffaAlgo, bhAlgo, ffaPsoAlgo, ffaRMSAlgo))

	pflag.StringVar(&srcExpt, "srcexpt", "", fmt.Sprintf("Instance files mask."))

	/*
		algo params
	*/

	//pflag.Float64SliceVar
	pflag.Float64Var(&algoSettings.alpha, "alpha", 0.1, "Param alpha for FFA algos.")
	pflag.Float64Var(&algoSettings.gamma, "gamma", 0.01, "Param gamma for FFA algos.")
	pflag.Float64Var(&algoSettings.betta, "betta", 1.0, "Param betta for FFA algos.")
	pflag.Float64Var(&algoSettings.stagn, "stagn", 1.0, "Param stang_persent for BHA algo.")

	pflag.StringVar(&algoSettings.collapse, "collapse", bhCollapseRand, fmt.Sprintf("Choose collapse for BHA [%s|%s]", bhCollapseRand, bhCollapseStand))
	pflag.StringVar(&algoSettings.move, "move", ffaMoveBest, fmt.Sprintf("Choose move_type for FFA [%s|%s]", ffaMoveBest, ffaMoveStand))
	pflag.StringVar(&algoSettings.CalcDist, "dist", ffaDistEuclid, fmt.Sprintf("Choose dist_type for FFA [%s|%s]", ffaDistEuclid, ffaDistManhat))

	// saving results

	pflag.StringVar(&fileSave, "res", "", fmt.Sprintf("File to save."))
	pflag.BoolVar(&saveTable, "stab", false, "Save in table style.")
	pflag.BoolVar(&save2Csv, "scsv", false, "Save in csv style.")

	pflag.Parse()

	configData = conf.ReadConfig(configName)

	if fileSave == "" {
		pflag.Usage()
		logrus.Fatal("Error! Miss filename")
		return
	}

	if !saveTable && !save2Csv {
		pflag.Usage()
		logrus.Fatal("Miss choice saving type")
		return
	}

	if srcExpt == "" {
		pflag.Usage()
		logrus.Fatal("Miss source instance filenames")
		return
	}

	// start working algo

	exptParams := scpexpt.NewExptParams(configData.PopSize, configData.MaxIter,
		configData.NumExpt, supmath.NewBinarizer(supmath.GetTransferByStr(configData.TransferFun), supmath.GetDiscreteByStr(configData.DiscreteMethod)))

	exptMaker := scpexpt.NewScpExptMaker()
	var solver scpalgo.ScpSolver
	var data [][]string
	var headers []string

	if algoChoice == ffaAlgo || algoChoice == ffaPsoAlgo || algoChoice == ffaRMSAlgo {

		var moveType scpalgo.MoveType
		if algoSettings.move == ffaMoveBest {
			moveType = scpalgo.BestFFMove
		} else if algoSettings.move == ffaMoveStand {
			moveType = scpalgo.StandardMove
		} else {
			pflag.Usage()
			logrus.Fatal("Bad value for move type")
			return
		}

		if algoChoice == ffaAlgo {
			solver = scpalgo.NewFFASolver([]float64{algoSettings.alpha}, algoSettings.gamma,
				algoSettings.betta, 2, moveType, scpalgo.NoChange)
		} else if algoChoice == ffaPsoAlgo {
			solver = scpalgo.NewPSOSolver([]float64{algoSettings.alpha}, algoSettings.gamma,
				algoSettings.betta, 2, moveType, scpalgo.NoChange)
		} else {
			solver = scpalgo.NewRSMSolver([]float64{algoSettings.alpha}, algoSettings.gamma,
				algoSettings.betta, 2, moveType, scpalgo.NoChange)
		}

	} else if algoChoice == bhAlgo {
		var col scpalgo.CollapseType
		if algoSettings.collapse == bhCollapseRand {
			col = scpalgo.RandCollapse
		} else if algoSettings.collapse == bhCollapseStand {
			col = scpalgo.StandCollapse
		} else {
			pflag.Usage()
			logrus.Fatal("Bad value for collapse type")
			return
		}

		var norm scpalgo.NormType
		norm = scpalgo.MaxNorm

		solver = scpalgo.NewBHASolver(norm, col, algoSettings.stagn)
	} else {
		pflag.Usage()
		logrus.Fatal("Bad value algo choice")
		return
	}

	// get files by mask
	filesInstance, err := filepath.Glob(srcExpt)
	if err != nil {
		panic(err)
	}

	spin.Start()
	data, headers = exptMaker.TestSetInstanceFiles(filesInstance, exptParams, solver, parser_scp.ParseScp)
	spin.Stop()
	if save2Csv {
		file, err := os.Create(fileSave + ".csv")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)

		exptMaker.Save2Csv(writer, data, headers)
	}

	if saveTable {
		file, err := os.Create(fileSave)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		exptMaker.Save2File(file, data, headers)
	}

}
