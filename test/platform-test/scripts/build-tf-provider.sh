#!/usr/bin/env bash
# Copyright (C) 2015 The Gravitee team (http://gravitee.io)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build the Gravitee Terraform provider from source at a given git ref and
# stage it into a Terraform filesystem-mirror, so the platform-test e2e suite
# can exercise unreleased provider features (anything merged to main but not
# yet tagged + published to registry.terraform.io).
#
# WHY THIS EXISTS
# ---------------
# Our TF fixtures declare `source = "gravitee-io/apim"` with no version pin,
# so `terraform init` normally downloads the latest *registry* release. The
# provider is released manually (workflow_dispatch + goreleaser) and main
# routinely sits multiple unreleased versions ahead of the last tag, so any
# test of a recently-merged feature would otherwise be blocked on release.
#
# This script bypasses that gap by building the provider locally and putting
# the binary in the registry-mirror layout that TF understands. The companion
# wiring in helpers/terraform.ts picks up the generated .terraformrc via the
# TF_CLI_CONFIG_FILE env var (the standard TF-CLI override mechanism).
#
# Usage:
#   build-tf-provider.sh [REF]
#
# REF defaults to the TF_PROVIDER_REF env var, or "main" if neither is set.
# REF can be a branch, tag, or commit SHA.
#
# Prereqs:
#   - go (>= the provider's go.mod minimum)
#   - git
#
# Idempotency:
#   - Re-running overwrites the mirror layout for the same target version.
#   - Different host OS/arch entries coexist, so a darwin_arm64 dev box and
#     a linux_amd64 CI runner can share the same mirror root.

set -euo pipefail

REF="${1:-${TF_PROVIDER_REF:-main}}"
MIRROR="${TF_PROVIDER_MIRROR:-${HOME}/.terraform.d/gko-e2e-mirror}"

# Synthetic version that's higher than any real release. TF picks the highest
# non-prerelease version it can see, so this guarantees the local build wins
# over any registry version when fixtures don't pin a version. Picked high
# enough that it will never collide with a real release.
TARGET_VERSION="99.0.0"

WORK="$(mktemp -d -t tf-provider-build-XXXX)"
trap 'rm -rf "$WORK"' EXIT

echo "[build-tf-provider] cloning gravitee-io/terraform-provider-apim @ ${REF}"
# Shallow clone of the default branch, then fetch the requested ref. Two-step
# because a single --branch=<sha> doesn't work for arbitrary SHAs on most
# git servers.
git clone --depth 50 --no-tags https://github.com/gravitee-io/terraform-provider-apim.git "$WORK/repo" >/dev/null
(
  cd "$WORK/repo"
  if git rev-parse --verify --quiet "${REF}^{commit}" >/dev/null; then
    # REF already present locally (the default branch, or a SHA within the
    # shallow history) — check it out by name.
    git checkout --quiet "$REF"
  else
    # REF is a tag, non-default branch, or out-of-history SHA. An explicit
    # bare-ref fetch only updates FETCH_HEAD: with --no-tags it creates no
    # local refs/tags/<tag> and no remote-tracking branch, so `git checkout
    # "$REF"` would fail. Check out FETCH_HEAD instead, which holds exactly
    # the fetched commit regardless of whether REF is a tag/branch/SHA.
    git fetch --depth 50 --no-tags origin "$REF" >/dev/null
    git checkout --quiet FETCH_HEAD
  fi
)

GOOS="$(go env GOOS)"
GOARCH="$(go env GOARCH)"
BIN_NAME="terraform-provider-apim_v${TARGET_VERSION}"
TARGET_DIR="${MIRROR}/registry.terraform.io/gravitee-io/apim/${TARGET_VERSION}/${GOOS}_${GOARCH}"

echo "[build-tf-provider] go build -> ${TARGET_DIR}/${BIN_NAME}"
mkdir -p "$TARGET_DIR"
(
  cd "$WORK/repo"
  go build -o "${TARGET_DIR}/${BIN_NAME}" .
)
chmod +x "${TARGET_DIR}/${BIN_NAME}"

# Write the CLI config. `direct { exclude = ... }` keeps the public registry
# in play for any other provider a future fixture might add.
TF_CONFIG="${MIRROR}/.terraformrc"
cat > "$TF_CONFIG" <<EOF
provider_installation {
  filesystem_mirror {
    path    = "${MIRROR}"
    include = ["registry.terraform.io/gravitee-io/apim"]
  }
  direct {
    exclude = ["registry.terraform.io/gravitee-io/apim"]
  }
}
EOF

# Stamp a marker the test harness can use to attribute drift / behaviour to
# the local build rather than the registry release.
COMMIT_SHA="$(git -C "$WORK/repo" rev-parse --short HEAD)"
cat > "${MIRROR}/build-info.json" <<EOF
{
  "ref": "${REF}",
  "commit": "${COMMIT_SHA}",
  "target_version": "${TARGET_VERSION}",
  "host": "${GOOS}_${GOARCH}",
  "built_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
}
EOF

echo "[build-tf-provider] mirror ready: ${MIRROR}"
echo "[build-tf-provider] export TF_CLI_CONFIG_FILE=${TF_CONFIG}"
echo "[build-tf-provider] built provider ref ${REF} (${COMMIT_SHA}) as ${TARGET_VERSION} for ${GOOS}_${GOARCH}"
