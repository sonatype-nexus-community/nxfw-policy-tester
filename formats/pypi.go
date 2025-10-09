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

import "fmt"

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
		{Name: "django", Version: "1.6", SecurityLevel: SecurityCritical, Extension: "whl", Qualifier: "py2.py3-none-any"},
		{Name: "flask", Version: "0.12", SecurityLevel: SecurityHigh, Extension: "whl", Qualifier: "py2.py3-none-any"},
		{Name: "click", Version: "7.0", SecurityLevel: SecurityMedium, Extension: "whl", Qualifier: "py2.py3-none-any"},
	}
}

func (p PyPIFormat) ConstructURL(nexusURL, repoName string, pkg Package) string {
	// PyPI wheel format: package-version-qualifier.whl
	// Example: django-1.6-py2.py3-none-any.whl
	// Path: /repository/pypi-proxy/packages/django/1.6/django-1.6-py2.py3-none-any.whl
	//       /repository/pupy-proxy/packages/filelock/3.12.2/filelock-3.12.2-py3-none-any.whl

	filename := fmt.Sprintf("%s-%s-%s.%s", pkg.Name, pkg.Version, pkg.Qualifier, pkg.Extension)

	return fmt.Sprintf("%s/repository/%s/packages/%s/%s/%s",
		nexusURL, repoName, pkg.Name, pkg.Version, filename)
}

func (p PyPIFormat) FormatPackageName(pkg Package) string {
	if pkg.Qualifier != "" {
		return fmt.Sprintf("%s@%s (%s, .%s)", pkg.Name, pkg.Version, pkg.Qualifier, pkg.Extension)
	}
	return fmt.Sprintf("%s@%s (.%s)", pkg.Name, pkg.Version, pkg.Extension)
}
