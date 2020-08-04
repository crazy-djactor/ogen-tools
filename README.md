# Ogen development tools

> A group of tool for continious building/deploying Ogen.

## Important Note

This services are developed to work over Linux. There are no future plans to make it working for other OS.

It is recommended to create a cron process to clean Docker images and containers constantly, builds are done using heavy size docker images and might take a lot of disk usage.

There is an optional script to remove all docker images on `clean_docker.sh`

## Explanation

This repository contains two tools to build Olympus development scenarios easily.

* Compiler
* Launcher

### Compiler

The compiler is a restful API that continuously runs Olympus cross-compiling scripts with a POST call from Github webhooks.

### Launcher

the launcher is a daemon script that runs a test network automatically.