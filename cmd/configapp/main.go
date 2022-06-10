package main

import (
	"fmt"
	conf "scpmod/internal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const (
	optUpdate = "update"
	optCreate = "create"
)

func main() {
	var configData conf.ScpConfig

	var option string // option on config
	var configName string
	// flag for option parameter
	pflag.StringVarP(&option, "opt", "o", "", fmt.Sprintf("Choose option [%s|%s]", optUpdate, optCreate))

	// glag for file conf param
	pflag.StringVarP(&configName, "file", "f", "", fmt.Sprintf("Give name of config file (without extension)"))

	//	flag population
	pflag.IntVarP(&configData.PopSize, "popsize", "p", 0, "Population size.")

	// flag max iter
	pflag.IntVarP(&configData.MaxIter, "iters", "i", 0, "Num of iterations for algo.")

	// flag num expt
	pflag.IntVarP(&configData.NumExpt, "expt", "e", 0, "Num of experiences iterations.")

	// flag for transfer function
	pflag.StringVarP(&configData.TransferFun, "transfer", "t", "",
		"Transfer function.")
	//pflag.Lookup("transfer").NoOptDefVal = "s1"

	// flag for discretization method
	pflag.StringVarP(&configData.DiscreteMethod, "disrete", "d", "",
		"Discretization method.")
	//pflag.Lookup("disrete").NoOptDefVal = "standard"

	pflag.Parse()

	if option == "" || (option != optUpdate && option != optCreate) {
		pflag.Usage()
		logrus.Fatal("Error option type")
		return
	}

	if configName == "" {
		pflag.Usage()
		logrus.Fatal("Don`t get filename")
		return
	}

	if option == optCreate && configData.HaveNil() {
		pflag.Usage()
		logrus.Fatal("Using [create] requires all parameters")
		return
	}

	if option == optCreate {
		conf.WriteConfig(configData, configName)
	} else if option == optUpdate {
		memConfig := conf.ReadConfig(configName)
		if configData.DiscreteMethod != "" {
			memConfig.DiscreteMethod = configData.DiscreteMethod
		}
		if configData.TransferFun != "" {
			memConfig.TransferFun = configData.TransferFun
		}
		if configData.MaxIter > 0 {
			memConfig.MaxIter = configData.MaxIter
		}
		if configData.PopSize > 0 {
			memConfig.PopSize = configData.PopSize
		}
		if configData.NumExpt > 0 {
			memConfig.NumExpt = configData.NumExpt
		}
		conf.WriteConfig(memConfig, configName)
	}
}
