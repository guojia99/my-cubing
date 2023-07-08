/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/1 下午4:08.
 *  * Author: guojia(https://github.com/guojia99)
 */

function syncAllProjectBestScores() {
    $.ajax({
        url: "./../api/score/report/all_project_best",
        type: 'GET',
        async: true,
        timeout: 5000,
        success: function (response) {
            console.log(response)
            const body = $("#best_table_body")
            const data = response["Data"]
            for (let i = 0; i < data.length; i++) {
                const project = data[i]["Project"]

                if (project === "速可乐" || project === "菊爆浩浩") {
                    continue
                }

                if (data[i]["BestScore"] === 0) {
                    continue
                }
                if (project === "多盲") {
                    body.append(`<tr>
                        <th scope="row">${project}</th>
                         <td>${data[i]["BestPlayer"]}</td>
                       <td>${data[i]["MBFScore"]} (${formatTimeByProject(data[i]["BestScore"])})</td>
                        <td>-</td>
                        <td>-</td> 
                    </tr>`)
                    continue
                }


                let bestTd = `
                        <td>${data[i]["BestPlayer"]}</td>
                        <td>${formatTimeByProject(data[i]["BestScore"], project)}</td>
                `
                if (data[i]["BestPlayer"] === "-") {
                    bestTd = `<td>-</td><td>-</td>`
                }

                let avgTd = `
                        <td>${formatTimeByProject(data[i]["AvgScore"])}</td>
                        <td>${data[i]["AvgPlayer"]}</td> 
                `
                if (data[i]["AvgPlayer"] === "-") {
                    avgTd = `<td>-</td><td>-</td>`
                }

                body.append(`
                   <tr>
                        <th scope="row">${project}</th>
                        ${bestTd}
                        ${avgTd}
                    </tr>
                `)
            }
        },
    });
}


function syncAllProjectScores() {
    $.ajax({
        url: "./../api/score/report/all_project_score",
        type: 'GET',
        async: true,
        timeout: 5000,
        success: function (response) {
            const projectList = response["ProjectList"]
            const best = response["Best"]
            const avg = response["Avg"]
            let allProjectBody = $("#all_project_body")
            console.log(response)

            for (let i = 0; i < projectList.length; i++) {
                // 获取所有的成绩
                const project = projectList[i], projectAvg = avg[project], projectBest = best[project]
                let maxLength = projectBest.length
                if (maxLength === 0) {
                    continue
                }
                let tableBody = ""

                if (project === "多盲"){
                    for (let i = 0; i < maxLength; i++) {
                        // 加入
                        let tr = `
                            <tr>
                                <td>${i}</td>
                                <td>${projectBest[i]["Player"]}</td>
                                <td>${projectBest[i]["R1"]} / ${projectBest[i]["R2"]}</td>
                                <td>${formatTimeByProject(projectBest[i]["R3"])}</td>
                            </tr>`
                        tableBody += tr
                    }
                    let table = `
                        <div class="col-md-6" style="margin-top: 30px">
                                <h3 class="text-center" style="margin-bottom: 15px"><strong>${project}排名</strong></h3>
                                <table class="table table-bordered table-striped" style="text-align:center">
                                    <thead>
                                    <tr>
                                        <th scope="col">排名</th>
                                        <th scope="col">选手</th>
                                        <th scope="col">还原</th>
                                        <th scope="col">时长</th>
                                    </tr>
                                    </thead>
                                    <tbody>${tableBody}</tbody>
                                </table>
                        </div>`
                    allProjectBody.append(table)
                    continue
                }

                let bestRoute = 0, avgRoute = 0
                let lastBestScore = 0, lastAvgScore = 0
                for (let i = 0; i < maxLength; i++) {
                    // 如果和上次成绩不同
                    if (projectBest[i]["Best"] !== lastBestScore) {
                        bestRoute = i + 1
                    }
                    lastBestScore = projectBest[i]["Best"]

                    // 这里因为只有平均才有可能小于最佳
                    let avgTd = `<td>-</td><td>-</td><td>-</td>`
                    if (i < projectAvg.length) {
                        if (projectAvg[i]["Avg"] !== lastAvgScore) {
                            avgRoute = i + 1
                        }
                        avgTd = `<td>${formatTimeByProject(projectAvg[i]["Avg"])}</td><td>${projectAvg[i]["Player"]}</td><td>${avgRoute}</td>`
                        lastAvgScore = projectAvg[i]["Avg"]
                    }

                    // 加入
                    let tr = `
                            <tr>
                                <td>${bestRoute}</td><td>${projectBest[i]["Player"]}</td><td>${formatTimeByProject(projectBest[i]["Best"], project)}</td>
                                ${avgTd}
                            </tr>`
                    tableBody += tr
                }


                let table = `
                <div class="col-md-6" style="margin-top: 30px">
                        <h3 class="text-center" style="margin-bottom: 15px"><strong>${project}排名</strong></h3>
                        <table class="table table-bordered table-striped" style="text-align:center">
                            <thead>
                            <tr>
                                <th scope="col">排名</th>
                                <th scope="col" colspan="2">单次</th>
                                <th scope="col" colspan="2">平均</th>
                                <th scope="col">排名</th>
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
                    <td>${i + 1}</td>
                    <td>${best[i]["Player"]}</td>
                    <td>${best[i]["Count"]}</td>
                    <td>${avg[i]["Count"]}</td>
                    <td>${avg[i]["Player"]}</td>
                    <td>${i + 1}</td>
                </tr>
                `)
            }
        },
    });
}