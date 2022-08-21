package linode

import (
	"encoding/json"
	"fmt"
)

type getDomainsResponse struct {
	Data []struct {
		Domain string `json:"domain"`
		ID     int    `json:"id"`
	} `json:"data"`
}

func (c *client) GetDomainID(domain string) (int, error) {
	body, err := c.doGetRequest("/v4/domains")
	if err != nil {
		return 0, fmt.Errorf("failed to get domain list from linode: %w", err)
	}

	var respData getDomainsResponse
	if err := json.Unmarshal(body, &respData); err != nil {
		return 0, fmt.Errorf("failed to unmarshal domains list from linode: %w", err)
	}

	for _, d := range respData.Data {
		if d.Domain == domain {
			return d.ID, nil
		}
	}

	return 0, ErrDomainNotFound
}

type getRecordsResponse struct {
	Data []RetrievedARecord `json:"data"`
}

type RetrievedARecord struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Target string `json:"target"`
	TTLSec int    `json:"ttl_sec"`
	Type   string `json:"type"`
}

func (c *client) GetARecords(domainID int, names []string) ([]RetrievedARecord, error) {
	body, err := c.doGetRequest(fmt.Sprintf("/v4/domains/%d/records", domainID))
	if err != nil {
		return nil, fmt.Errorf("failed to get record list from linode: %w", err)
	}

	var respData getRecordsResponse
	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal record list from linode: %w", err)
	}

	namesMap := make(map[string]struct{}, len(names))
	for _, name := range names {
		namesMap[name] = struct{}{}
	}

	records := make([]RetrievedARecord, 0)
	for _, r := range respData.Data {
		if r.Type != "A" {
			continue
		}

		if _, ok := namesMap[r.Name]; ok {
			records = append(records, r)
		}
	}

	return records, nil
}

type ARecordUpdate struct {
	Target string `json:"target"`
}

func (c *client) UpdateARecord(domainID int, recordID int, update ARecordUpdate) error {
	reqBody, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	_, err = c.doPutRequest(fmt.Sprintf("/v4/domains/%d/records/%d", domainID, recordID), reqBody)
	if err != nil {
		return fmt.Errorf("failed to update linode dns record: %w", err)
	}

	return nil
}
