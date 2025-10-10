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

// SecurityLevel represents the security classification of a package
type SecurityLevel string

const (
	SecurityCritical    SecurityLevel = "Security-Critical"
	SecurityHigh        SecurityLevel = "Security-High"
	SecurityMedium      SecurityLevel = "Security-Medium"
	SecurityLow         SecurityLevel = "Security-Low"
	IntegrityPending    SecurityLevel = "Integrity-Pending"
	IntegritySuspicious SecurityLevel = "Integrity-Suspicious"
	SecurityMalicious   SecurityLevel = "Security-Malicious"
)

// Package represents a package to be checked
type Package struct {
	Name          string
	Version       string
	SecurityLevel SecurityLevel
	Extension     string
	Qualifier     string // For PyPI wheel qualifiers like py2.py3-none-any
}

// PackageFormat represents a package format handler
type PackageFormat interface {
	GetName() string
	GetDisplayName() string
	GetPackages() []Package
	ConstructURL(nexusURL, repoName string, pkg Package) string
	FormatPackageName(pkg Package) string
}

// CheckResult represents the result of checking a package
type CheckResult struct {
	Package   Package
	Available bool
	HTTPCode  int
}
