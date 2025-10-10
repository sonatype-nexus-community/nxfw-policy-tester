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
package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	"github.com/sonatype-nexus-community/nxfw-policy-tester/formats"
)

var (
	currentRuntime string = runtime.GOOS
	commit                = "unknown"
	version               = "dev"
)

// Color codes
const (
	ColorRed     = "\033[0;31m"
	ColorGreen   = "\033[0;32m"
	ColorYellow  = "\033[1;33m"
	ColorCyan    = "\033[0;36m"
	ColorMagenta = "\033[0;35m"
	ColorBlue    = "\033[0;34m"
	ColorReset   = "\033[0m"
)

// getSecurityColor returns the color code for a security level
func getSecurityColor(level formats.SecurityLevel) string {
	switch level {
	case formats.SecurityCritical:
		return ColorRed
	case formats.SecurityHigh:
		return ColorMagenta
	case formats.SecurityMedium:
		return ColorYellow
	default:
		return ColorReset
	}
}

// readInput reads a line from stdin
func readInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// selectFormat prompts the user to select a package format
func selectFormat() formats.PackageFormat {
	formats := []formats.PackageFormat{
		// formats.CondaFormat{},
		formats.MavenFormat{},
		formats.NPMFormat{},
		formats.NuGetFormat{},
		formats.PyPIFormat{},
	}

	fmt.Printf("%sSelect package format:%s\n", ColorYellow, ColorReset)
	for i, format := range formats {
		fmt.Printf("%d) %s\n", i+1, format.GetDisplayName())
	}

	choice := readInput("Enter choice: ")
	choiceI, err := strconv.Atoi(choice)
	if err != nil {
		fmt.Printf("%sError: Invalid choice.%s\n", ColorRed, ColorReset)
		os.Exit(1)
	}

	if choiceI > len(formats) {
		fmt.Printf("%sError: Invalid choice.%s\n", ColorRed, ColorReset)
		os.Exit(1)
	}

	return formats[(choiceI - 1)]
}

// selectRepository prompts the user to select a repository from available proxies
func selectRepository(apiClient *v3.APIClient, ctx context.Context, formatName string) (string, error) {
	// Get all repositories using the RepositoryManagementAPI
	repos, _, err := apiClient.RepositoryManagementAPI.GetAllRepositories(ctx).Execute()
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

	fmt.Printf("\n%sAvailable %s proxy repositories:%s\n", ColorYellow, formatName, ColorReset)
	for i, repo := range proxyRepos {
		fmt.Printf("%d) %s\n", i+1, repo.GetName())
	}

	choice := readInput("\nSelect repository (enter number): ")

	var selectedIndex int
	_, err = fmt.Sscanf(choice, "%d", &selectedIndex)
	if err != nil || selectedIndex < 1 || selectedIndex > len(proxyRepos) {
		return "", fmt.Errorf("invalid selection")
	}

	return proxyRepos[selectedIndex-1].GetName(), nil
}

