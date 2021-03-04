package main // import "github.com/nixberg/ddns"

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/pelletier/go-toml"
)

type config struct {
	APIToken    string   `toml:"apiToken"`
	ZoneID      string   `toml:"zoneID"`
	RecordNames []string `toml:"recordNames"`
}

func main() {
	config, err := readConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	api, err := cloudflare.NewWithAPIToken(config.APIToken)
	if err != nil {
		fmt.Println("Error creating Cloudflare API client:", err)
		os.Exit(1)
	}

	if len(os.Args) == 2 && os.Args[1] == "validate-config" {
		_, err := api.DNSRecords(context.TODO(), config.ZoneID, cloudflare.DNSRecord{})

		if err != nil {
			fmt.Println("Could not validate config:", err)
			os.Exit(1)
		}

		return
	}

	ipAddress, err := getIPv4Address()
	if err != nil {
		fmt.Println("Error looking up IP address:", err)
		os.Exit(1)
	}

	for _, name := range config.RecordNames {
		err = updateRecord(api, config.ZoneID, name, ipAddress)
		if err != nil {
			fmt.Printf("Error updating A record \"%s\": %s\n", name, err)
		}
	}
}

func readConfig() (*config, error) {
	decoder := toml.NewDecoder(os.Stdin)
	decoder.Strict(true)

	config := config{}
	err := decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func updateRecord(api *cloudflare.API, zoneID, name, ipAddress string) error {
	filter := cloudflare.DNSRecord{Name: name, Type: "A"}

	records, err := api.DNSRecords(context.TODO(), zoneID, filter)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		_, err := api.CreateDNSRecord(context.TODO(), zoneID, cloudflare.DNSRecord{
			Type:    "A",
			Name:    name,
			Content: ipAddress,
		})

		if err != nil {
			return err
		}

		fmt.Printf("Created A record \"%s\" pointing to %s\n", name, ipAddress)

	} else {
		if records[0].Content == ipAddress {
			return nil
		}

		records[0].Content = ipAddress
		err := api.UpdateDNSRecord(context.TODO(), zoneID, records[0].ID, records[0])

		if err != nil {
			return err
		}

		fmt.Printf("Set A record \"%s\" to %s\n", name, ipAddress)
	}

	return nil
}

func getIPv4Address() (string, error) {
	response, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(data))

	parsed := net.ParseIP(ip)

	if parsed == nil {
		return "", errors.New("not an IP address")
	}

	if parsed.To4() == nil {
		return "", errors.New("not an IPv4 address")
	}

	return ip, nil
}
