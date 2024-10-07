import { api } from '/js/api.js';

window.addEventListener("DOMContentLoaded", function() {
    getProfile();

    document.getElementById("shift-preferred").addEventListener("click", confirmProfileRegistered);
});

const getProfile = async () => {
    try {
        const result = await api.get('account_profile');
        const form = document.getElementById("profile-form");
        document.getElementById('display_name').textContent = result.display_name;
        document.getElementById('account_role').textContent = result.account_role;
    } catch (e) {
        console.error(e)
    }
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