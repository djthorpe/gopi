# Networking

This section describes the networking features, which allow communication over 
[HTTP](https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol)
and RPC ([Remote procedure call](https://en.wikipedia.org/wiki/Remote_procedure_call))
protocols. You can also advertise and discover
[Services](https://en.wikipedia.org/wiki/Service_discovery) on your local network.

## Overview

These are the units you can embed into your application:

  * `gopi.Server` A HTTP or RPC server which can serve requests;
  * `gopi.ConnPool` A pool of connections to remote servers;
  * `gopi.ServiceDiscovery` A mechanism to either discovery available network services or register services;
  * `gopi.PingService` An RPC service which responds to requests with an empty response;
  * `gopi.InputService` As RPC service which emits input events (key presses, etc.);
  * `gopi.HttpStatic` A HTTP service which serves any file or folder on the filesystem.

These are examples you can look at which demonstate the features:

  * (`rpcping`)[https://github.com/djthorpe/gopi/tree/master/cmd/rpcping] is both a RPC client and server, 
    demonstrating the ping service.
  * (`hellohttp`)[https://github.com/djthorpe/gopi/tree/master/cmd/hellohttp] is a HTTP server which can
    serve static files.

