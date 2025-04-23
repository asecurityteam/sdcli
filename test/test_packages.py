import pytest

@pytest.mark.parametrize("name,version", [
    ("make", ""),
    ("git", ""),
    ("docker-ce-cli", "5:27.5.1"),
    ("docker-compose-plugin", "2.33.1"),
    ("docker-compose", "1.29.2"),
    ("yamllint", "1.29.0-1")
])
def test_packages(host, name, version):
    pkg = host.package(name)
    assert pkg.is_installed
    assert pkg.version.startswith(version)

