const http = require('http');

const requestListener = function (req, res) {
  if (req.url == '/v1/health') {
    res.writeHead(200);
    res.end('pong!ping!')
    return
  }
  res.writeHead(200);
  res.end('Hello, World! From FaaS-Native');
}

const server = http.createServer(requestListener);
server.listen(8000);
