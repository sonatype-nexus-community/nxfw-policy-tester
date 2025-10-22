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
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sonatype-nexus-community/nxfw-policy-tester/formats"
	"github.com/sonatype-nexus-community/nxfw-policy-tester/util"
)

// Print a line to CLI in a colour and reset
func PrintCliln(msg string, color string) {
	fmt.Printf("%s%s%s\n", color, msg, util.ColorReset)
}

// Prompt to select Format
func PromptSelectFormat(possibleFormats []formats.PackageFormat) formats.PackageFormat {
	PrintCliln("Select package format:", util.ColorYellow)
	for i, format := range possibleFormats {
		PrintCliln(fmt.Sprintf("%d) %s", i+1, format.GetDisplayName()), util.ColorReset)
	}

	choice := ReadInput("Enter choice: ")
	choiceI, err := strconv.Atoi(choice)
	if err != nil {
		PrintCliln("Error: Invalid choice.", util.ColorRed)
		os.Exit(1)
	}

	if choiceI > len(possibleFormats) {
		PrintCliln("Error: Invalid choice.", util.ColorRed)
		os.Exit(1)
	}

	return possibleFormats[(choiceI - 1)]
}

// ReadInput reads a line from stdin
func ReadInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
