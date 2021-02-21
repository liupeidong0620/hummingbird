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
$ sudo ./hummingbird-linux-amd64 -interface en0 -module config

# help
$ ./hummingbird-linux-amd64 -h

```

### config interface & route
 
 > linux.sh would take care of tun & routes.

```shell script
$ sh linux.sh start
```

  </details>

 <details>
    <summary><b>With MacOS</b></summary>

### start hummingbird

```sh
$ sudo ./hummingbird-darwin-amd64 -interface en0 -module config

# help
$ ./hummingbird-darwin-amd64 -h

```

### config interface & route
 
 > darwin.sh would take care of tun & routes.

```shell script
$ sh darwin.sh start
```
  </details>

## server example

[hummingbird-server](https://github.com/liupeidong0620/hummingbird-server.git).

## TODO
