#!/bin/bash

CSVGITPATH=/home/gitlab-runner/p1/gameconf
SRCPATH=/home/gitlab-runner/p1/myth
RUNDIR=/home/gitlab-runner/cehua/dev_linux

# 使用说明
function show_use() {
	echo -e "Usage: $0 (updcsv restart_server upd_restart)"
	echo -e "\t-- updcsv: 更新csv"
	echo -e "\t-- restart_server: 重启策划服服务器"
	echo -e "\t-- upd_restart: 更新配置并且重启服务器"
}

# 更新csv
function upd_csv() {
	echo "csv git path: $CSVGITPATH " 
	cd $CSVGITPATH
	git pull
	cd $RUNDIR
	echo "run dir: $RUNDIR"
	./update_csv.sh
	echo "build server exe on linux"
	cd $SRCPATH
	echo "src dir: $SRCDIR"
	git checkout sync1.5
	git config pull.ff only
	git pull
	#git fetch --all
	#git reset --hard origin/sync1.5
	go mod tidy
	make build
}

# 重启策划服
function restart_cehua_server() {
	echo "run dir: $RUNDIR"
	cd $RUNDIR
	./update_bin.sh
	sleep 1
	./notify_robot.sh stop
	echo "notify stopped\n"
	sleep 1
	./ctl_server_linux.sh stop
	sleep 1
	./ctl_server_linux.sh start
	sleep 1
	./notify_robot.sh finish
	sleep 1
	echo "notify finish\n"
}

if [ $# -eq 0 ]; then
	echo -e "\e[;31mparam error!!!!\e[0m"
	show_use
	exit
fi

if [ $1 == 'updcsv' ]; then
	echo -e "\e[;32m---------- upd csv --------------------\e[0m"	
	upd_csv
elif [ $1 == 'restart_server' ]; then 
	echo -e "\e[;32m---------- restart cehua server --------------------\e[0m"	
	restart_cehua_server
elif [ $1 == 'upd_restart' ]; then 
	echo -e "\e[;32m---------- upd csv and restart cehua server --------------------\e[0m"	
	upd_csv
	sleep 2
	restart_cehua_server
else
	echo -e "\e[;31m $1 is not supportted param!!!!\e[0m"	
fi

