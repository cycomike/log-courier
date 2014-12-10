# Logstash Integration

Log Courier is built to work seamlessly with [Logstash](http://logstash.net)
1.4.x.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](http://doctoc.herokuapp.com/)*

- [Installation](#installation)
  - [Logstash 1.5+ Plugin Manager](#logstash-15-plugin-manager)
  - [Manual installation](#manual-installation)
  - [Local-only Installation](#local-only-installation)
- [Configuration](#configuration)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation

### Logstash 1.5+ Plugin Manager

Logstash 1.5 introduces a new plugin manager that makes installing additional
plugins extremely easy.

Simply run the following commands to install the latest stable version of the
Log Courier plugins. If you are only receiving events, you only need to install
the input plugin.

		cd /path/to/logstash
		bin/logstash plugin install logstash-input-log-courier
		bin/logstash plugin install logstash-output-log-courier

Once the installation is complete, you can start using the plugins!

### Manual installation

For Logstash 1.4.x the plugins and dependencies need to be installed manually.

First build the Log Courier gem the plugins require. The file you will need will
be called log-courier-X.X.gem, where X.X is the version of Log Courier you have.

		git clone https://github.com/driskell/log-courier
		cd log-courier
		make gem

Switch to the Logstash installation directory and install it. Note that because
this is JRuby it may take a minute to finish the install. The ffi-rzmq-core and
ffi-rzmq gems bundled with Logstash will be upgraded during the installation,
which will require an internet connection.

		cd /path/to/logstash
		export GEM_HOME=vendor/bundle/jruby/1.9
		java -jar vendor/jar/jruby-complete-1.7.11.jar -S gem install /path/to/log-courier-X.X.gem

The remaining step is to manually install the Logstash plugins.

		cd /path/to/log-courier
		cp -rvf lib/logstash /path/to/logstash/lib

### Local-only Installation

If you need to install the gem and plugins on a server without an internet
connection, you can download the latest ffi-rzmq-core and ffi-zmq gems from the
rubygems site, transfer them across, and install them yourself. Once they are
installed, follow the instructions for Manual Installation and the process can
be completed without an internet connection.

* https://rubygems.org/gems/ffi-rzmq-core
* https://rubygems.org/gems/ffi-rzmq

## Configuration

The 'courier' input and output plugins will now be available. An example
configuration for the input plugin follows.

		input {
				courier {
						port            => 12345
						ssl_certificate => "/opt/logstash/ssl/logstash.cer"
						ssl_key         => "/opt/logstash/ssl/logstash.key"
				}
		}

The following options are available for the input plugin:

* transport - "tcp", "tls", "plainzmq" or "zmq" (default: "tls")
* address - Interface address to listen on (defaults to all interfaces)
* port - The port number to listen on
* ssl_certificate - Path to server SSL certificate (tls)
* ssl_key - Path to server SSL private key (tls)
* ssl_key_passphrase - Password for ssl_key (tls, optional)
* ssl_verify - If true, verifies client certificates (tls, default false)
* ssl_verify_default_ca - Accept client certificates signed by systems root CAs
(tls)
* ssl_verify_ca - Path to an SSL CA certificate to use for client certificate
verification (tls)
* curve_secret_key - CurveZMQ secret key for the server (zmq)
* max_packet_size - The maximum packet size to accept (default 10485760,
corresponds to Log Courier's `"spool max bytes"`)

The following options are available for the output plugin:

* addresses - Address to connect to in array format (only the first address will
be used at the moment)
* port - Port to connect to
* ssl_ca - Path to SSL certificate to verify server certificate
* ssl_certificate - Path to client SSL certificate (optional)
* ssl_key - Path to client SSL private key (optional)
* ssl_key_passphrase - Password for ssl_key (optional)
* spool_size - Maximum number of events to spool before a flush is forced
(default 1024)
* idle_timeout - Maxmimum time in seconds to wait for a full spool before
flushing anyway (default 5)

NOTE: The tcp, plainzmq and zmq transports are not implemented in the output
plugin at this time. It supports only the tls transport.
