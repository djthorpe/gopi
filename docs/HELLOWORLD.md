
## Getting Started with __gopi__

You can use the __gopi__ framework to develop a number of
different kinds of applications:

  * Command-line applications which don't interact with the user,
    or interact in a simple way;
  * Micro-services which interact with other services or user
    interfaces through the network;
  * Applications which respond to events (from input devices, timers,
    network activity, etc);
  * Graphical user applications which provide rich user interfaces
    and respond to user interactions.

The plan for __gopi__ is to support these on - eventually - a variety
of hardware platforms which can make the best use of the target
environment.

The prototypical `helloworld` application demonstrates how to develop
a simple command-line application:

```
package main

import (
	"fmt"
	"os"

  // Import Frameworks
	gopi "github.com/djthorpe/gopi"

	// Import modules
	_ "github.com/djthorpe/gopi/sys/logger"
)

func helloWorld(app *gopi.AppInstance, done chan<- struct{}) error {
	// If -name argument is used then use that, else output generic message
	if name, exists := app.AppFlags.GetString("name"); exists {
		fmt.Println("Hello,", name)
	} else {
		fmt.Println("Hello, World (use -name flag to specify your name)")
	}

	// If wait flag is set, then wait until CTRL+C is pressed to continue
	if wait, _ := app.AppFlags.GetBool("wait"); wait {
		fmt.Println("Press CTRL+C to exit")
		app.WaitForSignal()
	}

	// Signal that main thread is done
	done <- gopi.DONE
	return nil
}

func main() {
	// Create the configuration
	config := gopi.NewAppConfig()
	config.AppFlags.FlagString("name", "", "Your name")
	config.AppFlags.FlagBool("wait", false, "Wait for CTRL+C interrupt to end")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, helloWorld))
}
```

Here the __gopi__ framework is imported alongside an unnamed _module_ with
the name `sys/logger` which provides simple log output. Every __gopi__ application
requires a logger module.

If you run this command-line application you'll get the usual "Hello, World" output
from the `helloWorld` function, but you also get a few additional features for free.
Try running it with the `-help` flag:

```
bash$ cd "${GOPATH}/src/github.com/djthorpe/gopi"
bash$ go install cmd/helloworld/helloworld.go
bash$ helloworld -help

Usage of helloworld:
  -debug
    	Set debugging mode
  -log.append
    	When writing log to file, append output to end of file
  -log.file string
    	File for logging (default: log to stderr)
  -name string
    	Your name
  -verbose
    	Verbose logging
  -wait
    	Wait for CTRL+C interrupt to end
```

Try running the application with various flags including `-name` and `-wait`. Looking
at the source code in more detail, apart from the framework and module import statements
there are two functions:

  1. The `main` function creates an application configuration and adds two configuration
    flags. It then runs the application as a command line application with a single
    foreground task, `helloWorld`.
  2. The `helloWorld` function is a "foreground task" which when finished will terminate
    the program. It reads the configuration flags, outputs some information. If the
    `-wait` flag is used then it waits until the task is interrupted by keyboard (pressing
    CTRL and C keys) or interrupt signal (using the `kill` command). The `done <- gopi.DONE`
    line is not strictly necessary when there's only one foreground task; usually this will
    signal to background tasks that they should terminate.

As you can see with this example, there's several areas to explore:

  * The application configuration, including configuration flags
  * Running foreground and background tasks and co-ordinating between them
  * Importing and using application modules

The following sections introduce these concepts.

## Configuration

Any application starts with a _configuration_ for that application, which can determine
the environment in which the application runs. It's intended that the configuration may
include command-line flags, metadata about the application and links to resources such
as fonts, images and other information required in order to the application to run.
You create a configuration file using the function `gopi.NewAppConfig`
which takes a list of modules required for the running of the application as argument,
and returns an `gopi.AppConfig` object:

```
type AppConfig struct {
	// The set of modules which are required, including dependencies
	Modules  []*Module

	// The command-line arguments
	AppArgs  []string

	// The command-line flags
	AppFlags *Flags

	// Whether to log at debugging level
	Debug    bool

	// Whether to output verbose information
	Verbose  bool
}
```

The comma-separated list of modules you provide to the `gopi.NewAppConfig` function 
will be expanded to also include any modules where there are dependencies and 
a logging module is implicit to the list of modules. You can set the `Debug`
and `Verbose` variables explicitly, but it is overridden when the `-debug`
or `-verbose` flags exist on the command-line when invoking the application
from the command-line.

More information about modules is given in a future section, but for now you
can note to refer to modules by their explicit name or use a reserved word
to include a module by type rather than by name. Here's a list of reserved
words and how they map onto module types:

| Reserved word | Type                      | Description                 |
| -- | -- | -- |
| "logger"      | gopi.MODULE_TYPE_LOGGER   | Logging module (implicit)   |
|	"hw"          | gopi.MODULE_TYPE_HARDWARE | Hardware module             |
|	"display"     | gopi.MODULE_TYPE_DISPLAY  | Display                     |
|	"graphics"    | gopi.MODULE_TYPE_GRAPHICS | Graphics Manager            |
|	"fonts"       | gopi.MODULE_TYPE_FONTS    | Font Manager                |
|	"vector"      | gopi.MODULE_TYPE_VECTOR   | 2D Graphics Renderer        |
|	"opengl"      | gopi.MODULE_TYPE_OPENGL   | 3D Graphics Renderer        |
|	"layout"      | gopi.MODULE_TYPE_LAYOUT   | Box Layout                  |
|	"gpio"        | gopi.MODULE_TYPE_GPIO     | GPIO Hardware Interface     |
|	"i2c"         | gopi.MODULE_TYPE_I2C      | I2C Hardware Interface      |
|	"spi"         | gopi.MODULE_TYPE_SPI      | SPI Hardware Interface      |
|	"input"       | gopi.MODULE_TYPE_INPUT    | Input Manager               |
|	"mdns"        | gopi.MODULE_TYPE_MDNS     | RPC Service Discovery       |
|	"timer"       | gopi.MODULE_TYPE_TIMER    | Timer Manager               |
|	"lirc"        | gopi.MODULE_TYPE_LIRC     | Infrared Hardware Interface |


## Declaring Command Line Flags

##Foreground and Background tasks

## Using Application Modules

## What's next?



