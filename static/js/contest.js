
function getAllQueryParams() {
    const queryParams = {};
    const queryString = window.location.search.substring(1);
    const vars = queryString.split("&");
    for (let i = 0; i < vars.length; i++) {
        const pair = vars[i].split("=");
        const key = decodeURIComponent(pair[0]);
        const value = decodeURIComponent(pair[1]);
        // 如果该参数已存在，则将其转换为数组形式
        if (queryParams[key]) {
            if (Array.isArray(queryParams[key])) {
                queryParams[key].push(value);
            } else {
                queryParams[key] = [queryParams[key], value];
            }
        } else {
            queryParams[key] = value;
        }
    }
    return queryParams;
}


function syncContest() {
    let params = getAllQueryParams()
    if (params["contest_id"] === undefined) {
        return
    }

    $.ajax({
        url: `./../api/score/report/contest/${params["contest_id"]}`,
        type: 'GET',
        async: true,
        timeout: 5000,
        success: function (response) {
            const projectBody = $("#all_project_body")
            document.getElementById("contest_name").innerHTML = `${response["ContestName"]} 赛果`


            for (let i = 0; i < response["ProjectList"].length; i++) {
                let thBy23Project = `<th scope="col">还原2</th><th scope="col">还原3</th>`
                let thBy45Project = ` <th scope="col">还原4</th><th scope="col">还原5</th>`


                let project = response["ProjectList"][i], tableBody = ""
                const data = response["Data"][project]
                for (let j = 0; j < data.length; j++) {
                    let bestStyle = "", bestUpIcons = ""
                    if (data[j]["IsBest"]) {
                        bestUpIcons = `<i class="bi bi-graph-up-arrow"></i>`
                        bestStyle = "color:#dc3545;font-weight:bold"
                    }

                    let avgStyle = ""
                    let avgUpIcons = ""
                    if (data[j]["IsBestAvg"]) {
                        avgUpIcons = ` <i class="bi bi-graph-up-arrow"></i>`
                        avgStyle = "color:#dc3545;font-weight:bold"
                    }


                    // 这里为了区分多个三个的项目等
                    let trBy23Project = `
                            <td>${formatTimeByProject(data[j]["R2"], project)}</td>
                            <td>${formatTimeByProject(data[j]["R3"], project)}</td>
                    `

                    let trBy45Project = `
                            <td>${formatTimeByProject(data[j]["R4"], project)}</td>
                            <td>${formatTimeByProject(data[j]["R5"], project)}</td>
                    `


                    switch (project) {
                        case "最少步":
                        case "三盲":
                        case "四盲":
                        case "五盲":
                        case "六阶":
                        case "七阶":
                        case "多盲":
                            trBy45Project = ""
                            thBy45Project = ""
                            break
                        case "菊爆浩浩":
                        case "速可乐":
                            trBy23Project = ""
                            trBy45Project = ""
                            thBy23Project = ""
                            thBy45Project = ""
                            break
                        default:
                            break
                    }

                    let tr = `
                        <tr>
                            <td>${j + 1}</td>
                            <td>${data[j]["Player"]}</td>
                            <td style="${bestStyle}">${formatTimeByProject(data[j]["Best"],project)} ${bestUpIcons}</td>
                            <td style="${avgStyle}">${formatTimeByProject(data[j]["Avg"])} ${avgUpIcons}</td>
                            <td>${formatTimeByProject(data[j]["R1"],project)}</td>
                               ${trBy23Project}
                               ${trBy45Project}
                        </tr>
                    `
                    tableBody += tr
                }

                let table = `
                        <div style="margin-top: 25px">
                            <div class="col-md-12">
                                <h3><strong>${project}</strong></h3>
                                    <table class="table table-bordered table-striped" style="text-align:center">
                                        <thead>
                                            <tr>
                                                <th scope="col">排名</th>
                                                <th scope="col">选手</th>
                                                <th scope="col">单次</th>
                                                <th scope="col">平均</th>
                                                <th scope="col">还原1</th>
                                                ${thBy23Project}
                                                ${thBy45Project}
                                            </tr>
                                        </thead>
                                    <tbody>${tableBody}</tbody>
                                </table>
                            </div>
                        </div>
`
                projectBody.append(table)
            }
        },
    });
}