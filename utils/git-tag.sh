#!/usr/bin/env bash

set -o nounset
set -o errexit

deduceBumpLevel() {
  commitMessage=${1:?sdeduceBumpLevel requires first argument to be the latest commit message}

  bumpLevel="nonknown"
  if [[ "${commitMessage}" =~ ^\s*\[(major|breaking)\] || "${commitMessage}" =~ ^\s*(major|breaking)\b ]]; then
    bumpLevel="major"
  elif [[ "${commitMessage}" =~ ^\s*\[(minor|feature)\] || "${commitMessage}" =~ ^\s*(minor|feature)\b ]]; then
    bumpLevel="minor"
  elif [[ "${commitMessage}" =~ ^\s*\[(patch|fix)\] || "${commitMessage}" =~ ^\s*(patch|fix)\b ]]; then
    bumpLevel="patch"
  elif [[ "${commitMessage}" =~ ^\s*\[(doc|ci)\] || "${commitMessage}" =~ ^\s*(doc|ci)\b ]]; then
    bumpLevel="nonfunctional"
  fi

  if [[ "${commitMessage}" =~ ^\s*\[no\s+release\] ]]; then
    bumpLevel="nonrelease"
  fi

  echo "${bumpLevel}"
}

semverBump() {
  currentVersion=${1:?semverBump requires first argument to be the current semver}
  bumpLevel=${2:?semverBump requires the second argument to be the bump level: major, breaking, minor, feature, patch, fix}

  if [[ "${bumpLevel}" == "major" || "${bumpLevel}" == "breaking" ]]; then
    echo $(echo -n "${currentVersion}" | awk -F. '{printf "%d.%d.%d", ($1 + 1), 0, 0}')
  elif [[ "${bumpLevel}" == "minor" || "${bumpLevel}" == "feature" ]]; then
    echo $(echo -n "${currentVersion}" | awk -F. '{printf "%d.%d.%d", $1, ($2 + 1), 0}')
  elif [[ "${bumpLevel}" == "patch" || "${bumpLevel}" == "fix" ]]; then
    echo $(echo -n "${currentVersion}" | awk -F. '{printf "%d.%d.%d", $1, $2, ($3 + 1)}')
  else
    echo $(echo -n "${currentVersion}")
  fi
}

lastTag=$(git describe --tags --abbrev=0 | sed 's/^v//')
echo "Previous release: v${lastTag}"

latestCommitSHA=$(git log --pretty=%h -1)
echo "Latest commit SHA: ${latestCommitSHA}"

latestCommitMessage=$(git log --format=%s -1 ${latestCommitSHA})
echo "Latest commit message: ${latestCommitMessage}"

# latestCommitDecorations=$(git log --format=%s -1 ${latestCommitSHA})
# echo "Latest commit decorations: ${latestCommitDecorations}"

deducedBumpLevel=$(deduceBumpLevel ${latestCommitMessage})
echo "Deduced bump level: ${deducedBumpLevel}"

# TODO: Should not proceed if latestCommitDecorations contains a tag
if [[ "${deducedBumpLevel}" =~ ^non ]]; then
  echo "Skipping release (${deducedBumpLevel})"
else
  nextTag=$(semverBump ${lastTag} ${deducedBumpLevel})

  if [[ "${lastTag}" == "${nextTag}" ]]; then
    echo "Error: release v${nextTag} already exists (${deducedBumpLevel})"
  else
    echo "Creating release v${nextTag}"
    # git tag -m "v${nextTag}" v${nextTag}
    # git push origin v${nextTag}

    hub release create -m "v${nextTag}" -m "${latestCommitMessage}" "v${nextTag}"
  fi
fi
