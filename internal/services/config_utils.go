package services

import (
	"fmt"
	"net/http"

	Config "dotcomfy/internal/config"
	Log "dotcomfy/internal/logger"
)

func DownloadConfigFile(repo_url, branch string) (string, error) {
	LOGGER = Log.GetLogger()

	raw_url = fmt.Sprintf("%s/-/raw/%s/dotcomfy/config.yaml", repo_url, branch)
	res, err := http.Get(raw_url)
	if err != nil {
		LOGGER.Infof("Error getting config file from repo: %v", err)
	}

	defer res.Body.Close()
}
