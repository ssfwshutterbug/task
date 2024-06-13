package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	_ "modernc.org/sqlite"
)

const database_name = "task.sqlite3"
const table_name = "task_info"
const column_length = 55

type Task struct {
	id         string
	c_time     string
	task       string
	result     string
	f_time     string
	status     string
	task_len   string
	result_len string
}

// specify sqlite database table columns
var header = Task{"id", "create_time", "task", "result", "finish_time", "status", "task_len", "result_len"}

func main() {

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// create task database default dir ~/.local/share/task if not exist
	database_path := filepath.Join(homedir, ".local/share/task")
	err = os.MkdirAll(database_path, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// sqlite database full path
	database := filepath.Join(database_path, database_name)

	db, err := sql.Open("sqlite", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var args = os.Args

	// if no args, print help msg and exit
	check_parameter(args, 1)

	// otherwise check args
	var operation = args[1]

	switch operation {
	case "init", "--init":
		create_table(db, table_name)
	case "list", "--list":
		query := fmt.Sprintf(`select * from %s where %s = '0'`, table_name, header.status)
		query_task(db, query, "part")
	case "list-all", "--list-all":
		query := fmt.Sprintf(`select * from %s`, table_name)
		query_task(db, query, "full")
	case "list-done", "--list-done":
		query := fmt.Sprintf(`select * from %s where %s = '1'`, table_name, header.status)
		query_task(db, query, "full")
	case "add", "--add":
		check_parameter(args, 2)
		content := args[2]
		add_task(db, content)
	case "finish", "--finish":
		check_parameter(args, 3)
		id := args[2]
		content := args[3]
		finish_task(db, id, content)
	case "delete", "--delete":
		check_parameter(args, 2)
		id := args[2]
		delete_task(db, id)
	case "help", "--help":
		help()
	default:
		help()
	}

}

// create table (init)
func create_table(db *sql.DB, table string) {
	createTableQuery := fmt.Sprintf(
		`create table if not exists %s (
        %s integer primary key autoincrement,
        %s text not null,
        %s text not null,
        %s text,
        %s text,
        %s integer not null,
        %s integer default 0,
        %s integer default 0
    );
    `, table, header.id, header.c_time, header.task, header.result, header.f_time, header.status, header.task_len, header.result_len)

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

// add task
func add_task(db *sql.DB, task string) {
	var c_time = time.Now().Format("2006-01-02 15:04:05")
	var status = "0"
	var task_len = len(task)

	insertQuery := fmt.Sprintf(
		`insert into %s (%s, %s, %s, %s)
    values ('%s', '%s', '%s', '%v');
    `, table_name, header.c_time, header.task, header.status, header.task_len, c_time, task, status, task_len)

	_, err := db.Exec(insertQuery)
	if err != nil {
		log.Fatal(err)
	}
}

// finish task
func finish_task(db *sql.DB, id string, content string) {
	var c_time = time.Now().Format("2006-01-02 15:04:05")
	var status = "1"
	var result_len = len(content)

	updateQuery := fmt.Sprintf(
		`update %s set %s = '%s', %s = '%s', %s = %s, %s = %v where %s = %s;
        `, table_name, header.result, content, header.f_time, c_time, header.status, status, header.result_len, result_len, header.id, id)

	_, err := db.Exec(updateQuery)
	if err != nil {
		log.Fatal(err)
	}
}

// delete task
func delete_task(db *sql.DB, id string) {
	deleteQuery := fmt.Sprintf(`delete from %s where %s = %s`, table_name, header.id, id)

	_, err := db.Exec(deleteQuery)
	if err != nil {
		log.Fatal(err)
	}
}

// query data
func query_task(db *sql.DB, query string, style string) {

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	// get table header
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	// store healder to a list
	var list = [8]string{}
	for i, j := range columns {
		list[i] = j
	}

	// customise color
	header_color := color.New(color.BgGreen, color.FgBlack)
	row_color := color.New(color.FgBlue)

	// according to style parameter, adjust output columns, with full will output all columns of table
	// but with part parameter, only output mainly three columns.
	if style == "full" {
		header_color.Printf("%-4s %-20s %-*s    %-*s    %-20s\n", list[0], list[1], column_length, list[2], column_length, list[3], list[4])

		// show task row
		// use flag to adjust output rows front color, every two rows will render once, good for looking
		var flag = 0
		for rows.Next() {
			var row = Task{}
			var total_lines int

			rows.Scan(&row.id, &row.c_time, &row.task, &row.result, &row.f_time, &row.status, &row.task_len, &row.result_len)

			row.task, row.result, total_lines = fill_blank(row.task, row.result)
			for i := 0; i < total_lines; i++ {
				task := row.task[column_length*i : column_length*(i+1)]
				result := row.result[column_length*i : column_length*(i+1)]
				if i == 0 {
					if flag%2 == 0 {
						fmt.Printf("%-4s %-20s %-*s    %-*s    %-20s\n", row.id, row.c_time, column_length, task, column_length, result, row.f_time)
					} else {
						row_color.Printf("%-4s %-20s %-*s    %-*s    %-20s\n", row.id, row.c_time, column_length, task, column_length, result, row.f_time)
					}
				} else {
					if flag%2 == 0 {
						fmt.Printf("  %+*s    %-*s    \n", column_length+4+20, task, column_length, result)
					} else {
						row_color.Printf("  %+*s    %-*s    \n", column_length+4+20, task, column_length, result)
					}
				}
			}
			flag = flag + 1
		}
	} else if style == "part" {
		header_color.Printf("%-4s %-20s %-*s\n", list[0], list[1], column_length, list[2])

		var flag = 0
		for rows.Next() {
			var row = Task{}
			var total_lines int

			rows.Scan(&row.id, &row.c_time, &row.task, &row.result, &row.f_time, &row.status, &row.task_len, &row.result_len)

			row.task, row.result, total_lines = fill_blank(row.task, row.result)
			for i := 0; i < total_lines; i++ {
				task := row.task[column_length*i : column_length*(i+1)]

				if i == 0 {
					if flag%2 == 0 {
						fmt.Printf("%-4s %-20s %-*s\n", row.id, row.c_time, column_length, task)
					} else {
						row_color.Printf("%-4s %-20s %-*s\n", row.id, row.c_time, column_length, task)
					}
				} else {
					if flag%2 == 0 {
						fmt.Printf("  %+*s\n", column_length+4+20, task)
					} else {
						row_color.Printf("  %+*s\n", column_length+4+20, task)
					}
				}
			}
			flag = flag + 1
		}
	}
}

// in order to show text in multiple rows, start with the second row, we need add extra
// blank space before task column.
func fill_blank(task, result string) (ntask string, nresult string, total int) {
	len_task, len_result := len(task), len(result)
	var total_lines int

	if len_task > len_result {
		total_lines = (len_task + (column_length - 1)) / column_length
	} else if len_task < len_result {
		total_lines = (len_result + (column_length - 1)) / column_length
	} else {
		total_lines = len_task % column_length
	}

	task = task + strings.Repeat(" ", total_lines*column_length-len_task)
	result = result + strings.Repeat(" ", total_lines*column_length-len_result)
	return task, result, total_lines
}

// check parameter
func check_parameter(args []string, parameter_len int) {
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
--update    update      update task result
--delete    delete      delete task
    `

	fmt.Println(msg)
}
