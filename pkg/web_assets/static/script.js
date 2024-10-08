document.addEventListener('alpine:init', () => {
  Alpine.store('job', {
    spec: null,
    jobName: null,
    jobRun: null,
    runId: null,

    fetchSpec: async function () {
      try {
        const response = await fetch(`/api/jobs/${this.jobName}`);
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        this.spec = await response.json();
      } catch (error) {
        console.error('Fetch error:', error);
      }
    },
    fetchJobRun: async function (runId) {
      try {
        const response = await fetch(`/api/jobs/${this.jobName}/runs/${runId}`);
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        this.jobRun = await response.json();
        this.runId = this.jobRun.id // update runId to the actual runId
      } catch (error) {
        console.error('Fetch error:', error);
      }
    },
    async init() {
      // get jobname from last part of url
      const { jobName, runId } = parseJobUrl(window.location.href);
      this.jobName = jobName;
      this.runId = runId === "latest" ? -1 : runId;

      this.fetchSpec();
      this.fetchJobRun(this.runId)
    }


  })

  Alpine.store('jobs', {
    jobs: null,


    fetchJobs: async function () {
      try {
        const response = await fetch('/api/jobs/');
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        this.jobs = await response.json();
        console.log(this.jobs, 777)

      } catch (error) {
        console.error('Fetch error:', error);
      }
    },

    init() {
      this.fetchJobs();

    }

  })

  Alpine.store('selectedJob', {
    jobName: null,


    set: function (jobName) {
      this.jobName = jobName;
    },
  })

  // alpine data component
  Alpine.data('coreLogs', () => ({

    logs: null,
    fetchLogs: async function () {
      try {
        const response = await fetch('/api/core/logs');
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        this.logs = await response.json();
      } catch (error) {
        console.error('Fetch error:', error);
      }
    },
    init() {
      this.fetchLogs();
    }

  }))


})


function triggerJob(jobName) {
  fetch(`/api/jobs/${jobName}/trigger`, {
    method: 'POST',
  }).then(response => {
    if (response.ok) {
      console.log(`Job ${jobName} triggered!`);
    } else {
      console.error(`Job ${jobName} could not be triggered!`);
    }
  });
}

function parseJobUrl(url) {
  // Using a regular expression to extract jobName and runId
  const regex = /\/jobs\/([^\/]+)\/([^\/]+)/;
  const match = url.match(regex);

  if (match && match.length >= 3) {
    return {
      jobName: match[1],
      runId: match[2]
    };
  } else {
    return {
      error: "Invalid URL format"
    };
  }
}


function truncateDateTime(dateTimeStr) {
  // Regular expression to match the date and time up to the minute
  const regex = /^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2})/;

  // Extract the matched part
  const match = dateTimeStr.match(regex);
  return match ? match[1] : null;
}

// Function to fetch version from the server
async function fetchVersion() {
  try {
    const response = await fetch('/api/version');
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    const data = await response.json();
    document.getElementById('version').textContent = `Version: ${data.version}`;
  } catch (error) {
    console.error('Fetch error:', error);
  }
}

// Fetch version when the DOM is fully loaded
document.addEventListener('DOMContentLoaded', () => {
  fetchVersion(); // Call the function when the DOM is fully loaded
});