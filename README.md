# Beel-platform

Beel-platform is a project to provide instant and disposable environments
for local and remote use. It's composed of two main tools: _blue_ and _beel-cloud_.

___

## Description

__beel__ is a tool to automate the process of creating local environments, and
to make those environments distributable. For that purpose blue is based on two
tools: _Habitat_ and _Docker_.

* __Habitat__ is a tool that provides a consistent way to build and run
applications in a cloud native manner. It centers application configuration,
management, and behavior around the application itself, not the infrastructure
that the application runs on. This allows Habitat to be deployed and run on
various infrastructure environments, such as bare metal, virtual machines,
containers, and platform as a service.

* __Docker__ is a tool that provides containers, an additional layer of
abstraction and automation of _operating system level_ virtualization on Linux,
Windows and MacOS. By using Docker, __blue__ it's capable of distributing
environments to development teams in a standardize and reliable way.

The approach that __blue__ follows allows the operations team to focus on the
entire platform instead of having to fix and maintain specific environments.
That guarantees continuous improvement on the platform and leads to a reliable
service.

__beel-cloud__ is still work in progress and its goal is to provide the same
service as __beel__, but instead of focusing on local environments it will be
able to utilize cloud services like Amazon Web Services, Google Cloud, etc. to
distribute and deploy environments, simplifying the application's delivering
process.

---

## Installation

A guide to install and use __beel__ is provided within its own [README](beel/) file.

---

## Contact

For more information contact the current maintainer: bbh@briceno.mx

---

## Credits

Initial research, development and proof of concept sponsored by [Bluespark](https://www.bluespark.com)
