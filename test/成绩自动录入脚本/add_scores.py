#   * Copyright (c) 2023 guojia99 All rights reserved.
#   * Created: 2023/7/4 下午6:29.
#   * Author: guojia(https://github.com/guojia99)
import csv

import requests


def score_parse(r1, r2, r3, r4, r5):
    out = []
    for r in [r1, r2, r3, r4, r5]:
        if r == "DNF" or r == "DNS" or r == "/":
            r = 0
        if ":" in f"{r}":
            minutes, seconds = r.split(":")
            minutes = int(minutes)
            seconds = float(seconds)
            total_seconds = minutes * 60 + seconds
            out.append(float(total_seconds))
            continue
        out.append(float(r))
    return out


def parse_csv(file_path):
    data = []
    with open(file_path, 'r', encoding='utf-8') as csv_file:
        csv_reader = csv.DictReader(csv_file)
        for row in csv_reader:
            data.append({
                'project': row["项目"],
                'name': row['名称'],
                'result': score_parse(row['#1'], row['#2'], row['#3'], row['#4'], row['#5']),
            })
    return data


def run(d: dict):
    data = parse_csv(d["csv_file_path"])

    resp = requests.request("POST", "http://127.0.0.1:14023/api/contests", json={
        "Name": d["contest_name"],
        "Description": d["contest_name"],
    })

    contest_id = resp.json()["id"]
    for s in data:
        val = {
            "PlayerName": s["name"],
            "ContestID": contest_id,
            "RouteNumber": 1,
            "ProjectName": s["project"],
            "Results": s["result"]
        }
        resp = requests.request("POST", "http://127.0.0.1:14023/api/score", json=val)
        if resp.status_code != 200:
            print(resp.json())

    if d.get("end"):
        requests.request("POST", f"http://127.0.0.1:14023/api/score/report/contest/{contest_id}/end")


if __name__ == "__main__":
    data = [
        "蛋糕",
        "Cuber浩",
        "羽妹妹",
        "Cuber-MN",
        "MIB T",
        "郭嘉",
        "兔兔",
        "Clansey",
        "石头",
        "老司机",
        "串串香",
        "Showball",
        "clkcj",
        "跳跳",
        "火花",
        "云翮",
        "北宅",
        "Justin",
        "小情子",
        "小马哥",
        "诚诚",
        "北凉",
        "孤烟往事",
    ]

    for i in data:
        requests.request("POST", url="http://127.0.0.1:14023/api/players", json={
            "Name": i
        })

    contests = [
        {
            "csv_file_path": './scores/魔缘20230525第一期群赛赛果.csv',
            "contest_name": "魔缘2023第一期群赛",
            "end": True,
        },
        {
            "csv_file_path": './scores/魔缘20230605第二期群赛赛果.csv',
            "contest_name": "魔缘2023第二期群赛",
            "end": True,
        },
        {
            "csv_file_path": './scores/魔缘20230612第三期群赛赛果.csv',
            "contest_name": "魔缘2023第三期群赛",
            "end": True,
        },
        {
            "csv_file_path": './scores/魔缘20230618第四期群赛赛果.csv',
            "contest_name": "魔缘2023第四期群赛",
            "end": True,
        },
        {
            "csv_file_path": './scores/魔缘20230626第五期群赛赛果.csv',
            "contest_name": "魔缘2023第五期群赛",
            "end": True,
        },
        {
            "csv_file_path": './scores/魔缘20230702第六期群赛赛果.csv',
            "contest_name": "魔缘2023第六期群赛",
            "end": False,
        },
    ]

    for i in contests:
        run(i)
