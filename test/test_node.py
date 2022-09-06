import pytest


def test_node_installed(host):
    assert host.exists('node')


def test_node_version(host):
    c = host.run('node -v')
    assert c.rc == 0
    assert c.stdout.startswith('v12')
