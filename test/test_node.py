import pytest


def test_node_installed(host):
    assert host.exists('node')


def test_npm_installed(host):
    assert host.exists('npm')


def test_node_version(host):
    c = host.run('node -v')
    assert c.rc == 0
    assert c.stdout.startswith('v20')


def test_npm_version(host):
    c = host.run('npm -v')
    assert c.rc == 0
    assert c.stdout.startswith('10.')
