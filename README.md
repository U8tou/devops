# SHH_SERVER

#### 介绍
​	作为个人开发者，当项目、服务多了之后，部署起来就会感觉非常繁琐。市面上有 Jenkents、Jpome优秀的自动化部署框架，但需要一个更轻量点的方案。最好是直接双击直接运行，而不需准备运行环境。
​	所以就写了这么个小工具，仿照n8n工作流，能够实现脚本管理，配置好后可以一键按流程执行。它既是一个自动化部署工具，同时也是个Admin后台模板。
​	若对你有所帮助，希望不吝啬一个start。

#### 软件架构

后端：golang xorm

前端：vue3 element-plus (基于优秀前端模板https://www.artd.pro/#/dashboard/console)

#### 安装教程

1.  下载压缩包（或自己打包构建）
2.  直接运行
    - window双击 shh-windows-amd64.exe，linux 执行
    - linux 直接执行 shh-linux-amd64
3.  出现运行成功控制台，浏览器访问 http://localhost:3001
4.  自行配置配置文件 setting.yaml. 若首次启动，则需要注册账户，首个注册的账户默认为 超级管理员，直接具备系统所有权限

#### 使用说明

![](https://bsyimg.luoca.net/imgtc/20260427/ebe9465390487ce81308c7eb51fa1984.webp)

![](https://bsyimg.luoca.net/imgtc/20260427/7cb657381426907f1314fc3934864650.webp)
