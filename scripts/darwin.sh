#!/bin/sh

TUN_IP=198.18.0.1
TUN_MASK=255.255.255.0
TUN_GATEWAY=198.18.0.1
TUN_NAME=utun123

tun_up(){
    echo "tun_up() run ..."
    sudo ifconfig $TUN_NAME $TUN_IP netmask $TUN_MASK $TUN_GATEWAY up
}

route_add() {
    echo "route_add() run ..."
    sudo route -n add -net 128.0.0.0 -netmask 128.0.0.0 $TUN_GATEWAY
    sudo route -n add -net 64.0.0.0 -netmask 192.0.0.0 $TUN_GATEWAY
    sudo route -n add -net 32.0.0.0 -netmask 224.0.0.0 $TUN_GATEWAY
    sudo route -n add -net 16.0.0.0 -netmask 240.0.0.0 $TUN_GATEWAY
    sudo route -n add -net 8.0.0.0 -netmask 248.0.0.0 $TUN_GATEWAY
    sudo route -n add -net 4.0.0.0 -netmask 252.0.0.0 $TUN_GATEWAY
    sudo route -n add -net 2.0.0.0 -netmask 254.0.0.0 $TUN_GATEWAY
    sudo route -n add -net 1.0.0.0 -netmask 255.0.0.0 $TUN_GATEWAY
}

route_del(){
    echo "route_del() run ..."
    sudo route -n delete -net 128.0.0.0 -netmask 128.0.0.0 $TUN_GATEWAY
    sudo route -n delete -net 64.0.0.0 -netmask 192.0.0.0 $TUN_GATEWAY
    sudo route -n delete -net 32.0.0.0 -netmask 224.0.0.0 $TUN_GATEWAY
    sudo route -n delete -net 16.0.0.0 -netmask 240.0.0.0 $TUN_GATEWAY
    sudo route -n delete -net 8.0.0.0 -netmask 248.0.0.0 $TUN_GATEWAY
    sudo route -n delete -net 4.0.0.0 -netmask 252.0.0.0 $TUN_GATEWAY
    sudo route -n delete -net 2.0.0.0 -netmask 254.0.0.0 $TUN_GATEWAY
    sudo route -n delete -net 1.0.0.0 -netmask 255.0.0.0 $TUN_GATEWAY
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
