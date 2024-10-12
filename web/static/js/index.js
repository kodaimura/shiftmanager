import { api } from '/js/api.js';
import { getJaTime } from '/js/script.js';

let holidayCache = {};
let targetDate;

window.addEventListener("DOMContentLoaded", async () => {
    targetDate = new Date(getJaTime().getFullYear(), getJaTime().getMonth() + 1);
    getProfile();

    await fetchHolidays();
    renderCalendar();
    renderYearMonth();
    renderLinks();
    getShiftPreferred();
    
    document.getElementById('prev-month').addEventListener('click', handlePrevMonth);
    document.getElementById('next-month').addEventListener('click', handleNextMonth);
    document.getElementById("shift-preferred").addEventListener("click", confirmProfileRegistered);
});

const fetchHolidays = async () => {
    const year = targetDate.getFullYear();
    const month = targetDate.getMonth() + 1;
    const url = `https://api.national-holidays.jp/${year}${String(month).padStart(2, '0')}`;
    try {
        holidayCache[`${year}-${month}`] = [];
        const response = await fetch(url);
        if (response.ok) {
            holidayCache[`${year}-${month}`] = await response.json();
        }
    } catch (error) {
        console.error('Error fetching holidays:', error);
    }
};

const isHoliday = (year, month, day) => {
    if (!holidayCache[`${year}-${month}`]) {
        console.error(`祝日データがロードされていません: ${year}-${month}`);
        return false;
    }
    return holidayCache[`${year}-${month}`].some((holiday) => {
        return holiday.date === `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
    });
};

const getProfile = async () => {
    try {
        const result = await api.get('account_profiles/me');
        document.getElementById('display_name').textContent = result.display_name;
        document.getElementById('account_role').textContent = result.account_role;
    } catch (e) {
        console.error(e)
    }
}

const renderCalendar = async () => {
    const year = targetDate.getFullYear();
    const month = targetDate.getMonth() + 1;

    const calendarBody = document.querySelector('#calendar tbody');
    calendarBody.innerHTML = '';

    const firstDay = new Date(year, month - 1, 1);
    const lastDay = new Date(year, month, 0);

    let row = document.createElement('tr');
    for (let i = 0; i < firstDay.getDay(); i++) {
        row.appendChild(document.createElement('td'));
    }

    for (let day = 1; day <= lastDay.getDate(); day++) {
        const div1 = document.createElement('div');
        const div2= document.createElement('div');
        div1.classList.add('day');
        div2.classList.add('names');

        const dayOfWeek = new Date(year, month - 1, day).getDay();
        if (dayOfWeek === 0 || isHoliday(year, month, day)) {
            div1.classList.add('holiday');
        } else if (dayOfWeek === 6) {
            div1.classList.add('saturday');
        } else {
            div1.classList.add('weekday');
        }
        div1.textContent = day;
        div2.dataset.day = day;

        const cell = document.createElement('td');
        cell.appendChild(div1);
        cell.appendChild(div2);

        row.appendChild(cell);
        if ((firstDay.getDay() + day) % 7 === 0) {
            calendarBody.appendChild(row);
            row = document.createElement('tr');
        }
    }

    if (row.children.length > 0) {
        calendarBody.appendChild(row);
    }
};

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
                    cell.innerHTML += `${accountId}&nbsp;`;
                }
            };
        }
    } catch (e) {
        console.error(e);
    }
}

 const renderYearMonth = () => {
    const year = targetDate.getFullYear();
    const month = targetDate.getMonth() + 1;

    document.getElementById('year').textContent = year;
    document.getElementById('month').textContent = month
 };

const renderLinks = () => {
    const year = targetDate.getFullYear();
    const month = targetDate.getMonth() + 1;

    document.getElementById('shift-preferred').href = `/shift_preferreds/me/${year}/${month}`
    document.getElementById('shift-generate').href = `/shift/${year}/${month}/generate`
    document.getElementById('shift-edit').href = `/shift/${year}/${month}/edit`
};

const handlePrevMonth = async () => {
    targetDate.setMonth(targetDate.getMonth() - 1);
    await fetchHolidays();
    renderCalendar();
    renderYearMonth();
    renderLinks();
    getShiftPreferred();
};

const handleNextMonth = async () => {
    targetDate.setMonth(targetDate.getMonth() + 1);
    await fetchHolidays();
    renderCalendar();
    renderYearMonth();
    renderLinks();
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