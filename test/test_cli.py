import pytest


@pytest.fixture(autouse=True)
def dir(tmpdir):
    """Fixture to execute asserts before and after a test is run"""
    # Setup: fill with any logic you want

    yield # this is where the testing happens


def test_help(host):
    c = host.run('/usr/bin/sdcli')
    assert c.rc == 0
    assert 'Available commands' in c.stdout

def test_go_dep(host):
    c = host.run('/usr/bin/sdcli go dep')
    assert c.rc != 0
    assert 'go.mod not found' in c.stderr

def test_python_dep(host):
    c = host.run('/usr/bin/sdcli python dep')
    assert c.rc != 0
    assert 'Usage: pipenv' in c.stderr


def test_yaml_lint(host):
    c = host.run('/usr/bin/sdcli yaml lint')
    assert c.rc == 0 # no yaml files in test directory, so no error and no output expected
    assert '' == c.stdout
    assert '' == c.stderr