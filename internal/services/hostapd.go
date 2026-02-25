/*
	This file is for handling hostapd configuration
*/

package services

import (
	"fmt"
	"os"
)

// TODO - let users edit the hostapd config
type HostapdConf struct {
	Interface             string
	Driver                string
	Ssid                  string
	Hw_mode               string
	Channel               string
	Auth_algs             string
	Wpa                   string
	Wpa_passphrase        string
	Wpa_key_mgmt          string
	Wpa_pairwise          string
	Rsn_pairwise          string
	Country               string
	Macaddr_acl           string
	Ignore_broadcast_ssid string
	Wpa_group_rekey       string
	Beacon_int            string
	Dtim_period           string
}

func (c HostapdConf) String() string {
	return fmt.Sprintf(`interface=%s
	driver=%s
	ssid=%s
	hw_mode=%s
	channel=%s
	auth_algs=%s
	wpa=%s
	wpa_passphrase=%s
	wpa_key_mgmt=%s
	wpa_pairwise=%s
	rsn_pairwise=%s
	country=%s
	macaddr_acl=%s
	ignore_broadcast_ssid=%s
	wpa_group_rekey=%s
	beacon_int=%s
	dtim_period=%s
	`,
		c.Interface,
		c.Driver,
		c.Ssid,
		c.Hw_mode,
		c.Channel,
		c.Auth_algs,
		c.Wpa,
		c.Wpa_passphrase,
		c.Wpa_key_mgmt,
		c.Wpa_pairwise,
		c.Rsn_pairwise,
		c.Country,
		c.Macaddr_acl,
		c.Ignore_broadcast_ssid,
		c.Wpa_group_rekey,
		c.Beacon_int,
		c.Dtim_period,
	)
}

func (c HostapdConf) WriteConf(ssid string, passphrase string, country string) {
	confDefault := HostapdConf{
		Interface:             "wlan0",
		Driver:                "nl80211",
		Ssid:                  ssid,
		Hw_mode:               "g",
		Channel:               "7",
		Auth_algs:             "1",
		Wpa:                   "2",
		Wpa_passphrase:        passphrase,
		Wpa_key_mgmt:          "WPA-PSK",
		Wpa_pairwise:          "CCMP",
		Rsn_pairwise:          "CCMP",
		Country:               country,
		Macaddr_acl:           "0",
		Ignore_broadcast_ssid: "0",
		Wpa_group_rekey:       "86400",
		Beacon_int:            "100",
		Dtim_period:           "2",
	}

	os.WriteFile("/etc/hostapd/hostapd.conf", []byte(confDefault.String()), 0644)
}
