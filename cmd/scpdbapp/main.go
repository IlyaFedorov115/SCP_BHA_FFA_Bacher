package main

import (
	"database/sql"
	"fmt"
	"scpmod/pgsql"

	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

//"user=dev_test1 password=fanfanfan dbname=test1 sslmode=disable"
func main() {
	var sourceStr = "" // option on config
	pflag.StringVarP(&sourceStr, "source", "s", "", fmt.Sprintf("Source path file/dir"))

	var sourceFile = true
	var sourceDir = false
	var sourceTemp = false

	var connStr = ""
	var matrix_name = "matrix"
	var matrix_elem_name = "matrix_element"

	pflag.StringVarP(&connStr, "conct", "c", "", fmt.Sprintf("Get connection string ['user= password= dbname= sslmode=']"))
	pflag.BoolVar(&sourceDir, "dir", false, "Source is directory.")
	pflag.BoolVar(&sourceTemp, "temp", false, "Source by template.")

	pflag.StringVarP(&matrix_name, "mtrx", "m", "matrix", fmt.Sprintf("Name of table for matrix"))
	pflag.StringVarP(&matrix_elem_name, "mtrx_elem", "e", "matrix_element", fmt.Sprintf("Name of table for matrix elements"))

	pflag.Parse()

	if connStr == "" {
		pflag.Usage()
		logrus.Fatal("Error: need to get connection string!")
		return
	}

	if sourceStr == "" {
		pflag.Usage()
		logrus.Fatal("Error: need to get source string!")
		return
	}

	dbInfo := pgsql.NewDBInfo(matrix_name, matrix_elem_name)
	dbTool := pgsql.NewDBTool("", "", "", "", dbInfo)

	dbVar := dbTool.GetConnectionByStr(connStr)
	if dbVar == nil {
		logrus.Fatal("Can`t connect to db.")
		return
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logrus.Fatal(err)
			return
		}
	}(dbVar)

	count := 0
	if sourceTemp {

		m, err := filepath.Glob(sourceStr)
		if err != nil {
			logrus.Fatal(err)
			return
		}

		for _, val := range m {
			count++
			dbTool.AddMatrixFromFile(dbVar, val)
		}
		logrus.Info(fmt.Sprintf("Was added %v files.", count))

		//dbTool.AddMatrixFromFile())

	} else if sourceDir {

		files, err := filepath.Glob(sourceStr + "/*")
		if err != nil {
			logrus.Fatal(err)
			return
		}

		for _, file := range files {
			count++
			dbTool.AddMatrixFromFile(dbVar, file)
		}
		logrus.Info(fmt.Sprintf("Was added %v files.", count))

	} else if sourceFile {
		dbTool.AddMatrixFromFile(dbVar, sourceStr)
		logrus.Info(fmt.Sprintf("Was added %v files.", 1))
	} else {
		return
	}

}
