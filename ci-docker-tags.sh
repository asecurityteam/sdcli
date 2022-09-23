#!/bin/sh
if [ -n "${TRAVIS_TAG}" ]; then
  echo "${TRAVIS_TAG}"
  echo "${TRAVIS_TAG}" | grep -Eq "^v[0-9]+[.]" && echo "${TRAVIS_TAG}" | sed -r 's/^v([0-9]+)[.].*$/v\1/'
elif [ "master" = "${TRAVIS_BRANCH}" ]; then
  echo latest
fi
