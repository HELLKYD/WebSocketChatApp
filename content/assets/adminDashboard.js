function webWorkerOnMessage(event) {
    let eventData = event.data;
    let connectedUserList = document.getElementById("connectedUsers");
    if (connectedUserList.innerHTML != "") {
        connectedUserList.innerHTML = "";
    }
    connectedUserList.innerHTML = "<tr><th>Id</th><th>Username</th></tr>";
    for (let i = 0; i < eventData.length; i++) {
        let data = eventData[i].split(":");
        let id = data[0];
        let username = data[1];
        connectedUserList.innerHTML += "<tr><td>" + id + "</td><td>" + username + "</td></tr>";
    }
}

function startDataCollection() {
    let webWorker = new Worker("/content/fetchConnectedUsers.js");
    webWorker.onmessage = webWorkerOnMessage;
}

async function verifyLogin(username, password) {
    let url = "https://thatonedev.de/api/verifyUserLoginData/?loginData=".concat(username + ":" + password);
    let response = await fetch(url);
    let responseJson = await response.json();
    if(responseJson.isValid) {
        let userInput = document.getElementById('userInput');
        let passwordInput = document.getElementById('passwordInput');
        let submitButton = document.getElementById('submitButton');
        userInput.hidden = true;
        passwordInput.hidden = true;
        submitButton.hidden = true;
        startDataCollection();
    } else {
        window.alert('Wrong credentials');
    }
}

function login() {
    let userInput = document.getElementById('userInput');
    let passwordInput = document.getElementById('passwordInput');
    verifyLogin(userInput.value, passwordInput.value);
}