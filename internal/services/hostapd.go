/*
	This file is for handling hostapd configuration
*/

package services

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Bryce285/PiAP/internal/utils"
)

var countryCodes = [249]string{"AD", "AE", "AF", "AG", "AI", "AL", "AM", "AO",
	"AQ", "AR", "AS", "AT", "AU", "AW", "AX", "AZ",
	"BA", "BB", "BD", "BE", "BF", "BG", "BH", "BI",
	"BJ", "BL", "BM", "BN", "BO", "BQ", "BR", "BS",
	"BT", "BV", "BW", "BY", "BZ", "CA", "CC", "CD",
	"CF", "CG", "CH", "CI", "CK", "CL", "CM", "CN",
	"CO", "CR", "CU", "CV", "CW", "CX", "CY", "CZ",
	"DE", "DJ", "DK", "DM", "DO", "DZ", "EC", "EE",
	"EG", "EH", "ER", "ES", "ET", "FI", "FJ", "FK",
	"FM", "FO", "FR", "GA", "GB", "GD", "GE", "GF",
	"GG", "GH", "GI", "GL", "GM", "GN", "GP", "GQ",
	"GR", "GS", "GT", "GU", "GW", "GY", "HK", "HM",
	"HN", "HR", "HT", "HU", "ID", "IE", "IL", "IM",
	"IN", "IO", "IQ", "IR", "IS", "IT", "JE", "JM",
	"JO", "JP", "KE", "KG", "KH", "KI", "KM", "KN",
	"KP", "KR", "KW", "KY", "KZ", "LA", "LB", "LC",
	"LI", "LK", "LR", "LS", "LT", "LU", "LV", "LY",
	"MA", "MC", "MD", "ME", "MF", "MG", "MH", "MK",
	"ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS",
	"MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA",
	"NC", "NE", "NF", "NG", "NI", "NL", "NO", "NP",
	"NR", "NU", "NZ", "OM", "PA", "PE", "PF", "PG",
	"PH", "PK", "PL", "PM", "PN", "PR", "PS", "PT",
	"PW", "PY", "QA", "RE", "RO", "RS", "RU", "RW",
	"SA", "SB", "SC", "SD", "SE", "SG", "SH", "SI",
	"SJ", "SK", "SL", "SM", "SN", "SO", "SR", "SS",
	"ST", "SV", "SX", "SY", "SZ", "TC", "TD", "TF",
	"TG", "TH", "TJ", "TK", "TL", "TM", "TN", "TO",
	"TR", "TT", "TV", "TW", "TZ", "UA", "UG", "UM",
	"US", "UY", "UZ", "VA", "VC", "VE", "VG", "VI",
	"VN", "VU", "WF", "WS", "YE", "YT", "ZA", "ZM",
	"ZW"}

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

func (c HostapdConf) WriteHostapdConf(ssid, hw_mode, passphrase, country string, ignore_broadcast_ssid bool, channel uint) error {
	if len(ssid) > 32 {
		return errors.New("SSID cannot be longer than 32 bytes")
	}
	if hw_mode != "g" && hw_mode != "a" {
		return errors.New("Unrecognized input for hw_mode. Acceptable options are 'g' and 'a'")
	}

	if channel != 0 {
		if channel >= 1 && channel <= 13 {
			if hw_mode != "g" {
				return errors.New("Invalid channel for selected value of hw_mode")
			}
		}
		if channel >= 36 && channel <= 165 {
			if hw_mode != "a" {
				return errors.New("Invalid channel for selected value of hw_mode")
			}
		}
	}
	channelStr := strconv.Itoa(int(channel))

	if len(passphrase) < 8 || len(passphrase) > 63 {
		return errors.New("Passphrase must be between 8 and 63 characters in length")
	}

	countryFound := false
	for _, value := range countryCodes {
		if country == value {
			countryFound = true
			break
		}
	}

	if countryFound == false {
		return errors.New("Not a valid country code")
	}

	var ignoreBroadcastStr string
	if ignore_broadcast_ssid {
		ignoreBroadcastStr = "1"
	} else {
		ignoreBroadcastStr = "0"
	}

	confDefault := HostapdConf{
		Interface:             "wlan0",
		Driver:                "nl80211",
		Ssid:                  ssid,
		Hw_mode:               hw_mode,
		Channel:               channelStr,
		Auth_algs:             "1",
		Wpa:                   "2",
		Wpa_passphrase:        passphrase,
		Wpa_key_mgmt:          "WPA-PSK",
		Wpa_pairwise:          "CCMP",
		Rsn_pairwise:          "CCMP",
		Country:               country,
		Macaddr_acl:           "0",
		Ignore_broadcast_ssid: ignoreBroadcastStr,
		Wpa_group_rekey:       "86400",
		Beacon_int:            "100",
		Dtim_period:           "2",
	}

	info, err := os.Stat("/etc/hostapd/hostapd.conf")
	if errors.Is(err, os.ErrNotExist) {
		return os.WriteFile("/etc/hostapd/hostapd.conf", []byte(confDefault.String()), 0644)
	}
	if err != nil {
		fmt.Printf("Error checking /etc/hostapd/hostapd.conf file existence: %v\n", err)
		return err
	}
	if !info.IsDir() {
		removeErr := utils.RemoveFile("/etc/hostapd/hostapd.conf.backup")
		if removeErr != nil {
			return removeErr
		}

		renameErr := os.Rename("/etc/hostapd/hostapd.conf", "/etc/hostapd/hostapd.conf.backup")
		if renameErr != nil {
			return renameErr
		}

		return os.WriteFile("/etc/hostapd/hostapd.conf", []byte(confDefault.String()), 0644)
	}

	return errors.New("Error: could not create /etc/hostapd/hostapd.conf")
}
