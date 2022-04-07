document.addEventListener('DOMContentLoaded', (event) => {
    const year = document.getElementById('year').value
    const month = document.getElementById('month').value

    const startDate = new Date(year, month - 1, 1)
    const lastDate = new Date(year, month, 0)
    const lastDay = lastDate.getDate()

    document.getElementById('calendar').innerHTML = makeCalendar(startDate, lastDate)
    reflectHolidaysToCalendar(year, month, lastDay)    
    reflectShiftToCalendar()
    reflectStoreHolidaysToCalendar()
    countNames(lastDay)
    addEventToCellInput(lastDay)

    document.getElementById('submit').addEventListener('click',saveShift)
}) 


const makeCalendar = (startDate, lastDate) => {
    const startDayOfWeek = startDate.getDay()
    const lastDay = lastDate.getDate()
    let day = 1 
    let calendar = ''
    calendar += 
    `<table class="table"><thead class="table-secondary">
    <tr><th>日</th><th>月</th><th>火</th><th>水</th>
    <th>木</th><th>金</th><th>土</th></tr></thead>`

    for (let w = 0; w < 6; w++) {
        calendar += '<tr>'

        for (let d = 0; d < 7; d++) {
            if (w == 0 && d < startDayOfWeek || day > lastDay) {
                calendar += '<td></td>'
            } else {
                calendar += `<td id="d${day}">${day}
                <div><input type="text" id="s${day}"></div></td>`
                day++
            }
        }
        calendar += '</tr>'
    }
    return calendar + '</table>'
}


const isHoliday = (year, month, day) => {
    const date = new Date(year, month - 1, day)
    return JapaneseHolidays.isHoliday(date) !== undefined
}


const reflectHolidaysToCalendar = (year, month, lastDay) => {
    for (let i = 1; i <= lastDay; i++) {
        if (isHoliday(year, month, i)) {
            document.getElementById(`d${i}`).classList.add('holiday')
        }
    }  
} 


const reflectShiftToCalendar = () => {
    let shift = document.getElementById('shift').value.split(',')
    for ([i, s] of shift.entries()) {
        let input = document.getElementById(`s${i + 1}`)
        if (input !== null) {
            input.value = " " + s
        }
    }
}


const reflectStoreHolidaysToCalendar = () => {
    let storeholiday = document.getElementById('storeholiday').value.split(',')
    for (h of storeholiday) {
        let d = document.getElementById(`d${h}`)
        if (d !== null) {
            d.classList.add('storeholiday')
        }
    }
}


const countNames = (lastDay) => {
    let counts = {}
    let str = ""
    for (let i = 1; i <= lastDay; i++){
        let input = document.getElementById(`s${i}`)
        if (input !== null){
            str += " " + input.value
        }
    }

    let names = str.split(' ')

    for (let n of names){
        if (n !== "") {
            if (counts[n]){
                counts[n] += 1 
            }else{
                counts[n] = 1
            }
        }
    }
    let count = ''
    for (let n in counts){
        count += ` ${n}: ${counts[n]}/` 
    }
    document.getElementById('count').innerHTML = count
}


const addEventToCellInput = (lastDay) => {
    let candidate = document.getElementById('candidate').value.split(',')

    for (let i = 1; i <= lastDay; i++) {
        let target = document.getElementById(`s${i}`)
        target.addEventListener('focus', (e) => {
            document.getElementById('message').innerHTML = "書き換え候補: " + candidate[i - 1]
        })
        target.addEventListener('focusout', (e) => {
            document.getElementById('message').innerHTML = ""
        })
        target.addEventListener('input', (e) => {
            countNames(lastDay)
            let shift = document.getElementById('shift')
            let ls = shift.value.split(',')
            e.target.value = e.target.value.replace("　", " ").replace(",", " ")
            ls[i - 1] = e.target.value
            shift.value = ls.join(',')

            let button = document.getElementById("submit")
            button.value = "保存"
            button.disabled = false
        })
    }
}


const saveShift = (event) => {
    const year = document.getElementById('year').value
    const month = document.getElementById('month').value
    const shift = document.getElementById('shift').value
    const storeholiday = document.getElementById('storeholiday').value
    fetch(`/api/shift/${year}/${month}`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({shift, storeholiday})
    })
    .then(response => {
        if (response.status === 200) {
             event.target.value = "保存済"
             event.target.disabled = true
        } else {
            document.getElementById("error").innerHTML = "保存失敗"
        }
    })
}

