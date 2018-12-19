#!/bin/sh

start="200"
count="205"
managerAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
userinfoTempfile=".userinfo_temp"
addressInfo=".addr_info"
addrbalanceInfo=".addr_balance_info"
rpc_addr="http://localhost:8801"
config_file="exec_config"
leftCh="["
rightCh="]"

function ParseConfigFile()
{
    if [ ! -f ${config_file} ]; then
        echo "Can't find config file."
        exit 1
    fi


}

function CheckAddressInfoFile()
{
    if [ -f ${addressInfo} ]; then
        rm ${addressInfo}
        touch ${addressInfo}
    fi

    if [ -f ${addrbalanceInfo} ]; then
        rm ${addrbalanceInfo}
        touch ${addrbalanceInfo}
    fi
}

function CreateUser() 
{
    for ((i=$start; i < ${count}; i++))
    do 
        username="f3d_user"$i
        methodname="Chain33.NewAccount"
        paramsInfo="[{\"label\": \"${username}\"}]"
        res=`curl --data-binary "{\"jsonrpc\":\"2.0\", \"method\": \"${methodname}\", \"params\": ${paramsInfo} , \"id\":0}" -H 'content-type:text/plain;' ${rpc_addr} ` 
        if [[ "${res}" =~ "ErrLabelHasUsed" ]]; then
            continue
        else
            echo "${res}" > ${userinfoTempfile}
            cat ${userinfoTempfile}
            useraddr=`cat ${userinfoTempfile} | grep "addr" | awk -F '"' '{print $10}'`
            echo "${username}:${useraddr}" >> ${addressInfo}
            rm ${userinfoTempfile}
        fi
    done
}

function Transfer()
{
    for ((i=$start; i < ${count}; i++))
    do 
        username="f3d_user"$i
        GetAddressByLabel $username
        if [ "X${address}" == "X" ]; then
            continue
        fi
        TransferToUser ${address} "100000"
        TransferToExecFromAccount "90000" "f3d"
    done
}

function Buy()
{
    buystart=$1
    GetKeyInfo "Buy" "method"
    buyKeysmethodname="${value}"
    GetKeyInfo "Buy" "param"
    RefreshParamInt64 "${value}" "num" "1"
    buyKeysParamsInfo="[${param}]"
    for ((i=$buystart; i < ${count}; i++))
    do 
        username="f3d_user"$i
        GetAddressByLabel $username
        res=`curl --data-binary '{"jsonrpc":"2.0", "method": '"${buyKeysmethodname}"', "params": '"${buyKeysParamsInfo}"' , "id": 0}' -H 'content-type:text/plain;' ${rpc_addr}`
        if [[ "${res}" =~ "Err" ]]; then
            continue
        else
            unsignedTx=`echo ${res} | awk -F '"' '{print $6}'`
            Sign ${unsignedTx} ${address}
            sleep 0.5
        fi
    done
}

function GetAddressByLabel()
{
    label=$1
    address=`cat ${addressInfo} | grep "${label}:" | awk -F ":" '{print$2}'`
}

function CheckBalance() 
{
    checkstart=$1
    for ((i=$checkstart; i < ${count}; i++))
    do 
        username="f3d_user"$i
        GetAddressByLabel $username
        ./chain33-cli account balance -a ${address} > ${addrbalanceInfo}
        balanceinfotime=`cat ${addrbalanceInfo} | grep balance | wc -l`
        # 账户中余额，且合约地址中也有余额
        if [ ${balanceinfotime} == 2 ]; then
            rm ${addrbalanceInfo}
            continue
        # 账户中有余额，但是合约地址中没有
        elif [ ${balanceinfotime} == 1 ]; then
            TransferToExecFromAccount "1000" "f3d"
            # sleep 0.5
        # 账户中没有钱，同样合约地址中也没有
        elif [ ${balanceinfotime} == 0 ]; then
            TransferToUser ${address} "2000"
            # sleep 0.5
            TransferToExecFromAccount "1000" "f3d"
            # sleep 0.5
        else
            echo "wrong balance info"
        fi

        balance=`cat ${addrbalanceInfo} | grep balance | awk -F '"' '{print $4}'`
        rm ${addrbalanceInfo}
        # sleep 0.5
    done
}

