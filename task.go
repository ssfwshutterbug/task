package main

import (
	"database/sql"
	"fmt"
	"io"
	_ "modernc.org/sqlite"
	"os"
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

	var args = os.Args

	// if no args, print help msg and exit
	checkParameter(args, 1)

	// otherwise check args
	var operator = args[1]

	switch operator {
	case "init", "--init":
		task.CreateTable(connection)
	case "list", "--list":
		query := fmt.Sprintf(`select * from %s where %s = '0'`, sqldatabase.Tablename, db.Header.Status)
		task.QueryTask(connection, query, &color, COLUMNLENGTH)
	case "list-all", "--list-all":
		query := fmt.Sprintf(`select * from %s`, sqldatabase.Tablename)
		task.QueryTask(connection, query, &color, COLUMNLENGTH)
	case "list-done", "--list-done":
		query := fmt.Sprintf(`select * from %s where %s = '1'`, sqldatabase.Tablename, db.Header.Status)
		task.QueryTask(connection, query, &color, COLUMNLENGTH)
	case "add", "--add":
		checkParameter(args, 2)
		content := args[2]
		task.AddTask(connection, content)
	case "done", "--done":
		checkParameter(args, 2)
		id := args[2]
		task.FinishTask(connection, id)
	case "delete", "--delete":
		checkParameter(args, 2)
		id := args[2]
		task.DeleteTask(connection, id)
	case "help", "--help":
		help()
	default:
		help()
	}

}

// check parameter
func checkParameter(args []string, parameter_len int) {
	if len(args) == parameter_len {
		help()
		os.Exit(1)
	}
}

// help
func help() {
	msg := `
usage: task operation [parameter]

available options:
--help      help        show help message
--init      init        create sqlite table
--list      list        list unfinished tasks
--list-all  list-all    list all tasks
--list-done list-done   list accomplished tasks
--add       add         add a new task
--done      done        update task result
--delete    delete      delete task
    `

	io.WriteString(os.Stdout, msg+"\n")
}
