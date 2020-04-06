package main

import (
	"bufio"
	"fmt"
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

func indexPro(c *cli.Context, topic string) (err error) {
	logClient, err := initLogClient(topic)
	if err != nil {
		return
	}

	println("Start Pro")

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
		logClient.PubLog(line)
	}

	err = cmd.Wait()
	logClient.CloseLog(err == nil)

	return nil
}

func index(c *cli.Context) (err error) {
	recordFlag := c.Bool("pro")
	fullCmd := strings.Join(c.Args().Slice(), " ")

	if recordFlag {
		return indexPro(c, fullCmd)
	}

	client, err := initClient()
	if err != nil {
		return
	}

	cmd := consCmd(c)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	client.Email(fmt.Sprintf("✅ %s", fullCmd), "")
	return nil
}
