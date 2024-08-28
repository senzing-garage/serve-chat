# serve-chat

If you are beginning your journey with [Senzing],
please start with [Senzing Quick Start guides].

You are in the [Senzing Garage] where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: serve-chat is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The serve-chat repository serves as a starting point for new repositories hosting Go code.
It also shows best practices that can be retro-fitted into existing repositories hosting Go code.

[![Go Reference Badge]][Package reference]
[![Go Report Card Badge]][Go Report Card]
[![License Badge]][License]
[![go-test-linux.yaml Badge]][go-test-linux.yaml]
[![go-test-darwin.yaml Badge]][go-test-darwin.yaml]
[![go-test-windows.yaml Badge]][go-test-windows.yaml]

[![golangci-lint.yaml Badge]][golangci-lint.yaml]

## Overview

### Install

1. Visit [Releases](https://github.com/senzing-garage/serve-chat/releases) page.
1. For the desired versioned release, in the "Assets" section,
   download the appropriate installation package.
    1. Use `.deb` file for Debian, Ubuntu and
       [others](https://en.wikipedia.org/wiki/List_of_Linux_distributions#Debian-based)
    1. Use `.rpm` file for Red Hat, CentOS, openSuse and
       [others](https://en.wikipedia.org/wiki/List_of_Linux_distributions#RPM-based).

1. :pencil2: Example installation for `.deb` file:

    ```console
    sudo apt install ./serve-chat-0.0.0.deb
    ```

1. :pencil2: Example installation for `.rpm` file:

    ```console
    sudo yum install ./serve-chat-0.0.0.rpm
    ```

### Using Docker

1. Run Docker image against local SQLite database.
   Example:

    ```console
    docker run \
      --env SENZING_TOOLS_ENABLE_ALL=true \
      --env SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db \
      --publish 8262:8262 \
      --rm \
      --volume /tmp/sqlite:/tmp/sqlite \
      senzing/serve-chat

    ```

1. Open browser on [localhost:8252](http://localhost:8262)

## References

1. [API documentation]
1. [Development]
1. [Errors]
1. [Examples]

[API documentation]: https://pkg.go.dev/github.com/senzing-garage/serve-chat
[Development]: docs/development.md
[Errors]: docs/errors.md
[Examples]: docs/examples.md
[Go Reference Badge]: https://pkg.go.dev/badge/github.com/senzing-garage/serve-chat.svg
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/senzing-garage/serve-chat
[Go Report Card]: https://goreportcard.com/report/github.com/senzing-garage/serve-chat
[go-test-darwin.yaml Badge]: https://github.com/senzing-garage/serve-chat/actions/workflows/go-test-darwin.yaml/badge.svg
[go-test-darwin.yaml]: https://github.com/senzing-garage/serve-chat/actions/workflows/go-test-darwin.yaml
[go-test-linux.yaml Badge]: https://github.com/senzing-garage/serve-chat/actions/workflows/go-test-linux.yaml/badge.svg
[go-test-linux.yaml]: https://github.com/senzing-garage/serve-chat/actions/workflows/go-test-linux.yaml
[go-test-windows.yaml Badge]: https://github.com/senzing-garage/serve-chat/actions/workflows/go-test-windows.yaml/badge.svg
[go-test-windows.yaml]: https://github.com/senzing-garage/serve-chat/actions/workflows/go-test-windows.yaml
[golangci-lint.yaml Badge]: https://github.com/senzing-garage/serve-chat/actions/workflows/golangci-lint.yaml/badge.svg
[golangci-lint.yaml]: https://github.com/senzing-garage/serve-chat/actions/workflows/golangci-lint.yaml
[License Badge]: https://img.shields.io/badge/License-Apache2-brightgreen.svg
[License]: https://github.com/senzing-garage/serve-chat/blob/main/LICENSE
[Package reference]: https://pkg.go.dev/github.com/senzing-garage/serve-chat
[Senzing Garage]: https://github.com/senzing-garage
[Senzing Quick Start guides]: https://docs.senzing.com/quickstart/
[Senzing]: https://senzing.com/
