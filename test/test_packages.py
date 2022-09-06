import pytest


@pytest.mark.parametrize("name,version", [
    ("make", "4.2"),
    ("gcc", "4:8.3"),
    ("git", "1:2.20"),
    ("docker-ce-cli", "5:20.10"),
    ("docker-compose", "1.21"),
])
def test_packages(host, name, version):
    pkg = host.package(name)
    assert pkg.is_installed
    assert pkg.version.startswith(version)

