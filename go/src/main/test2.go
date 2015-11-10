package main

import (
	"fmt"
	"os"
)

func main() {
	attr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	p, err := os.StartProcess("/home/packerdir/packer", []string{"/home/packerdir/packer", "build", "/home/jsondir/centos66.json"}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
}
