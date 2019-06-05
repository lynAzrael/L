#!/bin/bash

MAIN_HTTP="http://localhost:9901"
execName="user.p.guodun.token"
txTempFile="/root/.txTempFile"
ownerTxTempFile="/root/.ownerTxTempFile"
ownerTempFile="/root/.ownerTempFile"
addrTxTempFile="/root/.addrTxTempFile"
addrTempFile="/root/.addrTempFile"
addrBalanceInfo="/root/.addrBalance"

function checkTempFile() {
    if [ -f "${txTempFile}" ]; then
        rm "${txTempFile}"
        
    fi

    if [ -f "${ownerTempFile}" ]; then
        rm "${ownerTempFile}"
        
    fi

    if [ -f "${ownerTxTempFile}" ]; then
        rm "${ownerTxTempFile}"
        
    fi

    if [ -f "${addrTempFile}" ];then
        rm "${addrTempFile}"
    fi

    if [ -f "${addrBalanceInfo}" ]; then
        rm "${addrBalanceInfo}"
    fi

    touch "${ownerTempFile}"
    touch "${txTempFile}"
    touch "${ownerTxTempFile}"
    touch "${addrTempFile}"
    touch "${addrBalanceInfo}"
}

function clearTempFile() {
    if [ -f "${txTempFile}" ]; then
        rm "${txTempFile}"
        
    fi

    if [ -f "${ownerTempFile}" ]; then
        rm "${ownerTempFile}"
        
    fi

    if [ -f "${ownerTxTempFile}" ]; then
        rm "${ownerTxTempFile}"
        
    fi

    if [ -f "${addrTempFile}" ];then
        rm "${addrTempFile}"
    fi
}

function queryTxTransaction() {
    txHash=$1
    validator=$2
    expectRes=$3
    res=$(curl -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"Chain33.QueryTransaction","params":[{"hash":"'"${txHash}"'"}]}' -H 'content-type:text/plain;' ${MAIN_HTTP} | jq -r "${validator}")
    if [ "${res}" != "${expectRes}" ]; then
        return 1
    else
        return 0
    fi
}

function getAddrByExecName() {
    execAddr=$(curl -ksd '{"method":"Chain33.ConvertExectoAddr","params":[{"execname":"'"${execName}"'"}]}' ${MAIN_HTTP} | jq -r ".result")
    if [ "${execName}" != "" ] && [ "${execAddr}" != "" ]; then
        return 0
    else
        return 1
    fi
}

function getAddrTxInfo() {
    addr=$1
    tempfile=$2
    flag=$3
    local res=""
    res=$(curl -k -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"Chain33.GetTxByAddr","params":[{"addr":"'"${addr}"'", "flag":'"${flag}"', "count":1000000, "direction":0, "height": -1, "index":0}]}' -H 'content-type:text/plain;' ${MAIN_HTTP} | jq -r ".result.txInfos")
    if [ "${res}" != "" ] && [ $(echo "${res}" | jq -r "length") -gt 0 ]; then 
        count=$(echo "${res}" | jq -r "length")
        for ((i=0;i<count;i++))
        do
            txhash=$(echo "${res}" | jq -r ".[$i].hash")
            if [ "${txhash}" != "" ]; then
                echo "${txhash}" >> ${tempfile}
            fi
        done
    fi
}

function addOwnerAddr() {
    if [ ${#owner} == 34 ]; then
        if [ $(grep -c ${owner} ${ownerTempFile}) -eq 0 ];then
            echo "${owner}" >> ${ownerTempFile}
        fi

        if [ $(grep -c ${owner} ${addrTempFile}) -eq 0 ];then
            echo "${owner}" >> ${addrTempFile}
        fi
    fi
}

function getOwnerAddrByTxHash() {
    for hash in $(cat "${txTempFile}")
    do
        local res=""
        res=$(curl -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"Chain33.QueryTransaction","params":[{"hash":"'"${hash}"'"}]}' -H 'content-type:text/plain;' ${MAIN_HTTP})
        if [ "${res}" != "" ]; then
            txType=$(echo "${res}" | jq -r ".result.tx.payload.Ty")
            case ${txType} in
                7)
            {
                owner=$(echo "${res}" | jq -r ".result.tx.payload.tokenPreCreate.owner")
                addOwnerAddr
            }
            ;;
                8)
            {
                owner=$(echo "${res}" | jq -r ".result.tx.payload.tokenFinishCreate.owner")
                addOwnerAddr
            }
            ;;
                9)
            {
                owner=$(echo "${res}" | jq -r ".result.tx.payload.tokenRevokeCreate.owner")
                addOwnerAddr
            }
            ;;
                *)
            {
                echo "Unknown tx type: ${txType}, and hash is ${hash}"
            }
            ;;
            esac
        fi
    done
}

function getAddrByOwnerTxHash() {
    for ownerTxHash in $(cat "${ownerTxTempFile}")
    do
        local res=""
        res=$(curl -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"Chain33.QueryTransaction","params":[{"hash":"'"${hash}"'"}]}' -H 'content-type:text/plain;' ${MAIN_HTTP})
        if [ "${res}" != "" ]; then
            echo "${res}"
        fi
    done
}

