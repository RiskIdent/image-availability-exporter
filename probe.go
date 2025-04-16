package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

// imageExistsInAnotherRegistry checks if an image exists in a registry other than dockerhub.
// It does work for dockerhub as well, but we reach the rate limit very quickly.
// It uses the docker command to check if the image exists in the registry which makes it slow.
// ToDo: Implement this function without using an external command
func imageExistsInAnotherRegistry(image string) (bool, error) {
	command := []string{"manifest", "inspect", image}
	if dockerConfigDir != "" {
		command = []string{"--config", dockerConfigDir, "manifest", "inspect", image}
	}
	cmd := exec.Command("docker", command...)
	var buf bytes.Buffer
	cmd.Stderr = &buf
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("error running docker manifest inspect %s: %v, %s", image, err, buf.String())
	}
	return true, nil
}

func imageExistsInDockerhub(image string) (bool, error) {
	if !strings.Contains(image, ":") {
		image = fmt.Sprintf("%s:latest", image)
	}
	if strings.HasPrefix(image, "docker.io/") {
		image = strings.Replace(image, "docker.io/", "", 1)
	}
	imageAndTag := strings.Split(image, ":")
	if !strings.Contains(image, "/") {
		imageAndTag[0] = fmt.Sprintf("library/%s", imageAndTag[0])
	}

	resp, err := http.Get(fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/%s/", imageAndTag[0], imageAndTag[1]))
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}
	b, _ := io.ReadAll(resp.Body)
	return string(b) != "\"Resource not found\"" && string(b) != "Tag not found", nil
}

func imageExistsInRegistry(image string) (bool, error) {
	registry := getRegistryFromImageName(image)
	if registry == "docker.io" {
		return imageExistsInDockerhub(image)
	}
	return imageExistsInAnotherRegistry(image)
}

func getRegistryFromImageName(image string) string {
	parts := strings.Split(image, "/")

	if len(parts) > 1 && strings.Contains(parts[0], ".") {
		return parts[0]
	}

	return "docker.io"
}
