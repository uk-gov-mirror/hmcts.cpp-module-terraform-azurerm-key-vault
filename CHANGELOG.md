# Changelog

All notable changes to this Terraform module are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.1] - 2026-07-03

## :bug: Bug Fixes

- fix: Improve module logic to deduplicate default roles where possible @dawidstrozak (#16)

## [2.0.0] - 2026-06-29

- DTSPO-31407: add SemVer versioning and automated release drafting @manohar-hmcts (#11)

## :boom: Breaking Changes

- Fix release-drafter workflow - copy in the latest release drafter workflow template @dawidstrozak (#14)
- feat!: Enable RBAC Authorisation by default in this module @dawidstrozak (#12)

## :bug: Bug Fixes

- fix: Adjust changelog action reference from local file to common actions repository @dawidstrozak (#15)

## [1.0.0] - 2026-06-23

### Added

- Initial tagged release of the `cpp-module-terraform-azurerm-key-vault` module.
- Semantic Versioning with automated release drafting (Release Drafter) and
  automatic CHANGELOG updates via the shared HMCTS update-changelog action.
