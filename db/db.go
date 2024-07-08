package db

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"database/sql"

	"task/colorstr"

	_ "modernc.org/sqlite"
)

type Task struct {
	Id       string
	C_time   string
	Task     string
	F_time   string
	Status   string
	Task_len string
}

type Engine struct {
	Enginename   string
	Databasepath string
	Databasename string
	Tablename    string
}

// specify sqlite database table columns
var Header = Task{
	Id:       "id",
	C_time:   "createTime",
	Task:     "task",
	F_time:   "finishTime",
	Status:   "status",
	Task_len: "taskLength",
}

type DB interface {
	ConnectDB() *sql.DB
	CreateTable()
	AddTask()
	DeleteTask()
	FinishTask()
	QueryTask()
}

type Color struct {
	HeadColorBg string
	HeadColorFg string
	Color1Bg    string
	Color1Fg    string
	Color2Bg    string
	Color2Fg    string
}

func (engine *Engine) ConnectDB() *sql.DB {

	// create db file if not exist
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	database_path := filepath.Join(homedir, strings.TrimLeft(engine.Databasepath, "~"))

	if err := os.MkdirAll(database_path, 0755); err != nil && !os.IsExist(err) {
		fmt.Println(err)
	}

	database := filepath.Join(database_path, engine.Databasename)

	// connect db
	db, err := sql.Open(engine.Enginename, database)
	if err != nil {
		fmt.Println(err)
	}
	//defer db.Close()

	return db
}

func (engine *Engine) CreateTable(db *sql.DB) {
	createTableQuery := fmt.Sprintf(
		`create table if not exists %s (
        %s integer primary key autoincrement,
        %s text not null,
        %s text not null,
        %s text,
        %s integer default 0,
        %s integer not null
    );
    `, engine.Tablename, Header.Id, Header.C_time, Header.Task, Header.F_time, Header.Status, Header.Task_len)

	_, err := db.Exec(createTableQuery)
	if err != nil {
		fmt.Println("error:", err)
	}
}

// add task
func (engine *Engine) AddTask(db *sql.DB, task string) {
	var c_time = time.Now().Format("2006-01-02")
	var status = "0"
	var task_len = len(task)

	insertQuery := fmt.Sprintf(
		`insert into %s (%s, %s, %s, %s)
    values ('%s', '%s', '%s', '%v');
    `, engine.Tablename, Header.C_time, Header.Task, Header.Status, Header.Task_len, c_time, task, status, task_len)

	_, err := db.Exec(insertQuery)
	if err != nil {
		fmt.Println(err)
	}
}

// finish task
func (engine *Engine) FinishTask(db *sql.DB, id string) {
	var c_time = time.Now().Format("2006-01-02")
	var status = "1"

	updateQuery := fmt.Sprintf(
		`update %s set %s = '%s', %s = %s where %s = %s;
        `, engine.Tablename, Header.F_time, c_time, Header.Status, status, Header.Id, id)

	_, err := db.Exec(updateQuery)
	if err != nil {
		fmt.Println(err)
	}
}

// delete task
func (engine *Engine) DeleteTask(db *sql.DB, id string) {
	deleteQuery := fmt.Sprintf(`delete from %s where %s = %s`, engine.Tablename, Header.Id, id)

	_, err := db.Exec(deleteQuery)
	if err != nil {
		fmt.Println(err)
	}
}

// query data
func (engine *Engine) QueryTask(db *sql.DB, query string, color *Color, columnlen int) {

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	// get table header
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
	}

	// store healder to a list
	var list = [6]string{}
	for i, j := range columns {
		list[i] = j
	}

	colorHeader := fmt.Sprintf("%-4s %-20s %-*s\n", list[0], list[1], columnlen, list[2])
	colorHeader = colorstr.ColorizeRgb(color.HeadColorFg, color.HeadColorBg, colorHeader)
	io.WriteString(os.Stdout, colorHeader)

	var flag = 0
	for rows.Next() {
		var row = Task{}

		rows.Scan(&row.Id, &row.C_time, &row.Task, &row.F_time, &row.Status, &row.Task_len)

		task, total_lines := fill_blank(row.Task, columnlen)
		for i := 0; i < total_lines; i++ {
			task := task[columnlen*i : columnlen*(i+1)]

			// first line, need to display id, createtime, task
			if i == 0 && flag%2 == 0 {
				color1row := fmt.Sprintf("%-4s %-20s %-*s\n", row.Id, row.C_time, columnlen, task)
				color1row = colorstr.ColorizeRgb(color.Color1Fg, color.Color1Bg, color1row)
				io.WriteString(os.Stdout, color1row)
			}
			if i == 0 && flag%2 != 0 {
				color2row := fmt.Sprintf("%-4s %-20s %-*s\n", row.Id, row.C_time, columnlen, task)
				color2row = colorstr.ColorizeRgb(color.Color2Fg, color.Color2Bg, color2row)
				io.WriteString(os.Stdout, color2row)
			}

			// more than one line, display task
			if i != 0 && flag%2 == 0 {
				color1row := fmt.Sprintf("  %+*s\n", columnlen+4+20, task)
				color1row = colorstr.ColorizeRgb(color.Color1Fg, color.Color1Bg, color1row)
				io.WriteString(os.Stdout, color1row)
			}
			if i != 0 && flag%2 != 0 {
				color2row := fmt.Sprintf("  %+*s\n", columnlen+4+20, task)
				color2row = colorstr.ColorizeRgb(color.Color2Fg, color.Color2Bg, color2row)
				io.WriteString(os.Stdout, color2row)
			}
		}
		flag = flag + 1
	}
}

// in order to show text in multiple rows, start with the second row, we need add extra
// blank space before task column.
func fill_blank(task string, columnlen int) (ntask string, total int) {
	tasklen := len(task)
	desirecolumn := tasklen/columnlen + 2

	taskchar := task + strings.Repeat(" ", desirecolumn*columnlen-tasklen)

	return taskchar, desirecolumn
}
