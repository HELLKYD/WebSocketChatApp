import { Component, createSignal } from "solid-js";
import { render } from "solid-js/web";
import { Login } from "./Login";

let i = 0;
let text = "Login";
let speed = 100;
function typeWriter() {
    if (i < text.length) {
        let header = document.getElementById("header") as HTMLElement;
        header.innerHTML += text.charAt(i);
        i++;
        setTimeout(typeWriter, speed);
    }
}

const Client: Component<{ username: string, password: string }> = (props) => {
    type MessageObject = {
        type: MessageType,
        content: string,
        sender: string
    };

    type MessageType = number;

    const CHAT_MESSAGE: MessageType = 2;
    const SYSTEM_MESSAGE: MessageType = 1;
    const LOGGEDIN: string = "LOGGEDIN";
    const NOTLOGGEDIN: string = "NOTLOGGEDIN";

    let messageId = 1;
    const [currentMessage, setCurrentMessage] = createSignal<string>("");
    const [websocketConnection] = createSignal<WebSocket>(new WebSocket("ws://thatonedev.de/ws"));
    console.log("Attempting Connection");

    websocketConnection().onopen = () => {
        console.log("Succesfully Connected");
        websocketConnection().send(props.username + ":" + props.password);
    };

    const displayMessage = (message: string, message_id: string, author: string) => {
        message = author + ": " + message;
        return "<div class=\"w3-container\" id=" + message_id + "><p class=\"w3-panel w3-leftbar w3-border-aqua w3-round\">" + message + "</p></div>";
    };

    document.addEventListener('keyup', (event: KeyboardEvent) => {
        if (event.key == "Enter") {
            console.log("pressed enter key")
            let message = getMessageInput();
            sendMessage(message);
        }
    });

    websocketConnection().onmessage = event => {
        let receivedMessageString: string = event.data as string;
        let messageObject: MessageObject = JSON.parse(receivedMessageString);
        if (messageObject.type == CHAT_MESSAGE) {
            console.log(messageObject.content);
            console.log(messageObject.type);
            setCurrentMessage(messageObject.content);
            let chatHistory = document.getElementById("messages") as HTMLElement;
            if (chatHistory.innerHTML == "No Messages") {
                chatHistory.innerHTML = "";
            }
            chatHistory = document.getElementById("messages") as HTMLElement;
            chatHistory.innerHTML += displayMessage(currentMessage(), messageId.toString(), messageObject.sender);
            document.getElementById(messageId.toString())?.scrollIntoView();
            messageId++;
        } else if (messageObject.type == SYSTEM_MESSAGE && messageObject.content == LOGGEDIN) {
            console.log("Logged in successfully");
        } else if (messageObject.type == SYSTEM_MESSAGE && messageObject.content == NOTLOGGEDIN) {
            console.log("could not login terminated connection");
            websocketConnection().close();
            window.alert("Could not login. Terminated connection");
            setTimeout(() => {
                console.log("waiting");
                let body = document.getElementById("root") as HTMLElement;
                body.innerHTML = "";
                render(() => <Login />, body);
                typeWriter();
            }, 1000);
        }
    };

    websocketConnection().onclose = event => {
        console.log("Socket Closed Connection: ", event);
        websocketConnection().send("Client closed");
    };

    websocketConnection().onerror = error => {
        console.log("Socket Error: ", error);
    }

    const sendMessage = (msg: string) => {
        if (msg != "") {
            websocketConnection().send(msg);
        }
    };

    const getMessageInput = () => {
        let x = document.getElementById("inp") as HTMLInputElement;
        let value = x.value;
        x.value = "";
        return value;
    };

    document.title = "ChatApp";
    window.onload = function () {
        let i = 0;
        let text = "Chat App";
        let speed = 100;
        function typeWriter() {
            if (i < text.length) {
                let header = document.getElementById("header") as HTMLElement;
                header.innerHTML += text.charAt(i);
                i++;
                setTimeout(typeWriter, speed);
            }
        }
        typeWriter();
    };

    let w = new Worker("/content/fetchConnectedUsers.js");
    w.onmessage = ev => {
        let message = ev.data as string[];
        let list = document.getElementById("connectedUsers") as HTMLUListElement;
        list.innerHTML = "";
        for (let i = 0; i < message.length; i++) {
            let userData = message[i].split(":");
            let id = userData[0];
            let username = userData[1];
            list.innerHTML += "<li id=\"" + id + "\">" + username + "</li>";
        }
    };

    const showConnectedUsers = () => {
        let connectedUserList = document.getElementById("connectedUserList") as HTMLElement;
        if (connectedUserList.style.display === "none") {
            connectedUserList.style.display = "block";
            let btn = document.getElementById("showUserBtn") as HTMLElement;
            btn.innerHTML = "Hide Connected Users";
        } else {
            connectedUserList.style.display = "none";
            let btn = document.getElementById("showUserBtn") as HTMLElement;
            btn.innerHTML = "Show Connected Users";
        }
    };

    return (
        <>
            <div class="w3-content w3-container">
                <h2 id="header"></h2>
                <div id="messages" class="w3-border w3-round w3-border-black" style="height:250px;width:100%;overflow-y:auto;"></div><br></br>
                <input class="w3-input w3-round" id="inp" placeholder="Enter Message" style="width:300px"></input>
                <button class="w3-button w3-indigo w3-round" onclick={() => { sendMessage(getMessageInput()); }}>Send</button>
                <button class="w3-button w3-indigo w3-round" id="showUserBtn" onClick={() => { showConnectedUsers(); }}>Show Connected Users</button>
                <div id="connectedUserList" style="display:none;">
                    <h3>Connected Users:</h3>
                    <div class="w3-border w3-border-black w3-round" style="height:100px;width:25%;overflow-y:auto;">
                        <ul class="w3-ul w3-large" id="connectedUsers"></ul>
                    </div>
                </div>
            </div>
        </>);
};

export { Client };