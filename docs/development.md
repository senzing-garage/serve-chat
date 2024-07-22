# serve-chat development

## Install Go

1. See Go's [Download and install].

## Install Senzing C library

Since the Senzing library is a prerequisite, it must be installed first.

1. Verify Senzing C shared objects, configuration, and SDK header files are installed.
    1. `/opt/senzing/g2/lib`
    1. `/opt/senzing/g2/sdk/c`
    1. `/etc/opt/senzing`

1. If not installed, see [How to Install Senzing for Go Development].

## Install Git repository

1. Identify git repository.

    ```console
    export GIT_ACCOUNT=senzing-garage
    export GIT_REPOSITORY=serve-chat
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Using the environment variables values just set, follow
   steps in [clone-repository] to install the Git repository.

## Dependencies

1. A one-time command to install dependencies needed for `make` targets.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make dependencies-for-make

    ```

1. Install dependencies needed for [Go] code.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make dependencies

    ```

## Development cycle

Instructions are at
[Ogen QuickStart](https://ogen.dev/docs/intro/).

1. Get latest version of [ogen](https://github.com/ogen-go/ogen) code generator.
   Do this frequently (i.e. daily), as code is changing constantly.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    go get -d github.com/ogen-go/ogen
    ```

1. View version.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    go list -m github.com/ogen-go/ogen
    ```

1. Modify
   [openapi.json](../senzingchatservice/openapi.json).
   **Note:** It must be `json`.  For some reason `yaml` doesn't work.
1. Generate code from
   [openapi.json](../senzingchatservice/openapi.json).
   Example:

    ```console
     cd ${GIT_REPOSITORY_DIR}
     make generate

    ```

1. Modify
   [senzingchatservice.go](../senzingchatservice/senzingchatservice.go)
   implementing method invocations seen in
   [oas_unimplemented_gen.go](../senzingchatapi/oas_unimplemented_gen.go)

1. Create clean SQLite test database.
   Example:

    ```console
   cd ${GIT_REPOSITORY_DIR}
   make clean

    ```

1. Test.

    ```console
   cd ${GIT_REPOSITORY_DIR}
   make test

    ```

## Build

1. Build the binaries.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean build

    ```

1. The binaries will be found in the `${GIT_REPOSITORY_DIR}/target` directory.
   Example:

    ```console
    tree ${GIT_REPOSITORY_DIR}/target

    ```

## Run

1. Run without a build.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make run

    ```

1. Open a web browser at [localhost:8262](http://localhost:8262).

1. Run the binary.
   Example:

    ```console
    ${GIT_REPOSITORY_DIR}/target/linux/serve-chat

    ```

1. Clean up.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean

    ```

## Lint

1. Run Go tests.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make lint

    ```

## Test

1. Run Go tests.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean setup test

    ```

## Coverage

Create a code coverage map.

1. Run Go tests.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean setup coverage

    ```

   A web-browser will show the results of the coverage.
   The goal is to have over 80% coverage.
   Anything less needs to be reflected in [testcoverage.yaml].

## Documentation

1. Start [godoc] documentation server.
   Example:

    ```console
     cd ${GIT_REPOSITORY_DIR}
     make clean documentation

    ```

1. If a web page doesn't appear, visit [localhost:6060].
1. Senzing documentation will be in the "Third party" section.
   `github.com` > `senzing` > `go-cmdhelping`

1. When a versioned release is published with a `v0.0.0` format tag,
the reference can be found by clicking on the following badge at the top of the README.md page.
Example:

    [![Go Reference](https://pkg.go.dev/badge/github.com/senzing-garage/serve-chat.svg)](https://pkg.go.dev/github.com/senzing-garage/serve-chat)

1. To stop the `godoc` server, run

    ```console
     cd ${GIT_REPOSITORY_DIR}
     make clean

    ```

## Docker

1. Use make target to run a docker images that builds RPM and DEB files.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make docker-build

    ```

1. Run docker container.
   Example:

    ```console
    docker run \
      --env SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db \
      --publish 8262:8262 \
      --rm \
      --volume /tmp/sqlite:/tmp/sqlite \
      senzing/serve-chat --enable-all

    ```

## Package

### Package RPM and DEB files

1. Use make target to run a docker images that builds RPM and DEB files.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make package

    ```

1. The results will be in the `${GIT_REPOSITORY_DIR}/target` directory.
   Example:

    ```console
    tree ${GIT_REPOSITORY_DIR}/target

    ```

### Test DEB package on Ubuntu

1. Determine if `serve-chat` is installed.
   Example:

    ```console
    apt list --installed | grep serve-chat

    ```

1. :pencil2: Install `serve-chat`.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}/target
    sudo apt install ./serve-chat-0.0.0.deb

    ```

1. Run command.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    serve-chat --database-url sqlite3://na:na@/tmp/sqlite/G2C.db --enable-all

    ```

1. Remove `serve-chat` from system.
   Example:

    ```console
    sudo apt-get remove serve-chat

    ```

[clone-repository]: https://github.com/senzing-garage/knowledge-base/blob/main/HOWTO/clone-repository.md
[Download and install]: https://go.dev/doc/install
[Go]: https://go.dev/
[godoc]: https://pkg.go.dev/golang.org/x/tools/cmd/godoc
[How to Install Senzing for Go Development]: https://github.com/senzing-garage/knowledge-base/blob/main/HOWTO/install-senzing-for-go-development.md
[localhost:6060]: http://localhost:6060/pkg/github.com/senzing-garage/serve-chat/
[testcoverage.yaml]: ../.github/coverage/testcoverage.yaml
