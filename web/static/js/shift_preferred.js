import { api } from '/js/api.js';

let holidayCache = {};
let selectedDays = [];

window.addEventListener("DOMContentLoaded", async () => {
    const year = parseInt(document.getElementById('year').value);
    const month = parseInt(document.getElementById('month').value); 
    await fetchHolidays(year, month);
    generateCalendar(year, month);
    const data = await getShiftPreferred(year, month);

    document.getElementById("save").addEventListener("click", save);
});

const fetchHolidays = async (year, month) => {
    const url = `https://api.national-holidays.jp/${year}${month}`;
    try {
        const response = await fetch(url);
        holidayCache[`${year}-${month}`] = await response.json();
    } catch (error) {
        console.error('Error fetching holidays:', error);
    }
};

const isHoliday = (year, month, day) => {
    if (!holidayCache[`${year}-${month}`]) {
        console.error(`祝日データがロードされていません: ${year}-${month}`);
        return false;
    }
    const targetDate = `${year}-${month}-${day}`;
    return holidayCache[`${year}-${month}`].some(holiday => holiday.date === targetDate);
};

const generateCalendar = (year, month) => {
    const calendarDiv = document.getElementById('calendar');
    calendarDiv.innerHTML = '';

    const firstDay = new Date(year, month - 1, 1);
    const lastDay = new Date(year, month, 0);

    const table = document.createElement('table');
    table.classList.add('calendar-table');

    const header = document.createElement('tr');
    ['日', '月', '火', '水', '木', '金', '土'].forEach(day => {
        const th = document.createElement('th');
        th.textContent = day;
        header.appendChild(th);
    });
    table.appendChild(header);

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

        cell.addEventListener('click', () => handleDateClick(cell, day));

        row.appendChild(cell);
        if ((firstDay.getDay() + day) % 7 === 0) {
            table.appendChild(row);
            row = document.createElement('tr');
        }
    }

    if (row.children.length > 0) {
        table.appendChild(row);
    }

    calendarDiv.appendChild(table);
};

const handleDateClick = (cell, day) => {
    if (!selectedDays.includes(day)) {
        selectedDays.push(day);
        cell.style.backgroundColor = 'yellow';
    } else {
        selectedDays = selectedDays.filter(d => d !== day);
        cell.style.backgroundColor = '';
    }

    const form = document.getElementById('shift-preferred-form');
    form.elements['dates'].value = selectedDays.join(',');
};

const getShiftPreferred = async (year, month) => {
    try {
        const result = await api.get(`shift_preferred/${year}/${month}`);
        const form = document.getElementById('shift-preferred-form');
        form.elements['dates'].value = result.dates;
        selectedDays = result.dates.split(',')
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
        const result = await api.post(`shift_preferred/${year}/${month}`, body);
        console.log(result);
    } catch (e) {
        console.error(e);
    }
};
