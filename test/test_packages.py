import pytest

@pytest.mark.parametrize("name,version", [
    ("make", "4.3-4.1"),
    ("git", "1:2.30.2-1"),
    ("docker-ce-cli", "5:27.5.1"),
    ("docker-compose-plugin", "2.33.1"),
])
def test_packages(host, name, version):
    pkg = host.package(name)
    assert pkg.is_installed
    assert pkg.version.startswith(version)

