# serve-chat

## :warning: WARNING: serve-chat is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

## Overview

### Install

1. Visit [Releases](https://github.com/Senzing/serve-chat/releases) page.
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

1. [Development](docs/development.md)
1. [Errors](docs/errors.md)
1. [Examples](docs/examples.md)
1. [Package reference](https://pkg.go.dev/github.com/senzing/serve-chat)
