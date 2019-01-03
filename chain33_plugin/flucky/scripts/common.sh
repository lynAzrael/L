#!/bin/sh

rpc_addr="http://localhost:8801"

function GetKeyInfo()
{
    section=$1
    key=$2
    value=`sed -n "/^\[${section}/,/^\[/p" ${config_file} | grep ${key} | awk -F '=' '{print $2}' | tr -d '\r\n'`
}

function Curl()
{
    method=$1
    params=$2
    res=`curl --data-binary '{"jsonrpc":"2.0", "method": '"${method}"', "params": '"${params}"' , "id": 0}' -H 'content-type:text/plain;' ${rpc_addr} -s`
    if [[ "${res}" =~ "Err" ]]; then 
        return 1
    else
        return 0
    fi
}

function GetLocalTime()
{
    current=`date "+%Y-%m-%d %H:%M:%S"`
    timeStamp=`date -d "$current" +%s`

    temp=`date "+%N"`
    if [ ${temp:0:1} == 0 ]; then
        temp=${temp:1}
    fi
    #将current转换为时间戳，精确到毫秒  
    currentTimeStamp=$((timeStamp*1000+${temp}/1000000))
    echo $currentTimeStamp
}

function IsNum()
{
    input=$1

    if [ `echo $input | grep -q '[^0-9]'` ]; then
        return 0
    else
        return 1
    fi
}

GetLocalTime
