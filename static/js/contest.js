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
                let project = response["ProjectList"][i];
                let tableBody = ""
                const data = response["Data"][project]
                for (let j = 0; j < data.length; j++){


                    let bestStyle = ""
                    let bestUpIcons = ""
                    if (data[j]["IsBest"]){
                        bestUpIcons =  ` <i class="bi bi-graph-up-arrow"></i>`
                        bestStyle = "color:#dc3545;font-weight:bold"
                    }

                    let avgStyle = ""
                    let avgUpIcons = ""
                    if (data[j]["IsBestAvg"]){
                        avgUpIcons = ` <i class="bi bi-graph-up-arrow"></i>`
                        avgStyle = "color:#dc3545;font-weight:bold"
                    }

                    let tr = `
                        <tr>
                            <td>${j+1}</td>
                            <td>${data[j]["Player"]}</td>
                            <td style="${bestStyle}">${data[j]["Best"]} ${bestUpIcons}</td>
                            <td style="${avgStyle}">${data[j]["Avg"]} ${avgUpIcons}</td>
                            <td>${data[j]["Result1"]}</td>
                            <td>${data[j]["Result2"]}</td>
                            <td>${data[j]["Result3"]}</td>
                            <td>${data[j]["Result4"]}</td>
                            <td>${data[j]["Result5"]}</td>
                        </tr>
                    `
                    tableBody += tr
                }

                let table = `
                        <div style="margin-top: 25px">
                            <div class="col-md-12">
                                <h3 class="text-center"><strong>${project}</strong></h3>
                                    <table class="table table-bordered table-striped" style="text-align:center">
                                        <thead>
                                            <tr>
                                                <th scope="col">排名</th>
                                                <th scope="col">选手</th>
                                                <th scope="col">单次</th>
                                                <th scope="col">平均</th>
                                                <th scope="col">还原1</th>
                                                <th scope="col">还原2</th>
                                                <th scope="col">还原3</th>
                                                <th scope="col">还原4</th>
                                                <th scope="col">还原5</th>
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