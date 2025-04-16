package main

import "testing"

func TestImageExistsInRegistry(t *testing.T) {
	type args struct {
		imageid string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"When using the docker registry and a valid image then it should return true",
			args{
				imageid: "alpine:latest",
			},
			true,
			false,
		},
		{
			"When using the github registry and a valid image then it should return true",
			args{
				imageid: "ghcr.io/riskident/jelease:v0.6.3",
			},
			true,
			false,
		},
		{
			"When using the docker registry and an invalid image then it should return false",
			args{
				imageid: "alpine:latest-invalid",
			},
			false,
			true,
		},
		{
			"When a valid image is used from a private registry then it should return true",
			args{
				imageid: "ghcr.io/riskident/image-availability-exporter/busybox:latest",
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := imageExistsInRegistry(tt.args.imageid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageExistsInRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ImageExistsInRegistry() got = %v, want %v", got, tt.want)
			}
		})
	}
}
