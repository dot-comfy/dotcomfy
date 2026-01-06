package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	Log "dotcomfy/internal/logger"

	"github.com/hashicorp/go-getter"
)

func DownloadConfigFile(repo_url, branch string) (string, error) {
	LOGGER = Log.GetLogger()

	// Remove .git extension if present for raw file URL
	repo_url = strings.TrimSuffix(repo_url, ".git")

	// Convert github.com URL to raw.githubusercontent.com for direct file access
	raw_base_url := strings.Replace(repo_url, "github.com/", "raw.githubusercontent.com/", 1)

	// Use GitHub's standard raw file URL format
	var raw_url string
	if branch != "" {
		raw_url = fmt.Sprintf("%s/%s/dotcomfy/config.yaml", raw_base_url, branch)
	} else {
		raw_url = fmt.Sprintf("%s/main/dotcomfy/config.yaml", raw_base_url)
	}

	// Validate URL format for GitHub
	if !strings.Contains(repo_url, "github.com") {
		LOGGER.Warnf("Config download optimized for GitHub, may not work for: %s", repo_url)
	}

	// Create the base temp directory if it doesn't exist
	base_temp_dir := "/tmp/dotcomfy"
	os.MkdirAll(base_temp_dir, 0755)

	temp_dir_path, err := os.MkdirTemp(base_temp_dir, "config-")
	if err != nil {
		LOGGER.Infof("Error creating temp directory: %v", err)
		// Continue with fallback directory
		temp_dir_path = base_temp_dir
		os.MkdirAll(temp_dir_path, 0755)
	}

	client := &getter.Client{
		Ctx:  context.Background(),
		Dst:  temp_dir_path + "/config.yaml",
		Src:  raw_url,
		Mode: getter.ClientModeFile,
	}
	LOGGER.Infof("Config file URL path: %v", raw_url)
	err = client.Get()
	if err != nil {
		LOGGER.Infof("Error getting config file from repo: %v", err)
		return "", err
	}

	// Verify file was downloaded successfully
	_, err = os.ReadFile(temp_dir_path + "/config.yaml")
	if err != nil {
		LOGGER.Errorf("Downloaded config file is not readable: %v", err)
	} else {
		LOGGER.Infof("Config file downloaded successfully to: %s/config.yaml", temp_dir_path)
	}

	return temp_dir_path, nil
}
