
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

| Reserved word | Type                        | Description                 |
| -- | -- | -- |
| "logger"      | `gopi.MODULE_TYPE_LOGGER`   | Logging module (implicit)   |
|	"hw"          | `gopi.MODULE_TYPE_HARDWARE` | Hardware module             |
|	"display"     | `gopi.MODULE_TYPE_DISPLAY`  | Display                     |
|	"graphics"    | `gopi.MODULE_TYPE_GRAPHICS` | Graphics Manager            |
|	"fonts"       | `gopi.MODULE_TYPE_FONTS`    | Font Manager                |
|	"vector"      | `gopi.MODULE_TYPE_VECTOR`   | 2D Graphics Renderer        |
|	"opengl"      | `gopi.MODULE_TYPE_OPENGL`   | 3D Graphics Renderer        |
|	"layout"      | `gopi.MODULE_TYPE_LAYOUT`   | Box Layout                  |
|	"gpio"        | `gopi.MODULE_TYPE_GPIO`     | GPIO Hardware Interface     |
|	"i2c"         | `gopi.MODULE_TYPE_I2C`      | I2C Hardware Interface      |
|	"spi"         | `gopi.MODULE_TYPE_SPI`      | SPI Hardware Interface      |
|	"input"       | `gopi.MODULE_TYPE_INPUT`    | Input Manager               |
|	"mdns"        | `gopi.MODULE_TYPE_MDNS`     | RPC Service Discovery       |
|	"timer"       | `gopi.MODULE_TYPE_TIMER`    | Timer Manager               |
|	"lirc"        | `gopi.MODULE_TYPE_LIRC`     | Infrared Hardware Interface |

If you declare the use of a module by passing it into `gopi.NewAppConfig`
then you also need to anonymously import the module as per the example
above. For example, to import the `logger` and `hw` modules on a Raspberry
Pi you would do the following:

```
import (
	// Import modules
	_ "github.com/djthorpe/gopi/sys/hw/rpi"
)
```

If your target platform is Linux you might want to do the following instead:

```
import (
	// Import modules
	_ "github.com/djthorpe/gopi/sys/hw/linux"
)
```

