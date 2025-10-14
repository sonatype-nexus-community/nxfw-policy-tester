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
- [Test Data Available](#test-data-available)
  - [Cargo (Rust)](#cargo-rust)
  - [Conda (conda-forge)](#conda-conda-forge)
  - [CRAN (R)](#cran-r)
  - [Golang (Go)](#golang-go)
  - [Huggingface.co (AI / ML)](#huggingfaceco-ai--ml)
  - [Maven (Java)](#maven-java)
  - [NPM (Javascript / Typescript)](#npm-javascript--typescript)
  - [Nuget (.NET)](#nuget-net)
  - [PyPi (Python)](#pypi-python)
- [Development](#development)
- [The Fine Print](#the-fine-print)

## Installation

Obtain the binary for your Operating System and Architecture from the [GitHub Releases page](https://github.com/sonatype-nexus-community/nxfw-policy-tester/releases).

## Usage

Set your Sonatype Nexus Repository credentials in two environment variables:

-   `NEXUS_USERNAME`
-   `NEXUS_PASSWORD`

```bash
./nxfw-policy-tester
```

Follow the prompts - you'll need the URL to your Sonatype Nexus Repository installation (https:// only supported).

## Test Data Available

Test data is aimed to validate the [Sonatype Reference Policy Set](https://help.sonatype.com/en/reference-policies.html).

If you have customised or custom Policies, plesae consider this.

> NOTES:
>
> ^ Does not rely on actually malicious package(s) for verification - uses Sonatype staged packages marked as Malicious.
>
> ± Includes packages in both Pending and Suspicious states, where staged test data is available.
>
> § See [Firewall Specific Policies](https://help.sonatype.com/en/security-policies.html#firewall-specific-policies) - unrealistic to test in a generic manner.

### Cargo (Rust)

| Policy Type | Reference Policy               | Available            |
| ----------- | ------------------------------ | -------------------- |
| Legal       | `License-Banned`               | ❌                   |
| Legal       | `License-None`                 | ❌                   |
| Legal       | `License-Copyleft`             | ❌                   |
| Legal       | `License-Threat Not Assigned`  | ❌                   |
| Legal       | `License-AI-ML`                | N/A                  |
| Legal       | `License-Non-Standard`         | ❌                   |
| Legal       | `License-Weak-Copyleft`        | ❌                   |
| Security    | `Security-Namespace Conflict`§ | ❌                   |
| Security    | `Security-Malicious`           | ⛔️ No safe testdata |
| Security    | `Integrity-Rating`             | ⛔️ No safe testdata |
| Security    | `Security-Critical`            | ✅                   |
| Security    | `Security-High`                | ✅                   |
| Security    | `Security-Medium`              | ❌                   |
| Security    | `Security-Low`                 | ❌                   |

### Conda (conda-forge)

| Policy Type | Reference Policy               | Available            |
| ----------- | ------------------------------ | -------------------- |
| Legal       | `License-Banned`               | ❌                   |
| Legal       | `License-None`                 | ❌                   |
| Legal       | `License-Copyleft`             | ❌                   |
| Legal       | `License-Threat Not Assigned`  | ❌                   |
| Legal       | `License-AI-ML`                | N/A                  |
| Legal       | `License-Non-Standard`         | ❌                   |
| Legal       | `License-Weak-Copyleft`        | ❌                   |
| Security    | `Security-Namespace Conflict`§ | ❌                   |
| Security    | `Security-Malicious`           | ⛔️ No safe testdata |
| Security    | `Integrity-Rating`             | ⛔️ No safe testdata |
| Security    | `Security-Critical`            | ✅                   |
| Security    | `Security-High`                | ✅                   |
| Security    | `Security-Medium`              | ❌                   |
| Security    | `Security-Low`                 | ✅                   |

### CRAN (R)

| Policy Type | Reference Policy               | Available            |
| ----------- | ------------------------------ | -------------------- |
| Legal       | `License-Banned`               | ❌                   |
| Legal       | `License-None`                 | ❌                   |
| Legal       | `License-Copyleft`             | ❌                   |
| Legal       | `License-Threat Not Assigned`  | ❌                   |
| Legal       | `License-AI-ML`                | N/A                  |
| Legal       | `License-Non-Standard`         | ❌                   |
| Legal       | `License-Weak-Copyleft`        | ❌                   |
| Security    | `Security-Namespace Conflict`§ | ❌                   |
| Security    | `Security-Malicious`           | ⛔️ No safe testdata |
| Security    | `Integrity-Rating`             | ⛔️ No safe testdata |
| Security    | `Security-Critical`            | ❌                   |
| Security    | `Security-High`                | ✅                   |
| Security    | `Security-Medium`              | ✅                   |
| Security    | `Security-Low`                 | ❌                   |

### Golang (Go)

| Policy Type | Reference Policy               | Available            |
| ----------- | ------------------------------ | -------------------- |
| Legal       | `License-Banned`               | ❌                   |
| Legal       | `License-None`                 | ❌                   |
| Legal       | `License-Copyleft`             | ❌                   |
| Legal       | `License-Threat Not Assigned`  | ❌                   |
| Legal       | `License-AI-ML`                | N/A                  |
| Legal       | `License-Non-Standard`         | ❌                   |
| Legal       | `License-Weak-Copyleft`        | ❌                   |
| Security    | `Security-Namespace Conflict`§ | ❌                   |
| Security    | `Security-Malicious`           | ⛔️ No safe testdata |
| Security    | `Integrity-Rating`             | ⛔️ No safe testdata |
| Security    | `Security-Critical`            | ❌                   |
| Security    | `Security-High`                | ✅                   |
| Security    | `Security-Medium`              | ✅                   |
| Security    | `Security-Low`                 | ❌                   |

### Huggingface.co (AI / ML)

| Policy Type | Reference Policy               | Available |
| ----------- | ------------------------------ | --------- |
| Legal       | `License-Banned`               | ✅        |
| Legal       | `License-None`                 | ❌        |
| Legal       | `License-Copyleft`             | ❌        |
| Legal       | `License-Threat Not Assigned`  | ❌        |
| Legal       | `License-AI-ML`                | ❌        |
| Legal       | `License-Non-Standard`         | ❌        |
| Legal       | `License-Weak-Copyleft`        | ❌        |
| Security    | `Security-Namespace Conflict`§ | ❌        |
| Security    | `Security-Malicious`^          | ✅        |
| Security    | `Integrity-Rating`±            | ✅        |
| Security    | `Security-Critical`            | ❌        |
| Security    | `Security-High`                | ❌        |
| Security    | `Security-Medium`              | ❌        |
| Security    | `Security-Low`                 | ❌        |

### Maven (Java)

| Policy Type | Reference Policy               | Available |
| ----------- | ------------------------------ | --------- |
| Legal       | `License-Banned`               | ❌        |
| Legal       | `License-None`                 | ❌        |
| Legal       | `License-Copyleft`             | ❌        |
| Legal       | `License-Threat Not Assigned`  | ❌        |
| Legal       | `License-AI-ML`                | N/A       |
| Legal       | `License-Non-Standard`         | ❌        |
| Legal       | `License-Weak-Copyleft`        | ❌        |
| Security    | `Security-Namespace Conflict`§ | ❌        |
| Security    | `Security-Malicious`^          | ✅        |
| Security    | `Integrity-Rating`±            | ✅        |
| Security    | `Security-Critical`            | ✅        |
| Security    | `Security-High`                | ✅        |
| Security    | `Security-Medium`              | ✅        |
| Security    | `Security-Low`                 | ✅        |

### NPM (Javascript / Typescript)

| Policy Type | Reference Policy               | Available |
| ----------- | ------------------------------ | --------- |
| Legal       | `License-Banned`               | ✅        |
| Legal       | `License-None`                 | ❌        |
| Legal       | `License-Copyleft`             | ❌        |
| Legal       | `License-Threat Not Assigned`  | ❌        |
| Legal       | `License-AI-ML`                | N/A       |
| Legal       | `License-Non-Standard`         | ❌        |
| Legal       | `License-Weak-Copyleft`        | ❌        |
| Security    | `Security-Namespace Conflict`§ | ❌        |
| Security    | `Security-Malicious`^          | ✅        |
| Security    | `Integrity-Rating`±            | ✅        |
| Security    | `Security-Critical`            | ✅        |
| Security    | `Security-High`                | ✅        |
| Security    | `Security-Medium`              | ✅        |
| Security    | `Security-Low`                 | ✅        |

### Nuget (.NET)

| Policy Type | Reference Policy               | Available            |
| ----------- | ------------------------------ | -------------------- |
| Legal       | `License-Banned`               | ❌                   |
| Legal       | `License-None`                 | ❌                   |
| Legal       | `License-Copyleft`             | ❌                   |
| Legal       | `License-Threat Not Assigned`  | ❌                   |
| Legal       | `License-AI-ML`                | N/A                  |
| Legal       | `License-Non-Standard`         | ❌                   |
| Legal       | `License-Weak-Copyleft`        | ❌                   |
| Security    | `Security-Namespace Conflict`§ | ❌                   |
| Security    | `Security-Malicious`           | ⛔️ No safe testdata |
| Security    | `Integrity-Rating`             | ⛔️ No safe testdata |
| Security    | `Security-Critical`            | ✅                   |
| Security    | `Security-High`                | ✅                   |
| Security    | `Security-Medium`              | ✅                   |
| Security    | `Security-Low`                 | ✅                   |

### PyPi (Python)

| Policy Type | Reference Policy               | Available |
| ----------- | ------------------------------ | --------- |
| Legal       | `License-Banned`               | ❌        |
| Legal       | `License-None`                 | ❌        |
| Legal       | `License-Copyleft`             | ❌        |
| Legal       | `License-Threat Not Assigned`  | ❌        |
| Legal       | `License-AI-ML`                | N/A       |
| Legal       | `License-Non-Standard`         | ❌        |
| Legal       | `License-Weak-Copyleft`        | ❌        |
| Security    | `Security-Namespace Conflict`§ | ❌        |
| Security    | `Security-Malicious`^          | ✅        |
| Security    | `Integrity-Rating`±            | ✅        |
| Security    | `Security-Critical`            | ✅        |
| Security    | `Security-High`                | ✅        |
| Security    | `Security-Medium`              | ✅        |
| Security    | `Security-Low`                 | ✅        |

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
