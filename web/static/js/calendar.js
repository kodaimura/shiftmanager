let holidayCache = null;

export const renderCalendar = async (year, month, tbody) => {
    await fetchHolidays(year, month);
    setCalendar(year, month, tbody);
};

const fetchHolidays = async (year, month) => {
    const url = `https://api.national-holidays.jp/${year}${String(month).padStart(2, '0')}`;
    try {
        holidayCache = null;
        const response = await fetch(url);
        if (response.ok) {
            holidayCache = await response.json();
        }
    } catch (error) {
        console.error('Error fetching holidays:', error);
    }
};

const isHoliday = (year, month, day) => {
    if (!holidayCache) {
        return false;
    }
    return holidayCache.some((holiday) => {
        return holiday.date === `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
    });
};

const setCalendar = async (year, month, tbody) => {
    tbody.innerHTML = '';

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
        const div2= document.createElement('div');
        div1.classList.add('cell-day');
        div2.classList.add('cell-body');

        const dayOfWeek = new Date(year, month - 1, day).getDay();
        if (dayOfWeek === 0 || isHoliday(year, month, day)) {
            div1.classList.add('holiday');
        } else if (dayOfWeek === 6) {
            div1.classList.add('saturday');
        } else {
            div1.classList.add('weekday');
        }
        cell.dataset.day = day;
        div1.textContent = day;
        div2.dataset.day = day;

        wrap.appendChild(div1);
        wrap.appendChild(div2);
        cell.appendChild(wrap);

        row.appendChild(cell);
        if ((firstDay.getDay() + day) % 7 === 0) {
            tbody.appendChild(row);
            row = document.createElement('tr');
        }
    }

    if (row.children.length > 0) {
        tbody.appendChild(row);
    }
};