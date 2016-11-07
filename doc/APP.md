
# Using the Application Framework

The application framework makes it easier to develop command-line and UI-based
applications. It manages resources, displays and events in the system, allowing
you to write more minimal code for your application. Here is a short example
of your `main()` function:

```go

import "github.com/djthorpe/gopi/app"

func main() {
	// Create the configuration
	config := app.Config(app.APP_GPIO)

	// Create the application
	myapp, err := app.NewApp(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer myapp.Close()

	// Run the application
	if err := myapp.Run(RunLoop); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}

```

Here, you are creating a configuration file which indicates you want to use
the GPIO resource. You can combine with other resource flags, for example,
`app.APP_GPIO | app.APP_I2C` could be used to indicate you want to use
both the GPIO and I2C peripheral interfaces.

Your code is then provided in a function called `RunLoop`. For example,

```go

func RunLoop(app *app.App) error {
	pins := []hw.GPIOPin{ app.GPIO.PhysicalPin(40) }
	led, err := gopi.Open(hw.LED{GPIO: app.GPIO, Pins: pins}, app.Logger)
	if err != nil {
		return err
	}
	defer led.Close()

	go func() {
		for {
			led.(hw.LEDDriver).On()
			time.Sleep(1 * time.Second)
			led.(hw.LEDDriver).Off()
			time.Sleep(1 * time.Second)
		}
	}()

	app.WaitUntilDone()

	// Return success
	return nil
}
```

In this example, it switches an LED signal wired to Physical Pin 40 on and off
at an interval of one second. The method `WaitUntilDone` will block execution
until a signal is received to terminate the application, usually through
receiving a `SIGTERM` or `SIGINT` signal.



