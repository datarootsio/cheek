const triggerJob = (jobName) => {
    const Http = new XMLHttpRequest();
    const url = '/trigger/' + jobName;
    Http.open("GET", url);
    Http.send();

    Http.onreadystatechange = (e) => {
        console.log(Http.responseText)
    }
}