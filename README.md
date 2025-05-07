# wsserver

send to client
```
curl -i -X POST http://192.168.110.20:80/api/send -d '{"type":"echo","client_id":["24c0c0793df0fe8fb2bbddf2d8ad59cf","ec2338f5ae989a1573b337c540df6d6e"],"data":{"seq":100}}'
```

ws.html
```html
<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Document</title>
</head>

<body>
<button id="linkBtn">连接ws</button>
<hr>
<button id="breakBtn">断开</button>

<script>
  let linkBtn = document.getElementById('linkBtn');
  let breakBtn = document.getElementById('breakBtn');
  let ws;
  let connected;

  linkBtn.addEventListener('click', function () {
      if (connected){
          return;
      }

      ws = new WebSocket('ws://192.168.110.20/api/ws');

      //监听是否连接成功
      ws.onopen = () => {
          connected = true;
          console.log('ws连接状态:' + ws.readyState);
      };

      // 接听服务器发回的信息并处理展示
      ws.onmessage = (e) => {
          console.log(e);
      };

      // 监听连接关闭事件
      ws.onclose = function () {
          connected = false;
          // 监听整个过程中websocket的状态
          console.log('ws连接状态：' + ws.readyState);
      };

      // 监听并处理error事件
      ws.onerror = function (error) {
          connected = false;
          console.log(error);
      };
  });

  breakBtn.addEventListener('click', function () {
      if (ws) {
          ws.close();
      }
  });
</script>
</body>

</html>
```
