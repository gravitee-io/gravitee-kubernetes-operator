# e2e CI lane configuration

This directory holds the cluster configuration for the `e2e` CircleCI job.
The e2e test suite itself lives in [`test/platform-test/`](../platform-test/).

- `coverage.yaml` — PersistentVolume/Claim (kind hostPath `/tmp/coverage`) that
  collects Go coverage from the operator binary while the e2e suite runs. The
  job converts it to `cover-e2e.out`, which feeds the SonarCloud analysis (see
  `sonar.go.coverage.reportPaths` in `sonar-project.properties`).
- `operator.values.yaml` — Helm values mounting the coverage volume into the
  operator manager pod.

The `conformance` CI lane keeps its equivalent pair under `test/conformance/`.
