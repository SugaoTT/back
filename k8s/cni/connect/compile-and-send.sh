#!/bin/bash

set -euo pipefail
#環境変数の設定
WORKER1_IPADDR=192.168.0.210
WORKER2_IPADDR=192.168.0.218
WORKER3_IPADDR=192.168.0.222

echo connect.goのコンパイル
GOOS=linux GOARCH=amd64 go build ./connect.go

echo Worker1にファイル転送
scp ./connect ubuntu@${WORKER1_IPADDR}:~/
ssh ubuntu@${WORKER1_IPADDR} sudo cp ./connect /opt/cni/bin/

echo Worker2にファイル転送
scp ./connect ubuntu@${WORKER2_IPADDR}:~/
ssh ubuntu@${WORKER2_IPADDR} sudo cp ./connect /opt/cni/bin/

echo Worker3にファイル転送
scp ./connect ubuntu@${WORKER3_IPADDR}:~/
ssh ubuntu@${WORKER3_IPADDR} sudo cp ./connect /opt/cni/bin/