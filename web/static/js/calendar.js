
document.addEventListener('DOMContentLoaded', (event) => {
    const year = document.getElementById('year').value
    const month = document.getElementById('month').value
    const startDate = new Date(year, month - 1, 1)
    const lastDate = new Date(year, month, 0)

    document.getElementById('calendar').innerHTML = makeCalendar(startDate, lastDate)
    reflectHolidaysToCalendar(year, month, lastDate.getDate())
    reflectSelectDaysToCalendar(getSelectDays())
    addEventToCalendarCell(lastDate.getDate())
    addInfoToCalendarCell(startDate)
}) 


const makeCalendar = (startDate, lastDate) => {
    const startDayOfWeek = startDate.getDay()
    const lastDay = lastDate.getDate()
    let day = 1 
    let calendar = ''
    calendar += 
    `<table><tr><th>日</th><th>月</th><th>火</th><th>水</th>
    <th>木</th><th>金</th><th>土</th></tr>`

    for (let w = 0; w < 6; w++) {
        calendar += '<tr>'

        for (let d = 0; d < 7; d++) {
            if (w == 0 && d < startDayOfWeek || day > lastDay) {
                calendar += '<td></td>'
            } else {
                calendar += `<td id="d${day}">` + day + 
                `<div id="info${day}" class="cellinfo">0</div></td>`
                day++
            }
        }
        calendar += '</tr>'
    }
    return calendar + '</table>'
}


const isHoliday = (year, month, day) => {
    const date = new Date(year, month - 1, day)
    let bool = false
    try {
        bool = JapaneseHolidays.isHoliday(date) !== undefined
    }finally {
        return bool
    }
}


const reflectHolidaysToCalendar = (year, month, lastDay) => {
    for (let i = 1; i <= lastDay; i++) {
        if (isHoliday(year, month, i)) {
            document.getElementById(`d${i}`).classList.add('holiday')
        }
    }  
}


const reflectSelectDaysToCalendar = (days) => {
    let dayCell
    let classList
    for (d of days){
        dayCell = document.getElementById(`d${d}`)
        classList = (dayCell === null)? null : dayCell.classList
        if (classList != null && !classList.contains('selected')) {
            classList.add('selected')
        }
    }
}


const getSelectDays = () => {
    days = document.getElementById('selectdays').value
    return (days === '')? [] : days.split(',')
}


const setSelectDays = (days) => {
    document.getElementById('selectdays').value = days.join(',')
}


const addEventToCalendarCell = (lastDay) => {
    for (let i = 1; i <= lastDay; i++) {
        document.getElementById(`d${i}`).addEventListener('click', (e) => {
            let days = getSelectDays()
            let dayCell = document.getElementById(`d${i}`)
            let classList = dayCell.classList
            if (classList.contains('selected')) {
                classList.remove('selected')
                days = days.filter(d => d !== dayCell.id.substring(1))
            }else{
                classList.add('selected')
                days.push(dayCell.id.substring(1))
            }
            setSelectDays(days)
        })
    }
}


const addInfoToCalendarCell = (startDate) => {
    const year = startDate.getFullYear()
    const month = startDate.getMonth() + 1
    fetch(`/api/group/workables/${year}/${month}`)
    .then(response => response.json())
    .then(data => reflectInfo(data))
}


const reflectInfo = (data) => {
    let workabledays
    let dayCell
    for (d of data){
        workabledays = (d.workabledays === '')? [] : d.workabledays.split(',')
        for (day of workabledays){
            infoCell = document.getElementById(`info${day}`)
            infoCell.innerHTML = parseInt(infoCell.innerHTML) + 1
        }
    }
}
