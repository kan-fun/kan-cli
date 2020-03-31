package main

import (
	"github.com/urfave/cli/v2"

	"log"
	"os/exec"
	"os"
	"fmt"
	"strings"
)

func index(c *cli.Context) error {
	recordFlag := c.Bool("record")
	if recordFlag {
		panic("Not Support")
	}

	first := c.Args().First()
	tail := c.Args().Tail()

	cmd := exec.Command(first, tail...)
	log.Printf("Running command and waiting for it to finish...")
	
	cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
	cmd.Run()
	
	fullCmd := strings.Join(c.Args().Slice(), " ")
	client.Email(fmt.Sprintf("✅ %s", fullCmd), "123")
	return nil
}