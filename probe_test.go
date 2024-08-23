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
				imageid: "docker-riskident.2rioffice.com/platform/aow-reminder:1.1.0",
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := imageExistsInAnotherRegistry(tt.args.imageid)
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

func TestImageExistsInDockerhub(t *testing.T) {
	type args struct {
		image string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"When using a valid image then it should return true",
			args{
				image: "alpine:latest",
			},
			true,
			false,
		},
		{
			"When using an invalid image then it should return false",
			args{
				image: "alpine:latest-invalid",
			},
			false,
			false,
		},
		{
			"When the name of the image does not contain a tag then the tag is latest",
			args{
				image: "alpine",
			},
			true,
			false,
		},
		{
			"When the image has a full name then the registry part is deleted",
			args{
				image: "docker.io/grafana/grafana:11.1.0",
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := imageExistsInDockerhub(tt.args.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageExistsInDockerhub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("imageExistsInDockerhub() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRegistryFromImageName(t *testing.T) {
	type args struct {
		image string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"When the image is from the docker registry then it should return docker.io",
			args{
				image: "alpine:latest",
			},
			"docker.io",
		},
		{
			"When the image is from the github registry then it should return ghcr.io",
			args{
				image: "ghcr.io/riskident/jelease:v0.6.3",
			},
			"ghcr.io",
		},
		{
			"When the image is from a private registry then it should return the registry",
			args{
				image: "docker-riskident.2rioffice.com/platform/aow-reminder:1.1.0",
			},
			"docker-riskident.2rioffice.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRegistryFromImageName(tt.args.image); got != tt.want {
				t.Errorf("getRegistryFromImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}
