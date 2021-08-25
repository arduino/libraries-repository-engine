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

test_data_path = pathlib.Path(__file__).resolve().parent.joinpath("testdata")


def test_help(run_command):
    """Test the command line help."""
    # Run the `help modify` command
    engine_command = [
        "help",
        "remove",
    ]
    result = run_command(cmd=engine_command)
    assert result.ok
    assert "help for remove" in result.stdout

    # --help flag
    engine_command = [
        "remove",
        "--help",
    ]
    result = run_command(cmd=engine_command)
    assert result.ok
    assert "help for remove" in result.stdout


def test_invalid_flag(configuration, run_command):
    """Test the command's handling of invalid flags."""
    invalid_flag = "--some-bad-flag"
    engine_command = [
        "remove",
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
        "remove",
        "--config-file",
        configuration.path,
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert "LIBRARY_NAME argument is required" in result.stderr


def test_database_file_not_found(configuration, run_command):
    """Test the command's handling of incorrect LibrariesDB configuration."""
    engine_command = [
        "remove",
        "--config-file",
        configuration.path,
        "SpacebrewYun",
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert "Database file not found at {db_path}".format(db_path=configuration.data["LibrariesDB"]) in result.stderr


def test_remove_basic(configuration, run_command):
    """Test the basic functionality of the `remove` command."""
    # Run the sync command to generate test data
    engine_command = [
        "sync",
        "--config-file",
        configuration.path,
        test_data_path.joinpath("test_remove", "test_remove_basic", "repos.txt"),
    ]
    result = run_command(cmd=engine_command)
    assert result.ok
    assert pathlib.Path(configuration.data["LibrariesDB"]).exists()

    # Release reference syntax with missing version
    library_name = "SpacebrewYun"
    engine_command = [
        "remove",
        "--config-file",
        configuration.path,
        f"{library_name}@",
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert f"Missing version for library name {library_name}" in result.stderr

    # LIBRARY_NAME argument not in DB
    nonexistent_library_name = "nonexistent"
    engine_command = [
        "remove",
        "--config-file",
        configuration.path,
        nonexistent_library_name,
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert f"{nonexistent_library_name} not found" in result.stderr

    # LIBRARY_NAME@VERSION argument not in DB
    library_name = "SpacebrewYun"
    version = "99.99.99"
    engine_command = [
        "remove",
        "--config-file",
        configuration.path,
        f"{library_name}@{version}",
    ]
    result = run_command(cmd=engine_command)
    assert not result.ok
    assert f"Library release {library_name}@{version} not found" in result.stderr


def test_remove(configuration, run_command, working_dir):
    """Test the the removal of an entire library."""
    # Run the sync command to generate test data
    engine_command = [
        "sync",
        "--config-file",
        configuration.path,
        test_data_path.joinpath("test_remove", "repos.txt"),
    ]
    result = run_command(cmd=engine_command)
    assert result.ok
    assert pathlib.Path(configuration.data["LibrariesDB"]).exists()

    def git_clone_path_exists(host, owner, repo_name):
        return pathlib.Path(configuration.data["GitClonesFolder"]).joinpath(host, owner, repo_name).exists()

    def release_archive_path_exists(host, owner, library_name, version):
        sanitized_library_name = library_name.replace(" ", "_")
        return (
            pathlib.Path(configuration.data["LibrariesFolder"])
            .joinpath(host, owner, f"{sanitized_library_name}-{version}.zip")
            .exists()
        )

    def db_has_library(library_name):
        with pathlib.Path(configuration.data["LibrariesDB"]).open(mode="r", encoding="utf-8") as library_db_file:
            library_db = json.load(fp=library_db_file)

        for library in library_db["Libraries"]:
            if library["Name"] == library_name:
                return True

        return False

    def db_has_release(library_name, version):
        with pathlib.Path(configuration.data["LibrariesDB"]).open(mode="r", encoding="utf-8") as library_db_file:
            library_db = json.load(fp=library_db_file)

        for release in library_db["Releases"]:
            if release["LibraryName"] == library_name and release["Version"] == version:
                return True

        return False

    # Verify the pre-command environment is as expected
    assert git_clone_path_exists(host="github.com", owner="arduino-libraries", repo_name="ArduinoCloudThing")
    assert git_clone_path_exists(host="github.com", owner="arduino-libraries", repo_name="ArduinoIoTCloudBearSSL")
    assert git_clone_path_exists(
        host="github.com", owner="arduino-libraries", repo_name="UnoWiFi-Developer-Edition-Lib"
    )
    assert git_clone_path_exists(host="github.com", owner="arduino-libraries", repo_name="SpacebrewYun")

    # Note: The "ArduinoCloudThing" library is used as a "canary" to make sure that other libraries are not affected by
    # the removal process, so I don't bother to check all of its many releases
    # (there is lack of selection of appropriate libraries for test data)
    assert release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="ArduinoCloudThing", version="1.3.1"
    )
    assert release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="ArduinoIoTCloudBearSSL", version="1.1.1"
    )
    assert release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="ArduinoIoTCloudBearSSL", version="1.1.2"
    )
    assert release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="Arduino Uno WiFi Dev Ed Library", version="0.0.3"
    )
    assert release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="SpacebrewYun", version="1.0.0"
    )
    assert release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="SpacebrewYun", version="1.0.1"
    )
    assert release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="SpacebrewYun", version="1.0.2"
    )

    assert db_has_library(library_name="ArduinoCloudThing")
    assert db_has_library(library_name="ArduinoIoTCloudBearSSL")
    assert db_has_library(library_name="Arduino Uno WiFi Dev Ed Library")
    assert db_has_library(library_name="SpacebrewYun")

    assert db_has_release(library_name="ArduinoCloudThing", version="1.3.1")
    assert db_has_release(library_name="ArduinoIoTCloudBearSSL", version="1.1.1")
    assert db_has_release(library_name="ArduinoIoTCloudBearSSL", version="1.1.2")
    assert db_has_release(library_name="Arduino Uno WiFi Dev Ed Library", version="0.0.3")
    assert db_has_release(library_name="SpacebrewYun", version="1.0.0")
    assert db_has_release(library_name="SpacebrewYun", version="1.0.1")
    assert db_has_release(library_name="SpacebrewYun", version="1.0.2")

    # Run a remove command
    engine_command = [
        "remove",
        "--config-file",
        configuration.path,
        "ArduinoIoTCloudBearSSL",
        "Arduino Uno WiFi Dev Ed Library@0.0.3",
        "SpacebrewYun@1.0.1",
    ]
    result = run_command(cmd=engine_command)
    assert result.ok

    # Verify the post-command environment is as expected
    assert git_clone_path_exists(host="github.com", owner="arduino-libraries", repo_name="ArduinoCloudThing")
    assert not git_clone_path_exists(host="github.com", owner="arduino-libraries", repo_name="ArduinoIoTCloudBearSSL")
    assert git_clone_path_exists(
        host="github.com", owner="arduino-libraries", repo_name="UnoWiFi-Developer-Edition-Lib"
    )
    assert git_clone_path_exists(host="github.com", owner="arduino-libraries", repo_name="SpacebrewYun")

    assert release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="ArduinoCloudThing", version="1.3.1"
    )
    assert not release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="ArduinoIoTCloudBearSSL", version="1.1.1"
    )
    assert not release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="ArduinoIoTCloudBearSSL", version="1.1.2"
    )
    assert not release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="Arduino Uno WiFi Dev Ed Library", version="0.0.3"
    )
    assert release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="SpacebrewYun", version="1.0.0"
    )
    assert not release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="SpacebrewYun", version="1.0.1"
    )
    assert release_archive_path_exists(
        host="github.com", owner="arduino-libraries", library_name="SpacebrewYun", version="1.0.2"
    )

    assert db_has_library(library_name="ArduinoCloudThing")
    assert not db_has_library(library_name="ArduinoIoTCloudBearSSL")
    assert db_has_library(library_name="Arduino Uno WiFi Dev Ed Library")
    assert db_has_library(library_name="SpacebrewYun")

    assert db_has_release(library_name="ArduinoCloudThing", version="1.3.1")
    assert not db_has_release(library_name="ArduinoIoTCloudBearSSL", version="1.1.1")
    assert not db_has_release(library_name="ArduinoIoTCloudBearSSL", version="1.1.2")
    assert not db_has_release(library_name="Arduino Uno WiFi Dev Ed Library", version="0.0.3")
    assert db_has_release(library_name="SpacebrewYun", version="1.0.0")
    assert not db_has_release(library_name="SpacebrewYun", version="1.0.1")
    assert db_has_release(library_name="SpacebrewYun", version="1.0.2")
