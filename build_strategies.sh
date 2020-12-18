#! /bin/sh

set -xe

echo "------> Building JS allocating strategies "
yarn --cwd  `dirname $0`/pools/allocating_strategies/strategies
#FIXME yarn --cwd  `dirname $0`/pools/allocating_strategies/strategies test
yarn --cwd  `dirname $random_s_int32_strategy.js:0`/pools/allocating_strategies/strategies generate:all
