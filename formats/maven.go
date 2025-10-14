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

// MavenFormat implements PackageFormat for Maven
type MavenFormat struct{}

func (m MavenFormat) GetName() string {
	return "maven2"
}

func (m MavenFormat) GetDisplayName() string {
	return "Maven"
}

func (m MavenFormat) GetPackages() []Package {
	return []Package{
		{Name: "com.amazonaws/aws-android-sdk-core", Version: "2.75.0", PolicyName: SecurityCritical, Extension: "aar"},
		{Name: "org.jsoup/jsoup", Version: "1.13.1", PolicyName: SecurityHigh, Extension: "jar"},
		{Name: "ant/ant", Version: "1.6.5", PolicyName: SecurityMedium, Extension: "jar"},
		{Name: "org.springframework/spring-context", Version: "6.2.3", PolicyName: SecurityLow, Extension: "jar"},
		{Name: "org.sonatype/maven-policy-demo", Version: "1.1.0", PolicyName: SecurityMalicious, Extension: "jar"},
		{Name: "org.sonatype/maven-policy-demo", Version: "1.2.0", PolicyName: IntegritySuspicious, Extension: "jar"},
		{Name: "org.sonatype/maven-policy-demo", Version: "1.3.0", PolicyName: IntegrityPending, Extension: "jar"},
	}
}

func (m MavenFormat) ConstructURL(nexusURL, repoName string, pkg Package) string {
	parts := strings.Split(pkg.Name, "/")
	if len(parts) != 2 {
		return ""
	}
	group := parts[0]
	artifact := parts[1]

	// Convert group dots to slashes (e.g., com.amazonaws -> com/amazonaws)
	groupPath := strings.ReplaceAll(group, ".", "/")

	return fmt.Sprintf("%s/repository/%s/%s/%s/%s/%s-%s.%s",
		nexusURL, repoName, groupPath, artifact, pkg.Version, artifact, pkg.Version, pkg.Extension)
}

func (m MavenFormat) FormatPackageName(pkg Package) string {
	return fmt.Sprintf("%s@%s (.%s)", pkg.Name, pkg.Version, pkg.Extension)
}
