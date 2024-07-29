#!/bin/bash

cd ..

cleanup() {
  echo "Cleaning up..."
  kill $ganache_pid 2>/dev/null
  sync
  sudo sh -c 'echo 3 > /proc/sys/vm/drop_caches'
}

handle_ctrl_c() {
  echo "Ctrl+C pressed. Cleaning up..."
  cleanup
  echo "Exiting..."
  exit 1
}

trap cleanup EXIT
trap handle_ctrl_c SIGINT


no_order=(2000)
no_matcher=(1 2 5 7 10)
no_send_to=(1 1 1 1 1)
for n in "${no_order[@]}"
do
    echo Running $n orders...
    for ((i=0; i<${#no_matcher[@]}; i++));
    do
        echo Send to ${no_send_to[i]}/${no_matcher[i]} matcher\(s\)
        mkdir -p ./data/result/test_1/price_curve_"${no_send_to[i]}"_"${no_matcher[i]}"_$n
        data_file="./data/real_orders/real_orders_$n.json"
        ganache -a 200 -m '' -e 99999999999 --chain.chainId 1337 --p 8545 > /dev/null 2>&1 &
        ganache_pid=$!
        sleep 2
        go run . 8545 1337 "${no_matcher[i]}" "${no_send_to[i]}" run $data_file ./data/result/test_1/price_curve_"${no_send_to[i]}"_"${no_matcher[i]}"_$n > ./data/test_2/price_curve_"${no_send_to[i]}"_"${no_matcher[i]}"_$n/log
        kill $ganache_pid 2>/dev/null
        wait $ganache_pid 2>/dev/null
        sync
        sudo sh -c 'echo 3 > /proc/sys/vm/drop_caches'
        echo
    done
done


no_order=(1000 1500 2000)
no_matcher=(10)
no_send_to=(2)
for n in "${no_order[@]}"
do
    echo Running $n orders...
    for ((i=0; i<${#no_matcher[@]}; i++));
    do
        echo Send to ${no_send_to[i]}/${no_matcher[i]} matcher\(s\)
        mkdir -p ./data/result/compare_with_clob/price_curve_"${no_send_to[i]}"_"${no_matcher[i]}"_$n
        data_file="./data/real_orders/real_orders_$n.json"
        ganache -a 200 -m '' -e 99999999999 --chain.chainId 1337 --p 8545 > /dev/null 2>&1 &
        ganache_pid=$!
        sleep 2
        go run . 8545 1337 "${no_matcher[i]}" "${no_send_to[i]}" run $data_file ./data/result/compare_with_clob/price_curve_"${no_send_to[i]}"_"${no_matcher[i]}"_$n > ./data/test_2/price_curve_"${no_send_to[i]}"_"${no_matcher[i]}"_$n/log
        kill $ganache_pid 2>/dev/null
        wait $ganache_pid 2>/dev/null
        sync
        sudo sh -c 'echo 3 > /proc/sys/vm/drop_caches'
        echo
    done
done


no_order=(1000 1500 2000)
no_matcher=(10)
no_send_to=(2)
for n in "${no_order[@]}"
do
    echo Running $n orders...
    for ((i=0; i<${#no_matcher[@]}; i++));
    do
        echo Send to ${no_send_to[i]}/${no_matcher[i]} matcher\(s\)
        mkdir -p ./data/result/compare_with_clob/price_curve_"${no_send_to[i]}"_"${no_matcher[i]}"_$n
        data_file="./data/real_orders/real_orders_$n.json"
        ganache -a 200 -m '' -e 99999999999 --chain.chainId 1337 --p 8545 > /dev/null 2>&1 &
        ganache_pid=$!
        sleep 2
        go run . 8545 1337 "${no_matcher[i]}" "${no_send_to[i]}" run $data_file ./data/result/compare_with_clob/price_curve_"${no_send_to[i]}"_"${no_matcher[i]}"_$n > ./data/test_2/price_curve_"${no_send_to[i]}"_"${no_matcher[i]}"_$n/log
        kill $ganache_pid 2>/dev/null
        wait $ganache_pid 2>/dev/null
        sync
        sudo sh -c 'echo 3 > /proc/sys/vm/drop_caches'
        echo
    done
done