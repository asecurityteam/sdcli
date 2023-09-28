import pytest

# TODO: create a new github tag for new image
# TODO: fix unit test
@pytest.mark.parametrize("name,version", [
    ("make", "4.3-4.1"),
    ("gcc", "4:10.2.1-1"),
    ("git", "1:2.30.2-1"),
    ("docker-ce-cli", "5:20.10.6"),
    ("docker-compose-plugin", "2.11.2"),
])
def test_packages(host, name, version):
    pkg = host.package(name)
    assert pkg.is_installed
    assert pkg.version.startswith(version)

