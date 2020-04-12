package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func index(c *cli.Context) (err error) {
	isPro := c.Bool("pro")

	first := c.Args().First()
	tail := c.Args().Tail()

	fullCmd := strings.Join(c.Args().Slice(), " ")

	logClient, err := initLogClient(fullCmd, isPro)
	if err != nil {
		return
	}

	if isPro {
		log.Println("====🚀 kan Pro 🚀====")
	} else {
		log.Println("====😃 kan Basic 😃====")
	}

	cmd := exec.Command(first, tail...)

	if isPro {
		var stdout io.ReadCloser

		cmd.Stderr = os.Stderr

		stdout, err = cmd.StdoutPipe()
		if err != nil {
			return
		}

		cmd.Start()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			println(line)
			logClient.PubLog(line)
		}

		err = cmd.Wait()
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	}

	logClient.CloseLog(err == nil)

	return nil
}
