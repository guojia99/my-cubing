# my-cubing
魔方成绩管理系统

> - 网站功能：
>   - 成绩录入，成绩展示，成绩排名, 个人个性主页，成绩汇总，成绩可视化
> - 机器人功能：
>   - 成绩查询， 成绩排名，比赛汇总信息



## Core
- 项目地址：https://github.com/guojia99/my-cubing-core
- core项目用于实现基础模型的定义， 方法查询，成绩录入等核心功能， 提供对外统一的API接口，是直接和数据库做交互的项目。
- 该项目被 api、robot所依赖，需要保持此两个项目统一



## API

- 项目地址：https://github.com/guojia99/my-cubing-api
- api项目基于go-gin实现接口功能。
- api项目给予core 提供统一的对外的后台api接口， 并提供拓展其他与core成绩无关的接口实现。
- api项目还实现了权限控制， 对于查询外的操作提供统一的权限管理。



## UI

- 项目地址： https://github.com/guojia99/my-cubing-ui

- UI是前端项目， 基于react + bootstarp实现前端界面，该项目依赖后端所提供的api接口，实现统一的对外界面实现



## Robot

- 项目地址： https://github.com/guojia99/my-cubing-robot
- 机器人项目基于qq + cqhttp实现，通过core项目的接口实现对应的群机器人查询功能。



## Getaway

- 项目地址： https://github.com/guojia99/my-cubing-getaway
- 网关项目是为了对外统一的接口、限流、集合API以及作为UI项目的启动器， 通过内网转发实现前后端分离， 统一端口输出等。



