# FastCGI-Serve - A webserver for FastCGI

This is a dedicated webserver to serve applications over HTTP using FastCGI
interface, allowing for example to run (and test) PHP applications with
different Zend PHP and HHVM side by side.

Why? I couldn't get hhvm and PHP to run side by side with Apache on my Ubuntu system.
There is probably a simple solution out there to do what this project does, but
in the interest of furthering my #golang experience I thought this was a good project
to built.

**WARNING: This server is not meant for production.**

## Installation

You need go installed on your system, then call:

    $ go get github.com/beberlei/fastcgi-serve

Add `$GOPATH/bin` to your `$PATH` in `~/.bashrc` to make the `fastcgi-serve` command available.

You need PHP-FPM or HHVM on your system as well and listening for fastcgi
connections on `127.0.0.1:9000` (it does that by default).

Don't forget to run either `/etc/init.d/php-fpm start` or `/etc/init.d/hhvm start`.

## Usage

You can keep your Apache/Nginx setup with this proxy and start a webserver (defaults to `localhost:8080`)
for a given document root (defaults to current working directory).

    $ fastcgi-serve --document-root=/var/www --listen=127.0.0.1:8080
    Listening on http://localhost:8080
    Document root is /var/www
    Press Ctrl-C to quit.

## Configuration

The following settings are available:

- `--document-root` - The document root to serve files from (default: current working directory)
- `--listen` - The webserver bind address to listen to (default:127.0.0.1)
- `--server` - The FastCGI server to listen to
- `--server-port` The FastCGI server port to listen to
- `--index` The default script to call when request path cannot be served with an existing file

There also support to load additional environment variables into each request.
Create an `.env` file in the document root with key vaue pairs of options:

    FOO=BAR

## Acknowledgements

- Bundles fastcgi client library by Junqing Tan (ivan@mysqlab.net) as dependency
- http://code.google.com/p/go-fastcgi-client/
