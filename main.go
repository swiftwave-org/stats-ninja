package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/swiftwave-org/stats_ninja/host"
)

var serviceName = "stats_ninja"

var serviceTemplate = `
[Unit]
Description=Stats Ninja
After=network.target

[Service]
Type=simple
ExecStart=/bin/stats_ninja run {{.Endpoint}} {{.Auth}}
Restart=on-failure
RestartSec=10
KillMode=process
User=root

[Install]
WantedBy=multi-user.target
`

func main() {
	// fetch first argument
	args := os.Args[1:]
	// run cmd
	if len(args) > 0 && args[0] == "run" {
		endpointFlag := ""
		authFlag := ""
		if len(args) > 2 {
			endpointFlag = args[1]
			authFlag = args[2]
			run(endpointFlag, authFlag)
		} else {
			fmt.Println("Usage: stats_ninja run <submission_endpoint> <authorization_header_value>")
		}

	} else if len(args) > 0 && args[0] == "enable" {
		// enable cmd
		endpointFlag := ""
		authFlag := ""
		if len(args) > 2 {
			endpointFlag = args[1]
			authFlag = args[2]
			enable(endpointFlag, authFlag)
		} else {
			fmt.Println("Usage: stats_ninja run <submission_endpoint> <authorization_header_value>")
		}
	} else if len(args) > 0 && args[0] == "disable" {
		// disable cmd
		disable()
	} else {
		fmt.Println("Usage: stats_ninja <run|enable|disable>")
		os.Exit(1)
	}
}

func enable(submissionEndpoint, authorizationHeaderVal string) {
	// update template
	service := serviceTemplate
	service = strings.Replace(service, "{{.Endpoint}}", submissionEndpoint, -1)
	service = strings.Replace(service, "{{.Auth}}", authorizationHeaderVal, -1)
	// write to file
	file, err := os.Create(fmt.Sprintf("/etc/systemd/system/%s.service", serviceName))
	if err != nil {
		log.Println("Error creating service file: ", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(service)
	if err != nil {
		log.Println("Error writing to service file: ", err)
		return
	}
	// reload systemd
	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		log.Println("Error reloading systemd: ", err)
		return
	}
	// enable service
	err = exec.Command("systemctl", "enable", serviceName).Run()
	if err != nil {
		log.Println("Error enabling service: ", err)
		return
	}
	// start service
	err = exec.Command("systemctl", "start", serviceName).Run()
	if err != nil {
		log.Println("Error starting service: ", err)
		return
	}
	log.Println("Service enabled successfully")
}

func disable() {
	// stop service
	err := exec.Command("systemctl", "stop", serviceName).Run()
	if err != nil {
		log.Println("Error stopping service: ", err)
		return
	}
	// disable service
	err = exec.Command("systemctl", "disable", serviceName).Run()
	if err != nil {
		log.Println("Error disabling service: ", err)
		return
	}
	log.Println("Service disabled successfully")
}

func run(submissionEndpoint, authorizationHeaderVal string) {
	_, _ = host.Stats() // intentionally called. just to initialize current network stats
	// create a new docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Println("Error creating docker client:")
		panic(err)
	}
	for {
		<-time.After(1 * time.Minute)
		// fetch stats
		statsData, err := fetchStats(dockerClient)
		if err != nil {
			log.Println("Error fetching stats: ", err)
			continue
		}
		// convert to json
		jsonData, err := statsData.JSON()
		if err != nil {
			log.Println("Error converting stats to json: ", err)
			continue
		}
		// send to endpoint
		err = sendStats(submissionEndpoint, authorizationHeaderVal, jsonData)
		if err != nil {
			log.Println("Error sending stats to endpoint: ", err)
			continue
		}
	}
}

// private functions
func sendStats(submissionEndpoint string, authorizationHeaderVal string, jsonData []byte) error {
	// convert jsonData to a reader
	body := bytes.NewReader(jsonData)
	req, err := http.NewRequest("POST", submissionEndpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	// set the authorization header
	if authorizationHeaderVal != "" {
		req.Header.Set("Authorization", authorizationHeaderVal)
	}
	// send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	// close the response body
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	return nil
}
