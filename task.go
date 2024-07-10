package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "modernc.org/sqlite"
	"task/db"
)

const COLUMNLENGTH = 55

var sqldatabase = db.DataBase{
	Enginename:   "sqlite",
	Databasepath: "~/.local/share/task",
	Databasename: "task.sqlite3",
	Tablename:    "taskInfo",
}

var color = db.Color{
	HeadColorBg: "#32283A",
	HeadColorFg: "#F38BA8",
	Color1Bg:    "#1F1F1F",
	Color1Fg:    "#CACCD4",
	Color2Bg:    "#161B22",
	Color2Fg:    "#CDD6F4",
}

func main() {
	var connect db.Connecter
	connect = &sqldatabase
	connection := connect.ConnectDB()

	var task db.Tasker
	task = &sqldatabase

	operation(connection, task)
	connection.Close()

}

func operation(connection *sql.DB, task db.Tasker) {

	flag.BoolFunc("init", "initial database", func(s string) error {
		task.CreateTable(connection)
		return nil
	})

	flag.BoolFunc("list", "list unfinished task", func(s string) error {
		query := fmt.Sprintf(`select * from %s where %s = '0'`, sqldatabase.Tablename, db.Header.Status)
		task.QueryTask(connection, query, &color, COLUMNLENGTH)
		return nil
	})

	flag.BoolFunc("list-all", "list all tasks", func(s string) error {
		query := fmt.Sprintf(`select * from %s`, sqldatabase.Tablename)
		task.QueryTask(connection, query, &color, COLUMNLENGTH)
		return nil
	})

	flag.BoolFunc("list-done", "list finished tasks", func(s string) error {
		query := fmt.Sprintf(`select * from %s where %s = '1'`, sqldatabase.Tablename, db.Header.Status)
		task.QueryTask(connection, query, &color, COLUMNLENGTH)
		return nil
	})

	add := flag.String("add", "nil", "add task")
	done := flag.String("done", "nil", "mark task has been done")
	deleteitem := flag.String("delete", "nil", "delete task")

	flag.Parse()

	if *add != "nil" {
		task.AddTask(connection, *add)
	}
	if *done != "nil" {
		task.FinishTask(connection, *done)
	}
	if *deleteitem != "nil" {
		task.DeleteTask(connection, *deleteitem)
	}
}
