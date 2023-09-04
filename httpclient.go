package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AzurePipelinesApiPoolNameResponse struct {
	// Because there could be multiple Azure Pipeline pools with the same name, the API returns an array. The objects
	// have many more fields, but we only care about the ID, and omit defining the other fields.
	Value []struct {
		ID int `json:"id"`
	} `json:"value"`
}

func getPoolIdFromName(pat, organizationUrl, poolName string, httpClient *http.Client) (int64, error) {
	url := fmt.Sprintf("%s/_apis/distributedtask/pools?poolName=%s", organizationUrl, poolName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.SetBasicAuth("", pat)
	if err != nil {
		return 0, err
	}

	response, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	bytesResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	if !(response.StatusCode >= 200 && response.StatusCode <= 299) {
		return 0, fmt.Errorf("Azure DevOps REST API returned error. url: %s status: %d response: %s", url, response.StatusCode, string(bytesResponse))
	}

	var result AzurePipelinesApiPoolNameResponse
	err = json.Unmarshal(bytesResponse, &result)
	if err != nil {
		return 0, err
	}

	count := len(result.Value)
	if count == 0 {
		return 0, fmt.Errorf("agent pool with name `%s` not found in response", poolName)
	}

	if count != 1 {
		return 0, fmt.Errorf("found %d agent pools with name `%s`", count, poolName)
	}

	poolId := int64(result.Value[0].ID)

	return poolId, nil
}

func registerFakeAgent(pat, organizationUrl, agentNamePrefix string, capabilities *map[string]string, poolId int64, httpClient *http.Client) (string, error) {
	url := fmt.Sprintf("%s/_apis/distributedtask/pools/%d/agents?api-version=7.0", organizationUrl, poolId)

	for {
		fakeAgentName := fmt.Sprintf("%s-%s", agentNamePrefix, randomString(8))

		requestBodyTemplate := `{
			"name": "%s",
			"version": "99.999.9",
			"osDescription": "Linux 5.15.49-linuxkit-pr #1 SMP PREEMPT Thu May 25 07:27:39 UTC 2023",
			"enabled": true,
			"status": "offline",
			"provisioningState": "Provisioned",
			"systemCapabilities": %s
		}`

		capabilitiesJsonStr, _ := json.Marshal(*capabilities)

		requestBody := fmt.Sprintf(requestBodyTemplate, fakeAgentName, capabilitiesJsonStr)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
		if err != nil {
			return "", err
		}

		req.SetBasicAuth("", pat)
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/json")

		response, err := httpClient.Do(req)
		if err != nil {
			return "", err
		}
		defer response.Body.Close()

		if response.StatusCode == 409 {
			fmt.Printf("Agent with name '%s' already exists, retrying with different name ...\n", fakeAgentName)
			continue // 409 = HTTP "conflict"
		}

		if response.StatusCode != 200 {
			return "", fmt.Errorf("Azure DevOps REST API returned error. url: %s status: %d", url, response.StatusCode)
		}

		return fakeAgentName, nil
	}
}
