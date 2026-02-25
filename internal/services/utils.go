/* This file is for utility functions relating to service configurations */

package services

import (
	"errors"
	"fmt"
	"os"
)

func RemoveFile(filepath string) error {
	info, err := os.Stat(filepath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		fmt.Printf("Error checking %s file existence: %v\n", filepath, err)
		return err
	}
	if !info.IsDir() {
		removeErr := os.Remove(filepath)
		return removeErr
	}

	return errors.New("Error: could not delete file")
}

/* To clear custom configuration files on config reset */
func ClearConf(filepath string) error {
	removeErr := RemoveFile(filepath)
	if removeErr != nil {
		return removeErr
	}

	info, err := os.Stat(filepath + ".backup")
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		fmt.Printf("Error checking %s.backup file existence: %v\n", filepath, err)
		return err
	}
	if !info.IsDir() {
		renameErr := os.Rename(filepath+".backup", filepath)
		if renameErr != nil {
			return renameErr
		}
	}

	return nil
}
