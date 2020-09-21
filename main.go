package main

import (
	"log"
	"os"
	"path"

	care_sdk "github.com/byte-care/care-sdk-go"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var configFilePath string
var version string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configFilePath = path.Join(homeDir, ".carerc.yml")
	if _, err := os.Stat(configFilePath); err == nil {

	} else if os.IsNotExist(err) {
		panic("Please use care-config to init the env. 😀")
	} else {
		panic(err)
	}

	viper.SetConfigName(".carerc") // name of config file (without extension)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME")   // call multiple times to add many search paths

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(err)
	}
}

func initClient() (client *care_sdk.Client, err error) {
	AccessKey := viper.GetString("access-key")
	SecretKey := viper.GetString("secret-key")

	client, err = care_sdk.NewClient(AccessKey, SecretKey)
	if err != nil {
		panic(err)
	}

	return
}

func initLogClient(topic string, isPro bool) (client *care_sdk.LogClient, err error) {
	AccessKey := viper.GetString("access-key")
	SecretKey := viper.GetString("secret-key")

	client, err = care_sdk.NewLogClient(AccessKey, SecretKey, topic, isPro)
	if err != nil {
		panic(err)
	}

	return
}

func main() {
	app := &cli.App{
		Name:     "care",
		Usage:    "👧💻 CLI for Care 💻👦",
		HelpName: "care",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "pro"},
		},
		Action:  index,
		Version: version,
	}
	app.UseShortOptionHandling = true

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
