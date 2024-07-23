#!/bin/bash

no_order=(1000 5000)
no_matcher=(1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20)
no_send_to=(1 2 3 4 5)

cd ..

cleanup() {
  echo "Cleaning up..."
  kill $ganache_pid 2>/dev/null
  sync
  echo 3 > /proc/sys/vm/drop_caches
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
      for ((j=0; j<${#no_send_to[@]}; j++));
      do
          echo Send to ${no_send_to[j]}/${no_matcher[i]} matcher\(s\)
          # mkdir -p ./data/priceCurve_"${no_send_to[j]}"_"${no_matcher[i]}"_$n
          data_file="./data/_real_orders_$n.json"
          ganache -a 200 -m '' -e 99999999999 --chain.chainId 1337 --p 8545 > /dev/null 2>&1 &
          ganache_pid=$!
          sleep 2
          go run . 8545 1337 "${no_matcher[i]}" "${no_send_to[j]}" run $data_file ./data/priceCurve_"${no_send_to[j]}"_"${no_matcher[i]}"_$n > ./data/priceCurve_"${no_send_to[j]}"_"${no_matcher[i]}"_$n/log
          kill $ganache_pid 2>/dev/null
          wait $ganache_pid 2>/dev/null
          sync
          echo 3 > /proc/sys/vm/drop_caches
          echo
        done 
    done
done