function getUserAddrByOwnerAddr() {
    echo "Begin to execute function getUsrAddrByOwnerAddr."
    for addr in $(cat "${ownerTempFile}")
    do
        getAddrTxInfo "${addr}" "${ownerTxTempFile}" 1
    done

    for txHash in $(cat "${ownerTxTempFile}")
    do 
        queryTxTransaction "${txHash}" ".result.receipt.tyName" "ExecOk"
        if [ $? -eq 0 ]; then
            local res=""
            res=$(curl -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"Chain33.QueryTransaction","params":[{"hash":"'"${txHash}"'"}]}' -H 'content-type:text/plain;' ${MAIN_HTTP})
            if [ "${res}" != "" ]; then
                local toAddr=""
                toAddr=$(echo "${res}" | jq -r ".result.tx.payload.transfer.to")
                while [ "${toAddr}" != "" ] && [ $(grep -c ${toAddr} ${addrTempFile}) -eq 0 ]
                do
                    echo "${toAddr}" >> ${addrTempFile}

                    res=$(curl -k -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"Chain33.GetTxByAddr","params":[{"addr":"'"${toAddr}"'", "flag":1, "count":1000000, "direction":0, "height": -1, "index":0}]}' -H 'content-type:text/plain;' ${MAIN_HTTP} | jq -r ".result.txInfos")
                    if [ "${res}" != "" ] && [ $(echo "${res}" | jq -r "length") -gt 0 ]; then 
                        count=$(echo "${res}" | jq -r "length")
                        for ((i=0;i<count;i++))
                        do
                            txhash=$(echo "${res}" | jq -r ".[$i].hash")
                            if [ "${txhash}" != "" ]; then
                                local txRes=""
                                txRes=$(curl -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"Chain33.QueryTransaction","params":[{"hash":"'"${txhash}"'"}]}' -H 'content-type:text/plain;' ${MAIN_HTTP})
                                if [ "${txRes}" != "" ]; then
                                    toAddr=$(echo "${txRes}" | jq -r ".result.tx.payload.transfer.to")
                                fi
                            fi
                        done
                    fi
                done
            fi
        fi
    done
}

function getBalanceInfo() {
    for addr in $(cat "${addrTempFile}")
    do
        local res=""
        res=$(curl -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"Chain33.Query","params":[{"execer": "'"${execName}"'","funcName":"GetAccountTokenAssets","payload": {"address":"'"${addr}"'", "execer": "token"}}]}' -H 'content-type:text/plain;' ${MAIN_HTTP})
        if [ "${res}" != "" ]; then
            local count=0
            count=$(echo "$res" | jq -r ".result.tokenAssets | length")
            for((i=0;i<count;i++));
            do
                symbol=$(echo "${res}" | jq -r ".result.tokenAssets[${i}].symbol")
                getTokenBalanceInfo "${symbol}" "${addr}"
        done
        fi
    done
}

function getTokenBalanceInfo() {
    symbol=$1
    addr=$2
    #local res=""
    # res=$(curl -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"token.GetTokenBalance","params":[{"addresses": ["'${addr}'"],"tokenSymbol":"'"${symbol}"'","execer": "'"${execName}"'"}]}' -H 'content-type:text/plain;' ${MAIN_HTTP})
    # if [ "${res}" != "" ]; then
    #     balance=$(echo "${res}" | jq -r ".result[0].balance")
    #     newBalance=$(awk 'BEGIN{printf "%0.2f\n",'${balance}'/100000000}')
    #     #echo "${res}" | jq -r ".result[0]" | sed  "/currency/i\ \ \"Token\":\"${symbol}\"" | sed 's/\(balance\":\ \)[0-9]*/\1 '${newBalance}'/'
    #     echo "${res}" | jq -r ".result[0]" | sed  "/currency/i\ \ \"Token\":\"${symbol}\"" 
    # fi

    ./chain33-cli token token_balance -a "${addr}" -e "user.p.guodun.token" -s "${symbol}" --rpc_laddr "${MAIN_HTTP}" | jq -r ".[0]" >> "${addrBalanceInfo}"
}

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

function getTxHashByAddr() {
    echo "Please input the hostip"
    read ip
    if [ "${ip}" != "" ]; then
        MAIN_HTTP=${ip}
    fi

    echo "Please input the execName"
    read exec
    if [ "${exec}" != "" ]; then
        execName="${exec}"
    fi

    echo "Begin to get ${execName} account info."

    # get executor addr by executor name
    getAddrByExecName
    if [ $? -ne 0 ]; then
        echo "get addr by executor name falied."
        exit 1
    fi

    # 
    getAddrTxInfo "${execAddr}" "${txTempFile}" 0
}

function main() {
    
    checkTempFile

    getTxHashByAddr

    getOwnerAddrByTxHash

    getUserAddrByOwnerAddr
     
    getBalanceInfo

    clearTempFile
}

main