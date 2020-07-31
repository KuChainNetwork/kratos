# Kuchain

## 0. Build

Just use make to build: 

```bash
git clone https://github.com/KuChainNetwork/kuchain.git
cd kuchain
go mod vendor
make all
```

`kucd` and `kucli` will in build path.

## 1. Start a testnet

To start a testnet, should gen a genesis json config for testnet,
There need set genesis accounts, genesis coins, and genesis assets for each accounts.

### **init config**

```bash
cd build
./kucd init --chain-id=testing testing
```

### **add genesis accounts**

```bash
./kucd genesis add-account kuchain $(./kucli keys show validator -a)
./kucd genesis add-account testacc1 $(./kucli keys show validator -a)
./kucd genesis add-account testacc2 $(./kucli keys show validator -a)
./kucd genesis add-account validator $(./kucli keys show validator -a)
```

### **add a genesis address account for geneis validator

```bash
./kucd genesis add-address $(./kucli keys show validator -a)
```

### **add coins**

```bash
./kucd genesis add-coin "10000000000000stake" "for staking"
./kucd genesis add-coin "10000000000000validatortoken" "for staking"
./kucd genesis add-account-coin $(./kucli keys show validator -a) "200000000stake,900000000validatortoken"
./kucd genesis add-account-coin "validator" "200000000000stake,900000000validatortoken"
./kucd genesis add-account-coin "testacc1" "200000000000stake,900000000validatortoken"
./kucd genesis add-account-coin "testacc2" "200000000000stake,900000000validatortoken"
./kucd genesis add-account-coin "kuchain" "200000000000stake,900000000validatortoken"
```

you can add other genesis coins:

```bash
./kucd genesis add-coin "10000000000000kuchain/eos" "eos map token"
./kucd genesis add-account-coin "testacc1" "200000000kuchain/eos"
```

### **gen genesis tx**

```bash
./kucd gentx validator --name validator
./kucd collect-gentxs
```

For all:

```bash
./kucd init --chain-id=testing testing
./kucd genesis add-account kuchain $(./kucli keys show validator -a)
./kucd genesis add-account testacc1 $(./kucli keys show validator -a)
./kucd genesis add-account testacc2 $(./kucli keys show validator -a)
./kucd genesis add-account validator $(./kucli keys show validator -a)
./kucd genesis add-address $(./kucli keys show validator -a)
./kucd genesis add-coin "10000000000000stake" "for staking"
./kucd genesis add-coin "10000000000000validatortoken" "for staking"
./kucd genesis add-account-coin $(./kucli keys show validator -a) "200000000stake,900000000validatortoken"
./kucd genesis add-account-coin "validator" "200000000000stake,900000000validatortoken"
./kucd genesis add-account-coin "testacc1" "200000000000stake,900000000validatortoken"
./kucd genesis add-account-coin "testacc2" "200000000000stake,900000000validatortoken"
./kucd genesis add-account-coin "kuchain" "200000000000stake,900000000validatortoken"
./kucd gentx validator --name validator
./kucd collect-gentxs
./kucd start --log_level "*:debug" --trace
```

If ok, the chain will produce blocks:

```bash
I[2020-03-27|19:26:34.490] starting ABCI with Tendermint                module=main 
I[2020-03-27|19:26:39.759] Executed block                               module=state height=1 validTxs=0 invalidTxs=0
I[2020-03-27|19:26:39.768] Committed state                              module=state height=1 txs=0 appHash=1B4D486D319E94A6F8DBBFF6DA80CFBB1A7BAB26821CE02F4FF5AAFCD3B282C0
I[2020-03-27|19:26:44.829] Executed block                               module=state height=2 validTxs=0 invalidTxs=0
```

## 2. Start testnet by script

In Kuchain we can start a testnet by script:

```bash
git clone https://github.com/KuChainNetwork/kuchain.git
cd kuchain
go mod vendor
make all
./scripts/boot-testnet.sh ./testchain testing
```

In this case, you can use tx like this:

```bash
./build/kucli tx asset transfer kuchain kuchain1jzjrqc7dftedr8k3jegm2ad8mlacunlueavkcx 100kuchain/kcs --chain-id testing --keyring-backend test --home ./testchain/testing/cli --from validator
```

This cmd transfer `100kuchain/kcs` from account `kuchain` to address `kuchain1jzjrqc7dftedr8k3jegm2ad8mlacunlueavkcx`.

Script will use test keyring, so this cmd need params:

- `--chain-id` testing
- `--keyring-backend` test
- `--home` ./testchain/testing/cli

most help info can get from `-h`

```bash
./build/kucli tx asset transfer -h
Transfer coins and sign a trx

Usage:
  kucli tx asset transfer [from] [to] [coins] [flags]

Flags:
  -a, --account-number uint      The account number of the signing account (offline mode only)
 ...

Global Flags:
      --chain-id string   Chain ID of tendermint node
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/Users/fy/.kucli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```