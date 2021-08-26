# Copyright 2021 ARDUINO SA (http://www.arduino.cc/)
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published
# by the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
#
# You can be released from the requirements of the above licenses by purchasing
# a commercial license. Buying such a license is mandatory if you want to
# modify or otherwise use the software for commercial activities involving the
# Arduino software without disclosing the source code of your own applications.
# To purchase a commercial license, send an email to license@arduino.cc.
#

import json
import pathlib

import pytest

test_data_path = pathlib.Path(__file__).resolve().parent.joinpath("testdata")


def test_help(run_command):
    """Test the command line help."""
    # Run the `help modify` command
    engine_command = [
        "help",
        "modify",
    ]
    result = run_command(cmd=engine_command)
    assert result.ok
    assert "help for modify" in result.stdout

    # --help flag
    engine_command = [
        "modify",
        "--help",
    ]
    result = run_command(cmd=engine_command)
    assert result.ok
    assert "help for modify" in result.stdout


def test_invalid_flag(configuration, run_command):
    """Test the command's handling of invalid flags."""
    invalid_flag = "--some-bad-flag"
    engine_command = [
        "modify",
        invalid_flag,
        "--config-file",
        configuration.path,
        "SpacebrewYun",
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert f"unknown flag: {invalid_flag}" in result.stderr


def test_missing_library_name_arg(configuration, run_command):
    """Test the command's handling of missing LIBRARY_NAME argument."""
    engine_command = [
        "modify",
        "--config-file",
        configuration.path,
        "--repo-url",
        "https://github.com/Foo/Bar.git",
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert "accepts 1 arg(s), received 0" in result.stderr


def test_multiple_library_name_arg(configuration, run_command):
    """Test the command's handling of multiple LIBRARY_NAME arguments."""
    engine_command = [
        "modify",
        "--config-file",
        configuration.path,
        "--repo-url",
        "https://github.com/Foo/Bar.git",
        "ArduinoIoTCloudBearSSL",
        "SpacebrewYun",
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert "accepts 1 arg(s), received 2" in result.stderr


def test_database_file_not_found(configuration, run_command):
    """Test the command's handling of incorrect LibrariesDB configuration."""
    engine_command = [
        "modify",
        "--config-file",
        configuration.path,
        "--repo-url",
        "https://github.com/Foo/Bar.git",
        "SpacebrewYun",
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert "Database file not found at {db_path}".format(db_path=configuration.data["LibrariesDB"]) in result.stderr


def test_repo_url_basic(configuration, run_command):
    """Test the basic functionality of the `--repo-url` modification flag."""
    # Run the sync command to generate test data
    engine_command = [
        "sync",
        "--config-file",
        configuration.path,
        test_data_path.joinpath("test_modify", "test_repo_url_basic", "repos.txt"),
    ]
    result = run_command(cmd=engine_command)
    assert result.ok
    assert pathlib.Path(configuration.data["LibrariesDB"]).exists()

    # Library not in DB
    nonexistent_library_name = "nonexistent"
    engine_command = [
        "modify",
        "--config-file",
        configuration.path,
        "--repo-url",
        "https://github.com/Foo/Bar.git",
        nonexistent_library_name,
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert f"{nonexistent_library_name} not found" in result.stderr

    # No local flag
    engine_command = [
        "modify",
        "--config-file",
        configuration.path,
        "SpacebrewYun",
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert "No modification flags" in result.stderr

    # Invalid URL format
    invalid_url = "https://github.com/Foo/Bar"
    engine_command = [
        "modify",
        "--config-file",
        configuration.path,
        "--repo-url",
        invalid_url,
        "SpacebrewYun",
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert f"{invalid_url} does not have a valid format" in result.stderr

    # Same URL as already in DB
    library_name = "SpacebrewYun"
    library_repo_url = "https://github.com/arduino-libraries/SpacebrewYun.git"
    engine_command = [
        "modify",
        "--config-file",
        configuration.path,
        "--repo-url",
        library_repo_url,
        library_name,
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert f"{library_name} already has URL {library_repo_url}" in result.stderr


@pytest.mark.parametrize(
    "name, releases, old_host, old_owner, old_repo_name, new_host, new_owner, new_repo_name",
    [
        (
            "Arduino Uno WiFi Dev Ed Library",
            ["0.0.3"],
            "github.com",
            "arduino-libraries",
            "UnoWiFi-Developer-Edition-Lib",
            "gitlab.com",
            "foo-owner",
            "bar-repo",
        ),
        (
            "SpacebrewYun",
            ["1.0.0", "1.0.1", "1.0.2"],
            "github.com",
            "arduino-libraries",
            "SpacebrewYun",
            "gitlab.com",
            "foo-owner",
            "bar-repo",
        ),
    ],
)
def test_repo_url(
    configuration,
    run_command,
    working_dir,
    name,
    releases,
    old_host,
    old_owner,
    old_repo_name,
    new_host,
    new_owner,
    new_repo_name,
):
    """Test the `--repo-url` modification flag in action."""
    sanitized_name = name.replace(" ", "_")
    old_library_release_archives_folder = pathlib.Path(configuration.data["LibrariesFolder"]).joinpath(
        old_host, old_owner
    )
    old_git_clone_path = pathlib.Path(configuration.data["GitClonesFolder"]).joinpath(
        old_host, old_owner, old_repo_name
    )
    new_repo_url = f"https://{new_host}/{new_owner}/{new_repo_name}.git"
    new_library_release_archives_folder = pathlib.Path(configuration.data["LibrariesFolder"]).joinpath(
        new_host, new_owner
    )
    new_git_clone_path = pathlib.Path(configuration.data["GitClonesFolder"]).joinpath(
        new_host, new_owner, new_repo_name
    )
    # The "canary" library is not modified and so all its content should remain unchanged after running the command
    canary_name = "ArduinoIoTCloudBearSSL"
    sanitized_canary_name = "ArduinoIoTCloudBearSSL"
    canary_release = "1.1.2"
    canary_host = "github.com"
    canary_owner = "arduino-libraries"
    canary_repo_name = "ArduinoIoTCloudBearSSL"
    canary_repo_url = f"https://{canary_host}/{canary_owner}/{canary_repo_name}.git"
    canary_release_filename = f"{sanitized_canary_name}-{canary_release}.zip"
    canary_release_archive_path = pathlib.Path(configuration.data["LibrariesFolder"]).joinpath(
        canary_host, canary_owner, canary_release_filename
    )
    canary_git_clone_path = pathlib.Path(configuration.data["GitClonesFolder"]).joinpath(
        canary_host, canary_owner, canary_repo_name
    )
    canary_release_archive_url = "{base}{host}/{owner}/{filename}".format(
        base=configuration.data["BaseDownloadUrl"],
        host=canary_host,
        owner=canary_owner,
        filename=canary_release_filename,
    )

    # Run the sync command to generate test data
    engine_command = [
        "sync",
        "--config-file",
        configuration.path,
        test_data_path.joinpath("test_modify", "test_repo_url", "repos.txt"),
    ]
    result = run_command(cmd=engine_command)
    assert result.ok
    assert pathlib.Path(configuration.data["LibrariesDB"]).exists()

    # Verify the pre-command environment is as expected
    def get_library_repo_url(name):
        with pathlib.Path(configuration.data["LibrariesDB"]).open(mode="r", encoding="utf-8") as library_db_file:
            library_db = json.load(fp=library_db_file)
        for library in library_db["Libraries"]:
            if library["Name"] == name:
                return library["Repository"]
        raise

    def get_release_archive_url(name, version):
        with pathlib.Path(configuration.data["LibrariesDB"]).open(mode="r", encoding="utf-8") as library_db_file:
            library_db = json.load(fp=library_db_file)
        for release in library_db["Releases"]:
            if release["LibraryName"] == name and release["Version"] == version:
                return release["URL"]
        raise

    assert old_git_clone_path.exists()
    assert not new_git_clone_path.exists()
    assert canary_git_clone_path.exists()
    assert get_library_repo_url(name=name) != new_repo_url
    assert get_library_repo_url(name=canary_name) == canary_repo_url
    for release in releases:
        assert old_library_release_archives_folder.joinpath(f"{sanitized_name}-{release}.zip").exists()
        assert not new_library_release_archives_folder.joinpath(f"{sanitized_name}-{release}.zip").exists()
        assert get_release_archive_url(name=name, version=release) == (
            "{base}{host}/{owner}/{name}-{release}.zip".format(
                base=configuration.data["BaseDownloadUrl"],
                host=old_host,
                owner=old_owner,
                name=sanitized_name,
                release=release,
            )
        )
    assert canary_release_archive_path.exists()
    assert get_release_archive_url(name=canary_name, version=canary_release) == canary_release_archive_url

    # Run the repository URL modification command
    engine_command = [
        "modify",
        "--config-file",
        configuration.path,
        "--repo-url",
        new_repo_url,
        name,
    ]
    result = run_command(cmd=engine_command)
    assert result.ok

    # Verify the effect of the command was as expected
    assert not old_git_clone_path.exists()
    assert new_git_clone_path.exists()
    assert canary_release_archive_path.exists()
    assert canary_git_clone_path.exists()
    assert get_library_repo_url(name=name) == new_repo_url
    assert get_library_repo_url(name=canary_name) == canary_repo_url
    for release in releases:
        assert not old_library_release_archives_folder.joinpath(f"{sanitized_name}-{release}.zip").exists()
        assert new_library_release_archives_folder.joinpath(f"{sanitized_name}-{release}.zip").exists()
        assert get_release_archive_url(name=name, version=release) == (
            "{base}{host}/{owner}/{name}-{release}.zip".format(
                base=configuration.data["BaseDownloadUrl"],
                host=new_host,
                owner=new_owner,
                name=sanitized_name,
                release=release,
            )
        )
    assert canary_release_archive_path.exists()
    assert get_release_archive_url(name=canary_name, version=canary_release) == canary_release_archive_url


def test_types(configuration, run_command):
    """Test the `--types` modification flag in action."""
    name = "SpacebrewYun"
    raw_old_types = "Arduino"
    raw_new_types = "Arduino, Retired"
    canary_name = "Arduino Uno WiFi Dev Ed Library"
    raw_canary_types = "Partner"
    # Run the sync command to generate test data
    engine_command = [
        "sync",
        "--config-file",
        configuration.path,
        test_data_path.joinpath("test_modify", "test_types", "repos.txt"),
    ]
    result = run_command(cmd=engine_command)
    assert result.ok
    assert pathlib.Path(configuration.data["LibrariesDB"]).exists()

    def assert_types(name, raw_types):
        with pathlib.Path(configuration.data["LibrariesDB"]).open(mode="r", encoding="utf-8") as library_db_file:
            library_db = json.load(fp=library_db_file)
        for release in library_db["Releases"]:
            if release["LibraryName"] == name and release["Types"] != [
                raw_type.strip() for raw_type in raw_types.split(sep=",")
            ]:
                return False
        return True

    # Verify the pre-command DB is as expected
    assert assert_types(name=name, raw_types=raw_old_types)
    assert assert_types(name=canary_name, raw_types=raw_canary_types)

    # Run the modification command with existing types
    engine_command = [
        "modify",
        "--config-file",
        configuration.path,
        "--types",
        raw_old_types,
        name,
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert f"{name} already has types {raw_old_types}" in result.stderr

    # Run the modification command with existing types
    engine_command = [
        "modify",
        "--config-file",
        configuration.path,
        "--types",
        raw_new_types,
        name,
    ]
    result = run_command(cmd=engine_command)
    assert result.ok

    # Verify the effect of the command was as expected
    assert assert_types(name=name, raw_types=raw_new_types)
    assert assert_types(name=canary_name, raw_types=raw_canary_types)
