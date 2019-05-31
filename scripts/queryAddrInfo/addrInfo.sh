#!/bin/bash

MAIN_HTTP="http://localhost:8801"
execName="coins"

function getAllAddrInfo() {
    local res=""
    res=$(curl -k -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"Chain33.GetAccounts","params":[]}' -H 'content-type:text/plain;' ${MAIN_HTTP})
    result=$(echo "$res" | jq -r ".error|not")
    if [ "${result}" == "true" ]; then
        #echo "getAllAddrInfo successfully."
        local count=0
        count=$(echo "$res" | jq -r ".result.wallets | length")
        for((i=0;i<count;i++));
        do
            accAddr=$(echo "${res}" | jq -r ".result.wallets[${i}].acc.addr")
            getAccBalanceInfo "${accAddr}" "${execName}"
        done
    fi
}

function getAccBalanceInfo() {
    accAddr=$1
    execName=$2
    local res=""
    res=$(curl -k -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"Chain33.GetAllExecBalance","params":[{"addr":"'"${accAddr}"'"}]}' -H 'content-type:text/plain;' ${MAIN_HTTP})
    result=$(echo "${res}" | jq -r ".error | not")
    if [ "${result}" == "true" ] && [ "$(echo "${res}" | jq -r ".result.execAccount")" != "null" ]; then
        #echo "getAccBalanceInfo successfully."
        local count=0
        count=$(echo "$res" | jq -r ".result.execAccount | length")
        for((j=0;j<count;j++));
        do
            accInfo=$(echo "${res}" | jq -r ".result.execAccount[${j}]")
            execer=$(echo "${accInfo}" | jq -r ".execer")
            if [ "${execer}" == "${execName}" ]; then
                echo "${accInfo/\"\"/\"${accAddr}\"}"
            fi
        done
    fi
}

function main() {
    echo "Please input the hostip"
    read ip
    if [ "${ip}" != "" ]; then
        MAIN_HTTP=${ip}
    fi

    echo "Please input the execName"
    read exec
    if [ "${execName}" != "" ]; then
        execName="${exec}"
    fi

    echo "Begin to get ${execName} account info."

    getAllAddrInfo

    #clearTempFile
}

main