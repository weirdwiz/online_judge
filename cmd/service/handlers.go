package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type CompileRequest struct {
	Code     string `json:"code"`
	Language string `json:"lang"`
}

type CompileResponse struct {
	Output string `json:"output"`
}

func CompileCode(w http.ResponseWriter, r *http.Request) {
	var request CompileRequest
	var compiler string

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err.Error())
	}

	switch request.Language {
	case "C":
		compiler = "gcc"
	case "c":
		compiler = "gcc"
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	file, _ := os.Create("test.c")
	file.WriteString(request.Code)
	defer file.Close()
	cmd := exec.Command(compiler, "test.c")
	stdout, err := cmd.Output()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	cmd = exec.Command("./a.out")
	stdout, err = cmd.Output()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := &CompileResponse{
		Output: string(stdout),
	}

	data, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
