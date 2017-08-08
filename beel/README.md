# Installation and use guide.

## Description and current status
__beel__ was originally designed to be written using the [Go](https://golang.org/)
language, its purpose is to automate two main tasks: to build custom _Habitat_
packages, and to create & distribute custom _Docker_ images & containers.

It was written as planned, and it's capable of installing _habitat_ and to build
custom packages. However, due to an issue with the Docker's Go SDK, the second
task has been completed using [Bash](https://www.gnu.org/software/bash/).

Because of that, the _Habitat_ automated tasks are still work in progress, until
the issue gets fixed. However, the project *it's usable* if the current available
packages are enough for your project.

## Available packages
The entire list of available packages is located at the [habitat depot](https://bldr.habitat.sh/#/origins/bluespark).
These are the packages that had been already tested:

* bluespark/php5
* bluespark/php7
* bluespark/httpd

## Requirements
* Install Docker
  * For MacOS you can follow the the [official installation guide](https://docs.docker.com/docker-for-mac/install/).
  * For Linux:
    * For CentOS follow the [official installation guide](https://docs.docker.com/engine/installation/linux/docker-ce/centos/)
    * For Ubuntu follow the [official installation guide](https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/)
* Install AWS Command Line Interface
  * For MacOS you can follow the [official installation guide](http://docs.aws.amazon.com/cli/latest/userguide/cli-install-macos.html).
  * For Linux
    ```bash
    curl https://bootstrap.pypa.io/get-pip.py | sudo python
    sudo pip install awscli
    ```
* Configure your AWS credentials. You can get these credentials from 1password, search for _AWS bsp_dev_.
  ```bash
  $ aws configure
  AWS Access Key ID [None]: <ACCESS_KEY_ID>
  AWS Secret Access Key [None]: <SECRET_ACCESS_KEY>
  Default region name [None]: us-east-1
  Default output format [None]: ENTER
  ```

## Installation
* Clone the repository:
  ```bash
  git clone git@github.com/beel-platform/beel-platform.git
  ```
* Create a test directory, copy the configuration and run __beel__.
  ```bash
  cd bluespark-platform
  mkdir myproject
  cp beel/sh/blue.cfg myproject/
  cd myproject
  ../beel/sh/blue.sh
  ```
* If everything works as expected, now you can create your own project and
configure your beel.cfg file according to your needs.

## Troubleshooting
* Permissions on docker's daemon socket file.
  * The `/var/run/docker.sock` requires write permissions to execute docker,
  however the permissions by default belong to `root` as `660`, even if we do
  not recommend to run __beel__ using `root`, the issue remains, there's a
  solution by adding your unprivileged user to the docker group.
    ```bash
    sudo usermod -aG docker <your_unprivileged_user>
    ```
  * Usually adding the unprivileged user to the `docker` group should fix this
  issue. However, if that doesn't work and you're running __beel__ in a local
  machine, you can cheat a bit by changing the mode to `666` and allow all users
  to write. It's not a best practice, but for local environments it solves the
  problem.
    ```bash
    sudo chmod 666 /var/run/docker.sock
    ```
