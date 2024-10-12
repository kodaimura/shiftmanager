import { api } from '/js/api.js';
import { getJaTime } from '/js/script.js';

let holidayCache = {};

window.addEventListener("DOMContentLoaded", async () => {
    const year = parseInt(document.getElementById('year').value);
    const month = parseInt(document.getElementById('month').value); 
    await fetchHolidays(year, month);
    await getShiftPreferred(year, month);
    

    renderCalendar(year, month);
    document.getElementById("save").addEventListener("click", save);
});

const fetchHolidays = async (year, month) => {
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

const renderCalendar = (year, month) => {
    const form = document.getElementById('shift-preferred-form');
    const selectedDays = form.elements['dates'].value.split(',');
    const calendarBody = document.querySelector('#calendar tbody');
    calendarBody.innerHTML = '';

    const firstDay = new Date(year, month - 1, 1);
    const lastDay = new Date(year, month, 0);

    let row = document.createElement('tr');
    for (let i = 0; i < firstDay.getDay(); i++) {
        row.appendChild(document.createElement('td'));
    }

    for (let day = 1; day <= lastDay.getDate(); day++) {
        const cell = document.createElement('td');
        cell.textContent = day;

        const dayOfWeek = new Date(year, month - 1, day).getDay();
        if (dayOfWeek === 0 || isHoliday(year, month, day)) {
            cell.classList.add('holiday');
        } else if (dayOfWeek === 6) {
            cell.classList.add('saturday');
        } else {
            cell.classList.add('weekday');
        }

        if (selectedDays.includes(String(day))) {
            cell.style.backgroundColor = 'yellow';
        }

        cell.addEventListener('click', () => handleDateClick(cell, day));

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

const handleDateClick = (cell, day) => {
    const form = document.getElementById('shift-preferred-form');
    let selectedDays = form.elements['dates'].value.split(',');
    if (!selectedDays.includes(day)) {
        selectedDays.push(day);
        cell.style.backgroundColor = 'yellow';
    } else {
        selectedDays = selectedDays.filter(d => d !== day);
        cell.style.backgroundColor = '';
    }
    form.elements['dates'].value = selectedDays.join(',');
};

const getShiftPreferred = async (year, month) => {
    try {
        const response = await api.get(`shift_preferred/${year}/${month}`);
        const dates = response.dates;
        if (dates) {
            const form = document.getElementById('shift-preferred-form');
            form.elements['dates'].value = dates;
        }
    } catch (e) {
        console.error(e);
    }
};

const save = async () => {
    const year = parseInt(document.getElementById('year').value);
    const month = parseInt(document.getElementById('month').value); 
    const form = document.getElementById('shift-preferred-form');
    const body = {
        dates: form.elements['dates'].value,
        notes: '',
    };
    try {
       await api.post(`shift_preferred/${year}/${month}`, body);
    } catch (e) {
        console.error(e);
    }
};
