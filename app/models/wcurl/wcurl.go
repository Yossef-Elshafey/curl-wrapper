package wcurl

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"wcurl/app/command"
	"wcurl/app/models/endpoiot"
	"wcurl/storage"
)

type WcurlWrapper struct {
	Data     map[string]endpoint.Endpoint `json:"data"`
	Listener command.CommandHandler       `json:"-"`
}

var (
	// TODO: Cache loaded data
	exposeInput string
)

func (w *WcurlWrapper) InputRecivier(inp string) {
	exposeInput = inp
}

func (w *WcurlWrapper) SetCommand() {
	w.Listener.Add("init", "initialize a new project", w.NewProject)
	w.Listener.Add("curl", "Write regular curl commands like any", w.CurlHandler)
}

func (w *WcurlWrapper) FilePath() string {
	return storage.GetAbsoluteJsonFilePath()
}

func (w *WcurlWrapper) Load() WcurlWrapper {
	var ww WcurlWrapper
	path := storage.GetAbsoluteJsonFilePath()
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		ww.Data = make(map[string]endpoint.Endpoint)
		return ww
	}

	if len(data) == 0 {
		ww.Data = make(map[string]endpoint.Endpoint)
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

	err = os.WriteFile(w.FilePath(), j, 0644)
	if err != nil {
		fmt.Println("WriteFile error:", err)
	}
}

func (w *WcurlWrapper) validate(hash string) bool {
	*w = w.Load()
	if _, ok := w.Data[hash]; ok {
		return true
	}
	return false
}

func (w *WcurlWrapper) NewProject() {
	h := storage.HashExecPath()
	*w = w.Load()

	if w.validate(h) {
		return
	}

	w.Data[h] = endpoint.Endpoint{Ep: map[string][]string{"": make([]string, 0)}}
	w.Write()
}

func (w *WcurlWrapper) extractEndpoint(s string) string {
	re := regexp.MustCompile(`^(?:https?:\/\/)?[^\/]+(\/.*)`)
	matches := re.FindAllStringSubmatch(s, -1)
	return matches[0][1]
}

func (w *WcurlWrapper) CurlHandler() {
	userInput := w.Listener.Get()
	h := storage.HashExecPath()
	ep := w.extractEndpoint(userInput)
	*w = w.Load()
	fmt.Println(w.Data[h].Ep[ep])
}
