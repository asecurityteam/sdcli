#!/bin/sh
if [ -n "${TRAVIS_TAG}" ]; then
  echo "${TRAVIS_TAG}"
  echo "${TRAVIS_TAG}" | grep -Eq "^v[0-9]+[.]" && echo "${TRAVIS_TAG}" | sed -r 's/^v([0-9]+)[.].*$/v\1/'
# the check for event type allows to only tag 'latest' when the code is pushed to master. TRAVIS_BRANCH is set to master for PR CI, but we do not want to tag/push on PR, only on merge
elif [ "master" = "${TRAVIS_BRANCH}" ] && [ "push" = "${TRAVIS_EVENT_TYPE}" ]; then
  echo latest
fi
