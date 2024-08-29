# image-availability-exporter

A Prometheus exporter exposing metrics about unsupported images in your cluster

## Description

This exporter will scan all images in your cluster and expose metrics about the 
availability of these images. It will check if the image is available in the 
configured registry every X time (default: 12h) and if it is not, it will expose
a metric with the image name and the tag.

If the metric `image_missing` is greater than 0, it means that there are images
in your cluster that are not available in the configured registry.
