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

var engine = db.Engine{
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

	database := engine.ConnectDB()
	operation(database)
	database.Close()

}

func operation(database *sql.DB) {

	var args = os.Args

	// if no args, print help msg and exit
	checkParameter(args, 1)

	// otherwise check args
	var operator = args[1]

	switch operator {
	case "init", "--init":
		engine.CreateTable(database)
	case "list", "--list":
		query := fmt.Sprintf(`select * from %s where %s = '0'`, engine.Tablename, db.Header.Status)
		engine.QueryTask(database, query, &color, COLUMNLENGTH)
	case "list-all", "--list-all":
		query := fmt.Sprintf(`select * from %s`, engine.Tablename)
		engine.QueryTask(database, query, &color, COLUMNLENGTH)
	case "list-done", "--list-done":
		query := fmt.Sprintf(`select * from %s where %s = '1'`, engine.Tablename, db.Header.Status)
		engine.QueryTask(database, query, &color, COLUMNLENGTH)
	case "add", "--add":
		checkParameter(args, 2)
		content := args[2]
		engine.AddTask(database, content)
	case "done", "--done":
		checkParameter(args, 2)
		id := args[2]
		engine.FinishTask(database, id)
	case "delete", "--delete":
		checkParameter(args, 2)
		id := args[2]
		engine.DeleteTask(database, id)
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
