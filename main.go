package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"time"
)

var (
	port            int
	interval        time.Duration
	dockerConfigDir string
)

func main() {
	app := cli.NewApp()
	app.Name = "Image Availability Prometheus Exporter"
	app.Usage = "Expose Prometheus metrics about image availability in your Kubernetes cluster"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "port",
			Value:       8080,
			Usage:       "Port to expose metrics on",
			Destination: &port,
			EnvVar:      "PORT",
		},
		cli.DurationFlag{
			Name:        "interval",
			Value:       12 * time.Hour,
			Usage:       "Interval to check images",
			Destination: &interval,
			EnvVar:      "INTERVAL",
		},
		cli.StringFlag{
			Name:        "docker-config-dir",
			Value:       "",
			Usage:       "Directory where the docker config file is located",
			Destination: &dockerConfigDir,
			EnvVar:      "DOCKER_CONFIG_DIR",
		},
	}
	app.Action = run
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	go func() {
		err := startMetricsServer(port)
		if err != nil {
			log.Fatalf("error starting metrics server %v", err)
		}
	}()

	log.Infof("Starting image availability check every %s", interval)

	if dockerConfigDir != "" {
		log.Infof("Using docker config dir %s", dockerConfigDir)
	}

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
				log.Errorf("Error checking image %s: %v", image, err)
			}
			if !exists {
				log.Errorf("Image %s not found in registry", image)
				resolveMissingTotal.WithLabelValues(image).Set(1)
			} else {
				log.Infof("Image %s found in registry", image)
				resolveMissingTotal.WithLabelValues(image).Set(0)
			}
		}
		time.Sleep(interval)
	}
}
