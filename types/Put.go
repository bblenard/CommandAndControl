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
		result.Result = base64.StdEncoding.EncodeToString([]byte(err.Error()))
	}

	fileData, err := base64.StdEncoding.DecodeString(p.Content)
	if err != nil {
		result.Result = base64.StdEncoding.EncodeToString([]byte(err.Error()))
	}

	bytesWritten, err := fileHandle.Write(fileData)
	if err != nil {
		resultString := fmt.Sprintf("Partial file write (%d bytes): %s", bytesWritten, err.Error())
		result.Result = base64.StdEncoding.EncodeToString([]byte(resultString))
	}
	result.Result = base64.StdEncoding.EncodeToString([]byte("Wrote file without error"))
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
