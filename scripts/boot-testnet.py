#!/usr/bin/env python

import argparse
import json
import os
import subprocess
import sys
import time
import logging
import re
import random

# Symbols
chainID = 'testing'
mainChainSymbol = 'kuchain'
coreCoinSymbol = 'sys'
coreCoinDenom = '%s/%s' % (mainChainSymbol, coreCoinSymbol)
cliCmd = 'kucli'
nodeCmd = 'kucd'
coinBase = 1000000000000000000

# auths for test
auths = {}
pubkeys = {}
conpubkeys = {}
keyHexs = {}

# node info for test
nodes = {}

logging.basicConfig(level=logging.DEBUG)

args = None

def run(args):
   logging.debug('%s', args)
   if subprocess.call(args, shell=True):
      logging.error('run \"%s\" error, exitting', args)
      sys.exit(1)

def sleep(t):
    time.sleep(t)

def run_output(args):
   logging.debug('%s', args)

   try:
      out_bytes = subprocess.check_output(args, shell=True)
   except subprocess.CalledProcessError as e:
      out_bytes = e.output
      out_text  = out_bytes.decode('utf-8')
      code      = e.returncode
      logging.error('run \"%s\" error by %d and %s', args, code, out_text)
      sys.exit(1)

   return out_bytes.decode('utf-8')

def cli(cmd, noKey = None):
   cliParams = "--home %s/cli/ --keyring-backend test" % (args.home)
   if noKey is not None:
      cliParams = "--home %s/cli/ " % (args.home)
   return run_output('%s/%s %s %s' % (args.build_path, cliCmd, cliParams, cmd))

def cliByHome(home, cmd):
   cliParams = "--home %s" % (home)
   return run_output('%s/%s %s %s' % (args.build_path, cliCmd, cliParams, cmd))

def getNodeHomePath(name):
   return "%s/nodes/%s/" % (args.home, name)

def node(name, cmd):
   cliParams = "--home %s" % (getNodeHomePath(name))
   cmdRun = '%s/%s %s %s' % (args.build_path, nodeCmd, cliParams, cmd)
   return run_output(cmdRun)

def nodeInBackground(name, logPath, cmd):
   cliParams = "--home %s" % (getNodeHomePath(name))
   cmdRun = '%s/%s %s %s' % (args.build_path, nodeCmd, cliParams, cmd)

   with open(logPath, mode='w') as f:
      f.write(cmdRun + '\n')
   subprocess.Popen(cmdRun + '    2>>' + logPath, shell=True)

def nodeByCli(name, cmd):
   cliParams = "--home-client %s/cli/ --keyring-backend test" % (args.home)
   return node(name, '%s %s' % (cliParams, cmd))

def coreCoin(amt):
   return '%s%s' % (amt, coreCoinDenom)

def initWallet():
   logging.debug("init wallet")

   run('rm -rf %s/cli' % (args.home))
   
   genAuth(mainChainSymbol)
   genAuth('test')

   return

def genAuth(name):
   addInfoJSON = cli('keys add ' + name)
   infos = addInfoJSON.splitlines()
   pubkey = infos[3].split(':')[1][1:]

   valAuth = cli('keys show %s -a' % (name))[:-1]
   keyHex = cli('parse --chain-id %s %s' % (chainID, pubkey), True)[:-1]
   keyHex = keyHex.splitlines()[1].split(':')[1][1:]

   keysInfos = cli('parse --chain-id %s %s' % (chainID, keyHex), True)

   conpubkey = keysInfos.splitlines()[6][2:]

   auths[name] = valAuth
   pubkeys[name] = pubkey
   conpubkeys[name] = conpubkey
   keyHexs[name] = keyHex

   logging.debug("gen key %s %s %s %s" % (valAuth, pubkey, conpubkey, keyHex))

   return valAuth

def getAuth(name):
   return auths[name]

def getNodeName(num):
   return "validator%d" % (num)

