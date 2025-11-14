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

// PyPIFormat implements PackageFormat for PyPI
type PyPIFormat struct{}

func (p PyPIFormat) GetName() string {
	return "pypi"
}

func (p PyPIFormat) GetDisplayName() string {
	return "PyPI"
}

func (p PyPIFormat) GetPackages() []Package {
	return []Package{
		// Security
		{Name: "Django", Version: "1.6", PolicyName: SecurityCritical, Extension: "whl", Qualifier: "py2.py3-none-any"},
		{Name: "Flask", Version: "0.12", PolicyName: SecurityHigh, Extension: "whl", Qualifier: "py2.py3-none-any"},
		{Name: "Click", Version: "7.0", PolicyName: SecurityMedium, Extension: "whl", Qualifier: "py2.py3-none-any"},
		{Name: "requests-toolbelt", Version: "1.0.0", PolicyName: SecurityLow, Extension: "tar.gz", Qualifier: ""},
		{Name: "python-policy-demo", Version: "1.1.0", PolicyName: SecurityMalicious, Extension: "tar.gz", Qualifier: ""},
		{Name: "python-policy-demo", Version: "1.2.0", PolicyName: IntegrityRating, Extension: "tar.gz", Qualifier: ""},
		{Name: "python-policy-demo", Version: "1.3.0", PolicyName: IntegrityRating, Extension: "tar.gz", Qualifier: ""},

		// Legal
		{Name: "nltk", Version: "3.9.2", PolicyName: LicenseBanned, Extension: "whl", Qualifier: "py3-none-any"},
		{Name: "gallery_dl", Version: "1.29.0", PolicyName: LicenseCopyLeft, Extension: "whl", Qualifier: "py3-none-any"},
		{Name: "dawdaw", Version: "0.1.2", PolicyName: LicenseNonStandard, Extension: "whl", Qualifier: "py2.py3-none-any"},
		{Name: "pypi-project-no-license", Version: "0.1.1", PolicyName: LicenseNone, Extension: "tar.gz", Qualifier: ""},

		// None
		{Name: "google-cloud-vision", Version: "3.5.0", PolicyName: None, Extension: "tar.gz", Qualifier: ""},
	}
}

func (p PyPIFormat) ConstructURL(nexusURL, repoName string, pkg Package) string {
	// PyPI uses the "simple" API format in Nexus
	// Format: /repository/{repo}/simple/{package}/{package}-{version}-{qualifier}.{extension}
	// Example: /repository/pypi-proxy/simple/django/django-1.6-py2.py3-none-any.whl

	// /repository/pupy-proxy/packages/click/6.6/click-6.6-py2.py3-none-any.whl

	// Normalize package name (PyPI uses lowercase with hyphens replaced)
	normalizedName := strings.ToLower(pkg.Name)
	// normalizedName = strings.ReplaceAll(normalizedName, "_", "-")

	filename := fmt.Sprintf("%s-%s-%s.%s", pkg.Name, pkg.Version, pkg.Qualifier, pkg.Extension)
	if pkg.Extension == "tar.gz" {
		filename = fmt.Sprintf("%s-%s.%s", pkg.Name, pkg.Version, pkg.Extension)
	}

	return fmt.Sprintf("%s/repository/%s/packages/%s/%s/%s",
		nexusURL, repoName, normalizedName, pkg.Version, filename)
}

func (p PyPIFormat) FormatPackageName(pkg Package) string {
	if pkg.Qualifier != "" {
		return fmt.Sprintf("%s@%s (%s, .%s)", pkg.Name, pkg.Version, pkg.Qualifier, pkg.Extension)
	}
	return fmt.Sprintf("%s@%s (.%s)", pkg.Name, pkg.Version, pkg.Extension)
}