You can in fact import both modules simultaneously and use 
[build tags](https://golang.org/pkg/go/build/) to choose which variant to use.
You can see the examples in the `cmd` folder to see how cross-platform
applications can be developed which make the best use of the target platform.

## Declaring Command Line Flags

__golang__ has a package called `flags` to define command-line flags but __gopi__
builds upon this to provide some additional mechanisms with `gopi.Flags`.
Here is the list of functions:

```
// Create a new flags object
func NewFlags(name string) Flags

type Flags interface {
    // Parse command line argumentsinto flags and pure arguments
    Parse(args []string) error

    // Parsed reports whether the command-line flags have been parsed
    Parsed() bool

    // Name returns the name of the flagset
    Name() string

    // Args returns the command line arguments as an array which aren't flags
    Args() []string

    // Flags returns the array of flags which were set on the command line
    Flags() []string

    // HasFlag returns a boolean indicating if a flag was set on the command line
    HasFlag(name string) bool

    // Define flags and return pointer to the flag value
    FlagString(name string, value string, usage string) *string
    FlagBool(name string, value bool, usage string) *bool
    FlagDuration(name string, value time.Duration, usage string) *time.Duration
    FlagInt(name string, value int, usage string) *int
    FlagUint(name string, value uint, usage string) *uint
    FlagFloat64(name string, value float64, usage string) *float64

    // Return flag values and boolean value which indicates presence on command line
    GetBool(name string) (bool, bool)
    GetString(name string) (string, bool)
    GetDuration(name string) (time.Duration, bool)
    GetInt(name string) (int, bool)
    GetUint(name string) (uint, bool)
    GetFloat64(name string) (float64, bool)
}
```

Ultimately the `gopi.NewAppConfig` function call will return a `gopi.Flags` object
into which you can define your own command-line flags. Modules you use will also
have the ability to add flags, which you can see from the example above. If the
flag `-help` is invoked then instead of your application running, it simply prints
out the usage information for the flags and exits.

## Foreground and Background tasks

Once you have your configuration object, you can create an application instance
and run your application code. In general an application instance is created
of type `gopi.AppInstance`, each module creates an instance of that module for
use in your tasks, and foreground and background tasks are started.

Your application is terminated when your foreground tasks returns with either
an error or with `nil` indicating successful completion. Before this, you can
choose to terminate your background tasks earlier before final cleanup. 

In a __Command Line Tool__ here is generally what a foreground task might 
look like:

```
func ForegroundTask(app *gopi.AppInstance, done chan<- struct{}) error {
    // ...Parse command-line arguments and check for validity
    // ...Perform any other initialization
    err := ...
    if err != nil {
        // Return error which is printed out on os.Stderr
        // and sets exit condition to -1
        return err
    }

    // ...If there are background tasks then pass information
    // onto them...

    // Continue processing until signalled to stop
    app.WaitForSignal()

	// Signal to background tasks that main thread is done
	done <- gopi.DONE

    // ...Perform any other cleanup

    // Return success (exit condition is 0)
	return nil
}
```

In comparison this is what a background task might look like:

```
func BackgroundTask(app *gopi.AppInstance, done chan<- struct{}) error {

    // Subscribe to events from modules
    chan1 := app.Module1.Subscribe()
    chan2 := app.Module2.Subscribe()
    chan3 := app.ModuleInstance('Module3').Subscribe()

    FOR_LOOP: for {
        select {
        case <-done:
                break FOR_LOOP
        case evt := <-chan1:
            // ... Process Module1 event
        case evt := <-chan2:
            // ... Process Module2 event
        case evt := <-chan3:
            // ... Process Module3 event        
        }
    }

    // Unsubscribe from channels
    app.Module1.Unsubscribe(chan1)
    app.Module2.Unsubscribe(chan2)
    app.ModuleInstance('Module3').Unsubscribe(chan3)

    // Return success (exit condition is 0)
	return nil
}
```

One or more background tasks are essentially event-driven, accepting
events from modules (or from the application itself), processing them
and then waiting for other events, one of which can be the `done`
signal which propogates from the foreground task.

More information on events is given in a future section, but for now
it's important to distingush between how a foreground task and a background
task operates in a __Command Line Tool__.

Finally your `main` function will invoke your tasks after set-up and
co-ordinate communications between them. The most simple `main` function
with one foreground task and two background tasks may look like this:

```
func main() {
	os.Exit(
        gopi.CommandLineTool(
          gopi.NewAppConfig("Module1", "Module2", "Module3"),
          ForegroundTask,
          BackgroundTask1,BackgroundTask2,
        )
    )
}
```

## Using Application Modules

As mentioned, you can use modules within your code by:

  1. Importing the variant of module you wish to use anonymously into your
     application
  2. Indicating in your configuration file which modules you wish to use
  3. Using the instance of the module within either your foreground or background
     tasks, cross-referencing the abstract interface for the methods that can be used.

Some modules can be referenced using the `app.Module` format, but others require
you to use the `app.ModuleInstance()` function, and then to cast them to your chosen
interface. For example:

```
func ForegroundTask(app *gopi.AppInstance, done chan<- struct{}) error {
  logger := app.Logger // implements the gopi.Logger interface
  display := app.ModuleInstance("display").(gopi.Display) // implements the gopi.Display interface
  // ... code here
}

Here is a list of some application modules, their "abstract interface" names
and the import path you would use. Note that for some, there are different
implementations. Information on how to choose and use each interface is detailed
in the rest of this guide.

| Name        | Use                 | Abstract Interface    | Import                                      |
| -- | -- | -- | -- |
| "logger"    | app.Logger          | `gopi.Logger`         | `github.com/djthorpe/gopi/sys/logger`       |
| "timer"     | app.Timer           | `gopi.Timer`          | `github.com/djthorpe/gopi/sys/timer`        |
| "hw"        | app.Hardware        | `gopi.Hardware`       | `github.com/djthorpe/gopi/sys/hw/rpi`       |
| "display"   | app.Display         | `gopi.Display`        | `github.com/djthorpe/gopi/sys/hw/rpi`       |
| "hw"        | app.Hardware        | `gopi.Hardware`       | `github.com/djthorpe/gopi/sys/hw/linux`     |
| "graphics"  | app.GraphicsManager | `gopi.SurfaceManager` | `github.com/djthorpe/gopi/sys/graphics/rpi` |
| "fonts"     | app.FontManager     | `gopi.FontManager`    | `github.com/djthorpe/gopi/sys/fonts/rpi`    |
| "input"     | app.InputManager    | `gopi.InputManager`   | `github.com/djthorpe/gopi/sys/input/linux`  |
| "gpio"      | app.GPIO            | `gopi.GPIO`           | `github.com/djthorpe/gopi/sys/hw/linux`     |
| "gpio"      | app.GPIO            | `gopi.GPIO`           | `github.com/djthorpe/gopi/sys/hw/rpi`       |
| "i2c"       | app.I2C             | `gopi.I2C`            | `github.com/djthorpe/gopi/sys/hw/linux`     |
| "spi"       | app.SPI             | `gopi.SPI`            | `github.com/djthorpe/gopi/sys/hw/linux`     |
| "lirc"      | app.LIRC            | `gopi.LIRC`           | `github.com/djthorpe/gopi/sys/hw/linux`     |

## What's next?



