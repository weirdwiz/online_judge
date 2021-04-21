package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/weirdwiz/online_judge/compile_microservice/pkg/sandbox"
)

func CompileCode(w http.ResponseWriter, r *http.Request) {
	var request sandbox.CompileRequest
	var compiler string
	if r.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		code := r.FormValue("code")
		lang := r.FormValue("lang")

		request.Code = code
		request.Language = lang
	}

	fileExtention := strings.ToLower(request.Language)
	switch fileExtention {
	case "c":
		compiler = "gcc"
	case "cpp":
		compiler = "g++"
	case "py":
		compiler = "python3"
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	file, _ := os.Create("test." + fileExtention)
	file.WriteString(request.Code)
	defer file.Close()

	cmd := exec.Command(compiler, file.Name())
	stdout, err := cmd.Output()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	if fileExtention == "c" || fileExtention == "cpp" {
		cmd = exec.Command("./a.out")
		stdout, err = cmd.Output()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	response := &sandbox.CompileResponse{
		Output: string(stdout),
	}

	data, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
