---
spec_version: '2.0'
task_id: auscaster-metadata
created: '2026-10-29T01:18:17Z'
updated: '2026-10-29T05:28:06Z'
status: completed
harden_status: passed
size: small
risk_level: low
---

# Migrate public metadata and registry publishing to auscaster

## Summary

Move scafld's public project and artifact links to the transferred auscaster
repository and auscaster.com. Use repository-bound OIDC for the PyPI publisher.
The authorized metadata release will validate the npm and PyPI bindings after
this source change is merged; the migration register owns that live readback.

## Objectives

- Published wrapper metadata and generated package manifests identify auscaster.
- Installer defaults download artifacts from auscaster/scafld.
- The release workflow publishes npm and PyPI through exact repository-bound OIDC.
- Existing Go imports and installed package identifiers keep working.

## Scope

- .github/workflows/release.yml
- README.md
- docs/
- install.sh
- package/
- scripts/
- test/release/packaging_test.go

## Dependencies

- GitHub transfer to auscaster/scafld is verified with the unchanged repository ID.
- Registry bindings name auscaster/scafld, release.yml, and no environment name.

## Assumptions

- Package names and the v2 Go module identity remain stable.
- The existing artifact builder, verified-artifact promotion, and registry jobs remain the release owners.

## Touchpoints

- .github/workflows/release.yml
- docs/sourcey.config.ts
- install.sh
- package/npm/lib/platform.js
- package/pypi/src/scafld_launcher/install.py
- package/homebrew/
- package/scoop/
- package/winget/
- package/docker/
- scripts/build-release-artifacts.sh
- scripts/audit-scafld-metadata.sh
- test/release/packaging_test.go

## Risks

- Cross-repository writer credentials require provider verification after ownership changes.
- A pending trusted publisher must be exercised before its activation deadline.
- Preserve the Go module path, schema IDs, historical copyright, and the registered WinGet package identifier.

## Acceptance

Profile: standard

## Phase 1: Update public identities and registry publishing

Status: completed
Dependencies: none

Objective: Publish metadata from the existing release boundaries with the new brand and repository owner.

Changes:
- Update README.md, docs/, install.sh, package/, and scripts/ public links and artifact download defaults.
- Update .github/workflows/release.yml to use transferred package repositories, the auscaster image namespace, and PyPI OIDC.
- Update test/release/packaging_test.go expectations while retaining the stable v2 module contract.

Acceptance:
- [x] `ac1` command - Existing checks validate the runtime, packaging, and release contracts
  - Command: `make check`
  - Expected kind: `exit_code_zero`
  - Status: pass
  - Evidence: exit code was 0
  - Source event: entry-35

## Rollback

- Revert this source change if verification fails; retain published versions and old URL redirects.
- Preserve the old registry binding until the replacement has a successful publication and provider readback.

## Review

Status: completed
Verdict: pass
Mode: verify
Provider: claude:claude-opus-5-5
Output: claude.mcp_submit_review
Summary: No defects found; the review passes. The release workflow now publishes to PyPI through OIDC: `password: ${{ secrets.PYPI_API_TOKEN }}` is gone (main checkout had it at release.yml:114), and `id-token: write` was added with no environment set. The npm job was already on OIDC. All the places the old owner name used to appear now say auscaster: the three installers, the Homebrew/Scoop/WinGet templates, the Dockerfile, the image tags, the tap and bucket checkouts, the checksum URL used by render, the audit script, the docs config, README and docs. The things the spec said to keep are kept: the Go module path `github.com/nilstate/scafld/v2` (go.mod, build ldflags, audit module, test), the schema `$id`, the WinGet `0state.scafld` identifier and the 0state copyright line. A repo-wide search found no leftover `nilstate/tap|homebrew|scoop`, `ghcr.io/nilstate`, non-module `github.com/nilstate/scafld` or `0state.com` references. The only exception is context7.json, which is outside this task's scope. No open findings were carried into this verify pass, and the fix introduced no new blockers.

