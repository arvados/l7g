Lightning Proxy
---

To run:

```
nodejs prox.js
```

This is a small application to allow for rudimentary authentication (via a passphrase)
and proxying of requests to the different services.

This is mostly a one-off and isn't meant to be used for anything more than demo-ing the
prototype in place.

Look in the `config.json` for the passphrase to use and some of the various proxy options.
The default authentication phrase is the string `ok`.

In `proxy.js`, you can fiddle with some of the static routes and files.

`srv.js` is a simple test server to test the proxy.

