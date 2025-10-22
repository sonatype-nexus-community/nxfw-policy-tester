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
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sonatype-nexus-community/nxfw-policy-tester/cli"
	"github.com/sonatype-nexus-community/nxfw-policy-tester/formats"
	"github.com/sonatype-nexus-community/nxfw-policy-tester/nxiq"
	"github.com/sonatype-nexus-community/nxfw-policy-tester/nxrm"
	"github.com/sonatype-nexus-community/nxfw-policy-tester/util"
)

var (
	allSupportedFormats = []formats.PackageFormat{
		formats.CargoFormat{},
		formats.CranFormat{},
		formats.CondaFormat{},
		formats.GolangFormat{},
		formats.HuggingFaceFormat{},
		formats.MavenFormat{},
		formats.NPMFormat{},
		formats.NuGetFormat{},
		formats.PyPIFormat{},
	}
	currentRuntime string = runtime.GOOS
	commit                = "unknown"
	version               = "dev"
)

// displaySummary displays the configuration summary
func displaySummary(nexusURL, repoName string, format formats.PackageFormat) {
	cli.PrintCliln("\n=== Configuration Summary ===\n", util.ColorYellow)
	cli.PrintCliln("Nexus Repository URL: "+nexusURL, util.ColorReset)
	cli.PrintCliln("Format: "+format.GetDisplayName(), util.ColorReset)
	cli.PrintCliln("Repository: "+repoName, util.ColorReset)
	cli.PrintCliln("\nPackages to check:", util.ColorReset)

	packages := format.GetPackages()
	for _, pkg := range packages {
		color := pkg.PolicyName.GetSecurityColor()
		cli.PrintCliln(
			fmt.Sprintf(
				"  - %s %s[%s]",
				format.FormatPackageName(pkg),
				color,
				pkg.PolicyName,
			),
			util.ColorReset,
		)
	}
	fmt.Println()
}

