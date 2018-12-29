package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TaskResult interface {
	MarshalJSON() ([]byte, error)
	ResultForTaskID(uuid.UUID) bool
}

type ClientTask interface {
	Execute() *Result
	IsIssuedTo(uuid.UUID) bool
}

type TaskMetaData struct {
	ID         string `json: "ID"`
	Target     string `json: "Target"`
	Type       string `json: "Type"`
	Completed  bool   `json: "Completed"`
	Title      string
	IssuedTime time.Time
}

type Result struct {
	ID     string
	Target string
	Result string
}

type Task struct {
	TaskMetaData
	Details json.RawMessage
}

type TaskReport struct {
	ID      string
	Target  string
	Type    string
	Title   string
	Details json.RawMessage
	Result  string
}

func (r TaskReport) String() string {
	template :=
		`Task Report:
Title: %s
Type: %s
Details:
  %s
Result:
  %s
`
	var detailsStr string
	switch r.Type {
	case "Get":
		gt := new(GetDetails)
		json.Unmarshal(r.Details, gt)
		detailsStr = gt.String()
	case "Put":
		pt := new(PutDetails)
		json.Unmarshal(r.Details, pt)
		detailsStr = pt.String()
	case "Execute":
		et := new(ExecuteDetails)
		json.Unmarshal(r.Details, et)
		detailsStr = et.String()
	}
	return fmt.Sprintf(template, r.Title, r.Type, detailsStr, r.Result)
}

func (t *Task) Init(target uuid.UUID, Type string) {
	t.ID = uuid.New().String()
	t.Target = target.String()
	t.Type = Type
	t.Completed = false
	t.IssuedTime = time.Now()
}

func (t *Task) Execute() (*Result, error) {
	// var err error
	var tr *Result
	if t.Completed {
		tr = new(Result)
		tr.Result = base64.StdEncoding.EncodeToString([]byte("Task already Completed"))
		return tr, nil
	}
	switch t.Type {
	case "Get":
		gt := new(GetDetails)
		err := json.Unmarshal(t.Details, gt)
		if err != nil {
			return nil, err
		}
		tr = gt.Execute()
		if err != nil {
			return nil, err
		}
	case "Put":
		pt := new(PutDetails)
		err := json.Unmarshal(t.Details, pt)
		if err != nil {
			return nil, err
		}
		tr = pt.Execute()
		if err != nil {
			return nil, err
		}
	case "Execute":
		et := new(ExecuteDetails)
		err := json.Unmarshal(t.Details, et)
		if err != nil {
			return nil, err
		}
		tr, err = et.Execute()
		if err != nil {
			return nil, err
		}
	}
	tr.ID = t.ID
	tr.Target = t.Target
	return tr, nil
}

func (t Task) String() string {
	var str string
	str = fmt.Sprintf("Task ID: %s\n", t.ID)
	str = fmt.Sprintf("%sTarget ID: %s\n", str, t.Target)
	str = fmt.Sprintf("%sType: %s\n", str, t.Type)
	switch t.Type {
	case "Get":
		details := new(GetDetails)
		err := json.Unmarshal(t.Details, details)
		if err != nil {
			fmt.Println("failed to unmarshal task details ", err.Error())
		}
		str = fmt.Sprintf("%sDetails: %v", str, details)
	case "Put":
		details := new(PutDetails)
		err := json.Unmarshal(t.Details, details)
		if err != nil {
			fmt.Println("failed to unmarshal task details ", err.Error())
		}
		str = fmt.Sprintf("%sDetails: %v", str, details)
	case "Execute":
		details := new(ExecuteDetails)
		err := json.Unmarshal(t.Details, details)
		if err != nil {
			fmt.Println("failed to unmarshal task details ", err.Error())
		}
		str = fmt.Sprintf("%sDetails: %v", str, details)
	}
	return str
}
