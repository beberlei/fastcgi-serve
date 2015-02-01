# HHVM Webserver

This is a dedicated webserver to serve PHP applications through HHVMs FastCGI
interface, allowing to run (and test) PHP applications with Zend PHP and HHVM
side by side.

Why? I couldn't get hhvm and PHP to run side by side with Apache on my Ubuntu system.
There is probably a simple solution out there to do what this project does, but
in the interest of furthering my #golang experience I thought this was a good project
to built.

**WARNING: This server is not meant for production. Use it to test your web-apps with HHVM.**

## Installation

You need go installed on your system, then call:

    $ go get github.com/beberlei/hhvm-serve

Add `$GOPATH/bin` to your `$PATH` in `~/.bashrc` to make the `hhvm-serve` command available.

You need HHVM on your system as well and listening for fastcgi connections on `127.0.0.1:9000` (it does that by default).
Don't forget to `/etc/init.d/hhvm start`.

## Usage

You can keep your Apache/Nginx setup with this proxy and start a webserver (defaults to `localhost:8080`)
for a given document root (defaults to current working directory).

    $ hhvm-serve --document-root=/var/www --listen=127.0.0.1:8080
    Listening on http://localhost:8080
    Document root is /var/www
    Press Ctrl-C to quit.

## Acknowledgements

- Bundles fastcgi client library by Junqing Tan (ivan@mysqlab.net) as dependency
- http://code.google.com/p/go-fastcgi-client/
