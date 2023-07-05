/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/6/23 上午12:22.
 *  * Author: guojia(https://github.com/guojia99)
 */


let projectValue = ""
let contestValue = ""
let contestsList = []
let playerList = []
let scoreList = []

function enableInputScoreNumber(numbers) {
    for (let i = 1; i < 6; i++) {
        let value = $("#input-score-number" + `${i}-value`)
        if (numbers.indexOf(i) >= 0) {
            value.removeAttr("disabled")
            continue
        }
        value.attr("disabled", true)
        value.val("") // disable的时候清空input的内容
    }
}

function validateTimeFormat(time) {
    // DNF
    if (time === 'DNF') {
        return true;
    }
    // 校验纯秒数格式
    if (/^\d+(\.\d+)?$/.test(time)) {
        return true;
    }
    // 校验分+秒格式
    if (/^\d{1,2}[:：][0-5]?\d(\.\d+)?$/.test(time)) {
        return true;
    }
    // 校验时+分+秒格式
    return (/^\d{1,2}[:：][0-5]?\d[:：]\d{2}(\.\d+)?$/.test(time));
}

function parseTimeToSeconds(time) {
    const data = function (time) {
        if (time === 'DNF') {
            return 0;
        }
        // 解析纯秒数格式
        if (/^\d+(\.\d+)?$/.test(time)) {
            return parseFloat(time);
        }
        // 解析分+秒格式
        if (/^\d{1,2}[:：]\d{2}(\.\d+)?$/.test(time)) {
            const [minutes, seconds] = time.split(/[:：]/);
            return parseFloat(minutes) * 60 + parseFloat(seconds);
        }
        // 解析时+分+秒格式
        if (/^\d{1,2}[:：]\d{2}[:：]\d{2}(\.\d+)?$/.test(time)) {
            const [hours, minutes, seconds] = time.split(/[:：]/);
            return parseFloat(hours) * 3600 + parseFloat(minutes) * 60 + parseFloat(seconds);
        }
        return 0;
    }(time)
    return parseFloat(data)
}

function syncScoresInput() {
    // 将所有的输入框确认
    for (let i = 1; i < 6; i++) {
        let value = $("#input-score-number" + `${i}-value`)
        if (value.attr("disabled") === "disabled") {
            continue
        }
        if (validateTimeFormat(value.val())) {
            continue
        }
        value.val("")
    }


    // 确认成绩
    let submit = $("#submit-button")

    submit.attr("disabled", true)
    submit.addClass("btn-secondary")
    submit.removeClass("btn-success")
    for (let i = 1; i < 6; i++) {
        let value = $("#input-score-number" + `${i}-value`)
        if (value.attr("disabled") !== "disabled" && value.val() === "") {
            return
        }
    }
    submit.removeAttr("disabled")
    submit.removeClass("btn-secondary")
    submit.addClass("btn-success")
}

function syncScoresList() {
    if ($("#contest-select").val() === "") {
        return
    }
    if ($("#user-data-list-input").val() === "") {
        return
    }
    $.ajax({
        url: `./../api/score/player/${$("#user-data-list-input").val()}/contest/${$("#contest-select").val()}`, type: "GET", async: false, dataType: "json", timeout: 5000, // 设置超时时间为 5000 毫秒 (5 秒)
        contentType: "application/json; charset=UTF-8", success: function (response) {
            console.log(response)
            console.log(111)
            let group = $("#user-scores-list-group")
            group.empty()
            for (let i = 0; i < response["data"].length; i++) {
                let data = response["data"][i]
                let result = [data["R1"], data["R2"], data["R3"], data["R4"], data["R5"]]
                group.append(`<li class="list-group-item">${data['Project']} 成绩 ${result} </li>`)
            }
        }, error: function (data, status) {
            console.log(data, status)
        }
    })
}

function syncProject(select) {
    projectValue = select.options[select.selectedIndex].value
    if (projectValue === "") {
        return
    }

    // 更新 disabled
    switch (projectValue) {
        case "333fm":
        case "333bf":
        case "444bf":
        case "555bf":
        case "666":
        case "777":
        case "333mbf":
            console.log("只有三个的项目", projectValue);
            enableInputScoreNumber([1, 2, 3]);
            break
        case "o_cola":
        case "jhh":
            console.log("只有一轮的项目", projectValue);
            enableInputScoreNumber([1]);
            break
        default:
            enableInputScoreNumber([1, 2, 3, 4, 5]);
            break
    }
    syncScoresInput()
}


function submitScores() {
    const result = [];
    let input1 = $("#input-score-number1-value")
    let input2 = $("#input-score-number2-value")
    let input3 = $("#input-score-number3-value")
    let input4 = $("#input-score-number4-value")
    let input5 = $("#input-score-number5-value")

    const data = {
        "PlayerName": $("#user-data-list-input").val(),
        "ContestID": parseInt($("#contest-select").val()),
        "RouteNumber": 1,
        "ProjectName": $("#project-select").val(),
        "Results": [parseTimeToSeconds(input1.val()), parseTimeToSeconds(input2.val()), parseTimeToSeconds(input3.val()), parseTimeToSeconds(input4.val()), parseTimeToSeconds(input5.val()),],
    }
    console.log(data)

    $.ajax({
        url: "./../api/score", type: "POST", async: false, dataType: "json", timeout: 5000, // 设置超时时间为 5000 毫秒 (5 秒)
        data: JSON.stringify(data), contentType: "application/json; charset=UTF-8", success: function (response) {
            alert("成功记录")
            input1.val("")
            input2.val("")
            input3.val("")
            input4.val("")
            input5.val("")
            syncByTabScore()
            syncScoresList()
            syncScoresInput()
        }, error: function (data, status) {
            alert(`记录失败 ${data} ${status}`)
        }
    })
}


function syncAllData() {
    $.ajax({
        url: "./../api/contests", type: 'GET', async: false, timeout: 5000, // 设置超时时间为 5000 毫秒 (5 秒)
        success: function (response) {
            contestsList = response["Contests"]
        }, error: function (xhr, status, error) {
            if (status === 'timeout') {
                console.error('请求超时');
            } else {
                console.error(error);
            }
        }
    });

    $.ajax({
        url: "./../api/players", type: 'GET', async: false, timeout: 5000, // 设置超时时间为 5000 毫秒 (5 秒)
        success: function (response) {
            playerList = response["Data"]
        }, error: function (xhr, status, error) {
            if (status === 'timeout') {
                console.error('请求超时');
            } else {
                console.error(error);
            }
        }
    });
}

// syncByTabScore 选择添加记录时需要同步的数据
function syncByTabScore() {
    syncAllData()
    // 同步比赛
    if (contestsList != null) {
        const contestSelect = $('#contest-select')
        contestSelect.empty()
        for (let i = 0; i < contestsList.length; i++) {
            const contest = contestsList[i]
            if (contest["IsEnd"]) {
                continue
            }
            const option = `<option value="${contest["ID"]}"> ${contest["Name"]} </option>`
            contestSelect.append(option)
        }
    }

    // 同步用户信息
    if (playerList != null) {
        const playerSelect = $('#user-datalistOptions')
        playerSelect.empty()
        for (let i = 0; i < playerList.length; i++) {
            const player = playerList[i]
            const option = `<option value="${player["Name"]}" id="player_${player["ID"]}">${player["ID"]} - ${player["WcaId"]}</option>`
            playerSelect.append(option)
        }
    }
}

// syncByTabScore 选择添加比赛时需要同步数据
function syncByTabContest() {
    syncAllData()
    if (contestsList != null) {
        const contestTabList = $("#add-contest-tab-user-list")
        contestTabList.empty()
        for (let i = 0; i < contestsList.length; i++) {
            const contest = contestsList[i]
            // todo 这里加结束比赛的按钮
            contestTabList.append(`<li class='list-group-item'> ${contest["Name"]} </li>`)
        }
    }
}

function syncContestScore(select) {
    contestValue = select.options[select.selectedIndex].value
}

function syncByTabPlayer() {
    syncAllData()

    const playerTabList = $("#add-user-tab-user-list")
    playerTabList.empty()
    if (playerList != null) {
        for (let i = 0; i < playerList.length; i++) {
            const player = playerList[i]

            let wcaId = player["WcaId"]
            if (wcaId === "") {
                wcaId = "无WcaID"
            }
            playerTabList.append(`<li class='list-group-item'> ${player["Name"]}( ${wcaId} )</li>`)
        }
    }
}

function submitContest() {
    const name = $("#add-contest-input-name")
    const description = $("#add-contest-input-description")
    if (name.val() === "") {
        alert("比赛名不能为空")
        return
    }

    $.ajax({
        url: "./../api/contests", type: 'POST', data: {
            "Name": name.val(), "Description": description.val()
        }, async: false, timeout: 5000, // 设置超时时间为 5000 毫秒 (5 秒)
        success: function (response) {
            name.val("")
            description.val("")
            alert("修改或创建成功");
            syncByTabContest()
        }, error: function (xhr, status, error) {
            alert(`创建失败 ${error} \n${xhr} \n${status}`)
        }
    });
}

function submitPlayers() {
    const name = $("#add-user-input-name")
    if (name.val() === "") {
        alert("名字不能为空")
        return
    }

    const wcaId = $("#add-user-input-wca_id")

    $.ajax({
        url: "./../api/players", type: 'POST', data: {
            "Name": name.val(), "WcaID": wcaId.val()
        }, async: false, timeout: 5000, // 设置超时时间为 5000 毫秒 (5 秒)
        success: function (response) {
            name.val("")
            wcaId.val("")
            alert("修改或创建成功");
            syncByTabPlayer()
        }, error: function (xhr, status, error) {
            alert(`创建失败 ${error} \n${xhr} \n${status}`)
        }
    });
}

(function () {
    // 在启动时执行的代码
    syncByTabScore()
})();