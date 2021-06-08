# This file is part of libraries-repository-engine.
#
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

import string
import re
import hashlib
import json
import pathlib
import platform
import typing

import invoke.context
import pytest

test_data_path = pathlib.Path(__file__).resolve().parent.joinpath("testdata")


def test_all(run_command, working_dir):
    working_dir_path = pathlib.Path(working_dir)
    configuration = {
        "BaseDownloadUrl": "http://www.example.com/libraries/",
        "LibrariesFolder": working_dir_path.joinpath("libraries").as_posix(),
        "LogsFolder": working_dir_path.joinpath("logs").as_posix(),
        "LibrariesDB": working_dir_path.joinpath("libraries_db.json").as_posix(),
        "LibrariesIndex": working_dir_path.joinpath("libraries", "library_index.json").as_posix(),
        "GitClonesFolder": working_dir_path.joinpath("gitclones").as_posix(),
        # I was unable to get clamdscan working in the GitHub Actions runner, but the tests should pass with this set to
        # False when run on a machine with ClamAV installed.
        "DoNotRunClamav": True,
        # Arduino Lint should be installed under PATH
        "ArduinoLintPath": "",
    }

    # Generate configuration file
    with working_dir_path.joinpath("config.json").open("w", encoding="utf-8") as configuration_file:
        json.dump(obj=configuration, fp=configuration_file, indent=2)

    libraries_repository_engine_command = [
        working_dir_path.joinpath("config.json"),
        test_data_path.joinpath("test_all", "repos.txt"),
    ]

    # Run the engine
    result = run_command(cmd=libraries_repository_engine_command)
    assert result.ok

    # Test fresh output
    check_libraries(configuration=configuration)
    check_logs(
        configuration=configuration,
        golden_logs_parent_path=test_data_path.joinpath("test_all", "golden", "logs", "generate"),
        logs_subpath=pathlib.Path("github.com", "arduino-libraries", "ArduinoCloudThing", "index.html"),
    )
    check_logs(
        configuration=configuration,
        golden_logs_parent_path=test_data_path.joinpath("test_all", "golden", "logs", "generate"),
        logs_subpath=pathlib.Path("github.com", "arduino-libraries", "SpacebrewYun", "index.html"),
    )
    check_index(configuration=configuration)

    # Run the engine again
    result = run_command(cmd=libraries_repository_engine_command)
    assert result.ok

    # Test the updated output
    check_libraries(configuration=configuration)
    check_logs(
        configuration=configuration,
        golden_logs_parent_path=test_data_path.joinpath("test_all", "golden", "logs", "update"),
        logs_subpath=pathlib.Path("github.com", "arduino-libraries", "ArduinoCloudThing", "index.html"),
    )
    check_logs(
        configuration=configuration,
        golden_logs_parent_path=test_data_path.joinpath("test_all", "golden", "logs", "update"),
        logs_subpath=pathlib.Path("github.com", "arduino-libraries", "SpacebrewYun", "index.html"),
    )
    check_index(configuration=configuration)


def check_libraries(configuration):
    """Run tests to determine whether the library release archives are as expected.

    Keyword arguments:
    configuration -- dictionary defining the libraries-repository-engine configuration
    """
    with pathlib.Path(configuration["LibrariesIndex"]).open(mode="r", encoding="utf-8") as libraries_index_file:
        libraries_index = json.load(fp=libraries_index_file)

    for release in libraries_index["libraries"]:
        release_archive_path = pathlib.Path(
            configuration["LibrariesFolder"],
            release["url"].removeprefix(configuration["BaseDownloadUrl"]),
        )

        assert release_archive_path.exists()

        assert release["size"] == release_archive_path.stat().st_size

        assert release["checksum"] == "SHA-256:" + hashlib.sha256(release_archive_path.read_bytes()).hexdigest()


def check_logs(configuration, golden_logs_parent_path, logs_subpath):
    """Run tests to determine whether the engine's logs are as expected.

    Keyword arguments:
    configuration -- dictionary defining the libraries-repository-engine configuration
    golden_logs_parent_path -- parent path for the golden master logs to compare the actual logs against
    logs_subpath -- sub-path for both the actual and golden master logs
    """
    logs = pathlib.Path(configuration["LogsFolder"], logs_subpath).read_text(encoding="utf-8")

    golden_logs_template = golden_logs_parent_path.joinpath(logs_subpath).read_text(encoding="utf-8")
    # Fill template with mutable content
    golden_logs = string.Template(template=golden_logs_template).substitute(
        git_clones_folder=configuration["GitClonesFolder"]
    )

    # Timestamps in the actual logs are not expected to match the golden logs, so replace with a placeholder
    timestamp_placeholder = "TIMESTAMP_PLACEHOLDER"
    timestamp_regex = re.compile(pattern=r"^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}", flags=re.MULTILINE)
    logs = re.sub(pattern=timestamp_regex, repl=timestamp_placeholder, string=logs)
    golden_logs = re.sub(pattern=timestamp_regex, repl=timestamp_placeholder, string=golden_logs)

    assert logs == golden_logs


