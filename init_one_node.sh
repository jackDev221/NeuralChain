#!/bin/bash
# NOTE
# the first you have to grants chmod allow can read this file 
# run ./init_one_node.sh 
# enjoy
#

echo "***********************************************"
echo "***********************************************"
echo "    INIT YOUR NODE MUST ENTER SOME INFOR BELOW"
datadir=~/evrynet
listAcc=$(./gev --datadir "$datadir" account list)
hasAccount=0
address=''
password=''

if [[ -z "$listAcc" ]]; then
    echo "    In the localhost there are not accounts, you must create or import an account"
else
    echo "    In the localhost there are accounts here: $listAcc "
    hasAccount=1
fi
echo "***********************************************"

read -p "    Do you have an account? (y|n): " yn
if [[ $yn = 'y' ]]; then
    if [[ $hasAccount = 1 ]]; then
        read -p "    Do you want to use one of accounts above? (y|n): " yn
        if [[ $yn = 'y' ]]; then
            read -p "    Enter your address: "  address
            read -p "    Enter your passphrase to unlock: " -s password
            ./gev --datadir "$datadir" --port 30311 --rpc --rpcaddr 'localhost' --networkid 15 --gasprice '0' --password <(echo "$password") --unlock "$address" --etherbase "$address" --mine --allow-insecure-unlock
            exit 1
        fi
    fi

    echo    "    starting import your account with private key..."
    read -p "    Enter your private key: "  privateKey
    read -p "    Enter your passphrase to unlock: " -s password
    result=$(./gev account import --datadir "$datadir" --password <(echo "$password") <(echo "$privateKey"))
    result="$(cut -d'{' -f2 <<<"$result")"
    address="$(cut -d'}' -f1 <<<"$result")"
    aLength="${#address}"
    if [[ $aLength != 40 ]]; then
        exit 1
    fi
else
    read -p "    Enter your passphrase for new account: " -s password
    ./gev account new --datadir "$datadir" --password <(echo "$password")
    read -p "    Please enter the address have been created above to unlock and run your node: " address
fi

echo "***********************************************"
read -p "    Do you want to set the "$address" as coinbase? (y|n): " yn
echo "***********************************************"

if [[ $yn = 'y' ]]; then
    ./gev --datadir "$datadir" --port 30311 --rpc --rpcaddr 'localhost' --networkid 15 --gasprice '0' --password <(echo "$password") --unlock "$address" --etherbase "$address" --mine --allow-insecure-unlock
else
    ./gev --datadir "$datadir" --port 30311 --rpc --rpcaddr 'localhost' --networkid 15 --gasprice '0' --password <(echo "$password") --unlock "$address" --mine --allow-insecure-unlock
fi
