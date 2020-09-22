import pathlib
from setuptools import setup

HERE = pathlib.Path(__file__).parent

README = (HERE / "README.md").read_text()

setup(
    name="remote_alsamixer_client",
    version="0.0.4",
    packages=["remote_alsamixer_client"],
    # url="",
    description="Client for Remote Alsamixer",
    long_description=README,
    long_description_content_type="text/markdown",
)