# Sonatype Repository Firewall Policy Test Suite

<!-- Badges Section -->

[![shield_gh-workflow-test]][link_gh-workflow-test]
[![shield_license]][license_file]

<!-- Add other badges or shields as appropriate -->

---

This is a simple test suite that will attempt to pull known _bad_ packages into Proxy repositores in Sonatype Nexus Repository
to allow customers to validate that Repository Firewall policies are operating as they expect.

This tool does not know the specific outcomes, as each Customer will have differing [Policy Actions](https://help.sonatype.com/en/policy-actions.html)
set.

This tool expects the [Reference Policy Set](https://help.sonatype.com/en/reference-policies.html) to be in use, but that is not mandatory.

- [Installation](#installation)
- [Usage](#usage)
- [Tests Available](#tests-available)
- [Development](#development)
- [The Fine Print](#the-fine-print)

## Installation

Obtain the binary for your Operating System and Architecture from the [GitHub Releases page](https://github.com/sonatype-nexus-community/nexus-repo-asset-lister/releases).

## Usage

Set your Sonatype Nexus Repository credentials in two environment variables:

-   `NEXUS_USERNAME`
-   `NEXUS_PASSWORD`

```bash
./nxfw-policy-tester
```

Follow the prompts - you'll need the URL to your Sonatype Nexus Repository installation (https:// only supported).

## Tests Available

| Format       | Reference Policy      | Available |
| ------------ | --------------------- | --------- |
| Cargo (Rust) | `Security-Critical`   | ✅        |
| Cargo (Rust) | `Security-High`       | ✅        |
| CRAN (R)     | `Security-Critical`   | ❌        |
| CRAN (R)     | `Security-High`       | ✅        |
| CRAN (R)     | `Security-Medium`     | ✅        |
| CRAN (R)     | `Security-Low`        | ❌        |
| Maven        | `Security-Critical`   | ✅        |
| Maven        | `Security-High`       | ✅        |
| Maven        | `Security-Medium`     | ✅        |
| Maven        | `Security-Low`        | ✅        |
| Maven        | `Security-Malicious`^ | ✅        |
| Maven        | `Integrity-Rating`±   | ✅        |
| NPM          | `Security-Critical`   | ✅        |
| NPM          | `Security-High`       | ✅        |
| NPM          | `Security-Medium`     | ✅        |
| NPM          | `Security-Low`        | ✅        |
| NPM          | `Security-Malicious`^ | ✅        |
| NPM          | `Integrity-Rating`±   | ✅        |
| Nuget        | `Security-Critical`   | ✅        |
| Nuget        | `Security-High`       | ✅        |
| Nuget        | `Security-Medium`     | ✅        |
| Nuget        | `Security-Low`        | ✅        |
| PyPi         | `Security-Critical`   | ✅        |
| PyPi         | `Security-High`       | ✅        |
| PyPi         | `Security-Medium`     | ✅        |
| PyPi         | `Security-Low`        | ❌        |
| PyPi         | `Security-Malicious`^ | ✅        |
| PyPi         | `Integrity-Rating`±   | ✅        |

> NOTES:
>
> ^ Does not rely on real malicious package for verification - uses Sonatype staged packages marked as Malicious, but that are not actually malicious in content.
>
> ± Includes packages in both Pending and Suspicious states.

## Development

See [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## The Fine Print

Remember:

This project is part of the [Sonatype Nexus Community](https://github.com/sonatype-nexus-community) organization, which is not officially supported by Sonatype. Please review the latest pull requests, issues, and commits to understand this project's readiness for contribution and use.

-   File suggestions and requests on this repo through GitHub Issues, so that the community can pitch in
-   Use or contribute to this project according to your organization's policies and your own risk tolerance
-   Don't file Sonatype support tickets related to this project— it won't reach the right people that way

Last but not least of all - have fun!

<!-- Links Section -->

[shield_gh-workflow-test]: https://img.shields.io/github/actions/workflow/status/sonatype-nexus-community/nxfw-policy-tester/build.yml?branch=main&logo=GitHub&logoColor=white 'build'
[shield_license]: https://img.shields.io/github/license/sonatype-nexus-community/nxfw-policy-tester?logo=open%20source%20initiative&logoColor=white 'license'
[link_gh-workflow-test]: https://github.com/sonatype-nexus-community/nxfw-policy-tester/actions/workflows/build.yml?query=branch%3Amain
[license_file]: https://github.com/sonatype-nexus-community/nxfw-policy-tester/blob/main/LICENSE
