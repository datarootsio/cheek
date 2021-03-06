<img src="https://storage.googleapis.com/better-unified/cheek265.png" alt="cheek" width="128" />

# cheek

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/datarootsio/cheek?label=version)
[![dataroots](https://dataroots.io/maintained.svg)](https://dataroots.io/) [![codecov](https://codecov.io/gh/datarootsio/cheek/branch/main/graph/badge.svg?token=011KCCGPE6)](https://codecov.io/gh/datarootsio/cheek) ![example workflow](https://github.com/datarootsio/cheek/actions/workflows/ci.yml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/datarootsio/cheek)](https://goreportcard.com/report/github.com/datarootsio/cheek) [![Go Reference](https://pkg.go.dev/badge/github.com/datarootsio/cheek.svg)](https://pkg.go.dev/github.com/datarootsio/cheek) ![triggers](https://img.shields.io/badge/dynamic/json?color=blueviolet&label=jobs-triggered&query=%24.triggered&url=https%3A%2F%2Fapi.dataroots.io%2Fv1%2Fcheek%2Ftriggered&style=flat-square)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)

`cheek`, of course, stands for `C`rontab-like sc`H`eduler for `E`ffective `E`xecution of tas`K`s. `cheek` is a KISS approach to crontab-like job scheduling. It was born out of a (/my?) frustration about the big gap between a lightweight crontab and full-fledged solutions like Airflow.

`cheek` aims to be a KISS approach to job scheduling. Focus is on the KISS approach not to necessarily do this in the most robust way possible.

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
tz_location: Europe/Brussels
jobs:
  foo:
    command: date
    cron: "* * * * *"
    on_success:
      trigger_job:
        - bar
  bar:
    command:
      - echo
      - bar
      - foo
  coffee:
    command: this fails
    cron: "* * * * *"
    retries: 3
    on_error:
      notify_webhook:
        - https://webhook.site/4b732eb4-ba10-4a84-8f6b-30167b2f2762
```

If your `command` requires arguments, please make sure to pass them as an array like in `foo_job`.

Note that you can set `tz_location` if the system time of where you run your service is not to your liking.

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

![](https://storage.googleapis.com/better-unified/ui-screenshot2.png)

## Configuration

All configuration options are available by checking out `cheek --help` or the help of its subcommands (e.g. `cheek run --help`).

Configuration can be passed as flags to the `cheek` CLI directly. All configuration flags are also possible to set via environment variables. The following environment variables are available, they will override the default and/or set value of their similarly named CLI flags (without the prefix): `CHEEK_PORT`, `CHEEK_SUPPRESSLOGS`, `CHEEK_LOGLEVEL`, `CHEEK_PRETTY`, `CHEEK_HOMEDIR`, `CHEEK_NOTELEMETRY`.

## Events

There are two types of event you can hook into: `on_success` and `on_error`. Both events materialize after an (attempted) job run. Two types of actions can be taken as a response: `notify_webhook` and `trigger_job`. See the example below. Definition of these event actions can be done on job level or at schedule level, in the latter case it will apply to all jobs.

```yaml
on_success:
  notify_webhook:
    - https://webhook.site/e33464a3-1a4f-4f1a-99d3-743364c6b10f
jobs:
  coffee:
    command: this fails # this will create on_error event
    cron: "* * * * *"
    on_error:
      notify_webhook:
        - https://webhook.site/e33464a3-1a4f-4f1a-99d3-743364c6b10f
  beans:
    command: echo grind # this will create on_success event
    cron: "* * * * *"
```

Webhook are a generic way to push notifications to a plethora of tools. You can use it for instance via Zapier to push messages to a Slack channel.

## Docker

Check out the `Dockerfile` for an example on how to set up `cheek` within the context of a Docker image.

## Available versions

If you want to pin your setup to a specific version of `cheek` you can use the following template to fetch your `cheek` binary:

- latest version: https://storage.googleapis.com/better-unified/{os}/{arch}/cheek
- tagged version: https://storage.googleapis.com/better-unified/{os}/{arch}/cheek-{tag}
- `main` branch builds: https://storage.googleapis.com/better-unified/{os}/{arch}/cheek-{shortsha}

Where:
- `os` is one of `linux`, `darwin`
- `arch` is one of `amd64`, `arm64`, `386`
- `tag` is one the [available tags](https://github.com/datarootsio/cheek/tags)
- `shortsha` is a 7-char SHA and most commits on `main` will be available

## Usage stats

By default `cheek` reports minimal usage stats. Each time a job is triggered a simple request that (only) contains your `cheek` version is send to our servers. Check out the exact implementation [here](https://github.com/datarootsio/cheek/blob/main/pkg/telemetry.go). Note that you can always opt-out of this by passing the `-no-telemetry` or `-n` flag.

## Acknowledgements

Thanks goes to:

- [gronx](https://github.com/adhocore/gronx): for allowing me not to worry about CRON strings.
- [Charm](https://www.charm.sh/): for their bubble-icious TUI libraries.
- [Sam](https://github.com/sdebruyn) & [Frederik](https://github.com/frederikdesmedt): for valuable code reviews / feedback.

<br/>
 
![GitHub Contributors](https://contrib.rocks/image?repo=datarootsio/cheek)