def check_index(configuration):
    """Run tests to determine whether the generated library index is as expected.

    Keyword arguments:
    configuration -- dictionary defining the libraries-repository-engine configuration
    """
    checksum_placeholder = "CHECKSUM_PLACEHOLDER"

    # Load generated index
    with pathlib.Path(configuration["LibrariesIndex"]).open(mode="r", encoding="utf-8") as library_index_file:
        library_index = json.load(fp=library_index_file)
    for release in library_index["libraries"]:
        # The checksum values in the index will be different on every run, so it's necessary to replace them with a
        # placeholder before comparing to the golden index
        release["checksum"] = checksum_placeholder

    # Load golden index
    golden_library_index_template = test_data_path.joinpath("test_all", "golden", "library_index.json").read_text(
        encoding="utf-8"
    )
    # Fill in mutable content
    golden_library_index_string = string.Template(template=golden_library_index_template).substitute(
        base_download_url=configuration["BaseDownloadUrl"], checksum_placeholder=checksum_placeholder
    )
    golden_library_index = json.loads(golden_library_index_string)

    # Order of releases in the index is arbitrary so a simply equality assertion is not possible
    assert len(library_index["libraries"]) == len(golden_library_index["libraries"])
    for release in library_index["libraries"]:
        assert release in golden_library_index["libraries"]


# The engine's Git code struggles to get a clean checkout of releases under some circumstances.
def test_clean_checkout(run_command, working_dir):
    working_dir_path = pathlib.Path(working_dir)
    configuration = {
        "BaseDownloadUrl": "http://www.example.com/libraries/",
        "LibrariesFolder": working_dir_path.joinpath("libraries").as_posix(),
        "LogsFolder": working_dir_path.joinpath("logs").as_posix(),
        "LibrariesDB": working_dir_path.joinpath("libraries_db.json").as_posix(),
        "LibrariesIndex": working_dir_path.joinpath("libraries", "library_index.json").as_posix(),
        "GitClonesFolder": working_dir_path.joinpath("gitclones").as_posix(),
        "DoNotRunClamav": True,
        # Arduino Lint should be installed under PATH
        "ArduinoLintPath": "",
    }

    # Generate configuration file
    with working_dir_path.joinpath("config.json").open("w", encoding="utf-8") as configuration_file:
        json.dump(obj=configuration, fp=configuration_file, indent=2)

    libraries_repository_engine_command = [
        working_dir_path.joinpath("config.json"),
        test_data_path.joinpath("test_clean_checkout", "repos.txt"),
    ]

    # Run the engine
    result = run_command(cmd=libraries_repository_engine_command)
    assert result.ok

    # Load generated index
    with pathlib.Path(configuration["LibrariesIndex"]).open(mode="r", encoding="utf-8") as library_index_file:
        library_index = json.load(fp=library_index_file)

    for release in library_index["libraries"]:
        # ssd1306@1.0.0 contains a .exe and so should fail validation.
        assert not (release["name"] == "ssd1306" and release["version"] == "1.0.0")


@pytest.fixture(scope="function")
def run_command(pytestconfig, working_dir) -> typing.Callable[..., invoke.runners.Result]:
    """Provide a wrapper around invoke's `run` API so that every test will work in the same temporary folder.

    Useful reference:
        http://docs.pyinvoke.org/en/1.4/api/runners.html#invoke.runners.Result
    """

    executable_path = pathlib.Path(pytestconfig.rootdir).parent / "libraries-repository-engine"

    def _run(
        cmd: list,
        custom_working_dir: typing.Optional[str] = None,
        custom_env: typing.Optional[dict] = None,
    ) -> invoke.runners.Result:
        if cmd is None:
            cmd = []
        if not custom_working_dir:
            custom_working_dir = working_dir
        quoted_cmd = []
        for token in cmd:
            quoted_cmd.append(f'"{token}"')
        cli_full_line = '"{}" {}'.format(executable_path, " ".join(quoted_cmd))
        run_context = invoke.context.Context()
        # It might happen that we need to change directories between drives on Windows,
        # in that case the "/d" flag must be used otherwise directory wouldn't change
        cd_command = "cd"
        if platform.system() == "Windows":
            cd_command += " /d"
        # Context.cd() is not used since it doesn't work correctly on Windows.
        # It escapes spaces in the path using "\ " but it doesn't always work,
        # wrapping the path in quotation marks is the safest approach
        with run_context.prefix(f'{cd_command} "{custom_working_dir}"'):
            return run_context.run(
                command=cli_full_line,
                echo=False,
                hide=True,
                warn=True,
                env=custom_env,
                encoding="utf-8",
            )

    return _run


@pytest.fixture(scope="function")
def working_dir(tmpdir_factory) -> str:
    """Create a temporary folder for the test to run in. It will be created before running each test and deleted at the
    end. This way all the tests work in isolation.
    """
    work_dir = tmpdir_factory.mktemp(basename="TestWorkingDir")
    yield str(work_dir)
