package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudfoundry/cli/plugin"
)

type mapevent struct{}

const eventUrl = "https://hottopic.apps.bogata.cf-app.com/map"

type mapEventPost struct {
	app   string
	topic string
}

func (m *mapevent) Run(cliConnection plugin.CliConnection, args []string) {
	if len(args) > 0 && args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}

	if len(args) != 3 {
		fmt.Println("Must supply APP_NAME and TOPIC_NAME")
		os.Exit(1)
	}

	appName := args[1]
	topicName := args[2]

	postBody := mapEventPost{app: appName, topic: topicName}
	fmt.Printf("%+v", postBody)

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(postBody)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	res, err := http.Post(eventUrl, "application/json", b)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	fmt.Println(res)

	if res.StatusCode == http.StatusCreated {
		fmt.Printf("Successfully mapped event topic %s to app %s\n", topicName, appName)
	} else {
		fmt.Println("this is hackday project, this shouldnt have happened")
	}
}

func (m *mapevent) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "map-event",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "map-event",
				HelpText: "map an event topic to an app",
				UsageDetails: plugin.Usage{
					Usage:   "cf map-event APP_NAME TOPIC_NAME",
					Options: map[string]string{},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(mapevent))
}
