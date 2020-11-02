#!/bin/bash

#启动脚本
#入参：可执行文件名

# PRG是被执行的脚本的名称，可以认为PRG=="deploy.sh"，也可能是某个符号链，指向该脚本。
PRG="$0"


# 处理了一下PRG
while [ -h "$PRG" ]; do
  ls=`ls -ld "$PRG"`
  link=`expr "$ls" : '.*-> \(.*\)$'`
  if expr "$link" : '.*/.*' > /dev/null; then
    PRG="$link"
  else
    PRG=`dirname "$PRG"`/"$link"
  fi
done


# PRGDIR是PRG的目录路径名称
PRGDIR=`dirname "$PRG"`
cd "$PRGDIR"

function start(){
        #mkdir -p log
        PIDS=(`ps -ef | grep zhibiaocalcsvr | grep -v grep|grep -v deploy.sh|grep -v salt | awk '{print $2}'`)
        for ((i=0;i<${#PIDS[@]};++i))
        do
                kill -9 ${PIDS[i]}
                echo ${PIDS[i]}
        done
        #./zhibiaocalcsvr
        #./zhibiaocalcsvr > /dev/null 2>&1 &
        nohup ./zhibiaocalcsvr > runinfo.txt 2>&1 &
}

#停止脚本
#程序停止脚本
function stop(){
        echo -n "start to kill zhibiaocalcsvr: "

        PIDS=(`ps -ef | grep zhibiaocalcsvr | grep -v grep| grep -v deploy.sh|grep -v salt| awk '{print $2}'`)
        for ((i=0;i<${#PIDS[@]};++i))
        do
                kill -9 ${PIDS[i]}
                echo ${PIDS[i]}
        done
}

#状态监测脚本
#入参：可执行程序名称
function status(){
        local EXEC_FILE='zhibiaocalcsvr'
        numproc=`ps -ef|grep $EXEC_FILE|grep -v grep| grep -v deploy.sh|grep -v salt| wc -l`
        if [ $numproc -gt 0 ]; then
                pid=`ps -ef|grep ${EXEC_FILE}|grep -v '${EXEC_FILE}'|awk '{print $2}'`
                echo "$EXEC_FILE Service is online...$pid"
                return 126
        else
                echo "$EXEC_FILE Service is stopped..."
                return 127
        fi
}

case "$1" in
start)
   start
   ;;
stop)
   stop
   ;;
status)
   status
   ;;
*)
   echo $"Usage: $0 {start|stop|status}"
   exit 1
esac