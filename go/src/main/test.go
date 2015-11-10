package main

import (
	"os/exec"
)

func main() {
	cmd := exec.Command("/home/code/mycode/go/src/main/test", "", "")
	stdout, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(stdout))
}
