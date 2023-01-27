import pytest

def test_help(host):
    c = host.run('/usr/bin/sdcli')
    assert c.rc == 0
    assert 'Available commands' in c.stdout