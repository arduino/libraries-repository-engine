# Source:
# https://github.com/arduino/tooling-project-assets/blob/main/workflow-templates/assets/test-integration/test_all.py
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
import os
import pathlib
import platform
import shutil
import typing

import invoke.context
import pytest


@pytest.fixture
def configuration(working_dir):
    """Create a libraries-repository-engine configuration file and return an object containing its data and path."""
    working_dir_path = pathlib.Path(working_dir)

    # This is based on the `Librariesv2` production job's config.
    data = {
        "BaseDownloadUrl": "https://downloads.arduino.cc/libraries/",
        "LibrariesFolder": working_dir_path.joinpath("libraries").as_posix(),
        "LogsFolder": working_dir_path.joinpath("ci-logs", "libraries", "logs").as_posix(),
        "LibrariesDB": working_dir_path.joinpath("db.json").as_posix(),
        "LibrariesIndex": working_dir_path.joinpath("libraries", "library_index.json").as_posix(),
        "GitClonesFolder": working_dir_path.joinpath("gitclones").as_posix(),
        # I was unable to get clamdscan working in the GitHub Actions runner, but the tests should pass with this set to
        # False when run on a machine with ClamAV installed.
        "DoNotRunClamav": True,
        # Arduino Lint should be installed under PATH
        "ArduinoLintPath": "",
    }

    # Generate configuration file
    path = working_dir_path.joinpath("config.json")
    with path.open("w", encoding="utf-8") as configuration_file:
        json.dump(obj=data, fp=configuration_file, indent=2)

    class Object:
        """Container for libraries-repository-engine configuration data.

        Keyword arguments:
        data -- dictionary of configuration data
        path -- path of the configuration file
        """

        def __init__(self, data, path):
            self.data = data
            self.path = path

    return Object(data=data, path=path)


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
    work_dir = tmpdir_factory.mktemp(basename="IntegrationTestWorkingDir")
    yield os.path.realpath(work_dir)
    shutil.rmtree(work_dir, ignore_errors=True)
