package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/bblenard/C2/endpoints"
	"github.com/bblenard/C2/logging"
	"github.com/bblenard/C2/storage"
	"github.com/bblenard/C2/types"
)

func saveTasks(w http.ResponseWriter, req *http.Request) {
	logging.Journal.Logf(logging.HTTP, "Endpoint: %s\n", req.URL)
	responseEncoder := json.NewEncoder(w)
	task := new(types.Task)
	jd := json.NewDecoder(req.Body)
	err := jd.Decode(task)
	if err != nil {
		logging.Journal.Logf(logging.Error, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		responseEncoder.Encode(err.Error())
		return
	}
	err = storage.DB.SaveTasks(*task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	responseEncoder.Encode("")
}

func saveTaskResults(w http.ResponseWriter, req *http.Request) {
	logging.Journal.Logf(logging.HTTP, "Endpoint: %s\n", req.URL)
	responseEncoder := json.NewEncoder(w)
	tr := new(types.Result)
	jd := json.NewDecoder(req.Body)
	err := jd.Decode(tr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		responseEncoder.Encode("") // I don't want to give info to the client.
		logging.Journal.Logf(logging.Error, err.Error())
		return
	}
	responseEncoder.Encode("")
	err = storage.DB.SaveTaskResults(*tr)
	if err != nil {

	}
}

func saveClient(w http.ResponseWriter, req *http.Request) {
	logging.Journal.Logf(logging.HTTP, "Endpoint: %s\n", req.URL)
	responseEncoder := json.NewEncoder(w)
	requestDecoder := json.NewDecoder(req.Body)
	c := new(types.Client)
	c.New()
	err := requestDecoder.Decode(&c.Details)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		responseEncoder.Encode("Error")
	}
	storage.DB.SaveClients(*c)
	responseEncoder.Encode(c.ID)
}

func getClients(w http.ResponseWriter, req *http.Request) {
	logging.Journal.Logf(logging.HTTP, "Endpoint: %s\n", req.URL)
	clients, err := storage.DB.GetClients()
	if err != nil {
		logging.Journal.Logf(logging.Error, err.Error())
	}
	err = json.NewEncoder(w).Encode(clients)
	if err != nil {
		logging.Journal.Logf(logging.Error, err.Error())
	}
}

func getClientByID(w http.ResponseWriter, req *http.Request) {
	logging.Journal.Logf(logging.HTTP, "Endpoint: %s\n", req.URL)
	responseEncoder := json.NewEncoder(w)
	cUUIDRaw := req.Header.Get("CID")
	_, err := uuid.Parse(cUUIDRaw)
	if err != nil {
		responseEncoder.Encode(fmt.Sprintf("error: failed to parse CID %s", err.Error()))
		return
	}
	client, err := storage.DB.GetClientByID(cUUIDRaw)
	if err != nil {
		responseEncoder.Encode(fmt.Sprintf("error: failed to retrieve client %s", err.Error()))
		return
	}
	err = responseEncoder.Encode(client)
	if err != nil {
		responseEncoder.Encode(err.Error())
		return
	}
}

func getCompletedTasksByClient(w http.ResponseWriter, req *http.Request) {
	logging.Journal.Logf(logging.HTTP, "Endpoint: %s\n", req.URL)
	responseEncoder := json.NewEncoder(w)
	cUUIDRaw := req.Header.Get("CID")
	_, err := uuid.Parse(cUUIDRaw)
	if err != nil {
		fmt.Printf("header parse: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		responseEncoder.Encode(fmt.Sprintf("error: failed to parse CID %s", err.Error()))
		return
	}
	taskReports, err := storage.DB.GetCompletedTasksByClient(cUUIDRaw)
	if err != nil {
		fmt.Printf("storage call: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		responseEncoder.Encode(fmt.Sprintf("error: failed to retrieve tasks: %s", err.Error()))
		return
	}
	err = responseEncoder.Encode(taskReports)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		responseEncoder.Encode(fmt.Sprintf("error: failed to encode response: %s", err.Error()))
		return
	}
}

func getPendingTasksByClient(w http.ResponseWriter, req *http.Request) {
	logging.Journal.Logf(logging.HTTP, "Endpoint: %s\n", req.URL)
	responseEncoder := json.NewEncoder(w)
	cUUIDRaw := req.Header.Get("CID")
	_, err := uuid.Parse(cUUIDRaw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseEncoder.Encode(fmt.Sprintf("error: failed to parse CID %s", err.Error()))
		return
	}
	tasks, err := storage.DB.GetPendingTasksByClient(cUUIDRaw)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		responseEncoder.Encode(fmt.Sprintf("error: failed to retrieve tasks: %s", err.Error()))
		return
	}

	err = responseEncoder.Encode(tasks)
	if err != nil {
		fmt.Println("Encoding tasks")
		w.WriteHeader(http.StatusInternalServerError)
		responseEncoder.Encode(fmt.Sprintf("error: failed to encode response: %s", err.Error()))
		return
	}
}

// func exportDB(w http.ResponseWriter, req *http.Request) {
// }

// func importDB(w http.ResponseWriter, req *http.Request) {
// }

func ServerMode() error {
	logging.NewLogger(logging.Stdout, logging.All)
	logging.Journal.Logf(logging.Stage, "Creating Storage")
	err := storage.NewStorage(storage.MySQL)
	if err != nil {
		fmt.Println(err)
	}
	logging.Journal.Logf(logging.Database, "Database Status at Start: %s", storage.DB)
	http.HandleFunc(endpoints.SAVETASK, saveTasks) // Internal
	http.HandleFunc(endpoints.SAVETASKRESULT, saveTaskResults)
	http.HandleFunc(endpoints.SAVECLIENT, saveClient)
	http.HandleFunc(endpoints.GETCLIENTS, getClients)       // Internal
	http.HandleFunc(endpoints.GETCLIENTBYID, getClientByID) // Internal
	http.HandleFunc(endpoints.GETPENDINGTASKSBYCLIENT, getPendingTasksByClient)
	http.HandleFunc(endpoints.GETCOMPLETEDTASKSBYCLIENT, getCompletedTasksByClient) // Internal
	return http.ListenAndServe("0.0.0.0:8888", nil)
}

func main() {
	err := ServerMode()
	fmt.Println("Server Mode Exiting")
	if err != nil {
		panic(err)
	}
}
