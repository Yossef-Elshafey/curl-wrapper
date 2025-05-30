package wcurl

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"wcurl/app/command"
	"wcurl/app/models/endpoiot"
	"wcurl/app/storage"
)

// NOTE: Endpoint is a map no order gurantee
// NOTE: Add command stack, work with arrows
// NOTE: Regex is maniac ?

type Project map[string]endpoint.Endpoint

type WcurlWrapper struct {
	Data           []Project              `json:"data"`
	CommandHandler command.CommandHandler `json:"-"`
	storage        storage.Storage
}

func (w *WcurlWrapper) Init() {
	w.CommandHandler.Add("init", "initialize a new project", w.NewProject)
	w.CommandHandler.Add("curl", "Write regular curl commands like any (init new project if doesn't exist)", w.CurlHandler)
	w.CommandHandler.Add("list", "list endpoints", w.ListProjectEndpoints)
	w.CommandHandler.Add("exec", "Execute a command (exec { endpoint }.{ command num })", w.Execute)
}

func (w *WcurlWrapper) Load() WcurlWrapper {
	var ww WcurlWrapper
	path := w.storage.GetAbsoluteJsonFilePath()
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		ww.Data = make([]Project, 0)
		return ww
	}

	if len(data) == 0 {
		ww.Data = make([]Project, 0)
		return ww
	}

	err = json.Unmarshal(data, &ww)
	if err != nil {
		fmt.Println("error unmarshaling file:", err)
	}
	return ww
}

func (w *WcurlWrapper) Write() {
	j, err := json.Marshal(w)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}

	err = os.WriteFile(w.storage.GetAbsoluteJsonFilePath(), j, 0644)
	if err != nil {
		fmt.Println("WriteFile error:", err)
	}
}

func (w *WcurlWrapper) validateExistProject() bool {
	*w = w.Load()

	for _, proj := range w.Data {
		if _, ok := proj[w.storage.ProjectID()]; ok {
			return true
		}
	}
	return false
}

func (w *WcurlWrapper) GetProjectEndpoint() endpoint.Endpoint {
	ep := endpoint.Endpoint{}

	for _, projects := range w.Data {
		if e, ok := projects[w.storage.ProjectID()]; ok {
			ep = e
			break
		} else {
			ep = endpoint.Endpoint{}
		}
	}

	return ep
}

func (w *WcurlWrapper) NewProject() {
	*w = w.Load()

	if w.validateExistProject() {
		fmt.Println("Project already exist")
		return
	}

	project := Project{}
	project[w.storage.ProjectID()] = endpoint.Endpoint{}

	w.Data = append(w.Data, project)
	w.Write()
}

func (w *WcurlWrapper) CurlHandler() {
	*w = w.Load()
	ep := w.GetProjectEndpoint().AddEndpoint()

	if !w.validateExistProject() || len(w.Data) == 0 {
		project := Project{}
		project[w.storage.ProjectID()] = endpoint.Endpoint{}
		w.Data = append(w.Data, project)
	}

	for _, projects := range w.Data {
		projects[w.storage.ProjectID()] = ep
	}
	w.ShellExcuter(w.CommandHandler.GetUserInput())
	w.Write()
}

func (w WcurlWrapper) ListProjectEndpoints() {
	w = w.Load()
	w.GetProjectEndpoint().ListEndPoints()
}

func (w *WcurlWrapper) getExecValues() ([]string, error) {
	command := strings.TrimPrefix(w.CommandHandler.GetUserInput(), "exec")
	inp := strings.Split(command, "->")
	if len(inp) != 2 {
		return nil, errors.New("Unrecognzied exec command (ex. exec endpoint->0)")
	}
	return inp, nil
}

func (w WcurlWrapper) Execute() {
	limit, err := w.getExecValues()
	if err != nil {
		fmt.Println(err)
		return
	}
	w = w.Load()
	ep := w.GetProjectEndpoint()

	targetEp := strings.TrimSpace(limit[0])
	targetCmd, err := strconv.Atoi(strings.TrimSpace(limit[1]))
	if err != nil {
		fmt.Println("Error while handling target endpoint command")
	}

	cmd := ep.Ep[targetEp][targetCmd]
	w.ShellExcuter(cmd)
}

func (w WcurlWrapper) ShellExcuter(cmd string) {
	var ex *exec.Cmd
	ex = exec.Command("sh", "-c", cmd+" -s ") // -s slient, ignore networking values
	output, err := ex.CombinedOutput()

	if err != nil {
		fmt.Println("Error executing command:", err)
		fmt.Println("Command output:", string(output))
		return
	}

	fmt.Println(string(output))
}
