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

# auth for test
mainAuth = None
testAuth = None

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
   cli('keys add ' + mainChainSymbol) # add for root auth
   cli('keys add test') # add for auth for test

   return

def initChain():
   logging.debug("init chain")

   run('rm -rf %s/node' % (args.home))
   node('init --chain-id %s %s' % (chainID, chainID))
   
   node('genesis add-address %s' % (mainAuth))
   node('genesis add-account %s %s' % (mainChainSymbol, mainAuth))
   node('genesis add-coin %s \"%s\"' % (coreCoin(1000000000000000000000000000000000000000), "main core"))
   node('genesis add-account-coin %s %s' % (mainAuth, coreCoin(100000000000000000000000000000000)))
   node('genesis add-account-coin %s %s' % (mainChainSymbol, coreCoin(100000000000000000000000000000000)))

   genesisAccounts = ['testacc1', 'testacc2']
   for genesisAccount in genesisAccounts:
      node('genesis add-account %s %s' % (genesisAccount, testAuth))
      node('genesis add-account-coin %s %s' % (genesisAccount, coreCoin(10000000000000000000000)))

   # gentx
   nodeByCli('gentx %s --name %s ' % (mainAuth, mainChainSymbol))
   node('collect-gentxs')

   return

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

args = parser.parse_args()
logging.debug("args %s", args)

# Start Chain
logging.info("start kuchain testnet by %s to %s", args.home, args.build_path)

initWallet()

mainAuth = cli('keys show %s -a' % mainChainSymbol)
mainAuth = mainAuth[:-1]
logging.debug("main auth : %s", mainAuth)

testAuth = cli('keys show test -a')
testAuth = testAuth[:-1]
logging.debug("test auth : %s", testAuth)

initChain()
startChainByOneNode()