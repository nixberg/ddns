package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	cloudflare "github.com/cloudflare/cloudflare-go"
)

const (
	sleepInterval time.Duration = 125 * time.Second
)

func getActualIP() (string, error) {
	response, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	ip, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	stringIP := strings.TrimSpace(string(ip))

	if net.ParseIP(stringIP) == nil {
		return "", errors.New("Invalid IP")
	}

	return stringIP, nil
}

type config struct {
	Email   string   `toml:"email"`
	APIKey  string   `toml:"apiKey"`
	ZoneID  string   `toml:"zoneID"`
	Records []string `toml:"records"`
}

func readConfig() (config, error) {
	config := config{}
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return config, err
	}
	metadata, err := toml.Decode(string(data), &config)
	if err != nil {
		return config, err
	} else if len(metadata.Undecoded()) != 0 {
		return config, errors.New("Invalid configuration")
	}
	return config, nil
}

func main() {
	config, err := readConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	api, err := cloudflare.New(config.APIKey, config.Email)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(os.Args) == 2 && os.Args[1] == "check" {
		os.Exit(0)
	}

	if config.Email == "email" {
		for {
			fmt.Println("Found dummy config. Exiting.")
			os.Exit(0)
		}
	}

	update := func() {
		actualIP, err := getActualIP()
		if err != nil {
			fmt.Println(err)
			return
		}

		dnsRecords, err := api.DNSRecords(config.ZoneID, cloudflare.DNSRecord{})
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, recordName := range config.Records {
			recordExists := false
			for _, dnsRecord := range dnsRecords {
				if recordName == dnsRecord.Name {
					if dnsRecord.Content != actualIP {
						dnsRecord.Content = actualIP
						err := api.UpdateDNSRecord(config.ZoneID, dnsRecord.ID, dnsRecord)
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println("Set", dnsRecord.Name, "to", dnsRecord.Content)
						}
					}
					recordExists = true
					break
				}
			}

			if !recordExists {
				if _, err := api.CreateDNSRecord(config.ZoneID, cloudflare.DNSRecord{
					Type:    "A",
					Name:    recordName,
					Content: actualIP,
					Proxied: true,
				}); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Created", recordName, "pointing to", actualIP)
				}
			}
		}
	}

	for {
		update()
		time.Sleep(sleepInterval)
	}
}
