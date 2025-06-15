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
	w.CommandHandler.Add("exec", "Execute a command, exec <endpoint> -> <command number>", w.Execute)
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
		}
	}

	return ep
}

func (w *WcurlWrapper) NewProject() {
	*w = w.Load()

	if w.validateExistProject() {
		fmt.Printf("\n\rProject already exist")
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
	// TODO: take more than one path or more than one command or comined, threading idiot
	// TODO: need to write a command line parser functionality something like flags.parse()

	raw := strings.TrimSpace(strings.TrimPrefix(w.CommandHandler.GetUserInput(), "exec"))
	parts := strings.SplitN(raw, "->", 2)

	if len(parts) != 2 {
		return nil, errors.New("Invalid exec command. Expect: exec <endpoint> -> <number>")
	}

	path, commandNum := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

	if path == "" || commandNum == "" {
		return nil, errors.New("Invalid exec command. Expect: exec <endpoint> -> <number>")
	}

	return []string{path, commandNum}, nil
}

func (w WcurlWrapper) Execute() {
	limit, err := w.getExecValues()
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	w = w.Load()
	ep := w.GetProjectEndpoint()

	targetEp := strings.TrimSpace(limit[0])
	targetCmd, err := strconv.Atoi(strings.TrimSpace(limit[1]))
	if err != nil {
		fmt.Printf("\n\rError while handling target endpoint command")
	}

	cmd := ep.Ep[targetEp][targetCmd]
	w.ShellExcuter(cmd)
}

func (w WcurlWrapper) ShellExcuter(cmd string) {
	var ex *exec.Cmd
	ex = exec.Command("sh", "-c", cmd+" -s ")
	output, err := ex.CombinedOutput()
	if err != nil {
		fmt.Printf("\n\rError executing command %s, make sure server is up?", cmd)
		return
	}

	fmt.Printf("\n\r%s", string(output))
}
