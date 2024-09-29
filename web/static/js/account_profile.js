import { getErrorStatus, handleResponse, handleError } from '/js/script.js';

window.addEventListener("DOMContentLoaded", function() {
    get();
    document.getElementById("save").addEventListener("click", save);
});

const get = () => {
    fetch('/api/account_profile', {
        method: 'GET',
        headers: {"Content-Type": "application/json"},
    })
    .then(handleResponse)
    .then((data) => {
        const form = document.getElementById("profile-form");
        form.elements['display_name'].value = data.display_name;
        form.elements['account_role'].value = data.account_role;
    })
    .catch(handleError)
}

const save = () => {
    const form = document.getElementById("profile-form");
    const display_name = form.elements['display_name'].value;
    const account_role = form.elements['account_role'].value;
    
    const body = {
        display_name: display_name,
        account_role: account_role
    };

    fetch('/api/account_profile', {
        method: 'POST',
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(body)
    })
    .then(handleResponse)
    .then(() => {
        window.location.replace('/');
    })
    .catch(handleError)
}