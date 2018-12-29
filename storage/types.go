package storage

import "github.com/bblenard/C2/types"

type DBDump struct {
	Tasks   []types.Task
	Clients []types.Client
	Results []types.Result
}
