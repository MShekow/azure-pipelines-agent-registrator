package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const version = "1.0"

var organizationUrl string
var poolName string
var pat string
var agentNamePrefix string
var capabilities string

// TODO consider implementing auto-retry logic?

func init() {
	flag.StringVar(&organizationUrl, "organization-url", "", "https-URL of your Azure DevOps organization")
	flag.StringVar(&poolName, "pool-name", "", "name of your Azure DevOps pool")
	flag.StringVar(&pat, "pat", "", "Azure DevOps Personal Access Token with 'Agent Pools Read&Manage' permission")
	flag.StringVar(&agentNamePrefix, "agent-name-prefix", "", "prefix for the name of the agent, the tool appends a short random string")
	flag.StringVar(&capabilities, "capabilities", "", "one or more key-value pairs (notation: key=value), separated by semicolon")
}

func main() {
	flag.CommandLine.SetOutput(os.Stdout) // ensure that flag.PrintDefaults() does NOT print to stderr by default
	flag.Usage = func() {
		fmt.Printf("Usage of Azure Pipelines agent registration tool %s:\n\n", version)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()

	if !strings.HasPrefix(organizationUrl, "https://") {
		log.Fatal("You must provide a valid organization URL")
	}
	organizationUrl = strings.TrimSuffix(organizationUrl, "/")
	if len(poolName) == 0 {
		log.Fatal("You must provide a -pool-name")
	}
	if len(pat) == 0 {
		log.Fatal("You must provide a valid PAT")
	}
	if len(agentNamePrefix) == 0 {
		log.Fatal("You must provide an -agent-name-prefix")
	}
	if len(capabilities) == 0 {
		log.Fatal("You must provide capabilities")
	}

	capabilitiesMap := GetCapabilitiesMapFromString(capabilities)

	timeout := 2 * time.Second
	httpClient := &http.Client{
		Timeout: timeout,
	}

	poolId, err := getPoolIdFromName(pat, organizationUrl, poolName, httpClient)
	if err != nil {
		fmt.Printf("Unable to retrieve pool ID from name: %v\n", err)
		os.Exit(1)
	}

	fakeAgentName, err := registerFakeAgent(pat, organizationUrl, agentNamePrefix, capabilitiesMap, poolId, httpClient)
	if err != nil {
		fmt.Printf("Unable to create fake agent: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully registered agent with name %s\n", fakeAgentName)
}
