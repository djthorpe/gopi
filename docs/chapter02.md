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
    app, err := app.NewCommandLineTool(Main)
    if err != nil {
        panic(err)
    }
	os.Exit(app.Run())
}
```

The features of this command-line tool are:

  * Importing _frameworks_ from `gopi`;
  * Importing _units_ from `gopi`;
  * A `Main` function which takes an application and the command-line arguments as parameters;
  * A `main` function which is the entry point for the command-line tool, this function creates the `app` object, and exits after the `app.Run()` method has been called.

  In the background, the `NewCommandLineTool` method creates _unit instances_ and parses any command-line parameters.

  If you compile and run this tool, it will print out the message on your console and wait for you to press CTRL+C (the `WaitForSignal` method blocks until the os.Interrupt signal is returned).
  

