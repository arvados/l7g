var http = require('http'),
    httpProxy = require('http-proxy'),
    url = require("url"),
    cookie = require("cookie"),
    program = require("commander"),
    fs = require("fs");

// load config file

/*
program
  .arguments("<file>")
  .option("-c, --config <config>", "Config file")
  .action(function(file) {
    init();
  });
  */

console.log(process.argv);

var args = process.argv.slice(2);
var cfg_fn = "./config.json";
if (args.length > 0) {
  cfg_fn = args[0];
}

var CFG = JSON.parse(fs.readFileSync(cfg_fn));
console.log(CFG);

function parseCookies (request) {
  var list = {},
      rc = request.headers.cookie;

  rc && rc.split(';').forEach(function( cke ) {
    var parts = cke.split('=');
    //list[parts.shift().trim()] = decodeURI(parts.join('='));
    list[parts[0].trim()] = decodeURI(parts.slice(1).join('='))
  });

  return list;
}


function serializeCookie(cke) {
  var cookie_list = [];
  for (var x in cke) {
    cookie_list.push(cookie.serialize(x, cke[x]));
  }
  //return cookie_list.join("; ");
  return cookie_list;
}

//
// Create a proxy server with custom application logic
//
//var proxy = httpProxy.createProxyServer({});
var proxy = httpProxy.createProxyServer({ prependPath:true });

function srv_redirect(res, ckeh) {
  ckeh = ((typeof ckeh === "undefined") ? {} : ckeh);
  var t = serializeCookie(ckeh);
  if (t.length>0) {
    res.writeHead(302, {
      'Set-Cookie': serializeCookie(ckeh),
      "Location": "/"
    });
  } else {
    res.writeHead(301, {"Location":"/"});
  }
  res.end();
}

function srv_file(fn, res, ckeh) {
  ckeh = ((typeof ckeh === "undefined") ? {} : ckeh);

  fs.readFile(fn, "binary", function(err, file) {
    if(err) {
      res.writeHead(500, {"Content-Type": "text/plain"});
      res.write(err + "\n");
      res.end();
      return;
    }

    var t = serializeCookie(ckeh);
    //console.log("srv_file: ", t, t.length);

    if (t.length>0) {
      res.writeHead(200, {
        'Set-Cookie': serializeCookie(ckeh)
      });
    } else {
      res.writeHead(200);
    }
    res.write(file, "binary");
    res.end();
  });

}

var g_staticFiles = {
  "/favicon.ico" : "./curoverse.ico",
  "/css/style.css" : "./assets/css/style.css",
  "/" : "./assets/login.html",
  "/index" : "./assets/login.html",
  "/index.html" : "./assets/login.html",
  "/login" : "./assets/login.html"
};

var g_staticFilesAuth = {
  "" : "./assets/info.html",
  "/" : "./assets/info.html",
  "/info" : "./assets/info.html",
  "/info.html" : "./assets/info.html"
}

/*
var g_proxySites = {
  "/" : "http://localhost:8081/i",
  "/variant" : "http://localhost:8082/i",
  "/cgf" : "http://localhost:8083/i",
  "/tile" : "http://localhost:8084/i",
  "/phenotype" : "http://localhost:8085/i",
};
*/

//
// Create your custom server and just call `proxy.web()` to proxy
// a web request to the target passed in the options
// also you can use `proxy.ws()` to proxy a websockets request
//
var server = http.createServer(function(req, res) {

  console.log("\n\nprox>>>", req.url, req.method);

  var auth_var_name = "auth_code";
  var auth_val = CFG.auth_key;

  //console.log("???", req);
  //console.log("???", res);


  var ckeh = {};
  if (("headers" in req) &&
      ("cookie" in req.headers) &&
      (req.headers.cookie.length > 0)) {
    ckeh = cookie.parse(req.headers.cookie);
  }

  var authenticated = false;
  if ((auth_var_name in ckeh) && (ckeh[auth_var_name] === auth_val)) {
    authenticated = true;
  }

  if (req.method === "GET") {
    if (req.url in g_staticFiles) {

      console.log("simple get");


      if (authenticated && (req.url in g_staticFilesAuth)) {
        console.log("  skipping simple get, passing down...");
      } else {
        srv_file(g_staticFiles[req.url],res,ckeh);
        return;
      }

    }
  }


  if (authenticated) {
  //if ((auth_var_name in ckeh) && (ckeh[auth_var_name] === auth_val)) {
    console.log(">>> authorized", req.url);

    if (req.url in g_staticFilesAuth) {

      srv_file(g_staticFilesAuth[req.url], res, ckeh);
      return;
    }

    //sendreq(req, res);
    //proxy.web(req, res, { target: 'http://127.0.0.1:8124' });
    //proxy.web(req, res, { target: CFG.proxy });

    var route_loookup = false;
    var url_parts = req.url.split("/");
    if ((url_parts.length>0) && (url_parts[0] === "")) {
      url_parts = url_parts.slice(1);
    }

    if (url_parts.length>0) {
      route_lookup = true;
    }


    //if (req.url in CFG.proxy) {
    if (route_lookup && (url_parts[0] in CFG.proxy)) {
      //console.log("???");
      //console.log("   ", req, res);

      console.log("   .. ", req.url);
      console.log("   ..>>>", url_parts);
      console.log("   .. ", CFG.proxy[req.url]);
      console.log("   .. ", CFG.proxy[url_parts[0]]);

      var dst_url = CFG.proxy[url_parts[0]] + "/" + url_parts.slice(1).join("/");

      console.log("   ..dst >>", dst_url);

      //proxy.web(req, res, { target: CFG.proxy[req.url], prependPath:true });
      //proxy.web(req, res, { target: dst_url, prependPath:true });
      proxy.web(req, res, { target: dst_url, prependPath:true, ignorePath:true });
    } else {
      console.log("404>>>");
      srv_file("./assets/404.html",res,ckeh);
    }

    return;

  } else if (!(auth_var_name in ckeh) || (ckeh[auth_var_name].length==0)) {

    //console.log("clearing");
    //ckeh["message"] = " !";
  }

  var ourl = url.parse(req.url);
  var path_parts = ourl.pathname.replace(/^\/*/, '').split('/');

  /*
  for (var k=0; k<path_parts.length; k++) {
    console.log("  prox.path[", k, "]>", path_parts[k]);
  }
  */

  //console.log("not auth");

  ckeh["message"] = "Not authorized...";
  //srv_file("/", res, ckeh);
  srv_redirect(res, ckeh);

});

//console.log("listening on port 5050")
//server.listen(5050);

console.log("listening on port " + CFG.port )
server.listen(CFG.port);
