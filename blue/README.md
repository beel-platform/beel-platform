# blue's installation and use guide.

## Description and current status

__blue__ was originally designed to be written using the [Go](https://golang.org/) language,
to automate two main tasks: to build custom _Habitat_ packages, and to create &
distribute custom _Docker_ images & containers.

It was written as planned, and it's capable of installing _habitat_ and to build
custom packages. However, due to an issue with the Docker' Go library the second
task, has been completed using [Bash](https://www.gnu.org/software/bash/).

Because of that, the _Habitat_ automated tasks are still work in progress, until
the issue gets fixed. However, the project *it's usable* if the current available
packages are enough for your project.

## Available packages

The entire list of available packages is located at the [habitat depot](https://bldr.habitat.sh/#/origins/bluespark).
There are the packages that had been already tested.

* bluespark/php5
* bluespark/php7
* bluespark/httpd

## Requirements
* Install Docker
  * For MacOS you can follow the the [official installation guide](https://docs.docker.com/docker-for-mac/install/).
  * For Linux
    * Make sure that your docker service is running.
    ```bash
    sudo service docker start
    ```
    * Linux __sudo issue__. There's an _"issue"_ with the docker's socket
    permissions, the `/var/run/docker.sock` requires write permissions to
    execute docker, however the permissions by default belong to `root` as `660`,
    even if we do not recommend to run __blue__ using `root`, the issue remains,
    there's a solution by adding your unprivileged user to the docker group
    ```bash
    sudo usermod -aG docker <your_unprivileged_user>
    ```
* Install AWS Command Line Interface
  * For MacOS you can follow the [official installation guide](http://docs.aws.amazon.com/cli/latest/userguide/cli-install-macos.html).
  * For Linux
  ```bash
  curl https://bootstrap.pypa.io/get-pip.py | sudo python
  sudo pip install awscli
  ```
  * Configure your AWS credentials. You can get these settings from 1password.
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
git clone git@github.com:BluesparkLabs/bluespark-platform.git
```
* Copy the configuration example to your project

## Troubleshooting
* Permissions on docker's daemon socket file. Usually adding the unprivileged
user to the `docker` group should fix this issue, however, if that doesn't works
and you're running __blue__ in a local machine, you can cheat a bit and
```bash
sudo chmod 666 /var/run/docker.sock
```
this is not a best practice, but for a local environment it solves the problem.
