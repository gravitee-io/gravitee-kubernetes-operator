/**
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { LOG, PROJECT_DIR } from "./index.mjs";

// Where generated reports and their submission metadata live.
export const REPORT_DIR = path.join(
  PROJECT_DIR,
  "test",
  "conformance",
  "kubernetes.io",
  "gateway-api",
  "report",
);

export const FORK_REPO = "gravitee-io-labs/gateway-api";
export const UPSTREAM_REPO = "kubernetes-sigs/gateway-api";

// The implementations list moved from site-src/implementations.md to a Docsy
// layout in kubernetes-sigs/gateway-api#4734 and may well move again. Probe the
// known locations rather than hard coding one, and fail loudly when none match
// so a silent no-op submission is impossible.
const IMPLEMENTATIONS_LIST_CANDIDATES = [
  "site/content/en/docs/implementations/list.md",
  "site-src/implementations.md",
];

export function reportPath(version) {
  return path.join(REPORT_DIR, `standard-${version}-default-report.yaml`);
}

export async function loadSubmissionMeta() {
  const raw = await fs.readFile(
    path.join(REPORT_DIR, "submission.yaml"),
    "utf8",
  );
  return YAML.parse(raw);
}

export async function loadReport(version) {
  const file = reportPath(version);
  if (!fs.pathExistsSync(file)) {
    throw new Error(
      `No conformance report for version ${version} at ${file}.\n` +
        `Run the suite for that version first, or pass a version that has one:\n` +
        (await fs.readdir(REPORT_DIR))
          .filter((f) => f.endsWith("-default-report.yaml"))
          .map(
            (f) => `  - ${f.replace(/^standard-|-default-report\.yaml$/g, "")}`,
          )
          .join("\n"),
    );
  }
  return YAML.parse(await fs.readFile(file, "utf8"));
}

/**
 * Reduce a generated report to everything the rendered documents need.
 * Only the prose lives in submission.yaml; versions, mode, channel and the
 * conformance verdict are read back off the report so they cannot drift.
 */
export function describeReport(report, meta) {
  const profile = (report.profiles ?? [])[0];

  if (!profile) {
    throw new Error("Report contains no conformance profile");
  }

  const levels = ["core", "extended"].map((l) => profile[l]).filter(Boolean);
  const skipped = levels.flatMap((l) => l.skippedTests ?? []);
  const failed = levels.reduce((n, l) => n + (l.statistics?.Failed ?? 0), 0);
  const passed = levels.reduce((n, l) => n + (l.statistics?.Passed ?? 0), 0);

  if (failed > 0) {
    throw new Error(
      `Report has ${failed} failing test(s). A report with failures is not submittable.`,
    );
  }

  const unexplained = skipped.filter((t) => !meta.skips?.[t]);
  if (unexplained.length > 0) {
    throw new Error(
      `Skipped tests with no explanation in submission.yaml: ${unexplained.join(", ")}.\n` +
        `Add an entry under 'skips:' phrased to read inside "It does not support <phrase>".`,
    );
  }

  const profileLabel = meta.profiles?.[profile.name] ?? profile.name;

  return {
    version: report.implementation.version,
    gatewayAPIVersion: report.gatewayAPIVersion,
    channel: report.gatewayAPIChannel,
    mode: report.mode,
    profileName: profile.name,
    profileLabel,
    conformant: skipped.length === 0,
    skipped,
    passed,
    failed,
  };
}

function joinPhrases(phrases) {
  if (phrases.length === 1) {
    return phrases[0];
  }
  return `${phrases.slice(0, -1).join(", ")} and ${phrases[phrases.length - 1]}`;
}

/**
 * The one sentence that states our conformance level. It appears verbatim in
 * both the report README and the upstream implementations list, so it is
 * generated once here.
 */
export function conformanceStatement(desc, meta) {
  const subject = `The ${meta.name} provides`;

  if (desc.conformant) {
    return (
      `${subject} full conformance for ${desc.profileLabel} features ` +
      `in version ${desc.version}.`
    );
  }

  const phrases = desc.skipped.map((t) => meta.skips[t]);
  const plural = phrases.length > 1 ? "These features" : "This feature";

  return (
    `${subject} partial conformance for ${desc.profileLabel} features ` +
    `in version ${desc.version}. It does not support ${joinPhrases(phrases)}. ` +
    `${plural} will be introduced in a future release.`
  );
}

/**
 * The README that ships next to the report, both in this repository and in the
 * upstream conformance/reports/<gatewayAPIVersion>/<directory>/ directory. The
 * two copies are byte identical by design.
 */
