import pathlib
from setuptools import setup, find_packages



HERE = pathlib.Path(__file__).parent

README = (HERE / "README.md").read_text()

setup(
    name="alsa_grpc_client",
    version="0.0.11",
    packages=find_packages(),
    # url="",
    description="Client for Remote Alsamixer",
    long_description=README,
    long_description_content_type="text/markdown",
)