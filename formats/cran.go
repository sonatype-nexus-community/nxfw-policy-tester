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

// CranFormat implements PackageFormat for PyPI
type CranFormat struct{}

func (p CranFormat) GetName() string {
	return "r"
}

func (p CranFormat) GetDisplayName() string {
	return "CRAN(R)"
}

func (p CranFormat) GetPackages() []Package {
	return []Package{
		{Name: "readxl", Version: "0.1.0", PolicyName: SecurityHigh, Extension: "tar.gz"},
		{Name: "xgboost", Version: "0.6-3", PolicyName: SecurityMedium, Extension: "tar.gz"},
	}
}

func (p CranFormat) ConstructURL(nexusURL, repoName string, pkg Package) string {
	// /repository/r-proxy/src/contrib/Archive/readxl/readxl_0.1.0.tar.gz

	// Normalize package name (PyPI uses lowercase with hyphens replaced)
	// normalizedName := strings.ToLower(pkg.Name)
	// normalizedName = strings.ReplaceAll(normalizedName, "_", "-")

	// filename := fmt.Sprintf("%s-%s-%s.%s", pkg.Name, pkg.Version, pkg.Qualifier, pkg.Extension)
	// if pkg.Extension == "tar.gz" {
	// 	filename = fmt.Sprintf("%s-%s.%s", pkg.Name, pkg.Version, pkg.Extension)
	// }

	return fmt.Sprintf("%s/repository/%s/src/contrib/Archive/%s/%s_%s.%s",
		nexusURL, repoName, pkg.Name, pkg.Name, pkg.Version, pkg.Extension)
}

func (p CranFormat) FormatPackageName(pkg Package) string {
	if pkg.Qualifier != "" {
		return fmt.Sprintf("%s@%s (%s, .%s)", pkg.Name, pkg.Version, pkg.Qualifier, pkg.Extension)
	}
	return fmt.Sprintf("%s@%s (.%s)", pkg.Name, pkg.Version, pkg.Extension)
}
