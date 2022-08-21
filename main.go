package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/jackharley/lindyn/linode"
	"github.com/jackharley/lindyn/myip"
)

func main() {
	godotenv.Load()

	httpClient := http.Client{Timeout: 5 * time.Second}
	ipRetriever := myip.NewClient(httpClient)
	linodeClient := linode.NewClient(httpClient, os.Getenv("LINODE_PERSONAL_ACCESS_TOKEN"))

	currentIPv4, err := ipRetriever.GetIPv4()
	if err != nil {
		log.Fatalf("Failed to determine current IPv4 address, error: %v, exiting", err)
	}
	log.Printf("Current IPv4 address detected: %s", currentIPv4)

	domainID, err := linodeClient.GetDomainID(os.Getenv("LINODE_DOMAIN"))
	if err != nil {
		log.Fatalf("Failed to get domain zone ID from Linode, error: %v, exiting", err)
	}

	records, err := linodeClient.GetARecords(domainID, strings.Split(os.Getenv("LINODE_A_RECORDS"), ","))
	if err != nil {
		log.Fatalf("Failed to get domain record IDs from Linode, error: %v, exiting", err)
	}

	for _, record := range records {
		if record.Target == currentIPv4 {
			log.Printf("A record [%s] found with correct IPv4 address [%s], skipping", record.Name, record.Target)
			continue
		}

		log.Printf("A record [%s] found with incorrect IPv4 address [%s], updating...", record.Name, record.Target)
		err := linodeClient.UpdateARecord(domainID, record.ID, linode.ARecordUpdate{Target: currentIPv4})
		if err != nil {
			log.Printf("Failed to update record [%s] with IPv4 address [%s], error: %v", record.Name, currentIPv4, err)
		} else {
			log.Printf("Successfully updated record [%s] with IPv4 address [%s]", record.Name, currentIPv4)
		}
	}
}