def addValidatorToGenesis(name):
   valAuth = genAuth(name)
   logging.info("add validator %s %s", name, valAuth)

   # add to genesis
   node(mainChainSymbol, 'genesis add-address %s' % (valAuth))
   node(mainChainSymbol, 'genesis add-account %s %s' % (name, valAuth))
   node(mainChainSymbol, 'genesis add-account-coin %s %s' % (valAuth, coreCoin(10000000000000)))
   node(mainChainSymbol, 'genesis add-account-coin %s %s' % (name, coreCoin(100000000000000000000000000000000)))

def initGenesis(nodeNum):
   mainAuth = getAuth(mainChainSymbol)
   testAuth = getAuth('test')

   node(mainChainSymbol, 'genesis add-address %s' % (mainAuth))
   node(mainChainSymbol, 'genesis add-account %s %s' % (mainChainSymbol, mainAuth))
   node(mainChainSymbol, 'genesis add-coin %s \"%s\"' % (coreCoin(0), "main core"))
   node(mainChainSymbol, 'genesis add-account-coin %s %s' % (mainAuth, coreCoin(100000000000000000000000000000000)))
   node(mainChainSymbol, 'genesis add-account-coin %s %s' % (mainChainSymbol, coreCoin(100000000000000000000000000000000)))

   genesisAccounts = ['testacc1', 'testacc2']
   for genesisAccount in genesisAccounts:
      node(mainChainSymbol, 'genesis add-account %s %s' % (genesisAccount, testAuth))
      node(mainChainSymbol, 'genesis add-account-coin %s %s' % (genesisAccount, coreCoin(10000000000000000000000)))

   for i in range(0, nodeNum):
      addValidatorToGenesis(getNodeName(i + 1))
   
   genTx()

def modifyNodeCfg(name, key, oldValue, newValue=None):
   file = "%s/config/config.toml" % (getNodeHomePath(name))

   patternStr = r"^%s = .+" % (key)
   newStr = '%s = %s' % (key, oldValue)

   if (newValue is not None):
      patternStr = r"^%s = %s" % (key, oldValue)
      newStr = '%s = %s' % (key, newValue)

   with open(file, "r") as f1, open("%s.bak" % file, "w") as f2:
      for line in f1:
         f2.write(re.sub(patternStr, newStr, line))

   os.remove(file)
   os.rename("%s.bak" % file, file)

def appendNodeCfg(name, key, value):
   cliByHome(getNodeHomePath(name), "config %s %s" % (key, value))

def mkNodeDatas(name, num, totalNum):
   if name is not mainChainSymbol:
      # cp genesis from main node
      run('cp %s/config/genesis.json %s/config/genesis.json' % (getNodeHomePath(mainChainSymbol), getNodeHomePath(name)))

      # bind ports, use 3XX56, 3XX57 and 3XX58
      modifyNodeCfg(name, 'proxy_app', '"tcp://127.0.0.1:26658"', '"tcp://127.0.0.1:3%02d58"' % (num))
      modifyNodeCfg(name, 'laddr',     '"tcp://127.0.0.1:26657"', '"tcp://127.0.0.1:3%02d57"' % (num))
      modifyNodeCfg(name, 'laddr',     '"tcp://0.0.0.0:26656"',   '"tcp://0.0.0.0:3%02d56"' % (num))

      # connect to root and next
      toNum = num + 1
      if (toNum > totalNum):
         toNum = 1

      peers = "%s@127.0.0.1:26656," % nodes[mainChainSymbol]['nodeID']
      peers += "%s@127.0.0.1:3%02d56" % (nodes[getNodeName(toNum)]['nodeID'], toNum)

      modifyNodeCfg(name, 'persistent_peers', '"%s"' % peers)

   modifyNodeCfg(name, 'max_num_outbound_peers', '128')
   modifyNodeCfg(name, 'allow_duplicate_ip', 'true')
   appendNodeCfg(name, 'chain-id', chainID)
   appendNodeCfg(name, 'trust-node', 'true')


