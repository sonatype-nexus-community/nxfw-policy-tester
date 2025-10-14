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

// CargoFormat implements PackageFormat for Cargo
type CargoFormat struct{}

func (c CargoFormat) GetName() string {
	return "cargo"
}

func (c CargoFormat) GetDisplayName() string {
	return "Cargo"
}

func (c CargoFormat) GetPackages() []Package {
	// Placeholder - packages will be provided later
	return []Package{
		{Name: "hyper", Version: "0.14.9", PolicyName: SecurityCritical},
		{Name: "abi_stable", Version: "0.8.4", PolicyName: SecurityHigh},
	}
}

func (c CargoFormat) ConstructURL(nexusURL, repoName string, pkg Package) string {
	// Cargo package URL format in Nexus
	// /repository/cargo-proxy/crates/abi_stable/0.8.4/download

	return fmt.Sprintf("%s/repository/%s/crates/%s/%s/download",
		nexusURL, repoName, pkg.Name, pkg.Version)
}

func (c CargoFormat) FormatPackageName(pkg Package) string {
	return fmt.Sprintf("%s@%s (.crate)", pkg.Name, pkg.Version)
}
