/* @refresh reload */
import { render } from 'solid-js/web';
import { Client } from './Client';

import { Login } from './Login';

function getCookie(cookieName: string): string {
    let cookies = document.cookie;
    let name = cookieName + "=";
    let ca = cookies.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

function isCookieSet(name: string): boolean {
    if(getCookie(name) == "") {
        return false;
    } else {
        return true;
    }
}

function renderLoginPage() {
    render(() => <Login />, document.getElementById('root') as HTMLElement);
}

function renderClientPage() {
    let logindata = getCookie("logindata");
    let passwordAndUsername = logindata.split(':');
    render(() => Client({username: passwordAndUsername[0], password: passwordAndUsername[1]}), document.getElementById('root') as HTMLElement);
}

if (isCookieSet("logindata")) {
    renderClientPage();
} else {
    renderLoginPage();
}
