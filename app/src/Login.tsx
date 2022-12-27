import { Component } from "solid-js";

const LoginHtml: Component = () => {
    document.title = "Login";
    return (
        <>
            <h1 id="header" class="w3-center"></h1>
            <div class="w3-content w3-container">
                <form action="http://thatonedev.de/" method="post">
                    <input type="text" placeholder="Enter Username" class="w3-input w3-round" id="username" name="username"/>
                    <br />
                    <input type="password" placeholder="Enter Password" class="w3-input w3-round" id="password" name="password"/>
                    <br />
                    <input type="submit" class="w3-button w3-round w3-indigo" value="Login"/>
                </form>
            </div>
        </>
    );
};

const Login: Component = () => {
    window.onload = function () {
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
        typeWriter();
    };
    return (
        <LoginHtml />
    );
};

export { Login };