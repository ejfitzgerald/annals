package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/ejfitzgerald/annals"
)

var commandMapping = map[string]string{
	"clang":   "/usr/bin/clang",
	"clang++": "/usr/bin/clang++",
	"___mock": "/usr/bin/clang++",
}

func submitMetadata(duration time.Duration) {

	metadata := annals.CompilationMetadata{
		duration,
		os.Args,
	}

	encodedData, err := json.Marshal(metadata)
	if err != nil {
		fmt.Println("Error encoding data: ", err)
		return
	}

	fmt.Println("Encoded Data", encodedData)

	// try and submit the data to the server
	url := "http://127.0.0.1:9100/compilation"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(encodedData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Unable to submit data: ", err)
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response from server")
	}
}

func main() {
	exitCode := 1

	command := os.Args[0]
	baseCommand := path.Base(command)

	if newCommand, ok := commandMapping[baseCommand]; ok {
		// get the arguments
		args := os.Args[1:]

		// capture the start time
		start := time.Now()

		// run the sub command
		cmd := exec.Command(newCommand, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		finished := time.Now()
		delta := finished.Sub(start)

		// try and submit the result
		submitMetadata(delta)

		if err == nil {

			exitCode = 0
		}
	} else {
		fmt.Println("Failed to lookup command for: ", command)
	}

	os.Exit(exitCode)
}
