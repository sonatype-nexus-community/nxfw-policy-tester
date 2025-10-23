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

package nxrm

import (
	"context"
	"fmt"
	"net/http"
	"os"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	"github.com/sonatype-nexus-community/nxfw-policy-tester/cli"
	"github.com/sonatype-nexus-community/nxfw-policy-tester/formats"
	"github.com/sonatype-nexus-community/nxfw-policy-tester/nxiq"
	"github.com/sonatype-nexus-community/nxfw-policy-tester/util"
)

type NxrmConnection struct {
	apiClient *v3.APIClient
	baseUrl   string
	ctx       *context.Context
	username  string
	password  string
}

// Validate credentials by making a test call
func (c *NxrmConnection) validateConnection() error {
	apiResponse, err := c.apiClient.StatusAPI.IsAvailable(*c.ctx).Execute()
	if err != nil {
		cli.PrintCliln("Error: Failed to authenticate with Nexus Repository. Please check your credentials and URL.", util.ColorRed)
		cli.PrintCliln(fmt.Sprintf("Details: %v", err), util.ColorRed)
		return err
	}

	if apiResponse.StatusCode != http.StatusOK {
		cli.PrintCliln(fmt.Sprintf("Error: Nexus Repository is not reporting as available - status code: %d", apiResponse.StatusCode), util.ColorRed)
		return fmt.Errorf("error: Nexus Repository is not reporting as available - status code: %d", apiResponse.StatusCode)
	}

	return nil
}

func (c *NxrmConnection) CheckPackages(repoName string, format formats.PackageFormat, nxiqConnection *nxiq.NxiqConnection) ([]formats.CheckResult, error) {
	packages := format.GetPackages()
	results := make([]formats.CheckResult, 0, len(packages))

	cli.PrintCliln("\n=== Checking Package Availability ===\n", util.ColorYellow)

	for _, pkg := range packages {
		color := pkg.PolicyName.GetSecurityColor()
		cli.PrintCliln(
			fmt.Sprintf("Checking %s %s[%s]%s...\n",
				format.FormatPackageName(pkg),
				color,
				pkg.PolicyName, util.ColorReset),
			util.ColorReset,
		)

		url := format.ConstructURL(c.baseUrl, repoName, pkg)
		httpCode, err := c.DownloadPackageAtUrl(url)

		result := formats.CheckResult{
			Package:                       pkg,
			HTTPCode:                      httpCode,
			Available:                     false,
			Failed:                        false,
			Quarantined:                   false,
			QuarantinedWithExpectedPolicy: false,
		}

		if err != nil {
			cli.PrintCliln(
				fmt.Sprintf(
					"✗ Error attempting download: %s [%s] (Error: %v)\n",
					format.FormatPackageName(pkg),
					pkg.PolicyName,
					err,
				),
				util.ColorRed,
			)
			result.Failed = true
		} else if httpCode == http.StatusOK {
			cli.PrintCliln(
				fmt.Sprintf(
					"✓ Package available: %s [%s]\n\n",
					format.FormatPackageName(pkg),
					pkg.PolicyName,
				),
				util.ColorGreen,
			)
			result.Available = true
		} else if httpCode == http.StatusForbidden {
			if nxiqConnection == nil {
				println("ERROR: NO Connection to IQ")
				os.Exit(1)
			}
			quarantined, policyTriggered, fwErr := nxiqConnection.RetrieveFWQuarantineStatus(
				pkg.Name, pkg.Version, repoName, string(pkg.PolicyName),
			)
			if fwErr != nil {
				cli.PrintCliln(fmt.Sprintf("Error checking Firewall Quarantine Status: %v", fwErr), util.ColorRed)
			}
			result.Quarantined = quarantined
			result.QuarantinedWithExpectedPolicy = policyTriggered

			cli.PrintCliln(
				fmt.Sprintf(
					"✗ Package Quarrantined and NOT available: %s [%s]\n\n",
					format.FormatPackageName(pkg),
					pkg.PolicyName,
				),
				util.ColorRed,
			)
		} else {
			cli.PrintCliln(
				fmt.Sprintf(
					"✗ Package NOT available: %s [%s] (response code %d)\n\n",
					format.FormatPackageName(pkg),
					pkg.PolicyName,
					httpCode,
				),
				util.ColorRed,
			)
			result.Available = false
		}

		results = append(results, result)
	}

	return results, nil
}

func (c *NxrmConnection) DownloadPackageAtUrl(url string) (int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.SetBasicAuth(c.username, c.password)

	resp, err := c.apiClient.GetConfig().HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log the error if needed, but don't override the main return value
			fmt.Fprintf(os.Stderr, "Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	return resp.StatusCode, nil
}

func (c *NxrmConnection) GetConnectedIqServer() (string, error) {
	iqConnection, apiResponse, err := c.apiClient.ManageSonatypeRepositoryFirewallConfigurationAPI.GetConfiguration(*c.ctx).Execute()

	if err != nil || apiResponse.StatusCode != http.StatusOK {
		cli.PrintCliln("Error: Nexus Repository is not connected to Repository Firewall, or there was an error.", util.ColorRed)
		cli.PrintCliln(fmt.Sprintf("Details: %v", err), util.ColorRed)
		return "", err
	}

	return *iqConnection.Url, nil
}

// SelectRepository prompts the user to select a repository from available proxies in given format
func (c *NxrmConnection) SelectRepository(formatName string) (string, error) {
	// Get all repositories using the RepositoryManagementAPI
	repos, _, err := c.apiClient.RepositoryManagementAPI.GetAllRepositories(*c.ctx).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to list repositories: %w", err)
	}

	// Filter for proxy repositories of the correct format
	var proxyRepos []v3.RepositoryXO
	for _, repo := range repos {
		// Check if it's a proxy and matches the format
		if repo.GetType() == "proxy" && repo.GetFormat() == formatName {
			proxyRepos = append(proxyRepos, repo)
		}
	}

	if len(proxyRepos) == 0 {
		return "", fmt.Errorf("no %s proxy repositories found", formatName)
	}

	cli.PrintCliln(fmt.Sprintf("\nAvailable %s proxy repositories:", formatName), util.ColorYellow)
	for i, repo := range proxyRepos {
		cli.PrintCliln(fmt.Sprintf("%2d) %s", i+1, repo.GetName()), util.ColorReset)
	}

	choice := cli.ReadInput("\nSelect repository (enter number): ")

	var selectedIndex int
	_, err = fmt.Sscanf(choice, "%d", &selectedIndex)
	if err != nil || selectedIndex < 1 || selectedIndex > len(proxyRepos) {
		return "", fmt.Errorf("invalid selection")
	}

	return proxyRepos[selectedIndex-1].GetName(), nil
}

func NewNxrmConnection(nexusURL, username, password string) (*NxrmConnection, error) {
	// Create API client configuration
	configuration := v3.NewConfiguration()
	configuration.Servers = v3.ServerConfigurations{
		{
			URL: nexusURL + "/service/rest",
		},
	}

	// Create API client
	apiClient := v3.NewAPIClient(configuration)

	// Create context with basic auth
	ctx := context.WithValue(context.Background(), v3.ContextBasicAuth, v3.BasicAuth{
		UserName: username,
		Password: password,
	})

	connection := &NxrmConnection{
		apiClient: apiClient,
		baseUrl:   nexusURL,
		ctx:       &ctx,
		username:  username,
		password:  password,
	}

	err := connection.validateConnection()
	if err != nil {
		return nil, err
	}

	return connection, nil
}
