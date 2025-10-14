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

// HuggingFaceFormat implements PackageFormat for HuggingFace
type HuggingFaceFormat struct{}

func (h HuggingFaceFormat) GetName() string {
	return "huggingface"
}

func (h HuggingFaceFormat) GetDisplayName() string {
	return "HuggingFace"
}

func (h HuggingFaceFormat) GetPackages() []Package {
	// Placeholder - packages will be provided later
	return []Package{
		// Security
		// {Name: "sonatype/huggingface-policy-demo", Version: "9f69193fe915031a1cb5be8adef4a40b43778e9a", PolicyName: IntegrityPending, Extension: "", Qualifier: "9f69193fe915031a1cb5be8adef4a40b43778e9a:pytorch_model.bin"},
		// {Name: "sonatype/huggingface-policy-demo", Version: "5793ec913638e247ac9311e7b085d43a74e80a03", PolicyName: IntegritySuspicious, Extension: "", Qualifier: "5793ec913638e247ac9311e7b085d43a74e80a03:pytorch_model.bin"},
		// {Name: "sonatype/huggingface-policy-demo", Version: "538f4075f93b173f75f10e505448c4d1ddb05515", PolicyName: SecurityMalicious, Extension: "", Qualifier: "538f4075f93b173f75f10e505448c4d1ddb05515:pytorch_model.bin"},

		// Legal
		{Name: "OuteAI/OuteTTS-0.2-500M-GGUF", Version: "ee3de04a4d6ca4b41d7f2598734636c08c82c713", PolicyName: LicenseBanned, Extension: "", Qualifier: "ee3de04a4d6ca4b41d7f2598734636c08c82c713:OuteTTS-0.2-500M-FP16.gguf"},
	}
}

func (h HuggingFaceFormat) ConstructURL(nexusURL, repoName string, pkg Package) string {
	// HuggingFace model URL format in Nexus
	// Format: /repository/{repo}/{model-id}/resolve/main/{file}
	// Example: /repository/huggingface-proxy/bert-base-uncased/resolve/main/pytorch_model.bin

	// HuggingFace models typically have a model ID and reference a specific file
	// The Qualifier field can store the branch/ref and filename
	// Format: "branch:filename" (e.g., "main:pytorch_model.bin")

	branch := "main"
	filename := pkg.Extension

	if pkg.Qualifier != "" {
		parts := strings.Split(pkg.Qualifier, ":")
		if len(parts) == 2 {
			branch = parts[0]
			filename = parts[1]
		} else {
			filename = pkg.Qualifier
		}
	}

	return fmt.Sprintf("%s/repository/%s/%s/resolve/%s/%s",
		nexusURL, repoName, pkg.Name, branch, filename)
}

func (h HuggingFaceFormat) FormatPackageName(pkg Package) string {
	if pkg.Qualifier != "" {
		parts := strings.Split(pkg.Qualifier, ":")
		if len(parts) == 2 {
			return fmt.Sprintf("%s@%s (file: %s)", pkg.Name, pkg.Version, parts[1])
		}
		return fmt.Sprintf("%s@%s (file: %s)", pkg.Name, pkg.Version, pkg.Qualifier)
	}
	return fmt.Sprintf("%s@%s", pkg.Name, pkg.Version)
}