// displayResults displays the check results summary
func displayResults(results []formats.CheckResult, format formats.PackageFormat) {
	availableCount := 0
	failedCount := 0
	quarantinedCount := 0

	for _, result := range results {
		if result.Available {
			availableCount++
		} else if result.Quarantined {
			quarantinedCount++
		} else if result.Failed {
			failedCount++
		}
	}

	cli.PrintCliln("===================== Test Summary =====================", util.ColorYellow)
	cli.PrintCliln(fmt.Sprintf("Downloadable:         %02d", availableCount), util.ColorGreen)
	cli.PrintCliln(fmt.Sprintf("Quarantined:          %02d", quarantinedCount), util.ColorCyan)
	cli.PrintCliln(fmt.Sprintf("Failure:              %02d", failedCount), util.ColorRed)

	cli.PrintCliln("----------------------- Details ------------------------", util.ColorYellow)
	for _, result := range results {
		color := result.Package.PolicyName.GetSecurityColor()
		var status = "UNKNOWN"
		if result.Available {
			status = fmt.Sprintf("%sAVAILABLE%s", util.ColorGreen, util.ColorReset)
		} else if result.Quarantined && result.QuarantinedWithExpectedPolicy {
			status = fmt.Sprintf("%sQUARANTINED%s", util.ColorCyan, util.ColorReset)
		} else if result.Quarantined && !result.QuarantinedWithExpectedPolicy {
			status = fmt.Sprintf("%sOOOPS%s", util.ColorRed, util.ColorReset)
		} else if result.Failed {
			status = fmt.Sprintf("%sFAILED%s", util.ColorRed, util.ColorReset)
		}

		cli.PrintCliln(
			fmt.Sprintf(
				"%25s: %45s - %15s",
				result.Package.PolicyName,
				format.FormatPackageName(result.Package),
				status,
			),
			color,
		)
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
	println("             ___  _____  __  __  __  __  __  __  _  _  ____  ____  _  _ ")
	println("            / __)(  _  )(  \\/  )(  \\/  )(  )(  )( \\( )(_  _)(_  _)( \\/ )       ")
	println("           ( (__  )(_)(  )    (  )    (  )(__)(  )  (  _)(_   )(   \\  /        ")
	println("            \\___)(_____)(_/\\/\\_)(_/\\/\\_)(______)(_)\\_)(____) (__)  (__)        ")
	println("")
	println(fmt.Sprintf("	Running on:		%s/%s", currentRuntime, runtime.GOARCH))
	println(fmt.Sprintf("	Version: 		%s (%s)", version, commit))
	println("")
	println(strings.Repeat("⬢⬡", 42))
	println("")

	// Get credentials from environment variables
	nxrmUsername := os.Getenv("NXRM_USERNAME")
	nxrmPassword := os.Getenv("NXRM_PASSWORD")
	nxiqUsername := os.Getenv("NXIQ_USERNAME")
	nxiqPassword := os.Getenv("NXIQ_PASSWORD")

	if nxrmUsername == "" || nxrmPassword == "" {
		fmt.Printf("%sError: NXRM_USERNAME and NXRM_PASSWORD environment variables must be set.%s\n", util.ColorRed, util.ColorReset)
		fmt.Println("Example: export NXRM_USERNAME='your_username'")
		fmt.Println("         export NXRM_PASSWORD='your_password'")
		os.Exit(1)
	}

	if nxiqUsername == "" || nxiqPassword == "" {
		fmt.Printf("%sError: NXIQ_PASSWORD and NXIQ_USERNAME environment variables must be set.%s\n", util.ColorRed, util.ColorReset)
		fmt.Println("Example: export NXIQ_PASSWORD='your_username'")
		fmt.Println("         export NXIQ_USERNAME='your_password'")
		os.Exit(1)
	}

	// Get Nexus URL
	fmt.Printf("\n%sEnter your Sonatype Nexus Repository URL:%s\n", util.ColorYellow, util.ColorReset)
	fmt.Println("(Example: https://nexus.example.com)")
	nexusURL := cli.ReadInput("")
	nexusURL = strings.TrimSuffix(nexusURL, "/")

	// Validate URL
	if !strings.HasPrefix(nexusURL, "http://") && !strings.HasPrefix(nexusURL, "https://") {
		fmt.Printf("%sError: Invalid URL format. URL must start with http:// or https://%s\n", util.ColorRed, util.ColorReset)
		os.Exit(1)
	}

	// NXRM Connection
	nxrmConnection, err := nxrm.NewNxrmConnection(nexusURL, nxrmUsername, nxrmPassword)
	if err != nil {
		os.Exit(1)
	}

	cli.PrintCliln("✓ Successfully authenticated with Sonatype Nexus Repository", util.ColorGreen)

	// Get IQ URL
	nxiqUrl, err := nxrmConnection.GetConnectedIqServer()
	if err != nil {
		os.Exit(1)
	}

	// NXIQ Connection
	nxiqConnection, err := nxiq.NewNxiqConnection(nxiqUrl, nxiqUsername, nxiqPassword)
	if err != nil {
		os.Exit(1)
	}

	cli.PrintCliln(fmt.Sprintf("✓ Successfully authenticated with Sonatype IQ Server (%s)", nxiqUrl), util.ColorGreen)

	// Select package format
	format := cli.PromptSelectFormat(allSupportedFormats)

	// Select Repository
	repoName, err := nxrmConnection.SelectRepository(format.GetName())
	if err != nil {
		cli.PrintCliln(fmt.Sprintf("Error: %v", err), util.ColorRed)
		os.Exit(1)
	}

	// Display summary
	displaySummary(nexusURL, repoName, format)

	// Confirm
	confirmation := cli.ReadInput("Proceed with checking packages? (y/n): ")
	if !strings.HasPrefix(strings.ToLower(confirmation), "y") {
		cli.PrintCliln("User cancelled.", util.ColorRed)
		os.Exit(0)
	}

	// Check packages
	results, err := nxrmConnection.CheckPackages(repoName, format, nxiqConnection)
	if err != nil {
		cli.PrintCliln(fmt.Sprintf("Unexpected failure: %v", err), util.ColorRed)
		os.Exit(1)
	}

	// Display results
	displayResults(results, format)
}
