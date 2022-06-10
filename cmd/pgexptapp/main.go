package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"scpmod/pgsql"
	"scpmod/scpalgo"
	"scpmod/scpexpt"
	"scpmod/scpfunc"
	"scpmod/supmath"
	"strconv"

	"os"

	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

type MenuExptParams struct {
	PopSize         int
	NumIter         int
	NumExpt         int
	TransType       string
	DiscreteType    string
	InstanceName    string
	FileToSave      string
	SaveType        string
	CurrAlgo        string
	repairSlice     []*scpfunc.SolutionRepairer
	findOneInstance bool
	instancesSlice  []string
}

func (menu *MenuExptParams) Clear() {
	menu.PopSize = 20
	menu.NumExpt = 10
	menu.NumIter = 200
	menu.DiscreteType = "standard"
	menu.TransType = "s1"
	menu.FileToSave = ""
	menu.InstanceName = ""
}

type MenuFFAParams struct {
	Alpha    float64
	Betta    float64
	Gamma    float64
	MoveType scpalgo.MoveType
}

type MenuBHAParams struct {
	Normalization scpalgo.NormType
	CollapseType  scpalgo.CollapseType
}

//params
var transferChoices = []string{"s1", "s2", "s3", "s12", "s13", "v1", "v2", "v3", "v4"}
var discreteChoices = []string{"standard", "elist"}
var algoChoices = []string{"BHA", "FFA"}
var styleSaveChoices = []string{"table", "csv"}
var moveChoices = []string{"Stand move", "Move best"}
var normChoices = []string{"Max norm", "Mean norm", "None norm"}
var solver scpalgo.ScpSolver
var ffaParams = MenuFFAParams{0.0, 1.0, 0.001, scpalgo.StandardMove}
var bhaParams = MenuBHAParams{scpalgo.NoneNorm, scpalgo.RandCollapse}
var exptParams = MenuExptParams{20, 500, 10, "s3", "standard", "", "", "csv", "bha", nil, true, nil}
var dbTool *pgsql.DBTool
var dbVar *sql.DB

//tui
var app = tview.NewApplication()
var formMenuInstance = tview.NewForm()   //Menu instance
var formMenuExptParams = tview.NewForm() //Menu expt
var formFFAParams = tview.NewForm()      //Menu ffa
var formBHAParams = tview.NewForm()      //Menu bha
var pages = tview.NewPages()

var Quit = func() {
	app.Stop()
}

//wait
var modalWait = tview.NewModal().SetText(fmt.Sprintf("Wait until algo don`t end working")).AddButtons([]string{"Start"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
	if buttonLabel == "Start" {
		runResult(solver, exptParams)
		pages.SwitchToPage("continue")
	}
})

func startModalWait() *tview.Modal {
	modalWait.SetText(fmt.Sprintf("Wait until [%s] don`t end working", exptParams.CurrAlgo))
	return modalWait
}

//continue
var modalContinue = tview.NewModal().
	SetText(fmt.Sprintf("Want to continue?")).
	AddButtons([]string{"Quit", "This instance", "New instance"}).
	SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Quit" {
			Quit()
		}
		if buttonLabel == "New instance" {
			startMenuInstance()
			pages.SwitchToPage("Menu instance")
		}
		if buttonLabel == "This instance" {
			startMenuExpt()
			pages.SwitchToPage("Menu expt")
		}
	})

/*
 модуль ошибки поиска экземпляра
*/
var modalErrFindInstance = tview.NewModal().
	SetText(fmt.Sprintf("Can`t find instance: [%s]", exptParams.InstanceName)).
	AddButtons([]string{"Quit", "Try again"}).
	SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Quit" {
			Quit()
		}
		if buttonLabel == "Try again" {
			startMenuInstance()
			pages.SwitchToPage("Menu instance")
		}
	})

/*
 окно для выбора экземпляров для экспериментов
*/

