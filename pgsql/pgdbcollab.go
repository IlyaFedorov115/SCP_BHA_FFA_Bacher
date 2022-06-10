package pgsql

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"scpmod/parser_scp"

	_ "github.com/lib/pq"
)

type matrixElement struct {
	matrix_id int
	row       int
	column    int
	value     float64
}

type DBInfo struct {
	matrixTableName    string
	matrixElementTable string
}

func NewDBInfo(matrixTableName string, matrixElementTable string) DBInfo {
	return DBInfo{matrixTableName: matrixTableName, matrixElementTable: matrixElementTable}
}

type DBTool struct {
	userName string
	password string
	dbName   string
	sslMode  string
	DBInfo
}

func NewDBTool(userName string, password string, dbName string, sslMode string, info DBInfo) *DBTool {
	return &DBTool{userName: userName,
		password: password, dbName: dbName, sslMode: sslMode, DBInfo: info}
}

func (dbTool *DBTool) getConnStr() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbTool.userName, dbTool.password, dbTool.dbName, dbTool.sslMode)
}

func (dbTool *DBTool) makeQueryGetTableById() string {
	return fmt.Sprintf("select * from %s where id = $1", dbTool.matrixTableName)
}

func (dbTool *DBTool) makeQueryGetIDByName() string {
	return fmt.Sprintf("select id from %s where matrix_name = $1", dbTool.matrixTableName)
}

func (dbTool *DBTool) makeQueryGetIDNameByTempl() string {
	return fmt.Sprintf("select id, matrix_name from %s where matrix_name LIKE $1 ||'%%'", dbTool.matrixTableName)
}

func (dbTool *DBTool) makeQueryGetElements() string {
	return fmt.Sprintf("select * from %s where matrix_id = $1", dbTool.matrixElementTable)
}

func (dbTool *DBTool) makeQueryPostTable() string {
	return fmt.Sprintf("insert into %s (num_rows, num_columns, matrix_name) values ($1, $2, $3) returning id", dbTool.matrixTableName)
}

func (dbTool *DBTool) makeQueryPostElements() string {
	return fmt.Sprintf("insert into %s (matrix_id, row, col, value) values ($1, $2, $3, $4)", dbTool.matrixElementTable)
}

func (dbTool *DBTool) GetConnectionByStr(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func (dbTool *DBTool) GetConnection() *sql.DB {
	// open connection
	connStr := dbTool.getConnStr()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error with connect table", err)
		return nil
	}
	return db
}

func (dbTool *DBTool) GetTableById(matrixId int, db *sql.DB) ([][2]int, []float64,
	map[int][]int, map[int][]int, error) {

	return dbTool.getTable(matrixId, db)
}

func (dbTool *DBTool) GetTableByName(matrixName string, db *sql.DB) ([][2]int, []float64,
	map[int][]int, map[int][]int, error) {

	// get index
	var matrixId int
	row := db.QueryRow(dbTool.makeQueryGetIDByName(), matrixName)
	err := row.Scan(&matrixId)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return dbTool.getTable(matrixId, db)
}

func (dbTool *DBTool) GetIdByTempName(templateName string, db *sql.DB) ([]int, []string) {
	matrixId := make([]int, 0)
	names := make([]string, 0)
	rows, err := db.Query(dbTool.makeQueryGetIDNameByTempl(), templateName)

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		matrixId = append(matrixId, id)
		names = append(names, name)
	}
	return matrixId, names
}

func (dbTool *DBTool) getTable(matrixId int, db *sql.DB) (table [][2]int, costs []float64,
	alpha, betta map[int][]int, problem error) {

	var numRows, numColumns int
	var matrixName string
	// get size of table
	row := db.QueryRow(dbTool.makeQueryGetTableById(), matrixId)
	err := row.Scan(&matrixId, &numRows, &numColumns, &matrixName)
	if err != nil {
		panic(err)
	}

	//initialize return values
	table = make([][2]int, 0, numColumns)
	costs = make([]float64, numColumns, numColumns)
	alpha = make(map[int][]int, numRows)
	betta = make(map[int][]int, numColumns)

	// get matrix
	rows, err := db.Query(dbTool.makeQueryGetElements(), matrixId)
	if err != nil {
		panic(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	for rows.Next() {
		currElement := matrixElement{}
		err := rows.Scan(&currElement.matrix_id, &currElement.row, &currElement.column, &currElement.value)
		if err != nil {
			fmt.Println("Problem with getting Table in line ", currElement, "err = ", err)
			panic(err)
		}
		_, ok := alpha[currElement.row]
		if !ok {
			alpha[currElement.row] = make([]int, 0)
		}
		alpha[currElement.row] = append(alpha[currElement.row], currElement.column)
		_, ok = betta[currElement.column]
		if !ok {
			betta[currElement.column] = make([]int, 0)
		}
		betta[currElement.column] = append(betta[currElement.column], currElement.row)
		table = append(table, [2]int{currElement.row, currElement.column})

		// attention, if start 0, need to change
		costs[currElement.column-1] = currElement.value /// attention
	}

	return
}

func (dbTool *DBTool) AddMatrixFromFile1(filePath string) {
	_, filename := filepath.Split(filePath)
	for i := len(filename) - 1; i > 0; i-- {
		if filename[i] == '.' {
			filename = filename[:i]
			break
		}
	}
	table, costs, alpha, betta, err := parser_scp.ParseScp(filePath)
	if err != nil {
		fmt.Println("Problem with reading file", filename, "for adding to DB!!")
		panic(err)
	}
	dbTool.postTable1(table, costs, len(alpha), len(betta), filename)
}

func (dbTool *DBTool) postTable1(table [][2]int, costs []float64, numRows, numColumns int, name string) {
	// open connection
	connStr := dbTool.getConnStr()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	//try to close connection
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	// insert new table and get id
	var newId int
	err = db.QueryRow(dbTool.makeQueryPostTable(), numRows, numColumns, name).Scan(&newId)
	if err != nil {
		fmt.Println("error with getting id from new table")
		panic(err)
	}

	// insert elements of matrix
	queryPost := dbTool.makeQueryPostElements()
	for i := range table {
		_, err := db.Exec(queryPost, newId, table[i][0], table[i][1], costs[table[i][1]-1]) /// attention if start with 0...
		if err != nil {
			fmt.Println("Problem with writing element of matrix")
			panic(err)
		}
	}

}

func (dbTool *DBTool) AddMatrixFromFile(db *sql.DB, filePath string) {
	_, filename := filepath.Split(filePath)
	for i := len(filename) - 1; i > 0; i-- {
		if filename[i] == '.' {
			filename = filename[:i]
			break
		}
	}
	table, costs, alpha, betta, err := parser_scp.ParseScp(filePath)
	if err != nil {
		fmt.Println("Problem with reading file", filename, "for adding to DB!!")
		panic(err)
	}
	dbTool.postTable(db, table, costs, len(alpha), len(betta), filename)
}

func (dbTool *DBTool) postTable(db *sql.DB, table [][2]int, costs []float64, numRows, numColumns int, name string) {

	// insert new table and get id
	var newId int
	err := db.QueryRow(dbTool.makeQueryPostTable(), numRows, numColumns, name).Scan(&newId)
	if err != nil {
		fmt.Println("error with getting id from new table")
		panic(err)
	}

	// insert elements of matrix
	queryPost := dbTool.makeQueryPostElements()
	for i := range table {
		_, err := db.Exec(queryPost, newId, table[i][0], table[i][1], costs[table[i][1]-1]) /// attention if start with 0...
		if err != nil {
			fmt.Println("Problem with writing element of matrix")
			panic(err)
		}
	}

}
