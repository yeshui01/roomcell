#!/bin/bash
ulimit -n 65535

#ALL=(zone dbserver rank social dblog game gateway)
SERVICE_ARRAY=(zone archive dbserver dblog rank social game gateway)
SERVICE_ARRAY_CLOSE=(gateway rank dblog social game dbserver archive zone)
num=`/bin/pwd |grep gameserver |grep bbqzs |awk -F'/' '{print $4}' |awk -F'_' '{print $2}'`
ROOTDIR=/data/gameserver/bbqzs_${num}
LOGPATH=/data/gameserver/bbqzs_${num}/logs

help() {
   printf "$num 控制参数输入错误
请输入执行命令参数start|stop|check 用来启动|停止|检查全进程
参数archive|dbserver|dblog|rank|socail|game|gateway start|stop 用来启动|停止对应子进程\n"
}
start_server() {
   process=$1
   server_proc=`/bin/ps aux | grep -w configPath | grep -w bbqzs_${num} | grep -w ${process} | awk '{print $2}'`
   time_start=$SECONDS
   if [ "${server_proc}" != "" ];then
      starttime=`/bin/ps aux | grep -w configPath | grep -w bbqzs_${num} | grep -w ${process} | awk '{print $9}'`
      printf "${process}进程已经存在，启动时间为：${starttime} \n"
      return
   else
      [ ! -d $LOGPATH/old/${process} ] && mkdir -p $LOGPATH/old/${process}
      [ -f $LOGPATH/${process}_console.log ] && mv ${LOGPATH}/${process}_console.log ${LOGPATH}/old/${process}/${process}_console.log.`date +%Y%m%d%H%M`
      nohup ${ROOTDIR}/server/bin/${process} --configPath ${ROOTDIR}/server/config/${process}/ > ${LOGPATH}/${process}_console.log 2>&1 &
      sleep 1
      PID=`/bin/ps aux | grep -w configPath | grep -w bbqzs_${num} | grep -w ${process} | awk '{print $2}'`
      port=`cat /data/gameserver/bbqzs_${num}/server/config/${process}/app.yaml | grep listenAddr |awk -F ":" '{print $3}' | awk -F "\"" '{print $1}'`   
      if [ "$PID" != "" ];then
         if [ ${process} == 'archive' ];then
            echo -e "[info] \t bbqzs_${num} ${process} start --------- [ OK ]\n"
         else
            connects=`/bin/netstat -ntlup|grep -w ${PID} |awk '{print $4}'|awk -F ":" '{print $2}'|grep -w ${port}`
            while [[ "$connects" == "" ]]
            do
               if [ ${process} == 'gateway' ];then
                  connects=`/bin/netstat -ntlup|grep -w ${PID} |awk '{print $4}'|awk -F ":" '{print $4}'|grep -w ${port}`
               else
                  connects=`/bin/netstat -ntlup|grep -w ${PID} |awk '{print $4}'|awk -F ":" '{print $2}'|grep -w ${port}`
               fi
               time_end=$SECONDS
               if [ `expr "$time_end" - "$time_start"` -ge 30 ];then
                  echo -e "[info] \t bbqzs_${num} ${process}  start --------- [ timeout ]\n"
                  exit 9999
               else
                  sleep 5
               fi
            done
#         [ ! -d ${ROOTDIR}/pid/ ] && mkdir -p ${ROOTDIR}/pid/
#         echo $PID >  ${ROOTDIR}/pid/${process}.pid
         echo -e "[info] \t bbqzs_${num} ${process} start --------- [ OK ]\n"
         fi
      else
         echo -e "[info] \t bbqzs_${num} ${process} start --------- [ Failed ]\n"
         exit 9999
      fi	
   fi
}

