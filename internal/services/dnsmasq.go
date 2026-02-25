/*
	This file is for handling dnsmasq configuration
*/

package services

import (
	"errors"
	"fmt"
	"os"
)

/* TODO - this dnsmasq configuration may need to be generalized a little */
func WriteDnsmasqConf() error {
	confData := "interface=wlan0\ndhcp-range=192.168.100.2,192.168.100.10,12h"

	info, err := os.Stat("/etc/dnsmasq.conf")
	if errors.Is(err, os.ErrNotExist) {
		return os.WriteFile("/etc/dnsmasq.conf", []byte(confData), 0644)
	}
	if err != nil {
		fmt.Printf("Error checking /etc/dnsmasq.conf file existence: %v\n", err)
		return err
	}
	if !info.IsDir() {
		removeErr := RemoveFile("/etc/dnsmasq.conf.backup")
		if removeErr != nil {
			return removeErr
		}

		renameErr := os.Rename("/etc/dnsmasq.conf", "/etc/dnsmasq.conf.backup")
		if renameErr != nil {
			return renameErr
		}

		return os.WriteFile("/etc/dnsmasq.conf", []byte(confData), 0644)
	}

	return errors.New("Error: could not create /etc/dnsmasq.conf")
}
