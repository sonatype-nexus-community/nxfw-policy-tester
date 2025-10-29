/**
 * Copyright (c) 2019-present Sonatype, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nxiq

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	nxiq "github.com/sonatype-nexus-community/nexus-iq-api-client-go"
	"github.com/sonatype-nexus-community/nxfw-policy-tester/cli"
	"github.com/sonatype-nexus-community/nxfw-policy-tester/util"
)

type NxiqConnection struct {
	apiClient *nxiq.APIClient
	iqBaseUrl string
	ctx       *context.Context
}

func (c *NxiqConnection) RetrieveFWQuarantineStatus(componentName, componentVersion, repositoryName, expectedPolicy, format, repoBaseUrl string) (bool, bool, error) {
	// Workaround for Maven
	if format == "maven2" {
		componentNameParts := strings.Split(componentName, "/")
		componentName = componentNameParts[1]
	}

	fwResult, apiResponse, err := c.apiClient.FirewallAPI.GetQuarantineList(*c.ctx).ComponentName(componentName).RepositoryPublicId((repositoryName)).Execute()
	if err != nil || apiResponse.StatusCode != http.StatusOK {
		cli.PrintCliln("Error: Failed to query Sonatype Repository Firewall Quarantine List. Please check your credentials and URL.", util.ColorRed)
		cli.PrintCliln(fmt.Sprintf("Details: %v", err), util.ColorRed)
		return false, false, err
	}

	var expectedPolicyTriggered = false
	var quarantined = false

	for _, r := range fwResult.Results {
		var coordinates = r.ComponentIdentifier.GetCoordinates()

		var packageName, packageName2 string
		format := r.ComponentIdentifier.Format
		switch *format {
		case "cargo", "conda", "golang":
			packageName = coordinates["name"]
		case "docker":
			// repo.phorton.eu.ngrok.io-dockerhub-proxy-sonatypecommunity-docker-policy-demo-Security-High
			packageName = strings.ReplaceAll(
				fmt.Sprintf("%s-%s-%s-%s", repoBaseUrl, repositoryName, componentName, componentVersion),
				".", "-",
			)
		case "hf-model":
			packageName = coordinates["repo_id"]
		case "maven":
			packageName = coordinates["artifactId"]
		case "pypi":
			packageName = coordinates["name"]
			packageName2 = strings.ReplaceAll(coordinates["name"], "-", "_")
		default:
			packageName = coordinates["packageId"]
		}

		version, ok := coordinates["version"]
		if !ok {
			continue
		}

		if (packageName == componentName || packageName2 == componentName) && version == componentVersion {
			if r.Quarantined != nil && *r.Quarantined {
				quarantined = true
				if *r.PolicyName == expectedPolicy {
					expectedPolicyTriggered = true
					break
				}
			}
		}
	}

	return quarantined, expectedPolicyTriggered, nil
}

// Validate credentials by making a test call
func (c *NxiqConnection) validateConnection() error {
	_, apiResponse, err := c.apiClient.UserTokensAPI.GetUserTokenExistsForCurrentUser(*c.ctx).Execute()
	if err != nil {
		cli.PrintCliln(fmt.Sprintf("Error: Failed to authenticate with Sonatype IQ Server (%s). Please check your credentials.", c.iqBaseUrl), util.ColorRed)
		cli.PrintCliln(fmt.Sprintf("Details: %v", err), util.ColorRed)
		return err
	}

	if apiResponse.StatusCode != http.StatusOK {
		cli.PrintCliln(fmt.Sprintf("Error: Sonatype IQ Server is not reporting as available - status code: %d", apiResponse.StatusCode), util.ColorRed)
		return fmt.Errorf("error: Sonatype IQ Server is not reporting as available - status code: %d", apiResponse.StatusCode)
	}

	return nil
}

func NewNxiqConnection(nxiqUrl, username, password string) (*NxiqConnection, error) {
	// Create API client configuration
	configuration := nxiq.NewConfiguration()
	configuration.Servers = nxiq.ServerConfigurations{
		{
			URL: strings.TrimSuffix(nxiqUrl, "/"),
		},
	}

	// Create API client
	apiClient := nxiq.NewAPIClient(configuration)

	// Create context with basic auth
	ctx := context.WithValue(context.Background(), nxiq.ContextBasicAuth, nxiq.BasicAuth{
		UserName: username,
		Password: password,
	})

	connection := &NxiqConnection{
		apiClient: apiClient,
		ctx:       &ctx,
		iqBaseUrl: nxiqUrl,
	}

	err := connection.validateConnection()
	if err != nil {
		return nil, err
	}

	return connection, nil
}
