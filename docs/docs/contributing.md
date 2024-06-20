# Contributing

Developers are expected to write and test their code on their local environments and in some cases remote environments
when available.

You should only open a Pull Request when one of the below requirements are met:

* You have been asked to provide a preview or demo of your work, in which case you should create a `draft` pr
* You have tested your code locally and are confident your work is ready for peer review and QA

## Project structure

- `deployment`: Deployment scripts and charts
    - `charts`: Dependency charts
    - `k3d`: Local kubernetes cluster configuration
- `e2e`: E2E tests for environment

## Testing

Creating a pull request for this project will test the following:

* The task commands run successfully
* K3D can be started
* The charts can be installed
* The dependencies are accessible and healthy

If you add a new dependency then you will need to create an E2E test for it, you must also update the `taskfile.yaml`
file in the root of the project to install and uninstall the new dependency.

The relevant tasks for this are `k3d:install:charts` and `k3d:delete:charts`, you can follow the existing examples.
