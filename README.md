# SocksV
A proxy server that supports socks5,socksv,vpn and so on...

## Motivation

This project is created because of the blocking of internet in some areas,
it's very hard for them to connect the target website fastly.For example,`Github`
is limited for speed in some ares...

We build this project to bypass the blocking rule, to create a fast and clean internet.

## How It Works
The core of socksv is to run a proxy server outside the blocking area, which is
connnected from inside.The network flow is:

```text
 -----------        --------------          ------------         -------------
|Web Browser| <---> |Local Machine| <----> |Proxy Server| <---> |Target Website|
 -----------        --------------         -------------         --------------
```


## How To Run
Run command `socksv -h` to get help:
```bash
Usage of socksv:
  -P string
        proxy server port. (default "8080")
  -l int
        log level.0-info;1-debug;2-trace;3-warn;4-error. (default 1)
  -p string
        socks5 server port. (default "1080")
  -s string
        relay server to connect.
```

### 1. run as server

Runs as server on your proxy machine(like aws ec2) and listen at port 1081
```bash
socksv -P 8080
```

### 2. run as client

Runs as client on your local machine.

The following command will run socksv as client,  connect to proxy server at [proxy_ip]:8080, and listening socks5 stream at port 1080.

```bash
socksv -s [proxy_ip]:8080
```

> The communication between socksv proxy server and client is encrypted.

### 3.config chrome

Config your chrome to connect to your local socks client.

Here you need a chrome plugin  [SwitchyOmega](https://chrome.google.com/webstore/detail/proxy-switchyomega/padekgcemlokbadohgkifijomclgjgif?utm_source=chrome-ntp-icon)  and config the proxy as:
* `protocol`: socks5
* `proxy server`: 127.0.0.1 (the socksv client running in your local machine)
* `port`: 1080

> After the 3 steps configuration, your can visit whatever you want

## Plan TODO

`TODO` List:

 * udp support for socks5
 * chrome plugin client
 * ios client
 * android client
 * electron desktop

## Contributing

 If you are interested in this project and wanna contribute, please fork this project,
  modify, and submit a pull request.

 You are welcome to submit issues to help improve the code and experience.

 You can choose what to do in the `TODO` list.