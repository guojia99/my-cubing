/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/1 下午4:08.
 *  * Author: guojia(https://github.com/guojia99)
 */

function formatTime(seconds) {
    if (seconds === 0) {
        return "-";
    }

    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const remainingSeconds = (seconds % 60).toFixed(2);

    let formattedTime = "";
    if (hours > 0) {
        formattedTime += hours + ":";
    }
    if (minutes > 0) {
        formattedTime += minutes;
    }
    if (hours === 0 && minutes === 0) {
        formattedTime += remainingSeconds;
    } else if (remainingSeconds !== "0.00") {
        formattedTime += ":" + remainingSeconds;
    }

    return formattedTime;
}

function syncAllProjectBestScores() {
    $.ajax({
        url: "./../api/score/report/all_project_best",
        type: 'GET',
        async: true,
        timeout: 5000,
        success: function (response) {
            const body = $("#best_table_body")
            const data = response["Data"]
            for (let i = 1; i < data.length; i++) {
                body.append(`
                   <tr>
                        <th scope="row">${data[i]["Project"]}</th>
                        <td>${data[i]["BestPlayer"]}</td>
                        <td>${formatTime(data[i]["BestScore"])}</td>
                        <td>${formatTime(data[i]["AvgScore"])}</td>
                        <td>${data[i]["AvgPlayer"]}</td>
                    </tr>
                `)
            }
        },
    });
}


function syncAllProjectScores() {
    $.ajax({
        url: "./../api/score/report/all_project",
        type: 'GET',
        async: true,
        timeout: 5000,
        success: function (response) {
            const projectList = response["ProjectList"]
            const best = response["Best"]
            const avg = response["Avg"]
            let allProjectBody = $("#all_project_body")

            for (let i = 0; i < projectList.length; i++) {
                const project = projectList[i]

                const projectAvg = avg[project]
                const projectBest = best[project]
                let maxLength = projectAvg.length
                if (projectBest.length > maxLength) {
                    maxLength = projectAvg.length
                }
                if (maxLength === 0) {
                    continue
                }


                let tableBody = ""
                for (let i = 0; i < maxLength; i++) {

                    let avgPlayer = "-"
                    let avgScore = 0.0
                    let bestPlayer = "-"
                    let bestScore = 0.0

                    if (projectBest.length >= maxLength) {
                        bestPlayer = projectBest[i]["Player"]
                        bestScore = projectBest[i]["Score"]
                    }
                    if (projectAvg.length >= maxLength) {
                        avgPlayer = projectAvg[i]["Player"]
                        avgScore = projectAvg[i]["Score"]
                    }


                    let tr = `
                            <tr>
                                <td>${i}</td>
                                <td>${bestPlayer}</td>
                                <td>${formatTime(bestScore)}</td>
                                <td>${formatTime(avgScore)}</td>
                                <td>${avgPlayer}</td>
                            </tr>`
                    tableBody += tr
                }

                let table = `
                <div class="col-md-6">
                        <h3 class="text-center"><strong>${project}排名</strong></h3>
                        <table class="table table-bordered table-striped" style="text-align:center">
                            <thead>
                            <tr>
                                <th scope="col">排名</th>
                                <th scope="col" colspan="2">单次</th>
                                <th scope="col" colspan="2">平均</th>
                            </tr>
                            </thead>
                            <tbody>${tableBody}</tbody>
                        </table>
                </div>`
                allProjectBody.append(table)
            }
        },
    });
}

function syncSorScores() {
    $.ajax({
        url: "./../api/score/report/all_sor",
        type: 'GET',
        async: true,
        timeout: 5000,
        success: function (response) {
            let best = response["Best"]
            let avg = response["Avg"]
            let tableBody = $("#sor_table_body")

            for (let i = 0; i < best.length; i++) {
                tableBody.append(`
                <tr>
                    <td>${i}</td>
                    <td>${best[i]["Player"]}</td>
                    <td>${best[i]["Count"]}</td>
                    <td>${avg[i]["Count"]}</td>
                    <td>${avg[i]["Player"]}</td>
                    <td>${i}</td>
                </tr>
                `)
            }
        },
    });
}