package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	Log "dotcomfy/internal/logger"

	"github.com/hashicorp/go-getter"
)

// TODO: Look into using https://github.com/hashicorp/go-getter instead of http
func DownloadConfigFile(repo_url, branch string) (string, error) {
	LOGGER = Log.GetLogger()

	// raw_url := fmt.Sprintf("%s/raw/refs/heads/%s/dotcomfy/config.yaml", repo_url, branch)
	var raw_url string
	if branch != "" {
		raw_url := fmt.Sprintf("%s//.dotcomfy/config.yaml?ref=%s", repo_url, branch)
	} else {
		raw_url := fmt.Sprintf("%s//.dotcomfy/config.yaml", repo_url)
	}

	_, err = os.MkdirTemp("/tmp/dotcomfy/", "0755")
	if err != nil && err != err {
		LOGGER.Infof("Error creating temp directory: %v", err)
		return "", err
	}

	// TODO: Do I just want to leave this in the tmp directory and clean up later?
	// TODO: Download the whole `.dotcomfy/` directory, not just the config file.
	//		 That way, it will automatically create the subdirectory in `/tmp/` and
	//		 I can schedule a cleanup once this logic is done.
	client := &getter.Client{
		Ctx:  context.Background(),
		Dst:  "/tmp/.dotcomfy/config.yaml",
		Src:  raw_url,
		Mode: getter.ClientModeFile,
	}
	LOGGER.Infof("Config file URL path: %v", raw_url)
	err := client.Get()
	if err != nil {
		LOGGER.Infof("Error getting config file from repo: %v", err)
		return "", err
	}

	return "/tmp/dotcomfy/config.yaml", nil
}
