package main

import (
	"testing"
)

func Test_getImagesFromAllPods(t *testing.T) {
	images, err := getImagesFromAllPods()
	if err != nil {
		t.Errorf("getImagesFromAllPods() error = %v", err)
	}
	if len(images) == 0 {
		t.Errorf("getImagesFromAllPods() = %v, want more than 0", images)
	}
}
