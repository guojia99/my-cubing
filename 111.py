for i in range(100):
    print(f'''INSERT INTO mycube2.contests (id, created_at, name, description, is_end, round_ids, start_time, end_time) VALUES ({i + 9}, '2023-07-26 11:23:14.000', '测试比赛{i}', '测试比赛', 0, '[1]', '2023-07-26 11:23:29.000', '2023-07-26 11:23:31.000');''')
