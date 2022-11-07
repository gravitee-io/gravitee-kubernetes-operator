const branch = process.env.CIRCLE_BRANCH;
const prereleaseBranchName = config.branches.some(
  (it) => it === branch || (it.name === branch && !it.prerelease)
)
  ? branch
  : null;

const config = {
  branches: [{ name: "alpha", prerelease: true }, "master"],
  tagFormat: "${version}",
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/changelog",
    [
      "@semantic-release/exec",
      {
        prepareCmd:
          "IMG=graviteeio/kubernetes-operator:${nextRelease.version} BUNDLE_IMG=graviteeio/kubernetes-operator-bundle:${nextRelease.version} make docker-build docker-push bundle-standalone bundle-build bundle-push lint-fix",
      },
    ],
    [
      "@semantic-release/github",
      {
        assets: [{ path: "bundle.yml", label: "Operator resources bundle" }],
      },
    ],
    [
      "@semantic-release/git",
      {
        assets: [
          `CHANGELOG${
            prereleaseBranchName ? `_${prereleaseBranchName}` : ""
          }.md`,
        ],
        message: "chore(release): ${nextRelease.version} [skip ci]",
      },
    ],
  ],
};

module.exports = config;
