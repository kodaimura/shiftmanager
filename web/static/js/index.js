import { getErrorStatus, handleResponse, handleError } from '/js/script.js';

window.addEventListener("DOMContentLoaded", function() {
    getProfile();

    document.getElementById("shift-preferred").addEventListener("click", confirmProfileRegistered);
});

const getProfile = () => {
    fetch('/api/account_profile', {
        method: 'GET',
        headers: {"Content-Type": "application/json"},
    })
    .then(handleResponse)
    .then((data) => {
        document.getElementById('display_name').textContent = data.display_name;
        document.getElementById('account_role').textContent = data.account_role;
    })
    .catch(handleError)
}

const confirmProfileRegistered = (event) => {
    if (
        document.getElementById('display_name').textContent === '' ||
        document.getElementById('account_role').textContent === ''
    ) {
        event.preventDefault();
        alert('先にプロフ登録を行ってください。');
    }
}