package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func imageExistsInRegistry(image string) (bool, error) {
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