def initChain(nodeNum):
   logging.debug("init chain")

   initNode(mainChainSymbol, 0)
   initGenesis(nodeNum)

   for i in range(0, nodeNum):
      initNode(getNodeName(i + 1), i + 1)

   for i in range(0, nodeNum):
      mkNodeDatas(getNodeName(i + 1), i + 1, nodeNum)

   mkNodeDatas(mainChainSymbol, 0, nodeNum)
   return

def genTx():
   genTxCmd = 'gentx %s %s --name %s ' % (mainChainSymbol, getAuth(mainChainSymbol), mainChainSymbol)
   if args.sign:
      genTxCmd += '--sign '
   nodeByCli(mainChainSymbol, genTxCmd)
   node(mainChainSymbol, 'collect-gentxs')

def initNode(name, num):
   run('rm -rf %s' % (getNodeHomePath(name)))
   node(name, 'init --chain-id %s %s' % (chainID, name))

   nodeID = node(name, 'tendermint show-node-id')
   logging.info('init node %s, id: %s', name, nodeID)

   nodes[name] = {
      'name'   : name,
      'nodeID' : nodeID[:-1],
      'num'    : num
   }

def startChainNode(name):
   logPath = '%s/logs/%s.log' % (args.home, name)
   run('mkdir -p %s/logs/' % args.home)

   bootParams = 'start --log_level "%s"' % (args.log_level)
   if( args.trace is not None ):
      bootParams += ' --trace'

   nodeInBackground(name, logPath, bootParams)

def message(msg, *args):
    for arg in args:
        print(msg + ' comes from ' + arg)

def tx(fromAcc, cmd, *params):
   cmd = 'tx ' + cmd
   cmd += (" --yes --chain-id=%s" % chainID)
   cmd += (" --from=%s" % fromAcc)
   cmd += (" --node=%s" % "tcp://localhost:30157")

   if len(params) % 2 is not 0:
      logging.error("tx params len should be div by 2")

   paramsKVlen = len(params) / 2
   for pi in range(0, paramsKVlen):
      key = params[pi * 2]
      value = params[pi * 2 + 1]
      cmd += (" --%s=%s" % (key, value))

   logging.debug("run tx: %s", cmd)
   
   cli(cmd)

def createValidator(name):
   # create validator
   cmd = 'kustaking create-validator %s' % name
   tx(name, cmd, 
      "pubkey", node(name, 'tendermint show-validator')[:-1],
      "commission-rate", "0.10",
      "moniker", name)

def voteValidator(acc, node, coin):
   cmd = 'kustaking delegate %s %s %s' % (acc, node, coin)
   tx(acc, cmd)

def regNodes(totalNum):
   logging.info("reg nodes")
   for i in range(0, totalNum):
      createValidator(getNodeName(i + 1))
   
   # wait for create
   sleep(5)

   for i in range(0, totalNum):
      coreCoinBase = 1000000000000000 * coinBase
      voteValidator(getNodeName(i + 1), getNodeName(i + 1), coreCoin(coreCoinBase))



# Parse args
parser = argparse.ArgumentParser()
parser.add_argument('--build-path', metavar='', help='Kuchain build path', default='../build')
parser.add_argument('--home', metavar='', help='testnet data home path', default='./testnet')
parser.add_argument('--trace', action='store_true', help='if --trace to kucd')
parser.add_argument('--log-level', metavar='', help='log level for kucd', default='*:info')
parser.add_argument('--node-num', type=int, metavar='', help='val node number', default=12)
parser.add_argument('--sign', type=bool, metavar='', help='if sign genesis trx auto', default=True)

args = parser.parse_args()
logging.debug("args %s", args)

# Start Chain
logging.info("start kuchain testnet by %s to %s", args.home, args.build_path)

initWallet()
initChain(int(args.node_num))


# Start main first
startChainNode(mainChainSymbol)
sleep(1)

for i in range(0, int(args.node_num)):
   startChainNode(getNodeName(i + 1))

sleep(5)

regNodes(int(args.node_num))