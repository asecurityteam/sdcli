import pytest


@pytest.mark.parametrize(
    "name,version,cmd", [
        ("docker", "20.10", "-v"),
        ("docker-compose", "2.11", "version")
    ])
def test_packages(host, name, version, cmd):
    c = host.run("{} {}".format(name, cmd))
    assert c.rc == 0
    assert version in c.stdout