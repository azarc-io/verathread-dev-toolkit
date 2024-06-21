# Command Line Usage

Below is a list of commands available, you can also view this list by running the `task` command in the root of the project.

#### Setup

Installs `tilt` cli, `k3d` cli and creates a `.env` file in the root of the project.

```shell
task setup
```

#### K3d

Creates a new cluster:

```shell
task k3d:create
```

Destroys the cluster:

```shell
task k3d:delete
```

Installs chart dependencies:

```shell
task k3d:install:charts
```

Uninstalls charts installed with the previous command:

```shell
task k3d:delete:charts
```

Installs python dependencies required to work on this documentation:

```shell
task docs:install:deps
```

Serves and hot reloads this documentation:

```shell
task docs:serve
```

Builds this documentation:

```shell
task docs:build
```
