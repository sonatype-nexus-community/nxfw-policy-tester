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
		{Name: "golang.org/x/crypto", Version: "v0.3.0", PolicyName: SecurityHigh, Extension: "zip"},
		{Name: "github.com/hashicorp/yamux", Version: "v0.1.1", PolicyName: SecurityMedium, Extension: "zip"},

		// Legal

	}
}

func (n GolangFormat) ConstructURL(nexusURL, repoName string, pkg Package) string {
	// /repository/golang.org/github.com/fatih/color/@v/v1.16.0.zip
	return fmt.Sprintf("%s/repository/%s/%s/%%40v/%s.%s", nexusURL, repoName, pkg.Name, pkg.Version, pkg.Extension)
}

func (n GolangFormat) FormatPackageName(pkg Package) string {
	return fmt.Sprintf("%s@%s", pkg.Name, pkg.Version)
}
