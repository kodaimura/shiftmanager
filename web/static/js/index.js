import { api } from '/js/api.js';
import { getJaTime } from '/js/script.js';
import { renderCalendar } from '/js/calendar.js';

let displayNameMap = {};
let targetDate;

window.addEventListener("DOMContentLoaded", async () => {
    targetDate = new Date(getJaTime().getFullYear(), getJaTime().getMonth() + 1);

    renderYearMonth();
    renderLinks();
    await renderCalendarCustom();
    await getDisplayNameMap();
    getProfile();
    getShiftPreferred();
    
    document.getElementById('prev-month').addEventListener('click', handlePrevMonth);
    document.getElementById('next-month').addEventListener('click', handleNextMonth);
    document.getElementById("shift-preferred").addEventListener("click", confirmProfileRegistered);
});

const renderCalendarCustom = async () => {
    await renderCalendar(
        targetDate.getFullYear(), 
        targetDate.getMonth() + 1, 
        document.querySelector('#calendar tbody')
    );
}

const renderYearMonth = () => {
    const year = targetDate.getFullYear();
    const month = targetDate.getMonth() + 1;

    document.getElementById('year').textContent = year;
    document.getElementById('month').textContent = month
 }

 const renderLinks = () => {
    const year = targetDate.getFullYear();
    const month = targetDate.getMonth() + 1;

    document.getElementById('shift-preferred').href = `/shift_preferreds/me/${year}/${month}`
    document.getElementById('shift').href = `/shifts/${year}/${month}`
};

const getProfile = async () => {
    const roleMap = {
        '1': 'キッチン',
        '2': 'ホール',
        '3': 'キッチン・ホール',
    }
    try {
        const result = await api.get('account_profiles/me');
        document.getElementById('display_name').textContent = result.display_name;
        document.getElementById('account_role').textContent = roleMap[result.account_role];
    } catch (e) {
        console.error(e)
    }
}

const getDisplayNameMap = async () => {
    try {
        const result = await api.get('account_profiles');
        for (let data of result) {
            const accountId = data.account_id;
            const displayName = data.display_name;
            displayNameMap[accountId] = displayName;
        }
    } catch (e) {
        console.error(e);
    }
}

const getShiftPreferred = async () => {
    const year = targetDate.getFullYear();
    const month = targetDate.getMonth() + 1;
    try {
        const result = await api.get(`shift_preferreds?year=${year}&month=${month}`);
        for (let data of result) {
            const accountId = data.account_id;
            const dates = data.dates.split(',').map(Number);
            for (let date of dates) {
                const cell = document.querySelector(`div[data-day='${date}']`);
                if (cell) {
                    if (!cell.classList.contains('highlight')) {
                        cell.classList.add('highlight');
                    }
                    cell.innerHTML += `${displayNameMap[accountId]}&nbsp;`;
                }
            };
        }
    } catch (e) {
        console.error(e);
    }
}

const handlePrevMonth = async () => {
    targetDate.setMonth(targetDate.getMonth() - 1);
    renderYearMonth();
    renderLinks();
    await renderCalendarCustom();
    getShiftPreferred();
};

const handleNextMonth = async () => {
    targetDate.setMonth(targetDate.getMonth() + 1);
    renderYearMonth();
    renderLinks();
    await renderCalendarCustom();
    getShiftPreferred();
};

const confirmProfileRegistered = (event) => {
    if (
        document.getElementById('display_name').textContent === '' ||
        document.getElementById('account_role').textContent === ''
    ) {
        event.preventDefault();
        alert('先にプロフ登録を行ってください。');
    }
}