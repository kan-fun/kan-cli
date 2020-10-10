package main

import (
	"bufio"
	"errors"
	"fmt"
	caresdk "github.com/byte-care/care-sdk-go"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	goupdate "github.com/inconshreveable/go-update"
	"github.com/urfave/cli/v2"
)

func getRemoteVersion() (remoteVersion string, err error) {
	osString := runtime.GOOS

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.bytecare.xyz/bin", nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("platform", osString)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	remoteVersion = string(respBody)

	return
}

func getBinary(remoteVersion string) (reader io.ReadCloser, err error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://bin.bytecare.xyz/%s/%s", runtime.GOOS, remoteVersion)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("network error")
	}

	reader = resp.Body

	return
}

func update() {
	remoteVersion, err := getRemoteVersion()
	if err != nil {
		return
	}

	if remoteVersion != version {
		println("ByteCare updating...")
		reader, err := getBinary(remoteVersion)
		if err != nil {
			if reader != nil {
				reader.Close()
			}
			return
		}
		defer reader.Close()

		var options goupdate.Options
		_ = goupdate.Apply(reader, options)
	}
}

func logChanHandler(logClient *caresdk.LogClient, logChan chan string, done chan int) {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	buffer := strings.Builder{}

	for {
		select {
		case <-ticker.C:
			{
				if buffer.Len() != 0 {
					logClient.PubLog(buffer.String())
					buffer.Reset()
				}
			}
		case s, ok := <-logChan:
			{
				buffer.WriteString(s)
				if buffer.Len() == 10 {
					logClient.PubLog(buffer.String())
					buffer.Reset()
				}

				if !ok {
					break
				}
			}
		}
		done <- 1
	}
}

func index(c *cli.Context) (err error) {
	accessKey := c.String("access-key")
	secretKey := c.String("secret-key")

	if accessKey != "" && secretKey != "" {
		viper.Set("access-key", accessKey)
		viper.Set("secret-key", secretKey)

		if err := viper.WriteConfigAs(configFilePath); err != nil {
			panic(err)
		}

		return
	}

	isUpdate := c.Bool("update")
	if isUpdate {
		update()
		return
	}

	if _, err := os.Stat(configFilePath); err == nil {

	} else if os.IsNotExist(err) {
		println("Please use --access-key and --secret-key to set token. 😀")
		return err
	} else {
		panic(err)
	}

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(err)
	}

	disablePub := c.Bool("disable-pub")
	isPro := !disablePub

	first := c.Args().First()
	tail := c.Args().Tail()

	fullCmd := strings.Join(c.Args().Slice(), " ")

	logClient, err := initLogClient(fullCmd, isPro)
	if err != nil {
		return
	}

	logChan := make(chan string, 15)

	done := make(chan int)
	go logChanHandler(logClient, logChan, done)

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
			logChan <- line
		}

		err = cmd.Wait()
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	}

	logClient.CloseLog(err == nil)
	close(logChan)

	update()
	<-done

	return nil
}
