import { getErrorStatus, handleResponse, handleError } from '/js/script.js';

window.addEventListener("DOMContentLoaded", function() {
    document.getElementById("login").addEventListener("click", login);
});


const login = () => {
    const form = document.getElementById("login-form");
    const account_name = form.elements['account_name'].value;
    const account_password = form.elements['account_password'].value;

    const body = {
        account_name: account_name,
        account_password: account_password
    };

    fetch('/api/login', {
        method: 'POST',
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(body)
    })
    .then(handleResponse)
    .then(() => {
        window.location.replace('/');
    })
    .catch(handleError)
    .catch((error) => {
        const status = getErrorStatus(error);
        document.getElementById("error").innerHTML = (status === 401)
        ? "ユーザ名またはパスワードが異なります。" 
        : "ログインに失敗しました。";
    });
}