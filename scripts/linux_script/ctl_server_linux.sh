#!/bin/bash
STATUS_KEY=0
WORK_DIR=$(pwd)
#SERVICE_ARRAY=(dbserver dblog gateway rank social game archive account gmtool front report)
SERVICE_ARRAY=(battle2 syncbattle gmtool account zone archive dbserver dblog rank social game gateway)
SERVICE_ARRAY_CLOSE=(gateway game social rank dblog dbserver archive zone account gmtool battle2 syncbattle)

function get_pid() {
  pid=$(pgrep -f "$WORK_DIR/bin/$SERVICE_NAME --configPath=$WORK_DIR/app_config/$SERVICE_NAME")
}

function start_service() {
  get_pid
  if [ -z "$pid" ]; then
    nohup "$WORK_DIR"/bin/"$SERVICE_NAME" --configPath="$WORK_DIR/app_config/$SERVICE_NAME" >> "$WORK_DIR"/logs/"$SERVICE_NAME".log 2>&1 &
    sleep 3
    get_pid
    if [ -z "$pid" ]; then
      echo -e "$SERVICE_NAME startup \e[;31mfailed\e[0m"
        (( STATUS_KEY=1 ))
    else
      echo -e "$SERVICE_NAME startup \e[;32mseccuss\e[0m -- pid $pid"
    fi
  else
    echo -e "$SERVICE_NAME \e[;32mAlready running\e[0m"
  fi
}

function check_service() {
  get_pid
  if [ -z "$pid" ]; then
    echo -e "$SERVICE_NAME \e[;31mnot running\e[0m"
    start_service
  else
    echo -e "$SERVICE_NAME \e[;32mrunning\e[0m"
  fi
}

function kill_service() {
  echo "$SERVICE_NAME stoping"
  get_pid
  if [ -n "$pid" ]; then
      kill $pid
      sleep 1
      while true; do
        # 等待完成
        get_pid
        if [ -z "$pid" ]; then
          echo -e "$SERVICE_NAME \e[;31mis shutdown\e[0m"
          break
        else
          echo "$SERVICE_NAME stoping"
          sleep 1
        fi
      done
  else
      echo "$SERVICE_NAME is not running"
  fi
}

function service_status() {
  get_pid
  if [ -n "$pid" ]; then
    echo -e "$SERVICE_NAME \e[;32mrunning\e[0m -- pid $pid"
  else
    echo -e "$SERVICE_NAME \e[;31mnot running\e[0m"
  fi
}


while true; do
  case "$1" in
    start)
      # 建立log目录
      mkdir -p "$WORK_DIR"/logs
      if [ $# -eq 1 ]; then
        for e in "${SERVICE_ARRAY[@]}"; do
          SERVICE_NAME="$e"
          start_service
        done
      else
        for e in "$@"; do
          if [ "$e" == "start" ]; then continue; fi
          SERVICE_NAME="$e"
          start_service
        done
      fi
      if [ $STATUS_KEY -eq 1 ]; then
        exit 1
      fi
      shift
      break
      ;;
    stop)
      if [ $# -eq 1 ]; then
        for e in "${SERVICE_ARRAY_CLOSE[@]}"; do
          SERVICE_NAME="$e"
          kill_service
				  sleep 2
        done
        break
      fi
      for e in "$@"; do
        if [ "$e" == "stop" ]; then continue; fi
        SERVICE_NAME=$e
        kill_service
				sleep 2
      done
      shift
      break
      ;;
    status)
      if [ $# -eq 1 ]; then
        for e in "${SERVICE_ARRAY[@]}"; do
          SERVICE_NAME="$e"
          service_status
        done
        break
      fi
      for e in "$@"; do
        if [ "$e" == "status" ]; then continue; fi
        SERVICE_NAME=$e
        service_status
      done
      shift
      break
      ;;
    check)
      if [ $# -eq 1 ]; then
        for e in "${SERVICE_ARRAY[@]}"; do
          SERVICE_NAME="$e"
          check_service
        done
        break
      fi
      for e in "$@"; do
        if [ "$e" == "check" ]; then continue; fi
        SERVICE_NAME=$e
        check_service
      done
      shift
      break
      ;;
    help)
      echo -e "Usage: $0 [options] [service name ...]\n"
      echo -e "Options:\n start: \t启动服务(默认为启动所有服务)\n stop:  \t关闭服务(默认为关闭所有服务)\n check: \t检查状态(默认为检查所有服务)\n status:\t服务状态(默认为查看所有服务)\n\n help：    \t帮助"
      exit 0
      ;;
    *)
      echo "error options, pls check or use [$0 help]"
      exit 1
      ;;
  esac
done
