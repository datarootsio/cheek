![](https://dataroots.io/butt.png)

# butt

`butt`, of course, stands for Better Unified Time-Driven Triggering. `butt` is a KISS approach to CRON-like job scheduling. `butt` was born out of a (/my?) frustration about the big gap between a lightweight crontab and full-fledged solutions like Airflow.

`butt` aims to be a KISS approach to job scheduling. Focus is on the KISS approach not to necessarily do this in the most robust way possible.


## Getting started

Install via ... TBD

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

It requires to the scheduler to be up and running. The UI allows to get a quick overview on jobs that have run, that error'd and their logs.

## Acknowledgements

Thanks goes to:
- [gronx](https://github.com/adhocore/gronx): for allowing me not to worry about CRON strings.
- [Charm](https://www.charm.sh/): for their bubbleicious TUI libraries.
- [Freddy](https://github.com/frederikdesmedt): for butt jokes.
- [Murilo](https://github.com/murilo-cunha): for asking me about my butt.