stop_server() {
   process=$1
   time_start=$SECONDS
   server_proc=`/bin/ps aux | grep -w configPath | grep -w bbqzs_${num} | grep -w ${process} | awk '{print $2}'`
   if [ "${server_proc}" == "" ];then
      echo "bbqzs_${num} ${process}没有启动"
   else
      cd ${ROOTDIR}/server
#      rm -f ${ROOTDIR}/pid/${process}.pid
      kill $server_proc
      sleep 1
      PID=`/bin/ps aux | grep -w configPath | grep -w bbqzs_${num} | grep -w ${process} | awk '{print $2}'`
      while [ "${PID}" != "" ]
      do
         time_end=$SECONDS
         if [ `expr "$time_end" - "$time_start"` -ge 10 ];then
            echo -e "[info] \t bbqzs_${num} ${process}  stop --------- [ FAIL ]\n"
            exit 9999
         else
            sleep 2
            PID=`/bin/ps aux | grep -w configPath | grep -w bbqzs_${num} | grep -w ${process} | awk '{print $2}'`
         fi
      done
      echo -e "[info] \t bbqzs_${num} ${process}  stop --------- [ OK ]\n"
   fi
}

check_server() {
   process=$1
   server_proc=`/bin/ps aux | grep -w configPath | grep -w bbqzs_${num} | grep -w ${process} | awk '{print $2}'`
   if [ "${server_proc}" == "" ]; then
      printf "[info] \t bbqzs_${num} ${process} no working\n"
   else
      printf "[info] \t bbqzs_${num} ${process} is working, pid is ${server_proc} \n"
   fi
}

#start_archive() {
#   process=$1
#   server_proc=`/bin/ps aux | grep -w configPath | grep -w bbqzs_${num} | grep -w ${process} | awk '{print $2}'`
#   time_start=$SECONDS
#   if [ "${server_proc}" != "" ];then
#      starttime=`/bin/ps aux | grep -w configPath | grep -w bbqzs_${num} | grep -w ${process} | awk '{print $9}'`
#      printf "${process}进程已经存在，启动时间为${starttime} \n"
#      return
#   else
#      [ ! -d $LOGPATH/old/${process} ] && mkdir -p $LOGPATH/old/${process}
#      mv ${LOGPATH}/${process}_console.log ${LOGPATH}/old/${process}/${process}_console.log.`date +%Y%m%d%H%M`
#      nohup ${ROOTDIR}/server/bin/${process} --configPath ${ROOTDIR}/server/config/${process}/ > ${LOGPATH}/${process}_console.log 2>&1 &
#      sleep 1
#      PID=`/bin/ps aux | grep -w configPath | grep -w bbqzs_${num} | grep -w ${process} | awk '{print $2}'`
#      if [ "$PID" != "" ];then
##         [ ! -d ${ROOTDIR}/pid/ ] && mkdir -p ${ROOTDIR}/pid/ 
##         echo $PID >  ${ROOTDIR}/pid/${process}.pid
##         echo $PID >  ${ROOTDIR}/${process}.pid
#         echo -e "[info] \t bbqzs_${num} ${process} start --------- [ OK ]\n"
#      else
#         echo -e "[info] \t bbqzs_${num} ${process} start --------- [ Failed ]\n"
#         exit 9999
#      fi
#   fi  
#}
if [ $# == 1 ];then
   case $1 in
   start)
     for server_name in "${SERVICE_ARRAY[@]}"
     do
        start_server ${server_name}
     done
   ;;
   stop)
     for server_name in "${SERVICE_ARRAY_CLOSE[@]}"
     do
        stop_server ${server_name}
     done
   ;;
   check)
     for server_name in "${SERVICE_ARRAY_CLOSE[@]}"
     do
        check_server ${server_name}
     done
   ;;
   *)
      help
      exit 999
   ;;
   esac
elif [ $# == 2 ];then
   for ctrl in ${SERVICE_ARRAY[@]}
   do
      if [ $1 == $ctrl ];then
         case $2 in
         start)
            start_server $ctrl
         ;;
         stop)
            stop_server $ctrl
         ;;
         *)
         printf "bbqzs_$num 控制参数2输入错误，你只能填写start|stop\n"
         exit 9999
         ;;
         esac
         exit 0
      else
         unset ctrl
      fi
   done
   help
   exit 9999
else
   help
   exit 9999
fi


