# HHVM Webserver

This webserver is an on-demand proxy to HHVMs FastCGI interface, allowing
to run (and test) PHP applications with Zend PHP and HHVM side by side.

Setting up HHVM and PHP beside each other in the same Apache or Nginx setup
is a cumbersome, the HHVM docs even suggest to disable modphp5 and enable FastCGI
to use HHVM.

You can keep your Apache/Nginx setup with this proxy and start a webserver (defaults to `localhost:8080`)
for a given document root (defaults to `/var/www`).

    $ hhvm-serve --document-root=/var/www --listen=127.0.0.1:8080
    Listening on http://localhost:8080
    Document root is /var/www
    Press Ctrl-C to quit.

**WARNING: This server is not meant for production. Use it to test your web-apps with HHVM.**

Acknowledgements
- Bundles fastcgi client library by Junqing Tan (ivan@mysqlab.net) as dependency
- http://code.google.com/p/go-fastcgi-client/
