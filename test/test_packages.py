import pytest

@pytest.mark.parametrize("name,version", [
    ("make", "4.3-4.1"),
    ("gcc", "4:10.2.1-1"),
    ("git", "1:2.30.2-1+deb11u3"),
    ("docker-ce-cli", "5:27.3.1"),
    ("docker-compose-plugin", "2.29.7"),
])
def test_packages(host, name, version):
    pkg = host.package(name)
    assert pkg.is_installed
    assert pkg.version.startswith(version)

