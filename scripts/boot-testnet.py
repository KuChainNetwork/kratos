#!/usr/bin/env python

import argparse
import json
import os
import subprocess
import sys
import time
import logging

# Keys
testKey = "kuchain1ysggxvhq3aqxp2dnzw8ucqf0hgjn44tqcp7l3s"
testKeyMnemonic = "warm law where bid turtle tenant story logic air ancient gesture way main rabbit sock enlist hollow wealth stereo position fiscal expand mosquito latin"

# Symbols
chainID = 'testing'
mainChainSymbol = 'kuchain'
coreCoinSymbol = 'sys'
coreCoinDenom = '%s/%s' % (mainChainSymbol, coreCoinSymbol)
cliCmd = 'kucli'
nodeCmd = 'kucd'

# auths for test
auths = {}

logging.basicConfig(level=logging.DEBUG)

args = None

def run(args):
   logging.debug('%s', args)
   if subprocess.call(args, shell=True):
      logging.error('run \"%s\" error, exitting', args)
      sys.exit(1)

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

def cli(cmd):
   cliParams = "--home %s/cli/ --keyring-backend test" % (args.home)
   return run_output('%s/%s %s %s' % (args.build_path, cliCmd, cliParams, cmd))

def node(cmd):
   cliParams = "--home %s/node/" % (args.home)
   return run('%s/%s %s %s' % (args.build_path, nodeCmd, cliParams, cmd))

def nodeByCli(cmd):
   cliParams = "--home %s/node/ --home-client %s/cli/ --keyring-backend test" % (args.home, args.home)
   return run('%s/%s %s %s' % (args.build_path, nodeCmd, cliParams, cmd))

def coreCoin(amt):
   return '%s%s' % (amt, coreCoinDenom)

def initWallet():
   logging.debug("init wallet")

   run('rm -rf %s/cli' % (args.home))
   
   genAuth(mainChainSymbol)
   genAuth('test')

   return

def genAuth(name):
   cli('keys add ' + name)
   valAuth = cli('keys show %s -a' % (name))
   valAuth = valAuth[:-1]
   auths[name] = valAuth
   return valAuth

def getAuth(name):
   return auths[name]

def addValidator(name):
   valAuth = genAuth(name)
   logging.info("add validator %s %s", name, valAuth)

   # add to genesis
   node('genesis add-address %s' % (valAuth))
   node('genesis add-account %s %s' % (name, valAuth))
   node('genesis add-account-coin %s %s' % (valAuth, coreCoin(10000000000000)))
   node('genesis add-account-coin %s %s' % (name, coreCoin(100000000000000000000000000000000)))

def initChain(nodeNum):
   logging.debug("init chain")

   run('rm -rf %s/node' % (args.home))
   node('init --chain-id %s %s' % (chainID, chainID))

   mainAuth = getAuth(mainChainSymbol)
   testAuth = getAuth('test')
   
   node('genesis add-address %s' % (mainAuth))
   node('genesis add-account %s %s' % (mainChainSymbol, mainAuth))
   node('genesis add-coin %s \"%s\"' % (coreCoin(1000000000000000000000000000000000000000), "main core"))
   node('genesis add-account-coin %s %s' % (mainAuth, coreCoin(100000000000000000000000000000000)))
   node('genesis add-account-coin %s %s' % (mainChainSymbol, coreCoin(100000000000000000000000000000000)))

   genesisAccounts = ['testacc1', 'testacc2']
   for genesisAccount in genesisAccounts:
      node('genesis add-account %s %s' % (genesisAccount, testAuth))
      node('genesis add-account-coin %s %s' % (genesisAccount, coreCoin(10000000000000000000000)))

   for i in range(1, nodeNum):
      addValidator("validator%d" % (i))

   return

def genTx():
   nodeByCli('gentx %s --name %s ' % (getAuth(mainChainSymbol), mainChainSymbol))
   node('collect-gentxs')

def startChainByOneNode():
   bootParams = 'start --log_level "%s"' % (args.log_level)
   if( args.trace is not None ):
      bootParams += ' --trace'

   node(bootParams)

# Parse args
parser = argparse.ArgumentParser()
parser.add_argument('--build-path', metavar='', help='Kuchain build path', default='../build')
parser.add_argument('--home', metavar='', help='testnet data home path', default='./testnet')
parser.add_argument('--trace', action='store_true', help='if --trace to kucd')
parser.add_argument('--log-level', metavar='', help='log level for kucd', default='*:debug')
parser.add_argument('--node-num', type=int, metavar='', help='val node number', default=5)

args = parser.parse_args()
logging.debug("args %s", args)

# Start Chain
logging.info("start kuchain testnet by %s to %s", args.home, args.build_path)

initWallet()
initChain(int(args.node_num))
genTx()
startChainByOneNode()