export function renderReadme(desc, meta) {
  const reportFile = `standard-${desc.version}-default-report.yaml`;
  const releaseURL = `${meta.releasesURL}/${desc.version}`;

  return `# Gravitee

## Table of Contents

| API channel  | Implementation version                    | Mode    | Report                                                 |
|--------------|-------------------------------------------|---------|--------------------------------------------------------|
| ${desc.channel.padEnd(12)} | [version-${desc.version}](${releaseURL}) | ${desc.mode.padEnd(7)} | [version-${desc.version} report](./${reportFile}) |

> ${conformanceStatement(desc, meta)}

## Prerequisites

The following binaries are assumed to be installed on your device

  - [docker](https://docs.docker.com/get-started/get-docker/)
  - [kubectl](https://kubernetes.io/docs/tasks/tools/)
  - [kind](https://github.com/kubernetes-sigs/kind)
  - [go](https://go.dev/learn/)

The reproducer has been tested on macOS and Linux only.

## Reproducer

1. Clone the Gravitee Kubernetes Operator repository

\`\`\`bash
git clone --depth 1 --branch ${desc.version} https://github.com/gravitee-io/gravitee-kubernetes-operator.git
\`\`\`

2. Start the Kubernetes cluster

\`\`\`bash
make start-conformance-cluster
\`\`\`

3. Run a local Load Balancer Service

> The make target runs [cloud-provider-kind](https://kind.sigs.k8s.io/docs/user/loadbalancer). If you are reproducing on a macOS device, the binary requires \`sudo\` privileges and you will be prompted for a password. For Linux devices, cloud-provider-kind will be run using Docker compose.

\`\`\`bash
make cloud-lb
\`\`\`

4. Run the operator

\`\`\`bash
make run
\`\`\`

5. Install the Gravitee GatewayClass

\`\`\`bash
kubectl apply -f ./test/conformance/gateway-class-parameters.report.yaml -f ./test/conformance/gateway-class.yaml
\`\`\`

6. Run the conformance tests

\`\`\`bash
make conformance
\`\`\`

7. Print report

\`\`\`bash
cat test/conformance/kubernetes.io/gateway-api/report/${reportFile}
\`\`\`
`;
}

function badge(desc, meta) {
  const level = desc.conformant ? "Conformance" : "Partial Conformance";
  const colour = desc.conformant ? "green" : "orange";
  const label = encodeURIComponent(
    `Gateway API ${level} ${desc.gatewayAPIVersion}`,
  ).replaceAll("%2F", "/");
  const impl = encodeURIComponent(meta.name);
  // The badge must link to the directory the report was actually committed to,
  // not to a patch-version symlink that may not resolve as a browsable path.
  const dir = desc.reportDirVersion ?? desc.gatewayAPIVersion;
  const reports = `https://github.com/${UPSTREAM_REPO}/blob/main/conformance/reports/${dir}/${meta.directory}`;

  return `[![Conformance](https://img.shields.io/badge/${label}-${impl}-${colour})](${reports})`;
}

/** Our `### <name>` section in the upstream implementations list. */
export function renderImplementationsSection(desc, meta) {
  return [
    `### ${meta.name}`,
    "",
    badge(desc, meta),
    "",
    meta.description,
    "",
    conformanceStatement(desc, meta),
    "",
    meta.support,
    "",
  ].join("\n");
}

/**
 * Which directory under conformance/reports/ this report belongs in.
 *
 * Upstream keeps one real directory per minor (`v1.6`) and symlinks each patch
 * version at it (`v1.6.0 -> v1.6`). Writing through the symlink lands in the
 * right place but records the real path in git, so a badge built from the
 * report's own gatewayAPIVersion would advertise a path that is not what was
 * committed. Resolve it against the checkout instead of guessing.
 */
export function resolveReportDirVersion(rootDir, gatewayAPIVersion) {
  const reports = path.join(rootDir, "conformance", "reports");
  const minor = gatewayAPIVersion.replace(/^(v\d+\.\d+)\.\d+$/, "$1");

  for (const candidate of [gatewayAPIVersion, minor]) {
    const dir = path.join(reports, candidate);
    if (fs.pathExistsSync(dir)) {
      return path.basename(fs.realpathSync(dir));
    }
  }

  // First implementation to report against a brand new Gateway API release.
  // The minor form is the current upstream convention, but say so out loud:
  // creating the wrong directory shape is worth a second pair of eyes.
  LOG.yellow(
    `  No directory for ${gatewayAPIVersion} under conformance/reports yet — ` +
      `creating ${minor}. Check that upstream still groups reports by minor.`,
  );
  return minor;
}

export function findImplementationsList(rootDir) {
  const found = IMPLEMENTATIONS_LIST_CANDIDATES.map((p) =>
    path.join(rootDir, p),
  ).find((p) => fs.pathExistsSync(p));

  if (!found) {
    throw new Error(
      `Could not find the implementations list in ${rootDir}. Tried:\n` +
        IMPLEMENTATIONS_LIST_CANDIDATES.map((p) => `  - ${p}`).join("\n") +
        `\nUpstream has moved it; add the new path to IMPLEMENTATIONS_LIST_CANDIDATES.`,
    );
  }

  return found;
}

const anchorOf = (name) =>
  `#${name
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "")}`;

const byName = (a, b) => a.toLowerCase().localeCompare(b.toLowerCase());

/** Index of the line matching `re`, searching only within [from, to). */
function lineIndex(lines, re, from = 0, to = lines.length) {
  for (let i = from; i < to; i++) {
    if (re.test(lines[i])) return i;
  }
  return -1;
}

