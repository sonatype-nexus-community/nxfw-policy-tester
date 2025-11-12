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

// NuGetFormat implements PackageFormat for NuGet
type NuGetFormat struct{}

func (n NuGetFormat) GetName() string {
	return "nuget"
}

func (n NuGetFormat) GetDisplayName() string {
	return "NuGet"
}

func (n NuGetFormat) GetPackages() []Package {
	return []Package{
		// Security
		{Name: "log4net", Version: "2.0.3", PolicyName: SecurityCritical},
		{Name: "Newtonsoft.Json", Version: "6.0.4", PolicyName: SecurityHigh},
		{Name: "Microsoft.Owin", Version: "2.1.0", PolicyName: SecurityMedium},
		{Name: "Microsoft.AspNet.SignalR.Core", Version: "2.0.3", PolicyName: SecurityLow},

		// License
		{Name: "LigerShark.WebOptimizer.Core", Version: "3.0.344", PolicyName: LicenseNone},
		{Name: "MySql.Data", Version: "8.0.27", PolicyName: LicenseCopyLeft},
		{Name: "PayPalCheckoutSdk", Version: "1.0.3", PolicyName: LicenseCommercial},

		// None
		{Name: "Microsoft.AspNetCore.Mvc.NewtonsoftJson", Version: "5.0.3", PolicyName: None},
	}
}

func (n NuGetFormat) ConstructURL(nexusURL, repoName string, pkg Package) string {
	// NuGet v3 API format in Nexus
	// Format: /repository/{repo}/download/{package-id-lowercase}/{version}/{package-id-lowercase}.{version}.nupkg
	// Example: /repository/nuget-proxy/download/log4net/2.0.3/log4net.2.0.3.nupkg
	// /repository/nuget-proxy/v3/content/0/log4net/2.0.3/log4net.<version>.nupkg
	// NuGet uses lowercase package IDs in the download URL
	packageIDLower := strings.ToLower(pkg.Name)

	filename := fmt.Sprintf("%s.%s.nupkg", packageIDLower, pkg.Version)

	return fmt.Sprintf("%s/repository/%s/v3/content/0/%s/%s/%s",
		nexusURL, repoName, packageIDLower, pkg.Version, filename)
}

func (n NuGetFormat) FormatPackageName(pkg Package) string {
	return fmt.Sprintf("%s@%s (.nupkg)", pkg.Name, pkg.Version)
}
