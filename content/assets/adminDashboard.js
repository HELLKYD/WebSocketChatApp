function webWorkerOnMessage(event) {
    let eventData = event.data;
    let connectedUserList = document.getElementById("connectedUsers");
    if(connectedUserList.innerHTML != "") {
        connectedUserList.innerHTML = "";
    }
    connectedUserList.innerHTML = "<tr><th>Id</th><th>Username</th></tr>";
    for(let i = 0; i < eventData.length; i++) {
        let data = eventData[i].split(":");
        let id = data[0];
        let username = data[1];
        connectedUserList.innerHTML += "<tr><td>" + id + "</td><td>" + username + "</td></tr>";
    }
}

let webWorker = new Worker("/content/fetchConnectedUsers.js");
webWorker.onmessage = webWorkerOnMessage;