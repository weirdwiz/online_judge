package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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
		compiler = "python"
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	file, err := ioutil.TempFile("/tmp", compiler+"*."+fileExtention)
	if err != nil {
		panic(err)
	}

	defer os.Remove(file.Name())

	file.WriteString(request.Code)
	defer file.Close()

	stdout, err := compile(file, compiler)

	//	cmd := exec.Command(compiler, file.Name())
	//	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	//	if fileExtention == "c" || fileExtention == "cpp" {
	//		cmd = exec.Command("./a.out")
	//		stdout, err = cmd.Output()
	//		if err != nil {
	//			w.WriteHeader(http.StatusBadRequest)
	//			return
	//		}
	//	}

	response := &sandbox.CompileResponse{
		Output: stdout,
	}

	data, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func compile(file *os.File, compiler string) (string, error) {
	buf := new(strings.Builder)

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return buf.String(), err
	}

	var Cmd []string

	_, fileName := filepath.Split(file.Name())
	switch compiler {
	case "gcc":
		Cmd = []string{"/bin/sh", "-c", "gcc " + fileName + " && ./a.out"}
	case "python":
		Cmd = []string{"python", fileName}
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: compiler,
		Tty:   false,
		Cmd:   Cmd,
	}, nil, nil, nil, "")

	if err != nil {
		return buf.String(), err
	}

	defer cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})

	//	err = cli.CopyToContainer(ctx, resp.ID, "/"+fileName, file, types.CopyToContainerOptions{})
	//	if err != nil {
	//		return buf.String(), err
	//	}

	cmd := exec.Command("docker", "cp", file.Name(), resp.ID+":"+fileName)
	_, err = cmd.Output()

	if err != nil {
		return buf.String(), err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return buf.String(), err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return buf.String(), err
		}
	case <-statusCh:
	}

	//	if compiler == "gcc" {
	//		execResp, err := cli.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
	//			Cmd: []string{"./a.out"},
	//		})
	//		if err != nil {
	//			return buf.String(), err
	//		}
	//		cli.ContainerExecStart(ctx, execResp.ID, types.ExecStartCheck{})
	//	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return buf.String(), err
	}

	fmt.Println(buf)
	io.Copy(buf, out)

	return buf.String(), nil
}
