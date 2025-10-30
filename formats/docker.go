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

// DockerFormat implements PackageFormat for Maven
type DockerFormat struct{}

func (m DockerFormat) GetName() string {
	return "docker"
}

func (m DockerFormat) GetDisplayName() string {
	return "Docker"
}

func (m DockerFormat) GetPackages() []Package {
	return []Package{
		// Security
		{Name: "sonatypecommunity/docker-policy-demo", Version: "Security-Critical", PolicyName: SecurityCritical, Extension: ""},
		{Name: "sonatypecommunity/docker-policy-demo", Version: "Security-High", PolicyName: SecurityHigh, Extension: ""},
		{Name: "sonatypecommunity/docker-policy-demo", Version: "Security-Medium", PolicyName: SecurityMedium, Extension: ""},
		{Name: "sonatypecommunity/docker-policy-demo", Version: "Security-Low", PolicyName: SecurityLow, Extension: ""},
		{Name: "sonatypecommunity/docker-policy-demo", Version: "Security-Malicious", PolicyName: SecurityMalicious, Extension: ""},
		{Name: "sonatypecommunity/docker-policy-demo", Version: "Integrity-Suspicious", PolicyName: IntegrityRating, Extension: ""},
		{Name: "sonatypecommunity/docker-policy-demo", Version: "Integrity-Pending", PolicyName: IntegrityRating, Extension: ""},

		// Legal
	}
}

func (m DockerFormat) ConstructURL(nexusURL, repoName string, pkg Package) string {
	// Docker registry v2 API format in Nexus
	// Format: /v2/{image-name}/manifests/{tag}
	// Example: /v2/sonatypecommunity/docker-policy-demo/manifests/Security-Malicious

	// For Docker, pkg.Name is the image name and pkg.Version is the tag
	// The image name might include namespace (e.g., sonatypecommunity/docker-policy-demo)

	return fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", nexusURL, repoName, pkg.Name, pkg.Version)
}

func (m DockerFormat) FormatPackageName(pkg Package) string {
	return fmt.Sprintf("%s:%s", pkg.Name, pkg.Version)
}
