package storage

import (
	"database/sql"
	"fmt"

	"github.com/bblenard/C2/types"
	_ "github.com/mattn/go-sqlite3"
)

const (
	CLIENT_QUERY      = "SELECT ID, Hostname FROM clients"
	TASKS_QUERY       = "SELECT ID, Target, Type, Details FROM tasks WHERE Target=? AND Completed=?"
	TASK_REPORT_QUERY = `
	SELECT DISTINCT * FROM
	(SELECT t.ID, t.Target, t.Type, t.Title, t.Details, r.Result
	FROM tasks t
	LEFT JOIN task_results r ON t.ID=r.ID AND t.Target=r.Target
	UNION ALL
	SELECT t.ID, t.Target, t.Type, t.Title, t.Details, r.Result
	FROM task_results r
	LEFT JOIN tasks t ON t.ID=r.ID AND t.Target=r.Target
	WHERE t.ID IS NOT NULL OR r.ID IS NOT NULL) WHERE Target=?`
	DELETE_TASK_QUERY = "DELETE FROM tasks WHERE ID=?"
)

type MySQLDB struct {
	db *sql.DB
}

func (m *MySQLDB) String() string {
	return ""
}

func (m *MySQLDB) Init() error {
	db, err := sql.Open("sqlite3", "./C2Database.sqlite")
	if err != nil {
		return err
	}
	m.db = db
	createTables := []string{
		"CREATE TABLE IF NOT EXISTS clients (ID TEXT PRIMARY KEY ON CONFLICT ABORT, Hostname TEXT);",
		"CREATE TABLE IF NOT EXISTS tasks (ID TEXT PRIMARY KEY ON CONFLICT ABORT, Target TEXT, Type INTEGER, Completed INTEGER, Title TEXT, IssuedTime timestamp, Details BLOB);",
		"CREATE TABLE task_results (ID TEXT, Target TEXT, Result TEXT, FOREIGN KEY(ID) REFERENCES tasks(ID), PRIMARY KEY (ID, Target));",
	}

	for _, table := range createTables {
		_, err = db.Exec(table)
		if err != nil {
			return err
		}
	}

	return err
}

func (m *MySQLDB) SaveTasks(tasks ...types.Task) error {
	stmt, err := m.db.Prepare("INSERT INTO tasks(ID, Target, Type, Completed, Title, IssuedTime, Details) values(?,?,?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("failed to create prepared statement: %s", err.Error())
	}
	for _, task := range tasks {
		_, err := stmt.Exec(task.ID, task.Target, task.Type, task.Completed, task.Title, task.IssuedTime, task.Details)
		if err != nil {
			return fmt.Errorf("failed to insert task into db: %s", err.Error())
		}
	}
	return nil
}

func (m *MySQLDB) SaveTaskResults(trs ...types.Result) error {
	insertStmt, err := m.db.Prepare("INSERT INTO task_results(ID, Target, Result) values(?,?,?)")
	if err != nil {
		return fmt.Errorf("failed to create prepared statement: %s", err.Error())
	}
	updateStmt, err := m.db.Prepare("UPDATE tasks SET Completed=1 WHERE ID=? AND Target=?")
	if err != nil {
		return fmt.Errorf("failed to create prepared statement: %s", err.Error())
	}
	for _, result := range trs {
		_, err := insertStmt.Exec(result.ID, result.Target, result.Result)
		if err != nil {
			return fmt.Errorf("failed to insert task result into db: %s", err.Error())
		}
		_, err = updateStmt.Exec(result.ID, result.Target)
		if err != nil {
			return fmt.Errorf("failed to update task to completed: %s", err.Error())
		}
	}
	return nil
}

func (m *MySQLDB) SaveClients(clients ...types.Client) error {
	stmt, err := m.db.Prepare("INSERT INTO clients(ID, Hostname) values(?,?)")
	if err != nil {
		return fmt.Errorf("failed to create prepared statement: %s", err.Error())
	}
	for _, client := range clients {
		_, err := stmt.Exec(client.ID, client.Details.Hostname)
		if err != nil {
			return fmt.Errorf("failed to insert client into db: %s", err.Error())
		}
	}
	return nil
}

func (m *MySQLDB) GetClients() ([]types.Client, error) {
	rows, err := m.db.Query("SELECT * FROM clients")
	if err != nil {
		return []types.Client{}, err
	}
	defer rows.Close()
	clients := new([]types.Client)
	for rows.Next() {
		client := types.Client{}
		err := rows.Scan(&client.ID, &client.Details.Hostname)
		if err != nil {
			return *clients, fmt.Errorf("failed to scan into client struct: %s", err)
		}
		*clients = append(*clients, client)
	}
	return *clients, nil
}

func (m *MySQLDB) GetClientByID(id string) (types.Client, error) {
	stmt, err := m.db.Prepare("SELECT * FROM clients WHERE ID=?")
	if err != nil {
		return types.Client{}, fmt.Errorf("failed to prepare statement: %s", err)
	}
	row := stmt.QueryRow(id)
	client := types.Client{}
	err = row.Scan(&client.ID, &client.Details.Hostname)
	if err != nil {
		return types.Client{}, fmt.Errorf("failed to scan row into client: %s", err)
	}
	return client, nil
}

func (m *MySQLDB) GetPendingTasksByClient(id string) ([]types.Task, error) {
	stmt, err := m.db.Prepare(TASKS_QUERY)
	if err != nil {
		return []types.Task{}, fmt.Errorf("failed to prepare tasks query: %s", err)
	}
	rows, err := stmt.Query(id, false)
	if err != nil {
		return []types.Task{}, fmt.Errorf("failed to query tasks for client: %s", err)
	}
	defer rows.Close()
	tasks := new([]types.Task)
	for rows.Next() {
		t := types.Task{}
		err := rows.Scan(&t.ID, &t.Target, &t.Type, &t.Details)
		if err != nil {
			return *tasks, fmt.Errorf("failed to scan into pending task: %s", err)
		}
		*tasks = append(*tasks, t)
	}
	return *tasks, nil
}

func (m *MySQLDB) GetCompletedTasksByClient(id string) ([]types.TaskReport, error) {
	stmt, err := m.db.Prepare(TASK_REPORT_QUERY)
	if err != nil {
		return []types.TaskReport{}, fmt.Errorf("failed to prepare TASK_REPORT_QUERY statement: %s", err)
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return []types.TaskReport{}, fmt.Errorf("failed to execute TASK_REPORT_QUERY statement: %s", err)
	}
	defer rows.Close()
	reports := new([]types.TaskReport)
	for rows.Next() {
		t := types.TaskReport{}
		err := rows.Scan(&t.ID, &t.Target, &t.Type, &t.Title, &t.Details, &t.Result)
		if err != nil {
			return *reports, fmt.Errorf("failed to scan into task report: %s", err)
		}
		*reports = append(*reports, t)
	}
	return *reports, nil
}

func (m *MySQLDB) DeleteTasks(taskIDs ...string) error {
	stmt, err := m.db.Prepare(DELETE_TASK_QUERY)
	if err != nil {
		return fmt.Errorf("failed to prepare DELETE_TASK_QUERY statement: %s", err)
	}
	var numberOfDeleteRows int64
	for _, id := range taskIDs {
		deleteResult, err := stmt.Exec(id)
		num, err := deleteResult.RowsAffected()
		if err != nil {
			return err
		}
		numberOfDeleteRows = numberOfDeleteRows + num
	}
	if numberOfDeleteRows != int64(len(taskIDs)) {
		return fmt.Errorf("failed to delete all tasks specified in request")
	}
	return nil
}
