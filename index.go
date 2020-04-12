package main

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func consCmd(c *cli.Context) *exec.Cmd {
	first := c.Args().First()
	tail := c.Args().Tail()

	cmd := exec.Command(first, tail...)

	return cmd
}

func index(c *cli.Context) (err error) {
	proFlag := c.Bool("pro")
	fullCmd := strings.Join(c.Args().Slice(), " ")

	logClient, err := initLogClient(fullCmd, proFlag)
	if err != nil {
		return
	}

	if proFlag {
		println("====🚀 kan Pro 🚀====")
	} else {
		println("====😃 kan Basic 😃====")
	}

	cmd := consCmd(c)
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		println(line)
		if proFlag {
			logClient.PubLog(line)
		}
	}

	err = cmd.Wait()
	logClient.CloseLog(err == nil)

	return nil
}
