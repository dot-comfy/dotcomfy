package services

import (
	"os"
	"path/filepath"
	"strings"
)

// TODO: write documentation
func RenameSymlinkUnix(old_dotfiles_dir, dotcomfy_dir, new_path string) (string, error) {
	// center_path represents the path of the directory entry
	// with the dotcomfy_path prefix removed.
	center_path := strings.TrimPrefix(new_path, dotcomfy_dir)
	old_path := old_dotfiles_dir + center_path
	// Want to check to see if new_entry has a corresponding entry
	// in old_dotfiles_path. If so, rename corresponding entry to
	// {corresponding_entry}.pre-dotcomfy, put new_entry symlink in its place.
	_, err := os.Stat(old_path)
	if err == nil {
		new_path = new_path + ".pre-dotcomfy"
		err = os.Rename(old_path, new_path)
		if err != nil {
			return "", err
		}
		err = os.Symlink(new_path, old_path)
		if err != nil {
			return "", err
		}
		return old_path, nil
	} else if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(old_path), 0755)
		if err != nil {
			return "", err
		}
		err = os.Symlink(new_path, old_path)
		if err != nil {
			return "", err
		}
		return old_path, nil
	} else {
		return "", err
	}
}
