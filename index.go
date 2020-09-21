package main

import (
	"bufio"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func index(c *cli.Context) (err error) {
	token := c.String("token")
	if token != "" {
		viper.Set("token", token)

		if err := viper.WriteConfigAs(configFilePath); err != nil {
			panic(err)
		}

		return
	}

	if _, err := os.Stat(configFilePath); err == nil {

	} else if os.IsNotExist(err) {
		panic("Please use --token to set token. 😀")
	} else {
		panic(err)
	}

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(err)
	}

	isPro := c.Bool("pro")

	first := c.Args().First()
	tail := c.Args().Tail()

	fullCmd := strings.Join(c.Args().Slice(), " ")

	logClient, err := initLogClient(fullCmd, isPro)
	if err != nil {
		return
	}

	if isPro {
		log.Println("====🚀 care Pro 🚀====")
	} else {
		log.Println("====😃 care Basic 😃====")
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
