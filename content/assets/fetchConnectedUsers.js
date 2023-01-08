async function getUsers() {
    let response = await fetch("https://thatonedev.de/api/connectedUsers");
    let users = await response.json();
    let out = [];
    for(let i = 0; i < users.length; i++) {
        out[i] = users[i].id + ":" + users[i].username;
    }
    postMessage(out);
    setTimeout("getUsers()", 5000);
}
getUsers();