#!/bin/bash

case $1 in
"help")
    # 定义表格数据
    helpData=(
        "命令名称\t\t功能说明\t\t参数说明"
        "1. start\t\t开启程序\t\t如果第二个参数=false则视为不开启守护进程，默认为true"
        "2. stop\t\t\t关闭程序\t\t如果第二个参数=true则视为强制关闭，默认为false"
        "3. status\t\t查看程序运行状态"
        "4. restart\t\t重新启动\t\t第二个参数相当于start和stop的第二个参数"
        "5. build\t\t编译且打包程序\t\t第二个参数代表版本，同时也代表out目录下的文件夹名称"
        "\t\t\t\t\t\t第三个参数代表平台，默认为当前平台（非交叉编译）"
        "\t\t\t\t\t\t第四个参数代表架构，默认为amd64"
        "6. git-push\t\t推送一个分支到主分支\t第二个参数代表分支名称，默认为dev"
        "7. script \t\t执行命令\t\t第二个参数代表-t参数"
    )

    # 打印表格数据行
    for row in "${helpData[@]}"; do
        echo $row
    done
    ;;
"start")
    # 启动程序 如果第二个参数=false则视为不开启守护进程
    daemonParam=true
    if [ "$2" = false ]; then
        daemonParam=false
    fi
    chmod 755 ./jericho-go && ./jericho-go -daemon=$daemonParam
    ;;
"stop")
    # 关闭程序 如果第二个参数=true则视为强制关闭
    killParamCode=15
    killParamText="关闭成功"
    if [ "$2" = true ]; then
        killParamCode=9
        killParamText="关闭成功（强制关闭）"
    fi
    ps aux | grep jericho-go | grep -v grep | awk '{print $2}' | xargs kill -$killParamCode && echo $killParamText
    ;;
"status")
    # 查看程序运行状态
    ps aux | grep jericho-go | grep -v grep | xargs
    ;;
"restart")
    # 重新启动 第二个参数相当于start和stop的第二个参数
    source manager.sh stop $2 && source manager.sh start $2
    ;;
"build")
    # 编译且打包程序
    # 第二个参数代表版本，同时也代表out目录下的文件夹名称
    # 第三个参数代表平台，默认为当前平台（非交叉编译）
    # 第四个参数代表架构，默认为amd64
    if [ -z $3 ]; then
        goosName="当前平台"
    else
        goosName=$3
    fi

    if [ -z $4 ]; then
        goarchCode="amd64"
    fi

    echo "开始打包……" &&
        echo "初始化打包输出目录" &&
        rm -rf "./out/$2" &&
        mkdir -p "./out/$2" &&
        echo "编译『程序：jericho-go』『版本：$2』『平台：$goosName』『架构：$goarchCode』" &&
        CGO_ENABLED=0 GOOS=$3 GOARCH=$goarchCode go build -a -race -v -o "./out/$2/" "jericho-go" &&
        echo "编译成功" && 
        echo "打包：模板文件" &&
        cp -r ./templates "./out/$2/" &&
        echo "打包配置文件" &&
        rm -rf "./out/$2/settings" &&
        mkdir "./out/$2/settings" &&
        cp ./settings/app.ini.exa "./out/$2/settings/" &&
        cp ./settings/db.ini.exa "./out/$2/settings/" &&
        cp ./manager.sh "./out/$2/" &&
        echo "打包静态文件" &&
        cp -r ./static "./out/$2" &&
        echo "编译完成" &&
        open ./out/$2
    ;;
"git-push")
    # 第二个参数代表分支名称，默认为dev
    if [ -z $2 ]; then
        branch="dev"
    fi

    git add --all &&
        git commit -mm &&
        git push origin $branch &&
        git checkout master &&
        git pull origin master &&
        git merge $branch &&
        git push origin master &&
        git checkout $branch
    ;;
"script")
    # 第二个参数代表需要执行的t参数
    ./jericho-go -t=$2
    ;;
*)
    echo "没有这个命令，请执行help命令查看帮助"
    ;;
esac
