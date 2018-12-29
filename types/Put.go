package types

import (
	"encoding/base64"
	"fmt"
	"os"
)

type PutDetails struct {
	Type    string
	Path    string
	Content string
}

func (p PutDetails) Execute() *Result {
	result := new(Result)

	fileHandle, err := os.Create(p.Path)
	defer fileHandle.Close()
	if err != nil {
		result.Result = err.Error()
	}

	fileData, err := base64.StdEncoding.DecodeString(p.Content)
	if err != nil {
		result.Result = err.Error()
	}

	bytesWritten, err := fileHandle.Write(fileData)
	if err != nil {
		result.Result = fmt.Sprintf("Partial file write (%d bytes): %s", bytesWritten, err.Error())
	}
	result.Result = "Wrote file without error"
	return result
}

func (pt *PutDetails) String() string {
	template :=
		`Type: %s
Path: %s
Content: %s
`
	return fmt.Sprintf(template, pt.Type, pt.Path, pt.Content)
}
