import pytest

@pytest.mark.parametrize(
    "name,version,cmd", [
        ("docker", "27.5", "-v"),
        ("docker", "2.33", "compose version"),
        ("docker-compose", "1.29", "-v")
    ])
def test_packages(host, name, version, cmd):
    c = host.run("{} {}".format(name, cmd))
    assert c.rc == 0
    assert version in c.stdout