package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
)

func TodoCMD(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	var command struct {
		Command string `json:"command"`
	}

	err := json.NewDecoder(r.Body).Decode(&command)
	if err != nil {
		http.Error(w, "failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	output, err := executeShellCommand(command.Command)
	if err != nil {
		http.Error(w, "exec failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"output": string(output),
	}
	json.NewEncoder(w).Encode(response)
}

func executeShellCommand(cmd string) ([]byte, error) {
	command := exec.Command("sh", "-c", cmd)
	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out

	err := command.Run()
	return out.Bytes(), err
}
