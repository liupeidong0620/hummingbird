#!/bin/sh

TUN_IP=198.18.0.1
TUN_MASK=15
TUN_NAME=utun123

PROXY_ADDR=$2
DEFAULT_GW=$3

tun_up(){
    echo "tun_up() run ..."
    sudo ip addr add $TUN_IP/$TUN_MASK dev $TUN_NAME
    sudo ip link set dev $TUN_NAME up
}

route_add() {
    echo "route_add() run ..."
    #echo "gw: $DEFAULT_GW proxy_addr: $PROXY_ADDR"
    sudo ip route add $PROXY_ADDR via $DEFAULT_GW
    sudo ip route add 0.0.0.0/1 dev $TUN_NAME src $TUN_IP
    sudo ip route add 128.0.0.0/1 dev $TUN_NAME src $TUN_IP
}

route_del(){
    echo "route_del() run ..."
    sudo ip route del 0.0.0.0/1 dev $TUN_NAME src $TUN_IP
    sudo ip route del 128.0.0.0/1 dev $TUN_NAME src $TUN_IP
}

main(){
    echo "$1"
    if [ "$1" = "start" ];then
        echo "start tun and load route..."
        tun_up
        route_add
    elif [ "$1" = "stop" ];then
        echo "del route ..."
        route_del
    fi
}

main $1
