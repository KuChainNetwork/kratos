#!/usr/bin/env bash

ROOT_DIR=$1
NULL=
CHAIN_ID="testing"

printf "params $1 $2\\n"

ROOT_DIR="${ROOT_DIR}/${CHAIN_ID}"

mkdir -p ${ROOT_DIR}
mkdir -p ${ROOT_DIR}/node/config
mkdir -p ${ROOT_DIR}/cli/

printf "start testnet for kuchain ${CHAIN_ID} in ${ROOT_DIR}\\n"

PARAMS="--home ${ROOT_DIR}/node/"
PARAMSCLI="--home ${ROOT_DIR}/cli/ --keyring-backend test"

./build/kucli ${PARAMSCLI} keys add kuchain

VALKEY=$(./build/kucli ${PARAMSCLI} keys show kuchain -a)

./build/kucli ${PARAMSCLI} keys add test

TESTVALKEY=$(./build/kucli ${PARAMSCLI} keys show test -a)

printf "current val key ${VALKEY}\\n"

./build/kucd init ${PARAMS} --chain-id=${CHAIN_ID} ${CHAIN_ID}
./build/kucd ${PARAMS} genesis add-account kuchain ${VALKEY}
./build/kucd ${PARAMS} genesis add-account testacc1 ${TESTVALKEY}
./build/kucd ${PARAMS} genesis add-account testacc2 ${TESTVALKEY}
./build/kucd ${PARAMS} genesis add-address ${VALKEY}
./build/kucd ${PARAMS} genesis add-coin "1000000000000000000000000000000000000000kuchain/sys" "main token"
#./build/kucd ${PARAMS} genesis add-coin "1000000000000000000000000000000000000000validatortoken" "for staking"
./build/kucd ${PARAMS} genesis add-account-coin ${VALKEY} "100000000000000000000000kuchain/sys"
./build/kucd ${PARAMS} genesis add-account-coin kuchain "100000000000000000000000kuchain/sys"

printf "./build/kucd ${PARAMS} gentx ${VALKEY} --keyring-backend test --name kuchain --home-client ${ROOT_DIR}/cli/\\n"
./build/kucd ${PARAMS} gentx ${VALKEY} --keyring-backend test --name kuchain --home-client ${ROOT_DIR}/cli/


./build/kucd ${PARAMS} collect-gentxs


if ["$3" == "$NULL"]; then
   printf "./build/kucd ${PARAMS} start --log_level \"*:debug\" --trace\\n"
   ./build/kucd ${PARAMS} start --log_level "*:debug" --trace
else
   PluginPath=$3
   printf "use plugin path ${PluginPath}\\n"
   printf "./build/kucd ${PARAMS} start --log_level \"*:debug\" --trace --plugin-cfg \"${PluginPath}\"\\n"
   ./build/kucd ${PARAMS} start --log_level "*:debug" --trace --plugin-cfg "${PluginPath}"
fi


