# Sonatype Repository Firewall Policy Test Suite

<!-- Badges Section -->
[![shield_gh-workflow-test]][link_gh-workflow-test]
[![shield_license]][license_file]
<!-- Add other badges or shields as appropriate -->

---

This is a simple test suite that will attempt to pull known *bad* packages into Proxy repositores in Sonatype Nexus Repository 
to allow customers to validate that Repository Firewall policies are operating as they expect.

This tool does not know the specific outcomes, as each Customer will have differing [Policy Actions](https://help.sonatype.com/en/policy-actions.html)
set.

This tool expects the [Reference Policy Set](https://help.sonatype.com/en/reference-policies.html) to be in use, but that is not mandatory.

- [Usage](#usage)
- [Development](#development)
- [The Fine Print](#the-fine-print)

## Usage

Currently this works for NPM format only.

```bash
./test-npm.sh
```

## Development

See [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## The Fine Print

Remember:

This project is part of the [Sonatype Nexus Community](https://github.com/sonatype-nexus-community) organization, which is not officially supported by Sonatype. Please review the latest pull requests, issues, and commits to understand this project's readiness for contribution and use.

* File suggestions and requests on this repo through GitHub Issues, so that the community can pitch in
* Use or contribute to this project according to your organization's policies and your own risk tolerance
* Don't file Sonatype support tickets related to this project— it won't reach the right people that way

Last but not least of all - have fun!

<!-- Links Section -->
[shield_gh-workflow-test]: https://img.shields.io/github/actions/workflow/status/sonatype-nexus-community/community-project-template/test.yml?branch=main&logo=GitHub&logoColor=white "build"
[shield_license]: https://img.shields.io/github/license/sonatype-nexus-community/community-project-template?logo=open%20source%20initiative&logoColor=white "license"

[link_gh-workflow-test]: https://github.com/sonatype-nexus-community/community-project-template/actions/workflows/test.yml?query=branch%3Amain
[license_file]: https://github.com/sonatype-nexus-community/community-project-template/blob/main/LICENSE