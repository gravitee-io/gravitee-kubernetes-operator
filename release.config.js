const config = {
  branches: [{ name: "alpha", prerelease: true }, "master"],
  tagFormat: "${version}",
};

const branch = process.env.CIRCLE_BRANCH;
const changelogFileName =
  branch === "master" ? "CHANGELOG.md" : `CHANGELOG-${branch.toUpperCase()}.md`;

const plugins = [
  "@semantic-release/commit-analyzer",
  "@semantic-release/release-notes-generator",
  [
    "@semantic-release/changelog",
    {
      changelogFile: changelogFileName,
    },
  ],
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
      assets: [changelogFileName],
      message: "chore(release): ${nextRelease.version} [skip ci]",
    },
  ],
];

module.exports = { ...config, plugins };
