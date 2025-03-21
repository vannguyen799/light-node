#!/bin/bash

mapfile -t privateKeys < privateKeys.txt
mapfile -t proxies < proxies.txt



for i in "${!privateKeys[@]}"; do
    echo "Starting executor with Proxy: ${proxies[$i]}"
    pk=${privateKeys[$i]}
    proxy=${proxies[$i]}
    env http_proxy="http://$proxy" https_proxy="http://$proxy" no_proxy="localhost,127.0.0.1" \
        GRPC_URL=grpc.testnet.layeredge.io:9090 \
        CONTRACT_ADDR=cosmos1ufs3tlq4umljk0qfe8k5ya0x6hpavn897u2cnf9k0en9jr7qarqqt56709 \
        ZK_PROVER_URL=https://layeredge.mintair.xyz \
        API_REQUEST_TIMEOUT=1000 \
        POINTS_API=https://light-node.layeredge.io \
        PRIVATE_KEY=$pk \
        ./light-node &
    sleep 1
done

wait
echo "All executors have finished."


