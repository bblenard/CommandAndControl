package types

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

type GetDetails struct {
	Type string
	Path string
}

func (gt *GetDetails) Execute() *Result {
	result := new(Result)
	data, err := ioutil.ReadFile(gt.Path)
	if err != nil {
		fmt.Println("failed to read file: ", err)
		result.Result = base64.StdEncoding.EncodeToString([]byte(err.Error()))
	} else {
		fmt.Printf("Read %s\n", data)
		result.Result = base64.StdEncoding.EncodeToString(data)
	}
	return result
}

func (gt *GetDetails) String() string {
	template :=
		`Type: %s
Path: %s
`
	return fmt.Sprintf(template, gt.Type, gt.Path)
}
