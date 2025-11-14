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
	"strings"
)

// CondaFormat implements PackageFormat for Conda
type CondaFormat struct{}

func (c CondaFormat) GetName() string {
	return "conda"
}

func (c CondaFormat) GetDisplayName() string {
	return "Conda"
}

func (c CondaFormat) GetPackages() []Package {
	return []Package{
		// Security
		{Name: "gettext", Version: "0.19.8.1", PolicyName: SecurityCritical, Extension: "tar.bz2", Qualifier: "main/linux-64/h9b4dc7a_1"},
		{Name: "setuptools", Version: "61.2.0", PolicyName: SecurityHigh, Extension: "tar.bz2", Qualifier: "main/linux-64/py310h06a4308_0"},
		// {Name: "gettext", Version: "0.21.1", PolicyName: SecurityLow, Extension: "tar.bz2", Qualifier: "main/linux-64/h27087fc_0"},

		// License
		// {Name: "glmnet", Version: "2.2.1", PolicyName: LicenseCopyLeft, Extension: "conda", Qualifier: "main/linux-64/py310h31179b7_6"},
	}
}

func (c CondaFormat) ConstructURL(nexusURL, repoName string, pkg Package) string {
	// Conda package URL format in Nexus
	// Format: /repository/{repo}/{channel}/{platform}/{package}-{version}-{build}.{extension}

	// gettext@0.19.8.1?build=h9b4dc7a_1&channel=main&subdir=linux-64&type=conda
	// /asn1crypto/0.24.0/download/linux-64/asn1crypto-0.24.0-py37_1003.tar.bz2

	// For Conda, the qualifier field will be used to store: channel/platform/build
	// Parse qualifier as "channel/platform/build"
	parts := strings.Split(pkg.Qualifier, "/")

	var channel, platform, build string
	if len(parts) >= 3 {
		channel = parts[0]
		platform = parts[1]
		build = parts[2]
	} else if len(parts) == 2 {
		channel = "main"
		platform = parts[0]
		build = parts[1]
	} else {
		// Default values if not properly formatted
		channel = "main"
		platform = "linux-64"
		build = "0"
	}

	filename := fmt.Sprintf("%s-%s-%s.%s", pkg.Name, pkg.Version, build, pkg.Extension)

	return fmt.Sprintf("%s/repository/%s/%s/%s/%s",
		nexusURL, repoName, channel, platform, filename)
}

func (c CondaFormat) FormatPackageName(pkg Package) string {
	if pkg.Qualifier != "" {
		// Parse qualifier to show channel/platform/build separately
		parts := strings.Split(pkg.Qualifier, "/")
		if len(parts) >= 3 {
			return fmt.Sprintf("%s@%s (channel: %s, platform: %s, build: %s, .%s)",
				pkg.Name, pkg.Version, parts[0], parts[1], parts[2], pkg.Extension)
		}
		return fmt.Sprintf("%s@%s (%s, .%s)", pkg.Name, pkg.Version, pkg.Qualifier, pkg.Extension)
	}
	return fmt.Sprintf("%s@%s (.%s)", pkg.Name, pkg.Version, pkg.Extension)
}
