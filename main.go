package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	go func() {
		err := startMetricsServer(8080)
		if err != nil {
			log.Fatalf("error starting metrics server %v", err)
		}
	}()

	for {
		log.Info("Getting images in pods")
		imageList, err := getImagesFromAllPods()
		if err != nil {
			log.Fatalf("error getting images from pods %v", err)
		}

		imagesCount := len(imageList)

		for i, image := range imageList {
			log.WithFields(map[string]interface{}{
				"progress": fmt.Sprintf("%d/%d", i+1, imagesCount),
				"image":    image,
			}).Info("Checking image")
			exists, err := imageExistsInRegistry(image)
			if err != nil {
				log.Errorf("error checking image %s: %v", image, err)
			}
			if !exists {
				log.Errorf("image %s not found in registry", image)
				resolveMissingTotal.WithLabelValues(image).Set(1)
			} else {
				log.Infof("image %s found in registry", image)
				resolveMissingTotal.WithLabelValues(image).Set(0)
			}
		}
		time.Sleep(12 * time.Hour)
	}
}
