# sharpshooter-tunnel
sharpshooter protocol tunel


## Introduce
sharpshooter-tunnel is a base on [sharpshooter](https://github.com/soyum2222/sharpshooter)  tunnel .

convert sharpshooter proto to TCP .

client to server conversation use sharpshooter base on UDP .

## build
1. install go
2. install make
3. run make in sharpshooter-tunnel dir

## Usage

    ./sharpshooter-client-linux-amd64 genconf 
    
get config module file 

    ./sharpshooter-client-linux-amd64 -c config.json

use config.json as the configuration file


## Config

##### client config:
```json
{
	"remote_addr": "",
	"local_addr": "",
	"key": "sharpshooter",
	"con_num": 10,
	"send_win": 1024,
	"rec_win": 1024,
	"interval": 100,
	"debug": false,
	"p_port": 8888
}
``` 
    remote_addr : sharpshooter tunnel server ip address eg:192.168.1.2:8888
    local_addr  : sharpshooter tunnel client listen ip addr eg: 0.0.0.0:8080
    key         : conversation aes encrypt key 
    con_num     : multiplexing create max conn
    send_win    : sharpshooter send windows minni size
    rec_win     : sharpshooter receive windows minni size
    interval    : sharpshooter send data interval , utils ms
    debug       : open debug
    p_port      : pprof port
    

#### server config:
```json
{
        "local_addr": "",
        "key": "sharpshooter",
        "send_win": 1024,
        "rec_win": 1024,
        "interval": 100,
        "listen_port": 0,
        "debug": false,
        "p_port": 8888
}
```
    local_addr  : be proxied service IP address
    listen_port : listen UDP port 
    key         : conversation aes encrypt key 
    send_win    : sharpshooter send windows minni size
    rec_win     : sharpshooter receive windows minni size
    interval    : sharpshooter send data interval , utils ms
    debug       : open debug
    p_port      : pprof port
 