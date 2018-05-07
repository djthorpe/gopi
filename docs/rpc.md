
## Microservice Architecture

Microservices allow you to implement your system as a number of discrete services,
and allow those services to reside anywhere on your network and co-operate. There
are probably plenty of books and articles about these available but I will define
them as follows:

  * A __service__ is a unit which provides "remote procedure calls": basically,
    procedures and functions which can be called and return values across your
    network;
  * A __server__ accepts remote procedure calls and routes them to one or more 
    hosted services;
  * A __client__ is some code which can call these services across the network;
  * The __communication protocol__ defines how the information is transferred
    between service and client;
  * A __naming service__ allows servers to register the services they host and
    allows clients to determine where to access the services.

The gopi framework implements services and clients as modules. To create
microservice you would:

  * Create a service as module type `gopi.MODULE_TYPE_SERVICE`
  * In your bootstrap (main) code, create a configuration including all the
    services you want to serve
  * Call `RPCServerTool` in order to run your server

Here is an example "helloworld" server which serves both the helloworld
service and the metrics service:

```
package main

import (
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc/grpc"
	_ "github.com/djthorpe/gopi/sys/rpc/mdns"

	// RPC Services
	_ "github.com/djthorpe/gopi/rpc/grpc/helloworld"
	_ "github.com/djthorpe/gopi/rpc/grpc/metrics"
)

func main() {
	config := gopi.NewAppConfig("rpc/service/helloworld:grpc", "rpc/service/metrics:grpc")
	config.Service = "helloworld"
	os.Exit(gopi.RPCServerTool(config))
}
```

Some notes on this server implementation:

  * The communication protocol used is [gRPC](https://grpc.io/), but it would be possible 
    to implement a different protocol, such as JSON or [Twerp](https://twitchtv.github.io/twirp/)
	for example;
  * The naming service used is [mDNS](https://tools.ietf.org/html/rfc6762) but it would be
    possible to use a different naming service, for example [Consul](https://www.consul.io/) seems
	like a good choice;
  * The RPC services are imported anonymously and then the specific services are included in the
    configuration;
  * The service is announced as "helloworld" which allows clients to discover and connect to
    the service.

More information on how to develop your own services will be provided later. Of course, you'll
want to communicate with your server using a client, which can either be another service, a
desktop or mobile app, or command line tool, for example. The gopi framework provides an
implementation for your clients as follows:

  * A __Naming Service__ watches for new announcements of services being registered on the
    local network;
  * The __Client Pool__ collects these announcements and allows connections to be made to
    remote services. It also sends out messages when services are registered and de-registered;
  * A __Client Connection__ is provided by the pool, allows information about remote services
    to be queried
  * The client pool can create a client for which remote procedures can be called.

Here's an example command line tool which connects to the ```helloworld`` service:

```
func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	pool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool)
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	if records, err := pool.Lookup(ctx, "", "", 1); err != nil {
		return err
	} else if len(records) == 0 {
		return gopi.ErrDeadlineExceeded
	} else if conn, err := pool.Connect(records[0], 0); err != nil {
		return err
	} else if client := pool.NewClient("mutablelogic.Helloworld", conn); client == nil {
		return gopi.ErrAppError
	} else if message, err := client.(*hw.Client).SayHello(); err != nil {
		return err
	} else {
		fmt.Printf("%v says '%v'\n\n", conn.Name(), message)
		pool.Disconnect(conn)
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/client/helloworld:grpc")
	config.Service = "helloworld"
	os.Exit(gopi.CommandLineTool(config, Main))
}
```

Some notes on this client implementation:

  * The `helloworld` client is imported anonymously and included in the configuration;
  * The service is also named as `helloworld` before invoking the application as
    a command line tool.
  * The `Lookup` function for the client pool returns exactly one registered service
 	on the network, a timeout of 100ms is used to wait for the information to be gathered,
    and the method returns `gopi.ErrDeadlineExceeded` if no services are available
	on the network.
  * A connection is made to the service, and a new client is created of name 
    `mutablelogic.Helloworld` - it's returned as a generic `gopi.RPCClient` interface
    which needs to be cast to the actual client interface.
  * The `SayHello` method is called remotely seralizing and deserializing through
    protocol buffers.

More information on developing your own clients will be provided later.

