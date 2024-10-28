import { api } from '/js/api.js';
import { renderCalendar } from '/js/calendar.js';

let displayNameMap = {};

window.addEventListener("DOMContentLoaded", async () => {
    const year = parseInt(document.getElementById('year').value);
    const month = parseInt(document.getElementById('month').value); 

    await renderCalendarCustom(year, month);
    await renderModalCalendar(year, month);
    await getDisplayNameMap();
    getShiftPreferred(year, month);
    getShift(year, month);

    document.getElementById("generate").addEventListener("click", async (event) => {
        event.target.disabled = true;
        await postShiftgenerate();
        event.target.disabled = false;
    });
    document.getElementById("save").addEventListener("click", async (event) => {
        event.target.disabled = true;
        await save();
        event.target.disabled = false;
    });
});

const renderCalendarCustom = async (year, month) => {
    await renderCalendar(year, month, document.querySelector('#calendar tbody'));

    for (let i = 1; i <= 31; i++) {
        const cell = document.querySelector(`#calendar .cell-body[data-day='${i}']`);
        if (cell) {
            const input = document.createElement('input');
            input.dataset.day = i;
            input.onchange = setCounter;

            input.setAttribute('data-bs-toggle', 'tooltip');
            input.setAttribute('aria-label', 'candidate');
            input.setAttribute('title', '');
            const tooltip = new bootstrap.Tooltip(input);
            input.focus = () => tooltip.show();
            input.blur = () => tooltip.hide();
            input.oninput = () => {
                input.value = input.value.replace(/　/g, ' ');
            };

            cell.appendChild(input);
        }
    }
};

const renderModalCalendar = async (year, month) => {
    await renderCalendar(year, month, document.querySelector('#modal-calendar tbody'));

    for (let i = 1; i <= 31; i++) {
        const cell = document.querySelector(`#modal-calendar td[data-day='${i}']`);
        if (cell) {
            cell.addEventListener('click', () => handleClickCell(cell, i));
        }
    }
};

const handleClickCell = (cell, day) => {
    const form = document.getElementById('generate-form');
    let storeHoliday = form.elements['store_holiday'].value.split(',').filter(item => item !== '');
    if (!storeHoliday.includes(String(day))) {
        storeHoliday.push(day)
        cell.style.backgroundColor = 'gray';
    } else {
        storeHoliday = storeHoliday.filter(d => d !== String(day));
        cell.style.backgroundColor = '';
    }
    form.elements['store_holiday'].value = storeHoliday.join(',');
};

const getShift = async (year, month) => {
    try {
        const result = await api.get(`shifts/${year}/${month}`);
        if (result.shift_data) {
            const ls = result.shift_data.split(',');
            for (let i = 1; i <= 31; i++) {
                const input = document.querySelector(`#calendar input[data-day='${i}']`);
                if (input) {
                    input.value = ls[i - 1];
                }
            };
        }
        if (result.store_holiday) {
            const form = document.getElementById('generate-form');
            form.elements['store_holiday'].value = result.store_holiday;

            const ls = result.store_holiday.split(',').filter(item => item !== '');
            for (let i = 1; i <= 31; i++) {
                const cell = document.querySelector(`#modal-calendar td[data-day='${i}']`);
                if (ls.includes(String(i))) {
                    cell.style.backgroundColor = 'gray';
                }
            };
        }
        setCounter();
    } catch (e) {
        console.error(e);
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

const getShiftPreferred = async (year, month) => {
    try {
        const result = await api.get(`shift_preferreds?year=${year}&month=${month}`);
        let candidate = new Array(31).fill([]);
        for (let data of result) {
            const dates = data.dates.split(',').map(Number);
            for (let date of dates) {
                const cell = document.querySelector(`#modal-calendar .cell-body[data-day='${date}']`);
                if (cell) {
                    cell.textContent = (parseInt(cell.textContent) || 0) + 1;
                }
                candidate[date - 1] = [...candidate[date - 1], displayNameMap[data.account_id]];
            };
        }
        for (let i = 1; i <= 31; i++) {
            const input = document.querySelector(`#calendar input[data-day='${i}']`);
            if (input) {
                input.setAttribute('title', candidate[i - 1].join(', '));
                new bootstrap.Tooltip(input);
            }
        }
    } catch (e) {
        console.error(e);
    }
}

const postShiftgenerate = async () => {
    const year = parseInt(document.getElementById('year').value);
    const month = parseInt(document.getElementById('month').value); 
    const form = document.getElementById('generate-form');
    const body = {
        store_holiday: form.elements['store_holiday'].value,
    };

    try {
        await api.post(`shifts/${year}/${month}/generate`, body);
        const modal = bootstrap.Modal.getInstance('#generate-modal');
        modal.hide();
        alert('保存しました。')
        getShift(year, month);
    } catch (e) {
        alert('保存に失敗しました。')
        console.error(e);
    }
}

const save = async () => {
    const year = parseInt(document.getElementById('year').value);
    const month = parseInt(document.getElementById('month').value);
    let shift = '';
    for (let i = 1; i <= 31; i++) {
        const input = document.querySelector(`#calendar input[data-day='${i}']`);
        if (input) {
            shift += input.value + ','
        } else {
            shift += ','
        }
    };
    shift = shift.replace(/,\s*$/, '')
    shift = shift.replace(/\u3000/g, ' ');
    shift = shift.replace(/\s+/g, ' ');
    shift = shift.trim();

    const form = document.getElementById('generate-form');
    const body = {
        shift_data: shift,
        store_holiday: form.elements['store_holiday'].value,
    };
    try {
        await api.post(`shifts/${year}/${month}`, body);
        alert('保存しました。')
        getShift(year, month);
    } catch (e) {
        alert('保存に失敗しました。')
        console.error(e);
    }
};

const setCounter = () => {
    let countMap = {};
    for (let i = 1; i <= 31; i++) {
        const input = document.querySelector(`#calendar input[data-day='${i}']`);
        if (input) {
            const names = input.value.split(/\s+/).filter(Boolean);
            for (const x of names) {
                if (countMap[x]) {
                    countMap[x] += 1
                } else {
                    countMap[x] = 1
                }
            }
        }
    };

    let str = ' ';
    for (const x in countMap) {
        str += `${x}:${countMap[x]}  `;
    }

    const tooltipIcon = document.getElementById('counter');
    tooltipIcon.setAttribute('title', str);
    new bootstrap.Tooltip(tooltipIcon);
}