Attack log:
- `.github/workflows/release.yml:86-115 vs main checkout release.yml:86-115 and docs/release.md:78-80`: Contract drift: does the PyPI job publish through exact repository-bound OIDC as the spec requires? -> clean (The baseline passed `password: ${{ secrets.PYPI_API_TOKEN }}`; the worktree removes it and adds `id-token: write # PyPI trusted publishing`. No `environment:` key is set, which matches the spec's 'no environment name' binding and the docs' claim that trust names auscaster/scafld + release.yml. The npm job keeps its OIDC permission and has no NPM_TOKEN or NODE_AUTH_TOKEN.)
- `go.mod:1, scripts/build-release-artifacts.sh:41,89, scripts/audit-scafld-metadata.sh:17,154, test/release/packaging_test.go:20,454,561, internal/adapters/corebundle/assets/core/schemas/spec.json:3, package/winget/*.tmpl:2`: Regression: was any stable import or package identity renamed along with the GitHub owner? -> clean (Kept unchanged: the module path, the ldflags `-X github.com/nilstate/scafld/v2/...` (still matches go.mod, so version stamping works), the deps.dev lookup, the schema $id, and the WinGet PackageIdentifier 0state.scafld. Tests still enforce both the module path and the WinGet identifier.)
- `rg for `github.com/nilstate/scafld[^/]`, `ghcr.io/nilstate`, `nilstate/(tap|homebrew|scoop)`, `0state.com` across the worktree`: Final product shape: search for leftover old-owner public identities across all files, including untracked ones -> clean (No matches outside the preserved module/schema identities. The remaining 0state/nilstate hits are all intentional: go.mod, LICENSE copyright, the WinGet identifier, the build-resource LegalCopyright, the scafld-verify.sh go install path (module), and plans/positioning.md, which is internal. context7.json is outside the declared scope and is an external library ID.)
- `install.sh:4, package/npm/lib/platform.js:46, package/pypi/src/scafld_launcher/install.py:51,58`: Installer default consistency: do all three download paths agree on the new owner and keep the same override? -> clean (All three default to auscaster/scafld behind the same SCAFLD_GITHUB_REPOSITORY override, and the SCAFLD_INSTALL_BASE_URL local path is untouched. That means the smoke-release-installers flow against dist/ is unaffected.)
- `package/homebrew/scafld.rb.tmpl:3-22, package/scoop/scafld.json.tmpl:4-27, package/winget/scafld.installer.yaml.tmpl:10,13, package/docker/Dockerfile:8,14, scripts/render-package-managers.sh:33, release.yml:308-309`: Package-manager render path: do templates, checksum fetch, and the OCI image all point at the same release assets? -> clean (Every asset URL uses github.com/auscaster/scafld/releases/download/v{{VERSION}}, and render fetches checksums from that same release. Scoop checkver changed from the bare `"github"` (which infers the repo from the now non-GitHub homepage) to an explicit GitHub URL, which is valid Scoop syntax. render still runs ruby -c, jq empty, and the leftover-marker check.)
- `package/npm/package.json:8-11 vs release.yml publish-npm job`: npm trusted-publishing provenance: does package.json repository.url match the repository that will mint the OIDC token? -> clean (repository.url is git+https://github.com/auscaster/scafld.git. That matches the transferred repository the workflow runs from, so the provenance repo check will not reject the publish. The baseline nilstate URL would have mismatched after the transfer.)
- `test/release/packaging_test.go:400-438, 518-535`: Testing gaps: would the tests catch a revert to token-based publishing or stale template owners? -> clean (The tests forbid PYPI_API_TOKEN, NPM_TOKEN and NODE_AUTH_TOKEN, and require the PyPI and npm id-token comments plus the auscaster tap/bucket fail-closed messages. They also require auscaster release URLs in the Homebrew, Scoop and WinGet installer templates. make check (ac1) passed.)
- `release.yml:184-191,232-239 and launcher defaults`: Configured invariants: secrets in source, legacy fallbacks, or dual paths introduced? -> clean (No secrets were added to source, and no dual-read or old-owner fallback URLs were added to the installers. The Homebrew and Scoop writers still fail closed on missing channel-specific secrets, with messages updated to the new repositories.)
- `scripts/build-release-artifacts.sh:41,63,67 and package/winget/scafld.yaml.tmpl:2-10`: Windows resource and WinGet publisher metadata: are copyright and identity kept as the spec requires? -> clean (CompanyName, Publisher and PublisherUrl moved to auscaster. The manifest identity `0state.scafld`, the PackageIdentifier and the historical `Copyright (c) 0state` are kept, as the spec's Risks section requires.)

Findings:
- none

## Origin

Created by: scafld
Source: plan

## Current State

Status: completed
Current phase: none
Next: none
Reason: finalization receipt passed
Blockers: none
Allowed follow-up command: `none`
Latest runner update: 2026-10-29T05:28:06Z
Review gate: not_started

## Harden Rounds

### round-1

Status: passed
Started: 2026-10-29T01:21:58Z
Ended: 2026-10-29T01:23:18Z
Spec digest: bf0d7bfbc9b1e02c43a702dbbd254fa3fec61d67cb96a030f14dbb83adbe37bb
Shape decision: keep
True shape: A bounded metadata and registry-identity cutover through the existing release owners. No-op leaves active public identities stale; shrinking away installer defaults leaves future packages bound to the old owner; a new deployment or package layer is unnecessary because the existing builders and publisher jobs already own these effects.
Minimal plan: Update the public identities and existing publisher jobs, retain stable import and package IDs, run the existing checks, and validate the new registry bindings through the authorized changed release after merge.
Shared owner: .github/workflows/release.yml and its existing artifact builder and package templates.
Adapter boundaries: npm and PyPI remain thin verified-native-binary launchers; package-manager templates render the existing release artifacts; docs describe their owning configuration.
Required spec edits: 

Observations:
- design
  - Result: clean
  - Anchor: code:.github/workflows/release.yml:49
  - Note: The existing job builds and smokes the exact version, uploads it, and the GitHub release job restores that artifact. Reuse this owner; another builder or publisher would duplicate it.
- scope
  - Result: clean
  - Anchor: code:go.mod:1
  - Note: The v2 module identity remains stable while GitHub repository URLs change. The registered WinGet identifier and immutable schema identities are retained; this would fail if the migration renamed an import contract or published package ID.
- path
  - Result: clean
  - Anchor: code:package/npm/lib/platform.js:46
  - Note: The npm launcher already owns its download repository default. Matching Python and shell defaults, package templates, docs config, and the release workflow were inspected and are in scope; no new runtime adapter is needed.
- command
  - Result: clean
  - Anchor: code:Makefile:42
  - Note: make check is the existing root command for formatting, architecture, packaging, tests and race checks. It ran successfully on the prepared metadata patch before the clean approval baseline was restored.
- timing
  - Result: clean
  - Anchor: code:.github/workflows/release.yml:46
  - Note: Source acceptance runs after applying the patch and before release. Live registry validation is a downstream release checkpoint, so the source change does not claim an unexecuted publication.
- rollback
  - Result: clean
  - Anchor: code:.github/workflows/release.yml:63
  - Note: Source repair is a revert or forward correction. Existing published artifacts and old URL redirects remain available, and the old registry binding stays until the new binding is valid.


