package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/swiftwave-org/stats_ninja/host"
)

func main() {
	/* Configure the environment variables
	* SUBMISSION_ENDPOINT: The endpoint to submit the stats to
	* AUTHORIZATION_HEADER_VAL: The value of the authorization header
	* DOCKER_HOST: unix or tcp socket to connect to
	* HOSTNAME: The hostname of the host
	* This will send the stats to the endpoint using the authorization header
	*
	* Configure Volume Mounts
	* <docker socket of host>:/var/run/docker.sock
	 */
	submissionEndpoint := os.Getenv("SUBMISSION_ENDPOINT")
	authorizationHeaderVal := os.Getenv("AUTHORIZATION_HEADER_VAL")
	hostname := os.Getenv("HOSTNAME")
	// reject if the submission endpoint is not set
	if submissionEndpoint == "" {
		panic("SUBMISSION_ENDPOINT is not set")
	}
	if os.Getenv("DOCKER_HOST") == "" {
		panic("DOCKER_HOST is not set")
	}
	// if hostname is not set, fetch it from the system
	if hostname == "" {
		panic("HOSTNAME is not set")
	}
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
		// set hostname
		statsData.Hostname = hostname
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