function TransferToUser()
{
    rcvrAddr=$1
    amount=$2
    res=`./chain33-cli send bty transfer -a ${amount} -k ${managerAddr} -t ${rcvrAddr}`
    while [ "$res" == "ErrTxExpire" ]; 
    do
        sleep 10
        res=`./chain33-cli send bty transfer -a ${amount} -k ${managerAddr} -t ${rcvrAddr}`
    done
}
    

function TransferToExecFromAccount()
{
    amount=$1
    executor_name=$2
    unsignedTx=`./chain33-cli bty send_exec -a ${amount} -e ${executor_name}`
    Sign ${unsignedTx} ${address}
}

function Exec() 
{
    cmd=$*
    res=`$cmd`
    Sign ${res} ${addr}
}

function Sign()
{
    unsignedTx=$1
    sendaddr=$2
    GetKeyInfo "Sign" "method"
    signMethodName="${value}"
    GetKeyInfo "Sign" "param"
    RefreshParamString "${value}" "addr" "${sendaddr}"
    RefreshParamString "${param}" "txHex" "${unsignedTx}"

    res=`curl --data-binary '{"method": '"${signMethodName}"', "params": '"[${param}]"', "id": 0}' -H 'content-type:text/plain;' ${rpc_addr}`
    if [[ "${res}" =~ "Err" ]]; then 
        echo "${res}"
    else
        signedTx=`echo ${res} | awk -F '"' '{print $6}'`
        Send ${signedTx}
    fi
   
}

function Send() 
{
    signedtx=$1
    GetKeyInfo "Send" "method"
    sendMethodname="${value}"
    GetKeyInfo "Send" "param"
    RefreshParamString "${value}" "data" "${signedtx}"

    res=`curl --data-binary '{"method": '"${sendMethodname}"', "params": '"[${param}]"' , "id": 0}' -H 'content-type:text/plain;' ${rpc_addr} `
    if [[ "${res}" =~ "Err" ]]; then
        errorInfo=`echo ${res} | awk -F '"' '{print $8}'`
        echo "${errorInfo}"
    fi
}

function GetKeyInfo()
{
    section=$1
    key=$2
    value=`sed -n "/^\[${section}/,/^\[/p" ${config_file} | awk 'NR>1 {print p} {p=$0}' | grep ${key} | awk -F '=' '{print $2}' | tr -d '\r'`
}

function RefreshParamString() 
{
    oldParam=$1
    refreshKey=$2
    refreshVal=$3

    # 使用jq后续curl指令执行有问题
    # param=`echo ${oldParam} | jq 'to_entries | map(if .key == "'${refreshKey}'" then . + {"value": "'${refreshVal}'"} else . end) | from_entries'`
    param=`echo "${oldParam}" | awk ' { for (i=1;i<=NF;i++) {if (match($i, "'${refreshKey}'")) {gsub(/inputParam/, "'${refreshVal}'", $(i+1)); print}} }'`
    echo $param
}

function RefreshParamInt64()
{
    oldParam=$1
    refreshKey=$2
    refreshVal=$3

    param=`echo "${oldParam}" | awk ' { for (i=1;i<=NF;i++) {if (match($i, "'${refreshKey}'")) {gsub(/\"inputParam\"/, '${refreshVal}', $(i+1)); print}} }'`
}

function ReadOperation()
{
    GetKeyInfo "Op" "support"
    ops=`echo ${value} | awk 'BEGIN{FS="[,\"]"} {for (i=1;i<NF;i++) {if ($i != "") print $i}}'`
    for op in ${ops}; do
        #statements
        `${op}`
    done
}

# op=$1
# if [ "X${op}" == "Xcreate" ]; then
#     CheckAddressInfoFile
#     CreateUser
# elif [ "X${op}" == "Xtransfer" ]; then
#     Transfer
# elif [ "X${op}" == "Xbuy" ]; then
#     buystart=$2
#     BuyKeys $buystart
# elif [[ "X${op}" == "Xcheck_balance" ]]; then
#     #statements
#     checkstart=$2
#     CheckBalance $checkstart
# else
#     echo "Invalid operation."
# fi