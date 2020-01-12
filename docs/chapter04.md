---
description: Starting application development.
---

# Developing Applications

In order to start developing your own application, you can take the following [template repository](https://github.com/djthorpe/gopi-app), which simply prints **Hello, World** on the screen:

```bash
bash% git clone https://github.com/djthorpe/gopi-app helloworld
bash% cd helloworld
bash% git remote remove origin
```

You should edit the following files in the repository in order to get started:

* `cmd/helloworld` Rename this folder to match the name of your application;
* `cmd/helloworld/units.go` Add and remove unit module imports and add any directly referenced units to the `UNITS` global variable;
* `cmd/helloworld/events.go` Add and remove event handlers in this file, and ensure any event handlers are added to the `EVENTS` global variable;
* `cmd/helloworld/main.go` Add command-line flags to the `main` function and edit the `Main` function to perform startup and shutdown of your application.

The `Makefile` in the repository contains targets for Linux, Darwin and Raspberry Pi. Invoke the `make` command with the appropriate target name:

```bash
bash% make linux
bash% make darwin
bash% make rpi
```

This will test, compile and install your command. Without modification the `helloworld` command is installed.