func startMenuInstance() *tview.Form {
	modalWait.SetText(fmt.Sprintf("Wait until  don`t end working"))
	formMenuInstance.Clear(true)

	formMenuInstance.AddInputField("Name Instance", "", 30, nil, func(text string) {
		exptParams.InstanceName = text
	})

	formMenuInstance.AddCheckbox("Find start with?", false, func(checked bool) {
		exptParams.findOneInstance = !checked
	})

	formMenuInstance.AddButton("Next", func() {
		if exptParams.findOneInstance {
			_, costs, alpha, betta, err := dbTool.GetTableByName(exptParams.InstanceName, dbVar)
			if err != nil {
				pages.SwitchToPage("error")
				return
			}
			repair := scpfunc.NewSolutionRepairer(alpha, betta, costs)
			exptParams.repairSlice = make([]*scpfunc.SolutionRepairer, 1)
			exptParams.repairSlice[0] = repair
			exptParams.instancesSlice = make([]string, 1)
			exptParams.instancesSlice[0] = exptParams.InstanceName
			//exptParams.instancesSlice = append(exptParams.instancesSlice, exptParams.InstanceName)
		} else {
			idSlice, namesSlice := dbTool.GetIdByTempName(exptParams.InstanceName, dbVar)
			if idSlice == nil || len(idSlice) == 0 {
				pages.SwitchToPage("error")
				return
			}
			exptParams.instancesSlice = namesSlice
			exptParams.repairSlice = make([]*scpfunc.SolutionRepairer, len(idSlice))
			for i, id := range idSlice {
				_, costs, alpha, betta, err := dbTool.GetTableById(id, dbVar)
				if err != nil {
					logrus.Error("Problem with id", id)
					continue
				}
				repair := scpfunc.NewSolutionRepairer(alpha, betta, costs)
				//exptParams.repairSlice = append(exptParams.repairSlice, repair)
				exptParams.repairSlice[i] = repair
			}
		}
		startMenuExpt()
		pages.SwitchToPage("Menu expt")
	})

	formMenuInstance.AddButton("Quit", func() {
		Quit()
	})

	formMenuInstance.SetBorder(true).SetTitle("Enter instance name").SetTitleAlign(tview.AlignLeft)
	return formMenuInstance
}

/*
окно для выбора параметров экспериментов
*/
func startMenuExpt() *tview.Form {
	exptParams.Clear()
	formMenuExptParams.Clear(true)

	formMenuExptParams.AddInputField("File results", "", 80, nil, func(text string) {
		exptParams.FileToSave = text
	})

	formMenuExptParams.AddDropDown("Style result", styleSaveChoices, 0, func(option string, optionIndex int) {
		exptParams.SaveType = option
	})

	formMenuExptParams.AddInputField("Population size", "20", 10, nil, func(text string) {
		popsize, err := strconv.Atoi(text)
		if err != nil {
			return
		}
		exptParams.PopSize = popsize
	})

	formMenuExptParams.AddInputField("Max iters", "200", 10, nil, func(text string) {
		num, err := strconv.Atoi(text)
		if err != nil {
			return
		}
		exptParams.NumIter = num
	})

	formMenuExptParams.AddInputField("Num experiments", "10", 10, nil, func(text string) {
		count, err := strconv.Atoi(text)
		if err != nil {
			return
		}
		exptParams.NumExpt = count
	})

	formMenuExptParams.AddDropDown("Select transfer func", transferChoices, 0, func(option string, optionIndex int) {
		exptParams.TransType = option
	})

	formMenuExptParams.AddDropDown("Select discrete func", discreteChoices, 0, func(option string, optionIndex int) {
		exptParams.DiscreteType = option
	})

	formMenuExptParams.AddDropDown("Select algo", algoChoices, 0, func(option string, optionIndex int) {
		exptParams.CurrAlgo = option
	})

	formMenuExptParams.AddButton("Next", func() {
		if exptParams.CurrAlgo == "FFA" {
			startFFAMenu()
			pages.SwitchToPage("Menu ffa")
		} else {
			startBHAMenu()
			pages.SwitchToPage("Menu bha")
		}
	})

	formMenuExptParams.AddButton("Back", func() {
		startMenuInstance()
		pages.SwitchToPage("Menu instance")
	})

	formMenuExptParams.AddButton("Quit", func() {
		Quit()
	})

	formMenuExptParams.SetBorder(true).SetTitle("Enter experiment params").SetTitleAlign(tview.AlignLeft)
	return formMenuExptParams
}

