#!/usr/bin/env bash

ROOT_DIR=$1
NULL=
CHAIN_ID="testing"

printf "params $1 $2\\n"

ROOT_DIR="${ROOT_DIR}/${CHAIN_ID}"

mkdir -p ${ROOT_DIR}
mkdir -p ${ROOT_DIR}/node/config
mkdir -p ${ROOT_DIR}/cli/

printf "start testnet for kratos ${CHAIN_ID} in ${ROOT_DIR}\\n"

PARAMS="--home ${ROOT_DIR}/node/"
PARAMSCLI="--home ${ROOT_DIR}/cli/ --keyring-backend test"

./build/ktscli ${PARAMSCLI} keys add kratos

VALKEY=$(./build/ktscli ${PARAMSCLI} keys show kratos -a)

./build/ktscli ${PARAMSCLI} keys add test

TESTVALKEY=$(./build/ktscli ${PARAMSCLI} keys show test -a)

printf "current val key ${VALKEY}\\n"

./build/ktsd init ${PARAMS} --chain-id=${CHAIN_ID} ${CHAIN_ID}
./build/ktsd ${PARAMS} genesis add-account kratos ${VALKEY}
./build/ktsd ${PARAMS} genesis add-account testacc1 ${TESTVALKEY}
./build/ktsd ${PARAMS} genesis add-account testacc2 ${TESTVALKEY}
./build/ktsd ${PARAMS} genesis add-address ${VALKEY}
./build/ktsd ${PARAMS} genesis add-coin "1000000000000000000000000000000000000000kratos/kts" "main token"
#./build/ktsd ${PARAMS} genesis add-coin "1000000000000000000000000000000000000000validatortoken" "for staking"
./build/ktsd ${PARAMS} genesis add-account-coin ${VALKEY} "100000000000000000000000kratos/kts"
./build/ktsd ${PARAMS} genesis add-account-coin kratos "100000000000000000000000kratos/kts"

printf "./build/ktsd ${PARAMS} gentx ${VALKEY} --keyring-backend test --name kratos --home-client ${ROOT_DIR}/cli/\\n"
./build/ktsd ${PARAMS} gentx ${VALKEY} --keyring-backend test --name kratos --home-client ${ROOT_DIR}/cli/


./build/ktsd ${PARAMS} collect-gentxs


if ["$3" == "$NULL"]; then
   printf "./build/ktsd ${PARAMS} start --log_level \"*:debug\" --trace\\n"
   ./build/ktsd ${PARAMS} start --log_level "*:debug" --trace
else
   PluginPath=$3
   printf "use plugin path ${PluginPath}\\n"
   printf "./build/ktsd ${PARAMS} start --log_level \"*:debug\" --trace --plugin-cfg \"${PluginPath}\"\\n"
   ./build/ktsd ${PARAMS} start --log_level "*:debug" --trace --plugin-cfg "${PluginPath}"
fi


