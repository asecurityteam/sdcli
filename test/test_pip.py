import pytest


def test_pip_installed(host):
    assert host.exists('pip3')


def test_python_broken_deps(host):
    assert host.pip.check(pip_path='pip3')


@pytest.mark.parametrize("name", [
    "setuptools",
    "cookiecutter",
    "flake8",
    "coverage",
    "pytest",
    "pytest-cov",
    "pipenv",
    "oyaml",
    "python-slugify",
    "yamllint"
])
def test_python_packages(host, name):
    pkg = host.pip(name, pip_path='pip3')
    assert pkg.is_installed