/** End of the block started at `from`, i.e. the next heading at `level` or shallower. */
function blockEnd(lines, from, level) {
  const heading = new RegExp(`^#{1,${level}} `);
  const next = lineIndex(lines, heading, from + 1);
  return next === -1 ? lines.length : next;
}

/**
 * Move our entry into the right status list under
 * "## Gateway Controller Implementation Status", keeping the list alphabetical.
 * Conformant and Partially Conformant are separate lists and we must appear in
 * exactly one of them.
 */
function patchStatusIndex(lines, meta, ref, conformant) {
  const statusStart = lineIndex(
    lines,
    /^## Gateway Controller Implementation Status/,
  );
  if (statusStart === -1) {
    throw new Error(
      "Could not find '## Gateway Controller Implementation Status' in the implementations list",
    );
  }
  const statusEnd = blockEnd(lines, statusStart, 2);

  const entry = `- [${meta.name}][${ref}]`;
  const existing = new RegExp(`^- \\[${meta.name}\\]\\[\\d+\\]\\s*$`);

  // Drop any existing entry from both lists before re-inserting into the right one.
  for (let i = statusEnd - 1; i >= statusStart; i--) {
    if (existing.test(lines[i])) lines.splice(i, 1);
  }

  const wantedHeading = conformant
    ? "### Conformant"
    : "### Partially Conformant";
  const end = blockEnd(lines, statusStart, 2);
  let headingAt = lineIndex(
    lines,
    new RegExp(`^${wantedHeading}\\s*$`),
    statusStart,
    end,
  );

  // Upstream drops a status list entirely when it empties out. Recreate it.
  if (headingAt === -1) {
    const other = lineIndex(lines, /^### /, statusStart, end);
    const insertAt = other === -1 ? statusStart + 1 : blockEnd(lines, other, 3);
    lines.splice(insertAt, 0, wantedHeading, "", entry, "");
    return lines;
  }

  const listEnd = blockEnd(lines, headingAt, 3);
  let insertAt = listEnd;
  for (let i = headingAt + 1; i < listEnd; i++) {
    const m = lines[i].match(/^- \[([^\]]+)\]\[\d+\]/);
    if (m && byName(m[1], meta.name) > 0) {
      insertAt = i;
      break;
    }
    if (m) insertAt = i + 1;
  }

  lines.splice(insertAt, 0, entry);
  return lines;
}

/** Ensure a `[<ref>]:#anchor` definition exists, allocating a ref if needed. */
function ensureLinkRef(lines, anchor) {
  const defined = lineIndex(lines, new RegExp(`^\\[(\\d+)\\]:${anchor}\\s*$`));
  if (defined !== -1) {
    return { ref: lines[defined].match(/^\[(\d+)\]/)[1], lines };
  }

  const refs = lines
    .map((l) => l.match(/^\[(\d+)\]:#/))
    .filter(Boolean)
    .map((m) => Number(m[1]));

  if (refs.length === 0) {
    throw new Error("Could not find any link reference definitions to extend");
  }

  const ref = String(Math.max(...refs) + 1);
  const lastRefAt = lines.reduce(
    (last, l, i) => (/^\[\d+\]:#/.test(l) ? i : last),
    -1,
  );

  lines.splice(lastRefAt + 1, 0, `[${ref}]:${anchor}`);
  return { ref, lines };
}

/**
 * Replace our `### <name>` section under "## Implementations", or insert it in
 * alphabetical order when it is missing. It has been missing upstream since the
 * Docsy migration, so insertion is the expected path, not an edge case.
 */
function patchImplementationsSection(lines, meta, section) {
  const start = lineIndex(lines, /^## Implementations\s*$/);
  if (start === -1) {
    throw new Error(
      "Could not find '## Implementations' in the implementations list",
    );
  }
  const end = blockEnd(lines, start, 2);

  const ours = lineIndex(
    lines,
    new RegExp(`^### ${meta.name}\\s*$`),
    start,
    end,
  );
  const body = section.split("\n");

  if (ours !== -1) {
    lines.splice(ours, blockEnd(lines, ours, 3) - ours, ...body);
    return lines;
  }

  let insertAt = end;
  for (let i = start + 1; i < end; i++) {
    const m = lines[i].match(/^### (.+?)\s*$/);
    if (m && byName(m[1], meta.name) > 0) {
      insertAt = i;
      break;
    }
  }

  lines.splice(insertAt, 0, ...body);
  return lines;
}

/**
 * Apply every edit the upstream implementations list needs for a submission:
 * the status index entry, the link reference and the section body.
 */
export function patchImplementationsList(content, desc, meta) {
  let lines = content.split("\n");
  const anchor = anchorOf(meta.name);

  const { ref, lines: withRef } = ensureLinkRef(lines, anchor);
  lines = withRef;
  lines = patchStatusIndex(lines, meta, ref, desc.conformant);
  lines = patchImplementationsSection(
    lines,
    meta,
    renderImplementationsSection(desc, meta),
  );

  return lines.join("\n");
}
