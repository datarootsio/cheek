<img src="https://storage.googleapis.com/better-unified/cheek265.png" alt="cheek" width="128" />

# cheek

[![codecov](https://codecov.io/gh/datarootsio/cheek/branch/main/graph/badge.svg?token=011KCCGPE6)](https://codecov.io/gh/datarootsio/cheek) ![example workflow](https://github.com/datarootsio/cheek/actions/workflows/ci.yml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/datarootsio/cheek)](https://goreportcard.com/report/github.com/datarootsio/cheek) [![Go Reference](https://pkg.go.dev/badge/github.com/datarootsio/cheek.svg)](https://pkg.go.dev/github.com/datarootsio/cheek) [![dataroots](https://dataroots.io/maintained.svg)](https://dataroots.io/)

`cheek`, of course, stands for `C`rontab-like sc`H`eduler for `E`ffective `E`xecution of tas`K`s. `cheek` is a KISS approach to crontab-like job scheduling. It was born out of a (/my?) frustration about the big gap between a lightweight crontab and full-fledged solutions like Airflow.

`cheek` aims to be a KISS approach to job scheduling. Focus is on the KISS approach not to necessarily do this in the most robust way possible.

## TOC

- [cheek](#cheek)
  - [TOC](#toc)
  - [Getting started](#getting-started)
  - [Scheduler](#scheduler)
  - [UI](#ui)
  - [Configuration](#configuration)
  - [Docker](#docker)
  - [Acknowledgements](#acknowledgements)


## Getting started

Fetch the latest version for your system below.

[darwin-arm64](https://storage.googleapis.com/better-unified/darwin/arm64/cheek) |
[darwin-amd64](https://storage.googleapis.com/better-unified/darwin/amd64/cheek) |
[linux-386](https://storage.googleapis.com/better-unified/linux/386/cheek) |
[linux-arm64](https://storage.googleapis.com/better-unified/linux/arm64/cheek) |
[linux-amd64](https://storage.googleapis.com/better-unified/linux/amd64/cheek)

You can (for example) fetch it like below, make it executable and run it. Optionally put the `cheek` on your `PATH`.

```sh
curl https://storage.googleapis.com/better-unified/darwin/amd64/cheek -o cheek
chmod +x cheek
./cheek
```

Create a schedule specification using the below YAML structure:

```yaml
jobs:
  foo:
    command: date
    cron: "* * * * *"
    on_success:
      trigger_job:
        - bar
  bar:
    command:
      - /bin/bash
      - -c
      - "echo bar_foo"
  coffee:
    command: this fails
    cron: "* * * * *"
    retries: 3
    on_error:
      notify_webhook:
        - https://webhook.site/4b732eb4-ba10-4a84-8f6b-30167b2f2762
```

If your `command` requires arguments, please make sure to pass them as an array like in `foo_job`.

## Scheduler

The core of `cheek` consists of a scheduler that uses a schedule specified in a `yaml` file to triggers jobs when they are due.

You can launch the scheduler via:

```sh
cheek run ./path/to/my-schedule.yaml
```

Check out `cheek run --help` for configuration options.

## UI

`cheek` ships with a terminal ui you can launch via:

```sh
cheek ui
```

The UI allows to get a quick overview on jobs that have run, that error'd and their logs. It basically does this by fetching the state of the scheduler and by reading the logs that (per job) get written to `$HOME/.cheek/`. Note that you can ignore these logs, output of jobs will always go to stdout as well.

The UI requires the scheduler to be up and running.

![](https://storage.googleapis.com/better-unified/ui-screenshot2.png)

## Configuration

All configuration options are available by checking out `cheek --help` or the help of its subcommands (e.g. `cheek run --help`).

Configuration can be passed as flags to the `cheek` CLI directly. All configuration flags are also possible to set via environment variables. The following environment variables are available, they will override the default and/or set value of their similarly named CLI flags (without the prefix): `CHEEK_PORT`, `CHEEK_SUPPRESSLOGS`, `CHEEK_LOGLEVEL`, `CHEEK_PRETTY`, `CHEEK_HOMEDIR`.

## Docker

Check out the `Dockerfile` for an example on how to set up `cheek` within the context of a Docker image.

## Acknowledgements

Thanks goes to:

- [gronx](https://github.com/adhocore/gronx): for allowing me not to worry about CRON strings.
- [Charm](https://www.charm.sh/): for their bubble-icious TUI libraries.
- [Sam](https://github.com/sdebruyn) & [Frederik](https://github.com/frederikdesmedt): for valuable code reviews / feedback.

<br/>
 
![GitHub Contributors](https://contrib.rocks/image?repo=datarootsio/cheek)