/*
окно для параметров ffa
*/
func startFFAMenu() *tview.Form {

	formFFAParams.Clear(true)
	formFFAParams.AddInputField("Alpha", "0.1", 10, nil, func(text string) {
		alpha, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return
		}
		ffaParams.Alpha = alpha
	})

	formFFAParams.AddInputField("Betta", "1.0", 10, nil, func(text string) {
		betta, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return
		}
		ffaParams.Betta = betta
	})

	formFFAParams.AddInputField("Gamma", "0.001", 10, nil, func(text string) {
		gamma, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return
		}
		ffaParams.Gamma = gamma
	})

	formFFAParams.AddDropDown("Select move ff", moveChoices, 0, func(option string, optionIndex int) {
		if option == "Stand move" {
			ffaParams.MoveType = scpalgo.StandardMove
		} else {
			ffaParams.MoveType = scpalgo.BestFFMove
		}
	})

	formFFAParams.AddButton("Next", AfterAlgoFunc)
	formFFAParams.AddButton("Back", func() {
		startMenuExpt()
		pages.SwitchToPage("Menu expt")
	})
	formFFAParams.AddButton("Quit", func() {
		Quit()
	})

	formFFAParams.SetBorder(true).SetTitle("Enter FFA params").SetTitleAlign(tview.AlignLeft)
	return formFFAParams
}

func startBHAMenu() *tview.Form {
	formBHAParams.Clear(true)
	formBHAParams.AddDropDown("Select norm type", normChoices, 0, func(option string, optionIndex int) {
		if option == "Max norm" {
			bhaParams.Normalization = scpalgo.MaxNorm
		} else if option == "Mean norm" {
			bhaParams.Normalization = scpalgo.MeanNorm
		} else {
			bhaParams.Normalization = scpalgo.NoneNorm
		}
	})

	formBHAParams.AddButton("Next", AfterAlgoFunc)

	formBHAParams.AddButton("Back", func() {
		startMenuExpt()
		pages.SwitchToPage("Menu expt")
	})
	formBHAParams.AddButton("Quit", func() {
		Quit()
	})

	formBHAParams.SetBorder(true).SetTitle("Enter bha params").SetTitleAlign(tview.AlignLeft)
	return formBHAParams
}

func AfterAlgoFunc() {
	if exptParams.CurrAlgo == "FFA" {
		solver = scpalgo.NewFFASolver([]float64{ffaParams.Alpha}, ffaParams.Gamma, ffaParams.Betta, 2,
			ffaParams.MoveType, scpalgo.NoChange)
	} else {
		solver = scpalgo.NewBHASolver(bhaParams.Normalization, bhaParams.CollapseType, 1.0)
	}
	startModalWait()
	pages.SwitchToPage("wait")

}

func runResult(solver scpalgo.ScpSolver, params MenuExptParams) {

	filesave := exptParams.FileToSave
	var discr func(xVec []float64, items interface{})
	if params.DiscreteType == "elist" {
		discr = supmath.ElitistDiscrete
	} else {
		discr = supmath.StandardDiscrete
	}

	exptParams_ := scpexpt.NewExptParams(exptParams.PopSize, exptParams.NumIter, exptParams.NumExpt,
		supmath.NewBinarizer(supmath.GetTransferByStr(exptParams.TransType), discr))
	expt := scpexpt.NewScpExptMaker()

	data, headers := expt.TestSetInstance(exptParams.instancesSlice, exptParams.repairSlice, exptParams_, solver)
	if exptParams.SaveType == "table" {
		file, _ := os.Create(filesave)
		expt.Save2File(file, data, headers)
	} else {
		file, _ := os.Create(filesave)
		w := csv.NewWriter(file)
		expt.Save2Csv(w, data, headers)
	}
}

//constr := "user=dev_test1 password=fanfanfan dbname=test1 sslmode=disable"
func main() {
	// get connection to postgres database
	if len(os.Args) < 2 {
		logrus.Fatal("Need provide a PostgreSQL connect string.")
		return
	}

	dbInfo := pgsql.NewDBInfo("matrix", "matrix_element")
	dbTool = pgsql.NewDBTool("", "", "", "", dbInfo)

	dbVar = dbTool.GetConnectionByStr(os.Args[1])

	if dbVar == nil {
		logrus.Fatal("Can`t connect to db.")
		return
	}

	pingErr := dbVar.Ping()
	if pingErr != nil {
		logrus.Fatal("Bad connections string: ", pingErr)
		return
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(dbVar)

	//Menu instance
	//Menu expt
	//Menu ffa
	//Menu bha

	startMenuInstance()
	pages.AddPage("Menu instance", formMenuInstance, true, true)
	pages.AddPage("Menu expt", formMenuExptParams, true, false)
	pages.AddPage("Menu ffa", formFFAParams, true, false)
	pages.AddPage("Menu bha", formBHAParams, true, false)
	pages.AddPage("error", modalErrFindInstance, true, false)
	pages.AddPage("wait", modalWait, true, false)
	pages.AddPage("continue", modalContinue, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
