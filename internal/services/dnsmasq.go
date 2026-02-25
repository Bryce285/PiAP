/*
	This file is for handling dnsmasq configuration
*/

package services

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/Bryce285/PiAP/internal/utils"
)

/* Right now we are assuming a /24 network */
func ValidateDHCPRange(baseIPStr, startStr, endStr string) error {
	baseIP := net.ParseIP(baseIPStr).To4()
	startIP := net.ParseIP(startStr).To4()
	endIP := net.ParseIP(endStr).To4()

	if baseIP == nil || startIP == nil || endIP == nil {
		return errors.New("all addresses must be valid IPv4")
	}

	_, subnet, err := net.ParseCIDR(baseIPStr + "/24")
	if err != nil {
		return err
	}

	if !subnet.Contains(startIP) || !subnet.Contains(endIP) {
		return errors.New("DHCP range must be in same /24 subnet as base IP")
	}

	baseInt := utils.IpToUint32(baseIP)
	startInt := utils.IpToUint32(startIP)
	endInt := utils.IpToUint32(endIP)

	if startInt > endInt {
		return errors.New("range start must be less than or equal to range end")
	}

	if baseInt >= startInt && baseInt <= endInt {
		return errors.New("base IP cannot be inside DHCP range")
	}

	networkInt := utils.IpToUint32(subnet.IP)
	broadcastInt := networkInt | ^binary.BigEndian.Uint32(subnet.Mask)

	if startInt <= networkInt || endInt >= broadcastInt {
		return errors.New("range cannot include network or broadcast address")
	}

	return nil
}

func WriteDnsmasqConf(ipBaseAddr, dhcpStart, dhcpEnd string, leaseTime uint) error {
	leaseTimeStr := strconv.Itoa(int(leaseTime))

	if !utils.IsValidIPv4(dhcpStart) {
		return fmt.Errorf("Error: %s is not a valid IPv4 address.", dhcpStart)
	}
	if !utils.IsValidIPv4(dhcpEnd) {
		return fmt.Errorf("Error: %s is not a valid IPv4 address.", dhcpEnd)
	}

	dhcpRangeErr := ValidateDHCPRange(ipBaseAddr, dhcpStart, dhcpEnd)
	if dhcpRangeErr != nil {
		return dhcpRangeErr
	}

	confData := "interface=wlan0\ndhcp-range=" + dhcpStart + "," + dhcpEnd + "," + leaseTimeStr + "h"

	info, err := os.Stat("/etc/dnsmasq.conf")
	if errors.Is(err, os.ErrNotExist) {
		return os.WriteFile("/etc/dnsmasq.conf", []byte(confData), 0644)
	}
	if err != nil {
		fmt.Printf("Error checking /etc/dnsmasq.conf file existence: %v\n", err)
		return err
	}
	if !info.IsDir() {
		removeErr := utils.RemoveFile("/etc/dnsmasq.conf.backup")
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
