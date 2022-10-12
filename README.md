# wsserver
A goroutine-style WebSocket server based on `github.com/gorilla/websocket`, supports: Listener, Heart, Group...

# Install
```
go install github.com/go-the-way/wsserver@latest
```

## API Docs

### 1. Online client
```
curl $SERVER_ADDR/api/online_client
```

## Rpc Docs

### 1. Send to client
```
ServicePath: Sender
ServiceMethod: Send
Args: {"type":"echo","client_id":["client_a"],"data":{"seq":1000}}
Reply: {"code":200}
```

## Listener Docs

* Creator `when a new client connected, trigger creator listener` 
 
* Destroyer `when cached client closed, trigger destroyer listener`

## Code Styles
```
config                 -- App & Environment
handler                -- Handler routers
listener               -- Listeners
manager                -- Client manager
pkg                    -- Third-party pkg
rpc                    -- Rpc service
```

## Environment

### 1. SERVER_ADDR
*Http Server Address*
```
default val: :80
```

### 2. RPC_ADDR
*Rpc Server Address*
```
default val: :86
```

# Example

```
let ws = new WebSocket("ws://192.168.6.125:80/ws");
let seq = 1;
let INT;
ws.onopen = function () {
  console.log("已连接");
  INT = setInterval(function () {
    ws.send('{"type":"seq","data":{"seq":' + seq++ + "}}");
  }, 1000);
};
ws.onmessage = function (msg) {
  console.log("接收=>", msg.data);
};
ws.onclose = function () {
  console.log("已断开");
  clearInterval(INT);
};
```


