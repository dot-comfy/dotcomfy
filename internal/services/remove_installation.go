package services

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func RemoveInstallation(dotcomfy_dir, old_dotfiles_dir string) (err error) {
	err = filepath.WalkDir(dotcomfy_dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if strings.Contains(path, ".git") {
				fmt.Println("Skipping .git directory")
			} else if strings.Contains(path, dotcomfy_dir+"README.md") {
				fmt.Println("Skipping root level README.md")
			} else {
				center_path := strings.TrimPrefix(path, dotcomfy_dir)
				old_path := old_dotfiles_dir + center_path
				if strings.Contains(old_path, ".pre-dotcomfy") {
					old_name := strings.Replace(old_path, ".pre-dotcomfy", "", 1)
					// Remove symlink
					err = os.Remove(old_path)
					if err != nil {
						return err
					}
					err = os.Rename(old_path, old_name)
					if err != nil {
						return err
					}
				} else { // Just remove symlink
					err = os.Remove(old_path)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	config_path := filepath.Join(dotcomfy_dir, "config.toml")

	// Delete everything in ~/.dotcomfy
	dir, err := os.Open(dotcomfy_dir)
	if err != nil {
		os.Exit(1)
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		os.Exit(1)
	}

	for _, name := range names {
		file_path := filepath.Join(dotcomfy_dir, name)
		if file_path == config_path {
			continue
		}

		file_info, err := os.Stat(file_path)
		if err != nil {
			os.Exit(1)
		}

		if file_info.IsDir() {
			err = os.RemoveAll(file_path)
			if err != nil {
				log.Print(err)
				continue
			}
		} else {
			err = os.Remove(file_path)
			if err != nil {
				log.Print(err)
				continue
			}
		}
	}
	return nil
}
