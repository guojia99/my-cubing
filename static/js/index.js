/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/4 下午3:05.
 *  * Author: guojia(https://github.com/guojia99)
 */

function formatTime(timestamp){
    timestamp = timestamp * 1000

    let date = new Date(timestamp);
    let year = date.getFullYear();
    let month = ("0" + (date.getMonth() + 1)).slice(-2); // 月份从0开始，需要+1
    let day = ("0" + date.getDate()).slice(-2);
    let hours = ("0" + date.getHours()).slice(-2);
    let minutes = ("0" + date.getMinutes()).slice(-2);
    let seconds = ("0" + date.getSeconds()).slice(-2);
    return year + "/" + month + "/" + day + " " + hours + ":" + minutes + ":" + seconds
}

function syncContest() {
    $.ajax({
        url: "./../api/contests",
        type: 'GET',
        async: true,
        timeout: 5000,
        success: function (response) {
            console.log(response, response["Contests"])
            const data = response["Contests"]
            for (let i = 0; i < data.length; i++) {

                let status = "text-primary"
                if (data[i]["IsEnd"]){
                    status = "text-secondary"
                }

                $("#contest_body").append(`
                  <tr>
                    <td>${i+1}</td>
                     <td>${data[i]["Name"]}</td>
                    <td><i class="bi bi-boxes ${status}"></i></td>
                    <td>${formatTime(data[i]["StartTime"])} - ${formatTime(data[i]["EndTime"])}</td>
                    <td><a href="/contest?contest_id=${data[i]['ID']}" class="btn btn-success">前往</a></td>
                </tr>
                `)
            }
        },
    });
}