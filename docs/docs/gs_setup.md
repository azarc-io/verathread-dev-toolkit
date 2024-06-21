# Setup Environment

!!! note

    If you have already setup k3d, task.dev and other tooling in one of our other projects then
    you can skip this document.

### Task.dev CLI

This is a replacement for make, it provides a much simpler syntax and is well documented.

You can read the task.dev documentation [here](https://taskfile.dev/usage/)

##### OSX

```shell
brew install go-task
```

##### Snap

```shell
sudo snap install task --classic
```

##### NPM

```shell
npm install -g @go-task/cli
```

##### Go Install
```shell
go install github.com/go-task/task/v3/cmd/task@latest
```

!!! info

    You can run `task` on your command line to view a list of available tasks and their descriptions.

### Initial Setup Task

Once task is installed run 

```shell
task setup
```

This will create a .env file in the root of the project, install the tilt cli and install the K3D cli.

### Host Names

In order to make ingress work you will need to add some records to your hosts file, we will be adding two
entries. One is for talking to the ingress and the other is for hosting local docker images.

Start by editing your host file typically found at `/etc/hosts`.

Add these two records and set the ip address to `127.0.0.1` if you are running kubernetes locally, alternatively if you
are running kube on a separate machine you can use that machines ip address.

```text
127.0.0.1 dev.cluster.local
127.0.0.1 k3d-local-registry
```

!!! note

    If you are hosting kube on a separate machine you will need to add that machines ip address to the `cluster.yaml` file
    located in `deployment/k3d/cluster.yaml` 

    In here look for the section starting with:
    ```yaml
      k3s:
        extraArgs:
    ```
    And add another entry like:
    ```yaml
      - arg: --tls-san=127.0.0.1
        nodeFilters:
            - server:*
    ```

### Environment Files

Edit the `.env` file in the root of the project and update the values

Set your `NAMESPACE` to a value such as `<YOUR FIRST NAME>-dev` you must make sure that the same value is set in
the projects you are working on.

!!! warning

    You must set the same value for `NAMESPACE` in all projects
