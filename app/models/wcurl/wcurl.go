package wcurl

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"wcurl/app/command"
	"wcurl/app/models/endpoiot"
	"wcurl/storage"
)

type WcurlWrapper struct {
	Data     map[string]endpoint.Endpoint `json:"data"`
	Listener command.CommandHandler       `json:"-"`
}

var WcurlWrapperCache WcurlWrapper

func (w *WcurlWrapper) SetCommand() {
	w.Listener.Add("init", "initialize a new project", w.NewProject)
	w.Listener.Add("ep", "does something", w.NewProject)
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
	WcurlWrapperCache = ww
	return ww
}

func (w *WcurlWrapper) Write() {
	fmt.Println(w)
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

func (w *WcurlWrapper) NewProject() {
	fmt.Println("hello idiot")
	h := storage.HashExecPath()
	*w = w.Load()
	w.Data = make(map[string]endpoint.Endpoint)
	w.Data[h] = endpoint.Endpoint{Ep: map[string][]string{"admin/local": make([]string, 0)}}
	w.Write()
}
