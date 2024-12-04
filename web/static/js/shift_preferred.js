import { api } from '/js/api.js';
import { renderCalendar } from '/js/calendar.js';

window.addEventListener("DOMContentLoaded", async () => {
    const year = parseInt(document.getElementById('year').value);
    const month = parseInt(document.getElementById('month').value);
    await renderCalendar(year, month, document.querySelector('#calendar tbody'));
    addEventsToCalendar();
    await getShiftPreferreds(year, month);
    await getShiftPreferred(year, month);
    document.getElementById("save").addEventListener("click", save);
});

const addEventsToCalendar = () => {
    for (let i = 1; i <= 31; i++) {
        const cell = document.querySelector(`#calendar tbody td[data-day='${i}']`);
        if (cell) {
            cell.addEventListener('click', () => handleClickCell(cell, i));
        }
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
       alert('保存しました。')
       location.reload();
    } catch (e) {
        alert('保存に失敗しました。')
        console.error(e);
    }
};
