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
	"net/url"
	"strings"
)

// NPMFormat implements PackageFormat for NPM
type NPMFormat struct{}

func (n NPMFormat) GetName() string {
	return "npm"
}

func (n NPMFormat) GetDisplayName() string {
	return "NPM"
}

func (n NPMFormat) GetPackages() []Package {
	return []Package{
		{Name: "bson", Version: "1.0.9", SecurityLevel: SecurityCritical, Extension: "tgz"},
		{Name: "braces", Version: "1.8.5", SecurityLevel: SecurityHigh, Extension: "tgz"},
		{Name: "cookie", Version: "0.3.1", SecurityLevel: SecurityMedium, Extension: "tgz"},
		{Name: "react-dom", Version: "18.3.1", SecurityLevel: SecurityLow, Extension: "tgz"},
		{Name: "@sonatype/policy-demo", Version: "2.3.0", SecurityLevel: IntegrityPending, Extension: "tgz"},
		{Name: "@sonatype/policy-demo", Version: "2.2.0", SecurityLevel: IntegritySuspicious, Extension: "tgz"},
		{Name: "@sonatype/policy-demo", Version: "2.1.0", SecurityLevel: SecurityMalicious, Extension: "tgz"},
	}
}

func (n NPMFormat) ConstructURL(nexusURL, repoName string, pkg Package) string {
	filename := fmt.Sprintf("%s-%s.%s", pkg.Name, pkg.Version, pkg.Extension)
	if strings.Contains(pkg.Name, "/") {
		nameParts := strings.Split(pkg.Name, "/")
		filename = fmt.Sprintf("%s-%s.%s", nameParts[1], pkg.Version, pkg.Extension)
	}

	return fmt.Sprintf("%s/repository/%s/%s/-/%s",
		nexusURL, repoName, url.QueryEscape(pkg.Name), filename)
}

func (n NPMFormat) FormatPackageName(pkg Package) string {
	return fmt.Sprintf("%s@%s", pkg.Name, pkg.Version)
}
