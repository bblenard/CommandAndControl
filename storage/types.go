package storage

import "github.com/bblenard/CommandAndControl/types"

type DBDump struct {
	Tasks   []types.Task
	Clients []types.Client
	Results []types.Result
}
