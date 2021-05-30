# hummingbird

A hummingbird(tun2websocket) powered by gVisor TCP/IP stack

## How to Build

### build from source code

Go compiler version >= 1.15 is required

```text
$ git clone https://github.com/liupeidong0620/hummingbird.git
$ cd hummingbird
$ make
```

## QuickStart

 <details>
    <summary><b>With Linux</b></summary>

### start hummingbird

```sh
$ sudo ./hummingbird-linux-amd64 -interface en0 -proxy "ws://1.2.3.4:80"

# help
$ ./hummingbird-linux-amd64 -h

```

### config interface & route
 
 > scripts/linux.sh would take care of tun & routes.

```shell script
$ sh linux.sh start
```

  </details>

 <details>
    <summary><b>With MacOS</b></summary>

### start hummingbird

```sh
$ sudo ./hummingbird-darwin-amd64 -interface en0 -proxy "ws://1.2.3.4:80"

# help
$ ./hummingbird-darwin-amd64 -h

```

### config interface & route
 
 > scripts/darwin.sh would take care of tun & routes.

```shell script
$ sh darwin.sh start
```
  </details>

   <details>
    <summary><b>With Windows</b></summary>

### start hummingbird

> This runs on Windows, but you should install [wintun](https://www.wintun.net/)

```sh
# root authority
$ ./hummingbird-windows-amd64 -interface en0 -proxy "ws://1.2.3.4:80"

# help
$ ./hummingbird-windows-amd64 -h

```

### config interface & route

```shell script
netsh interface ip set address utun123 static 26.26.26.1 255.255.255.0

netsh interface ip set dns utun123 static 8.8.8.8

route add 0.0.0.0 MASK 128.0.0.0  26.26.26.1
```
  </details>

## server example

[hummingbird-server](https://github.com/liupeidong0620/hummingbird-server.git).

## TODO

* IPV6 test
