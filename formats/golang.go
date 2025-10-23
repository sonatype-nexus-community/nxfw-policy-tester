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

package formats

import (
	"fmt"
)

// GolangFormat implements PackageFormat for NPM
type GolangFormat struct{}

func (n GolangFormat) GetName() string {
	return "go"
}

func (n GolangFormat) GetDisplayName() string {
	return "Golang"
}

func (n GolangFormat) GetPackages() []Package {
	return []Package{
		// Security
		{Name: "github.com/tmc/langchaingo", Version: "v0.1.6", PolicyName: SecurityCritical, Extension: "zip"},
		{Name: "golang.org/x/crypto", Version: "v0.3.0", PolicyName: SecurityHigh, Extension: "zip"},
		{Name: "github.com/hashicorp/yamux", Version: "v0.1.1", PolicyName: SecurityMedium, Extension: "zip"},

		// Legal
		{Name: "github.com/lcomrade/lenpaste", Version: "v1.3.1", PolicyName: LicenseBanned, Extension: "zip"},
		{Name: "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3", Version: "v3.81.6", PolicyName: LicenseNone, Extension: "zip"},
		{Name: "go.wit.com/lib/cobol", Version: "v0.0.29", PolicyName: LicenseCopyLeft, Extension: "zip"},
		{Name: "github.com/unidoc/unipdf/v3", Version: "v3.69.0", PolicyName: LicenseNonStandard, Extension: "zip"},
	}
}

func (n GolangFormat) ConstructURL(nexusURL, repoName string, pkg Package) string {
	// /repository/golang.org/github.com/fatih/color/@v/v1.16.0.zip
	return fmt.Sprintf("%s/repository/%s/%s/%%40v/%s.%s", nexusURL, repoName, pkg.Name, pkg.Version, pkg.Extension)
}

func (n GolangFormat) FormatPackageName(pkg Package) string {
	return fmt.Sprintf("%s@%s", pkg.Name, pkg.Version)
}
