import pytest


@pytest.mark.parametrize("name", [
    "go",
    "gocov",
    "golangci-lint",
    "gocovmerge",
    "gocov-xml",
])
def test_go_utils(host, name):
    assert host.exists(name)
