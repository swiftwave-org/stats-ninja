package main

import (
	"bytes"
	"github.com/docker/docker/client"
	"github.com/swiftwave-org/stats_ninja/host"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	/* Configure the environment variables
	* SUBMISSION_ENDPOINT: The endpoint to submit the stats to
	* AUTHORIZATION_HEADER_VAL: The value of the authorization header
	* DOCKER_HOST: unix or tcp socket to connect to
	* This will send the stats to the endpoint using the authorization header
	*
	 */
	submissionEndpoint := os.Getenv("SUBMISSION_ENDPOINT")
	authorizationHeaderVal := os.Getenv("AUTHORIZATION_HEADER_VAL")
	// reject if the submission endpoint is not set
	if submissionEndpoint == "" {
		panic("SUBMISSION_ENDPOINT is not set")
	}
	if os.Getenv("DOCKER_HOST") == "" {
		panic("DOCKER_HOST is not set")
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

// private function to send stats to the endpoint
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
