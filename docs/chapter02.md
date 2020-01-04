# Getting Started

The simplest program prints out "Hello, World" on your terminal. It's quite a bit
longer than a normal golang helloworld (which would be about ten lines long):

```go
package main

import (
	"context"
	"fmt"
	// Frameworks
	"github.com/djthorpe/gopi/v2"
    "github.com/djthorpe/gopi/v2/app"   
    // Units
    _ "github.com/djthorpe/gopi/v2/unit/logger"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {

	// Print out name
	fmt.Println("Hello, World!")
	fmt.Println("Press CTRL+C to exit")

	// Wait for CTRL+C
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
    app, err := app.NewCommandLineTool(Main,nil)
    if err != nil {
        panic(err)
    }
	os.Exit(app.Run())
}
```

The features of this command-line tool are:

  * Importing __frameworks__ from `gopi`. For more about the `gopi` package, see below;
  * Importing __units__ from `gopi`. Note these are imported anonymously since all you are doing is registering their presence. For more about units, see below;
  * A `Main` function which takes an application and the command-line arguments as parameters;
  * A `main` function which is the entry point for the command-line tool, this function creates the `app` object, and exits after the `app.Run()` method has been called.

In the background, the `NewCommandLineTool` method creates _unit instances_ and parses any command-line parameters.

If you compile and run this tool, it will print out the message on your console and wait for you to press CTRL+C (the `WaitForSignal` method blocks until the os.Interrupt signal is returned).

## What is a Unit?

A __unit__ is a Go Language module which adheres to a particular interface, and has these other features:

  * It has a unique __Name__. The name is optional if the unit has a __Type__;
  * It has a __Type__. The type is optional is the unit has a __Name__;
  * It may depend on other units, which can be referred to by name or type;
  * It usually has a function defined to create a __Unit Instance__.

You can register units in your tool by importing them anonymously. For example, the
example above registers a logging unit:

```go
import (
    _ "github.com/djthorpe/gopi/v2/unit/logger"
)
```

If you look at the `init.go` file in the logging package (function bodies removed to aide reading):

```go
func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "gopi/logger",
		Type: gopi.UNIT_LOGGER,
        Config: func(app gopi.App) error { ... },
        New: func(app gopi.App) (gopi.Unit, error) { ... },
	})
}
```

The __Type__ parameter is tied to the interface for the returned unit, in this case the `New` function returns a unit which satisfies the `gopi.Logger` interface:

```go
type gopi.Logger interface {
	gopi.Unit

	Name() string 	// Name returns the name of the logger
	Clone(string) gopi.Logger 	// Returns a new logger with a different name
    Error(error) error 	// Error logs an error
	Debug(args ...interface{}) // Output debug message
	IsDebug() bool 	// IsDebug returns true if debugging is enabled
}
```

If you want to use a unit instance within your application, you need to 
do three things:

  1. Import the unit anonymously into your tool;
  2. Indicate you want to use the unit in your tool in an argument to `app.NewCommandLineTool`. This will also satisfy the unit dependencies
 so that unit instances are created in the right order, but you will need
 to make sure all the units are imported in step one;
  3. Use the `app.UnitInstance()` method to obtain the unit instance in your `Main` function. You will need to cast it to the correct interface in order to use it.

Here's a slightly expanded example which demonstrates using unitA and unitB in a command line application:

```go
package main

import (
	"context"
	"fmt"
	// Frameworks
	"github.com/djthorpe/gopi/v2"
    "github.com/djthorpe/gopi/v2/app"   
    // Use units A and B
    _ "github.com/djthorpe/gopi/v2/unit/A"
    _ "github.com/djthorpe/gopi/v2/unit/B"
)

type interfaceA interface {
    // ...
    gopi.Unit
}

type interfaceB interface {
    // ...
    gopi.Unit
}

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
    A := app.UnitInstance("A").(interfaceA)
    B := app.UnitInstance("B").(interfaceB)

    // ...

    // Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
    // Use units A and B
    app, err := app.NewCommandLineTool(Main,nil,"A","B")
    if err != nil {
        panic(err)
    }
	os.Exit(app.Run())
}
```

## Why Units?

Go interfaces provide excellent abstractions for encapsulating an object and methods, but creating complex cross-platform tools still requires some additional patterns and techniques.

Being able to include modules at runtime and for those modules to be magically created and automatically depend on others reduces run-time complexity.

## The Logger

You have already seen the logger unit used in the first "Hello, World" tool. In fact, every tool needs to have a logger defined and it's not necessary to declare the usage of the logger when calling `app.NewCommandLineTool`.

A convenience method `app.Log()` will return the logging unit. For example,

```go
func Main(app gopi.App, args []string) error {
    app.Log().Debug("In Main")
    // Return success
	return nil
}
```

You can only import one logger of type `gopi.UNIT_LOGGER` into your tool.
The unit you import as `github.com/djthorpe/gopi/v2/unit/logger` outputs
messages to `os.Stderr` and defines some command-line flags so that when you invoke your tool. For example, if you invoke it with the `-help` flag:

```bash
bash$ helloworld -help
  -debug
    	Debugging output
  -verbose
    	Verbose output (default true)
```

Other implementations of the `gopi.UNIT_LOGGER` could output messages to file or to the `syslog` service, for example.

