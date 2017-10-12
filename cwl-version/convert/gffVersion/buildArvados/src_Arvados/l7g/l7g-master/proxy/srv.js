/* 
 * http://stackoverflow.com/questions/3393854/get-and-set-a-single-cookie-with-node-js-http-server
 * Corey Hart
 */
var http = require('http');

function parseCookies (request) {
    var list = {},
        rc = request.headers.cookie;

    rc && rc.split(';').forEach(function( cookie ) {
        var parts = cookie.split('=');
        //list[parts.shift().trim()] = decodeURI(parts.join('='));
        list[parts[0].trim()] = decodeURI(parts.slice(1).join('='))
    });

    return list;
}


http.createServer(function (request, response) {

  // To Read a Cookie
  var cookies = parseCookies(request);

  console.log(".");
  for (var x in cookies) {
    console.log(">>>", x, cookies[x]);
  }

  // To Write a Cookie
  response.writeHead(200, {
    'Set-Cookie': 'mycookie=test',
    'Set-Cookie': 'mycookie2=test=foo=bar',
    'Content-Type': 'text/plain'
  });
  response.end('Hello World\n');
}).listen(8124);


http.createServer(function(req, resp) {
  console.log(">> 8081", req.url, req.method);
  resp.writeHead(200, { 'Content-Type': 'text/plain' });
  resp.end('8081\n');
}).listen(8081);

http.createServer(function(req, resp) {
  console.log(">> 8082", req.url, req.method);
  resp.writeHead(200, { 'Content-Type': 'text/plain' });
  resp.end('8082\n');
}).listen(8082);

http.createServer(function(req, resp) {
  console.log(">> 8083", req.url, req.method);
  resp.writeHead(200, { 'Content-Type': 'text/plain' });
  resp.end('8083\n');
}).listen(8083);

http.createServer(function(req, resp) {
  console.log(">> 8084", req.url, req.method);
  resp.writeHead(200, { 'Content-Type': 'text/plain' });
  resp.end('8084\n');
}).listen(8084);

http.createServer(function(req, resp) {
  console.log(">> 8085", req.url, req.method);
  resp.writeHead(200, { 'Content-Type': 'text/plain' });
  resp.end('8085\n');
}).listen(8085);

console.log('Server running');
