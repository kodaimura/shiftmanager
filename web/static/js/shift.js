import { api } from '/js/api.js';

let holidayCache = {};

window.addEventListener("DOMContentLoaded", async () => {
    const year = parseInt(document.getElementById('year').value);
    const month = parseInt(document.getElementById('month').value); 
    await fetchHolidays(year, month);
    renderCalendar(year, month);
    renderModalCalendar(year, month);
    getShiftPreferred(year, month);
    getShift(year, month);

    document.getElementById("generate").addEventListener("click", postShiftgenerate);
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
        const input = document.createElement('input');
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
        input.dataset.day = day;

        div2.appendChild(input);
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

const getShift = async (year, month) => {
    try {
        const result = await api.get(`shifts/${year}/${month}`);
        const dates = result.shift_data.split(',');
        for (let date of dates) {
            const input = document.querySelector(`#calendar input[data-day='${date}']`);
            if (input) {
                input.value = data;
            }
        };
    } catch (e) {
        console.error(e);
    }
}

const renderModalCalendar = (year, month) => {
    const calendarBody = document.querySelector('#modal-calendar tbody');
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
        div2.addEventListener('click', () => handleClickCell(div2, day));

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
    const form = document.getElementById('generate-form');
    let storeHoliday = form.elements['store_holiday'].value.split(',');
    if (!storeHoliday.includes(String(day))) {
        cell.style.backgroundColor = 'gray';
    } else {
        storeHoliday = storeHoliday.filter(d => d !== String(day));
        cell.style.backgroundColor = '';
    }
    form.elements['store_holiday'].value = storeHoliday.join(',');
};

const getShiftPreferred = async (year, month) => {
    try {
        const result = await api.get(`shift_preferreds?year=${year}&month=${month}`);
        for (let data of result) {
            const dates = data.dates.split(',').map(Number);
            for (let date of dates) {
                const cell = document.querySelector(`#modal-calendar div[data-day='${date}']`);
                if (cell) {
                    cell.textContent = (parseInt(cell.textContent) || 0) + 1;
                }
            };
        }
    } catch (e) {
        console.error(e);
    }
}

const postShiftgenerate = async () => {
    const year = parseInt(document.getElementById('year').value);
    const month = parseInt(document.getElementById('month').value); 
    const form = document.getElementById('generate-form');
    let storeHoliday = form.elements['store_holiday'].value.split(',');

    const body = {
        store_holiday: storeHoliday,
    };

    try {
        await api.post(`shifts/${year}/${month}`, body);
    } catch (e) {
        console.error(e);
    }
}