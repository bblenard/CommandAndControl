package storage

import (
	"fmt"

	"github.com/bblenard/C2/types"
)

type Type int

const (
	Memory Type = iota
	Mongo
	MySQL
)

type Storage interface {
	String() string
	Init() error
	SaveTasks(...types.Task) error
	SaveTaskResults(...types.Result) error
	SaveClients(...types.Client) error
	GetClients() ([]types.Client, error)
	GetClientByID(string) (types.Client, error)
	GetPendingTasksByClient(id string) ([]types.Task, error)
	GetCompletedTasksByClient(id string) ([]types.TaskReport, error)
}

var DB Storage

func NewStorage(t Type) error {

	switch t {
	case MySQL:
		fmt.Println("MySQL!!!")
		DB = new(MySQLDB)
		return DB.Init()
	}

	return nil
}
