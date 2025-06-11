---
title: Job Execution Flow
---

The following diagram illustrates how jobs flow from schedule to execution and how events trigger subsequent actions:

{{< mermaid >}}
graph TD
    A[Schedule YAML] --> B[Job Specs]
    
    B --> C[Cron Schedule]
    
    C --> D[JobRun Execution]
    F[Manual Trigger] --> D
    G[Job Trigger] --> D
    
    D --> H{Job Success?}
    
    H -->|Yes| I[on_success Events]
    H -->|No| J{Retries Left?}
    
    J -->|Yes| K[Retry After Delay]
    J -->|No| L[on_retries_exhausted Events]
    
    K --> D
    H -->|No| M[on_error Events]
    
    I --> N[Event Actions]
    M --> N
    L --> N
    
    N --> O[trigger_job]
    N --> P["notify_{type}_webhook"]
    
    O --> Q[New JobRun]
    Q --> R[Parent Context Added]
    R --> S[triggered_by_job_run field populated]
    
    P --> T[Webhook Payload]
    T --> U[Includes parent job data if triggered by job]
    
    G -.-> Q
    S -.-> D
    
    style A fill:#e1f5fe
    style B fill:#e3f2fd
    style D fill:#f3e5f5
    style Q fill:#f3e5f5
    style R fill:#e8f5e8
    style T fill:#fff3e0
    style U fill:#fff3e0
{{< /mermaid >}}

## Flow Explanation

1. **Schedule Definition**: Jobs are defined in the YAML configuration file
2. **Job Specs**: The scheduler parses job specifications
3. **Execution Triggers**: Jobs can be triggered by:
   - Cron schedules
   - Manual triggers
   - Other jobs via `trigger_job`
4. **Job Execution**: The job runs in its own process
5. **Success/Failure Handling**: Based on the outcome:
   - Success triggers `on_success` events
   - Failure checks for retries or triggers `on_error` events
   - Exhausted retries trigger `on_retries_exhausted` events
6. **Event Actions**: Events can trigger webhooks or other jobs
7. **Parent Context**: Triggered jobs include context from their parent job