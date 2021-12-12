![](https://dataroots.io/butt.png)

# butt

[![codecov](https://codecov.io/gh/Bart6114/butt/branch/main/graph/badge.svg?token=011KCCGPE6)](https://codecov.io/gh/Bart6114/butt) ![example workflow](https://github.com/bart6114/butt/actions/workflows/ci.yml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/bart6114/butt)](https://goreportcard.com/report/github.com/bart6114/butt) [![Go Reference](https://pkg.go.dev/badge/github.com/bart6114/butt.svg)](https://pkg.go.dev/github.com/bart6114/butt)


`butt`, of course, stands for Better Unified Time-Driven Triggering. `butt` is a KISS approach to crontab-like job scheduling. `butt` was born out of a (/my?) frustration about the big gap between a lightweight crontab and full-fledged solutions like Airflow.

`butt` aims to be a KISS approach to job scheduling. Focus is on the KISS approach not to necessarily do this in the most robust way possible.


## Getting started

Fetch the latest version for your system below.

[darwin-arm64](https://storage.googleapis.com/better-unified/darwin/arm64/butt) |
[darwin-amd64](https://storage.googleapis.com/better-unified/darwin/amd64/butt) |
[linux-386](https://storage.googleapis.com/better-unified/linux/386/butt) |
[linux-arm64](https://storage.googleapis.com/better-unified/linux/arm64/butt) |
[linux-amd64](https://storage.googleapis.com/better-unified/linux/amd64/butt)

You can (for example) fetch it like below, make it executable and run it. Optionally put the `butt` on your `PATH`.

```sh
curl https://storage.googleapis.com/better-unified/darwin/amd64/butt -o butt
chmod +x butt
./butt
```

Create a schedule specification using the below YAML structure:

```yaml
jobs:
  my_job:
    command: date
    cron: "* * * * *"
    triggers:
      - another_job
  another_job:
    command:
      - /bin/bash
      - -c
      - "sleep 2; echo bar"
  foo_job:
    command:
      - ls
      - .
    cron: "* * * * *"
  coffee_alert:
    command: this fails
    cron: "* * * * *"
    retries: 3
```

If your `command` requires arguments, please make sure to pass them as an array like in `foo_job`.

## Scheduler

The core of `butt` consists of a scheduler that uses a schedule specified in a `yaml` file to triggers jobs when they are due.

You can launch the scheduler via: 

```sh
butt run ./path/to/my-schedule.yaml
```

Check out `butt run --help` for configuration options.

## UI

`butt` ships with a terminal ui you can launch via:

```sh
butt ui
```

The UI allows to get a quick overview on jobs that have run, that error'd and their logs. It basically does this by fetching the state of the scheduler and by reading the logs that (per job) get written to `$HOME/.butt/`. Note that you can ignore these logs, output of jobs will always go to stdout as well.

The UI requires the scheduler to be up and running.

![](https://storage.googleapis.com/better-unified/ui-screenshot.png)



## Acknowledgements

Thanks goes to:
- [gronx](https://github.com/adhocore/gronx): for allowing me not to worry about CRON strings.
- [Charm](https://www.charm.sh/): for their bubble-icious TUI libraries.
- [Freddy](https://github.com/frederikdesmedt): for cracking them butt jokes.
- [Murilo](https://github.com/murilo-cunha): for asking me about my butt.