#!/bin/sh

start="1"
count="2000"
managerAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
userinfoTempfile=".userinfo_temp"
addressInfo=".addr_info"
addrbalanceInfo=".addr_balance_info"

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
        `./chain33-cli account create -l $username > ${userinfoTempfile}`
        cat ${userinfoTempfile}
        useraddr=`cat ${userinfoTempfile} | grep "addr" | awk -F '"' '{print $4}'`
        echo "${username}:${useraddr}" >> ${addressInfo}
        rm ${userinfoTempfile}
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

function BuyKeys()
{
    buystart=$1
    for ((i=$buystart; i < ${count}; i++))
    do 
        username="f3d_user"$i
        GetAddressByLabel $username
        unsignedTx=`./chain33-cli f3d game buy -n 1`
        Sign ${unsignedTx} ${address}
        sleep 0.5
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
    signedTx=`./chain33-cli wallet sign -d ${unsignedTx} -a ${sendaddr} -e 0`
    Send ${signedTx} 
}

function Send() 
{
    signedtx=$1
    ./chain33-cli wallet send -d ${signedtx}
}

op=$1
if [ "X${op}" == "Xcreate" ]; then
    CheckAddressInfoFile
    CreateUser
elif [ "X${op}" == "Xtransfer" ]; then
    Transfer
elif [ "X${op}" == "Xbuy" ]; then
    buystart=$2
    BuyKeys $buystart
elif [[ "X${op}" == "Xcheck_balance" ]]; then
    #statements
    checkstart=$2
    CheckBalance $checkstart
else
    echo "Invalid operation."
fi