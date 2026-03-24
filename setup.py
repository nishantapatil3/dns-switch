#!/usr/bin/env python3
"""Setup script for dns-switch"""

from setuptools import setup, find_packages
from pathlib import Path

# Read the README file
this_directory = Path(__file__).parent
long_description = (this_directory / "README.md").read_text()

setup(
    name="dns-switch",
    version="1.0.0",
    description="A user-friendly TUI for quickly switching between DNS configurations",
    long_description=long_description,
    long_description_content_type="text/markdown",
    author="Pinaka Team",
    author_email="",
    url="https://github.com/pinaka-io/dns-switch",
    packages=find_packages(),
    py_modules=["dns_switch"],
    install_requires=[
        "textual>=0.50.0",
        "pyyaml>=6.0",
    ],
    entry_points={
        "console_scripts": [
            "dns-switch=dns_switch:main",
        ],
    },
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "Intended Audience :: System Administrators",
        "License :: OSI Approved :: MIT License",
        "Operating System :: MacOS",
        "Operating System :: POSIX :: Linux",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        "Topic :: System :: Networking",
        "Topic :: Utilities",
    ],
    python_requires=">=3.8",
    keywords="dns tui terminal cli network",
    include_package_data=True,
    package_data={
        "": ["config.yaml"],
    },
)