// checkPackage checks if a package is available
func checkPackage(url, username, password string) (int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
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

// displaySummary displays the configuration summary
func displaySummary(nexusURL, repoName string, format formats.PackageFormat) {
	fmt.Printf("\n%s=== Configuration Summary ===%s\n", ColorYellow, ColorReset)
	fmt.Printf("Nexus URL: %s\n", nexusURL)
	fmt.Printf("Format: %s\n", format.GetDisplayName())
	fmt.Printf("Repository: %s\n", repoName)
	fmt.Println("\nPackages to check:")

	packages := format.GetPackages()
	for _, pkg := range packages {
		color := getSecurityColor(pkg.SecurityLevel)
		fmt.Printf("  - %s %s[%s]%s\n",
			format.FormatPackageName(pkg),
			color,
			pkg.SecurityLevel,
			ColorReset)
	}
	fmt.Println()
}

// checkPackages checks all packages and returns results
func checkPackages(nexusURL, repoName, username, password string, format formats.PackageFormat) []formats.CheckResult {
	packages := format.GetPackages()
	results := make([]formats.CheckResult, 0, len(packages))

	fmt.Printf("\n%s=== Checking Package Availability ===%s\n\n", ColorYellow, ColorReset)

	for _, pkg := range packages {
		color := getSecurityColor(pkg.SecurityLevel)
		fmt.Printf("Checking %s %s[%s]%s...\n",
			format.FormatPackageName(pkg),
			color,
			pkg.SecurityLevel,
			ColorReset)

		url := format.ConstructURL(nexusURL, repoName, pkg)
		httpCode, err := checkPackage(url, username, password)

		result := formats.CheckResult{
			Package:  pkg,
			HTTPCode: httpCode,
		}

		if err != nil {
			fmt.Printf("%s✗ Package not available: %s [%s]%s (Error: %v)\n\n",
				ColorRed,
				format.FormatPackageName(pkg),
				pkg.SecurityLevel,
				ColorReset,
				err)
			result.Available = false
		} else if httpCode == 200 {
			fmt.Printf("%s✓ Package available: %s [%s]%s\n\n",
				ColorGreen,
				format.FormatPackageName(pkg),
				pkg.SecurityLevel,
				ColorReset)
			result.Available = true
		} else {
			fmt.Printf("%s✗ Package not available: %s [%s]%s (HTTP %d)\n\n",
				ColorRed,
				format.FormatPackageName(pkg),
				pkg.SecurityLevel,
				ColorReset,
				httpCode)
			result.Available = false
		}

		results = append(results, result)
	}

	return results
}

// displayResults displays the check results summary
func displayResults(results []formats.CheckResult, format formats.PackageFormat) {
	successCount := 0
	failCount := 0

	for _, result := range results {
		if result.Available {
			successCount++
		} else {
			failCount++
		}
	}

	fmt.Printf("%s=== Check Summary ===%s\n", ColorYellow, ColorReset)
	fmt.Printf("%sAvailable: %d%s\n", ColorGreen, successCount, ColorReset)
	fmt.Printf("%sNot Available: %d%s\n", ColorRed, failCount, ColorReset)

	fmt.Printf("\n%s=== Security Level Breakdown ===%s\n", ColorYellow, ColorReset)
	for _, result := range results {
		color := getSecurityColor(result.Package.SecurityLevel)
		var status string
		if result.Available {
			status = fmt.Sprintf("%s[Available]%s", ColorGreen, ColorReset)
		} else {
			status = fmt.Sprintf("%s[Not Available]%s", ColorRed, ColorReset)
		}

		fmt.Printf("%s%s%s: %s %s\n",
			color,
			result.Package.SecurityLevel,
			ColorReset,
			format.FormatPackageName(result.Package),
			status)
	}
}

func main() {
	// Output Banner
	println(strings.Repeat("⬢⬡", 42))
	println("")
	println("	███████╗ ██████╗ ███╗   ██╗ █████╗ ████████╗██╗   ██╗██████╗ ███████╗  ")
	println(" 	██╔════╝██╔═══██╗████╗  ██║██╔══██╗╚══██╔══╝╚██╗ ██╔╝██╔══██╗██╔════╝  ")
	println("	███████╗██║   ██║██╔██╗ ██║███████║   ██║    ╚████╔╝ ██████╔╝█████╗    ")
	println(" 	╚════██║██║   ██║██║╚██╗██║██╔══██║   ██║     ╚██╔╝  ██╔═══╝ ██╔══╝    ")
	println(" 	███████║╚██████╔╝██║ ╚████║██║  ██║   ██║      ██║   ██║     ███████╗  ")
	println(" 	╚══════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝   ╚═╝      ╚═╝   ╚═╝     ╚══════╝  ")
	println("")
	println(fmt.Sprintf("	Running on:		%s/%s", currentRuntime, runtime.GOARCH))
	println(fmt.Sprintf("	Version: 		%s (%s)", version, commit))
	println("")
	println(strings.Repeat("⬢⬡", 42))
	println("")

	// Get credentials from environment variables
	username := os.Getenv("NEXUS_USERNAME")
	password := os.Getenv("NEXUS_PASSWORD")

	if username == "" || password == "" {
		fmt.Printf("%sError: NEXUS_USERNAME and NEXUS_PASSWORD environment variables must be set.%s\n", ColorRed, ColorReset)
		fmt.Println("Example: export NEXUS_USERNAME='your_username'")
		fmt.Println("         export NEXUS_PASSWORD='your_password'")
		os.Exit(1)
	}

	// Select package format
	format := selectFormat()

	// Get Nexus URL
	fmt.Printf("\n%sEnter your Sonatype Nexus Repository URL:%s\n", ColorYellow, ColorReset)
	fmt.Println("(Example: https://nexus.example.com)")
	nexusURL := readInput("")
	nexusURL = strings.TrimSuffix(nexusURL, "/")

	// Validate URL
	if !strings.HasPrefix(nexusURL, "http://") && !strings.HasPrefix(nexusURL, "https://") {
		fmt.Printf("%sError: Invalid URL format. URL must start with http:// or https://%s\n", ColorRed, ColorReset)
		os.Exit(1)
	}

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

	// Validate credentials by making a test call
	_, err := apiClient.StatusAPI.IsAvailable(ctx).Execute()
	if err != nil {
		fmt.Printf("%sError: Failed to authenticate with Nexus. Please check your credentials and URL.%s\n", ColorRed, ColorReset)
		fmt.Printf("Details: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s✓ Successfully authenticated with Nexus%s\n", ColorGreen, ColorReset)

	// Select repository
	repoName, err := selectRepository(apiClient, ctx, format.GetName())
	if err != nil {
		fmt.Printf("%sError: %v%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}

	// Display summary
	displaySummary(nexusURL, repoName, format)

	// Confirm
	confirmation := readInput("Proceed with checking packages? (y/n): ")
	if !strings.HasPrefix(strings.ToLower(confirmation), "y") {
		fmt.Println("Check cancelled.")
		os.Exit(0)
	}

	// Check packages
	results := checkPackages(nexusURL, repoName, username, password, format)

	// Display results
	displayResults(results, format)
}
