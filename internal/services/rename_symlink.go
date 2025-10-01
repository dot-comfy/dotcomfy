package services

import (
	"os"
	"path/filepath"

	Log "dotcomfy/internal/logger"
)

// TODO:
//
// write documentation
func RenameSymlinkUnix(old_dotfiles_dir, dotcomfy_dir, center_path string) (string, error) {
	LOGGER = Log.GetLogger()
	new_path := dotcomfy_dir + center_path
	old_path := old_dotfiles_dir + center_path
	// Want to check to see if new_entry has a corresponding entry
	// in old_dotfiles_path. If so, rename corresponding entry to
	// {corresponding_entry}.pre-dotcomfy, put new_entry symlink in its place.
	_, err := os.Stat(old_path)
	if err == nil {
		old_path_renamed := old_path + ".pre-dotcomfy"
		err = os.Rename(old_path, old_path_renamed)
		if err != nil {
			LOGGER.Error(err)
			return "", err
		}
		err = os.Symlink(new_path, old_path)
		if err != nil {
			LOGGER.Error(err)
			return "", err
		}
		return old_path, nil
	} else if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(old_path), 0755)
		if err != nil {
			LOGGER.Error(err)
			return "", err
		}
		err = os.Symlink(new_path, old_path)
		if err != nil {
			LOGGER.Error(err)
			return "", err
		}
		return old_path, nil
	} else {
		LOGGER.Error(err)
		return "", err
	}
}
