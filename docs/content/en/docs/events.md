---
title: Events & Notifications
---

There are three types of event you can hook into: `on_success`, `on_error`, and `on_retries_exhausted`. The first two events materialize after an (attempted) job run, while `on_retries_exhausted` fires only once when a job with retries configured fails all attempts.

## Event Types

- **on_success**: Triggered when a job completes successfully
- **on_error**: Triggered when a job fails (fires after each failed attempt)
- **on_retries_exhausted**: Triggered only once when all retries have been exhausted

## Action Types

Three types of actions can be taken as a response:
- `notify_webhook`: Send a generic webhook notification
- `notify_slack_webhook`: Send a Slack-compatible webhook notification  
- `notify_discord_webhook`: Send a Discord-compatible webhook notification
- `trigger_job`: Trigger another job to run

## Configuration Examples

Definition of these event actions can be done on job level or at schedule level, in the latter case it will apply to all jobs.

```yaml
on_success:
  notify_webhook:
    - https://webhook.site/e33464a3-1a4f-4f1a-99d3-743364c6b10f
jobs:
  coffee:
    command: this fails # this will create on_error event
    cron: "* * * * *"
    retries: 3 # retry up to 3 times before giving up
    on_error:
      notify_webhook: # fires after each failed attempt
        - https://webhook.site/e33464a3-1a4f-4f1a-99d3-743364c6b10f
    on_retries_exhausted:
      trigger_job: # only fires once when all retries fail
        - cleanup_job
      notify_webhook:
        - https://webhook.site/critical-alerts
  beans:
    command: echo grind # this will create on_success event
    cron: "* * * * *"
```

## Webhook Payloads

### Generic Webhook

The `notify_webhook` sends a JSON payload to your webhook url with the following structure:

```json
{
	"status": 0,
	"log": "I'm a teapot, not a coffee machine!",
	"name": "TeapotTask",
	"triggered_at": "2023-04-01T12:00:00Z",
	"triggered_by": "CoffeeRequestButton",
	"triggered": ["CoffeeMachine"], // this job triggered another one
	"retry_attempt": 2, // which retry attempt this was (0 = first attempt)
	"retries_exhausted": true, // true when all retries have been exhausted
	"triggered_by_job_run": { // parent job context when triggered by another job
		"status": 0,
		"log": "Parent job completed successfully",
		"name": "ParentJob",
		"triggered_at": "2023-04-01T11:59:00Z",
		"triggered_by": "cron"
	}
}
```

When a job is triggered by another job via `trigger_job`, the webhook payload includes a `triggered_by_job_run` field containing the complete context of the parent job that triggered it. This provides full visibility into the job execution chain and allows for more sophisticated workflow tracking and debugging.

### Slack Webhook

The `notify_slack_webhook` sends a JSON payload to your Slack webhook url with the following structure (which is Slack app compatible):

```json
{
	"text": "TeapotTask (exitcode 0):\nI'm a teapot, not a coffee machine!"
}
```

### Discord Webhook

The `notify_discord_webhook` sends a JSON payload to your Discord webhook url with the following structure (which is Discord app compatible):

```json
{
	"content": "TeapotTask (exitcode 0):\nI'm a teapot, not a coffee machine!"
}
```