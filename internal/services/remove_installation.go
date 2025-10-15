package services

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	Log "dotcomfy/internal/logger"
)

func RemoveInstallation(dotcomfy_dir, old_dotfiles_dir string) (err error) {
	LOGGER = Log.GetLogger()
	err = filepath.WalkDir(dotcomfy_dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			LOGGER.Error(err)
			return err
		}
		if !d.IsDir() {
			if !strings.Contains(path, ".git") && !strings.Contains(path, dotcomfy_dir+"README.md") {
				center_path := strings.TrimPrefix(path, dotcomfy_dir)
				old_path := old_dotfiles_dir + center_path
				_, err := os.Stat(old_path + ".pre-dotcomfy")
				if err == nil {
					// Remove symlink
					err = os.Remove(old_path)
					if err != nil {
						LOGGER.Warn(err)
					}
					err = os.Rename(old_path+".pre-dotcomfy", old_path)
					if err != nil {
						LOGGER.Warn(err)
					}
				} else { // Just remove symlink
					err = os.Remove(old_path)
					if err != nil {
						LOGGER.Warn(err)
					}
				}
			}
		}
		return nil
	})

	err = filepath.WalkDir(dotcomfy_dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			center_path := strings.TrimPrefix(path, dotcomfy_dir)
			old_path := old_dotfiles_dir + center_path

			not_empty, err := directoryContainsFiles(old_path)
			if err != nil {
				LOGGER.Warn(err)
			} else {
				if !not_empty {
					err = os.RemoveAll(old_path)
					if err != nil {
						LOGGER.Info(err)
					}
				}
			}

		}
		return nil
	})

	if err != nil {
		LOGGER.Error(err)
		return err
	}

	config_path := filepath.Join(dotcomfy_dir, "config.toml")

	// Delete everything in ~/.dotcomfy
	dir, err := os.Open(dotcomfy_dir)
	if err != nil {
		LOGGER.Fatal(err)
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		LOGGER.Fatal(err)
		os.Exit(1)
	}

	for _, name := range names {
		file_path := filepath.Join(dotcomfy_dir, name)
		if file_path == config_path {
			continue
		}

		file_info, err := os.Stat(file_path)
		if err != nil {
			LOGGER.Fatal(err)
			os.Exit(1)
		}

		if file_info.IsDir() {
			err = os.RemoveAll(file_path)
			if err != nil {
				LOGGER.Fatal(err)
				continue
			}
		} else {
			err = os.Remove(file_path)
			if err != nil {
				LOGGER.Fatal(err)
				continue
			}
		}
	}
	return nil
}

func directoryContainsFiles(dir_path string) (bool, error) {
	entries, err := os.ReadDir(dir_path)
	if err != nil {
		return false, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			return true, nil // Found a file
		}
	}
	return false, nil // No files found
}
