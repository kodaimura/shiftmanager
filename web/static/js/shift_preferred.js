import { api } from '/js/api.js';

let holidayCache = {};

window.addEventListener("DOMContentLoaded", async () => {
    const year = parseInt(document.getElementById('year').value);
    const month = parseInt(document.getElementById('month').value); 
    await fetchHolidays(year, month);
    renderCalendar(year, month);
    await getShiftPreferreds(year, month);
    await getShiftPreferred(year, month);
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
        const wrap = document.createElement('div');
        const div1 = document.createElement('div');
        const div2 = document.createElement('div');
        div1.classList.add('day');
        div2.classList.add('nums');

        const dayOfWeek = new Date(year, month - 1, day).getDay();
        if (dayOfWeek === 0 || isHoliday(year, month, day)) {
            cell.classList.add('holiday');
        } else if (dayOfWeek === 6) {
            cell.classList.add('saturday');
        } else {
            cell.classList.add('weekday');
        }

        cell.dataset.day = day;
        div1.textContent = day;
        div2.dataset.day = day;
        cell.addEventListener('click', () => handleClickCell(cell, day));

        wrap.appendChild(div1);
        wrap.appendChild(div2);
        cell.appendChild(wrap);

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

const handleClickCell = (cell, day) => {
    const form = document.getElementById('shift-preferred-form');
    let selectedDays = form.elements['dates'].value.split(',').filter(item => item !== '');
    if (!selectedDays.includes(String(day))) {
        selectedDays.push(day)
        cell.style.backgroundColor = 'yellow';
    } else {
        selectedDays = selectedDays.filter(d => d !== String(day));
        cell.style.backgroundColor = '';
    }
    form.elements['dates'].value = selectedDays.join(',');
};

const getShiftPreferreds = async (year, month) => {
    try {
        const result = await api.get(`shift_preferreds?year=${year}&month=${month}`);
        for (let data of result) {
            const dates = data.dates.split(',').map(Number);
            for (let date of dates) {
                const cell = document.querySelector(`#calendar div[data-day='${date}']`);
                if (cell) {
                    cell.textContent = (parseInt(cell.textContent) || 0) + 1;
                }
            };
        }
    } catch (e) {
        console.error(e);
    }
}

const getShiftPreferred = async (year, month) => {
    try {
        const response = await api.get(`shift_preferreds/me/${year}/${month}`);
        const dates = response.dates;
        if (dates) {
            const form = document.getElementById('shift-preferred-form');
            form.elements['dates'].value = dates;

            const selectedDays = dates.split(',');
            const cells = document.querySelectorAll('#calendar tbody td');
            cells.forEach(cell => {
                const day = cell.dataset.day;
                if (selectedDays.includes(day)) {
                    cell.style.backgroundColor = 'yellow';
                } else {
                    cell.style.backgroundColor = '';
                }
            });
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
       await api.post(`shift_preferreds/me/${year}/${month}`, body);
       window.location.replace('/');
    } catch (e) {
        console.error(e);
    }
};
