import pytest

@pytest.mark.parametrize("name,version", [
    ("make", ""),
    ("git", ""),
    ("docker-cli", ""),
    ("docker-compose", ""),
    ("yamllint", "")
])
def test_packages(host, name, version):
    pkg = host.package(name)
    assert pkg.is_installed
    assert pkg.version.startswith(version)

