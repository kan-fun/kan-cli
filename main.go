package main

import (
	"log"
	"os"
	"path"

	caresdk "github.com/byte-care/care-sdk-go"
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

	viper.SetConfigName(".carerc") // name of config file (without extension)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME")   // call multiple times to add many search paths
}

func initClient() (client *caresdk.Client, err error) {
	AccessKey := viper.GetString("access-key")
	SecretKey := viper.GetString("secret-key")

	client, err = caresdk.NewClient(AccessKey, SecretKey)
	if err != nil {
		panic(err)
	}

	return
}

func initLogClient(topic string, isPro bool) (client *caresdk.LogClient, err error) {
	AccessKey := viper.GetString("access-key")
	SecretKey := viper.GetString("secret-key")

	client, err = caresdk.NewLogClient(AccessKey, SecretKey, topic, isPro)
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
			&cli.BoolFlag{Name: "update"},
			&cli.StringFlag{Name: "access-key"},
			&cli.StringFlag{Name: "secret-key"},
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
