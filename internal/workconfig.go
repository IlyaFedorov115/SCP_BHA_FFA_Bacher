package internal

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ScpConfig struct {
	PopSize        int
	MaxIter        int
	NumExpt        int
	TransferFun    string
	DiscreteMethod string
}

func (config ScpConfig) HaveNil() bool {
	if config.DiscreteMethod == "" || config.TransferFun == "" ||
		config.MaxIter == 0 || config.NumExpt == 0 || config.PopSize == 0 {
		return true
	}
	return false
}

func ReadConfig(path string) ScpConfig {

	viper.SetDefault("pop_size", 20)
	viper.SetDefault("num_expt", 10)
	viper.SetDefault("max_iter", 400)
	viper.SetDefault("discrete", "standard")
	viper.SetDefault("transfer", "s1")

	viper.SetConfigName(path)          // name of config file (without extension)
	viper.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("../configs/") // path to look for the config file in
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal("Errir config file ", path)
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	var configData ScpConfig
	configData.PopSize = viper.GetInt("pop_size")
	configData.NumExpt = viper.GetInt("num_expt")
	configData.MaxIter = viper.GetInt("max_iter")
	configData.DiscreteMethod = viper.GetString("discrete")
	configData.TransferFun = viper.GetString("transfer")
	log.Infof("Config file: %v", path)

	return configData
}

func WriteConfig(configData ScpConfig, filename string) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.Set("pop_size", configData.PopSize)
	v.Set("num_expt", configData.NumExpt)
	v.Set("max_iter", configData.MaxIter)
	v.Set("discrete", configData.DiscreteMethod)
	v.Set("transfer", configData.TransferFun)
	err := v.WriteConfigAs(filename + ".yaml")
	if err != nil {
		log.Fatal("Error to write config file ", filename, " ", err)
	}
}
