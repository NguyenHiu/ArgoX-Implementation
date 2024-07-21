#!/bin/bash

no_order=(1000 5000 10000)
no_matcher=(1 2 5 10)
no_send_to=(1 1 2 5)

cd ..

cleanup() {
  echo "Cleaning up..."
  kill $ganache_pid 2>/dev/null
}

handle_ctrl_c() {
  echo "Ctrl+C pressed. Cleaning up..."
  cleanup
  echo "Exiting..."
  exit 1
}

trap cleanup EXIT
trap handle_ctrl_c SIGINT

for n in "${no_order[@]}"
do
    echo Running $n orders...
    for ((i=0; i<${#no_matcher[@]}; i++));
    do
        echo Send to ${no_send_to[i]}/${no_matcher[i]} matcher\(s\)
        data_file="./data/_real_orders_$n.json"
        ganache -a 200 -m '' -e 99999999999 --chain.chainId 1337 --p 8545 > /dev/null 2>&1 &
        ganache_pid=$!
        sleep 2
        go run . 8545 1337 "${no_matcher[i]}" "${no_send_to[i]}" run $data_file ./data/priceCurve_"${no_send_to[i]}"_"${no_matcher[i]}"_$n
        kill $ganache_pid 2>/dev/null
        echo
    done
done

