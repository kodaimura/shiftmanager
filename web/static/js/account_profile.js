import { api } from '/js/api.js';

window.addEventListener("DOMContentLoaded", function() {
    get();
    document.getElementById("save").addEventListener("click", save);
});

const get = async () => {
    try {
        const result = await api.get('account_profiles/me');
        const form = document.getElementById("profile-form");
        form.elements['display_name'].value = result.display_name;
        form.elements['account_role'].value = result.account_role;
    } catch (e) {
        console.error(e)
    }
}

const save = async () => {
    const form = document.getElementById("profile-form");
    const display_name = form.elements['display_name'].value;
    const account_role = form.elements['account_role'].value;
    
    const body = {
        display_name: display_name,
        account_role: account_role
    };

    try {
        await api.post('account_profiles/me', body);
        window.location.replace('/');
    } catch (e) {
        console.error(e)
    }
}