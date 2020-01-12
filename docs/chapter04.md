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

### Application Examples

There are a number of example applications in the `cmd` folder which you can examine and run. The following sections describe how to install and run these examples. Ultimately you can compile them all through the use of the `make` command:

* [cmd/helloworld](https://github.com/djthorpe/gopi/tree/v2/cmd/helloworld) is the canonical first program. It prints out your name and waits for you to send keyboard interrupt \(CTRL+C\);
* [cmd/timers](https://github.com/djthorpe/gopi/tree/v2/cmd/timers) prints out messages as a ticker fires once every second;
* [cmd/hwinfo](https://github.com/djthorpe/gopi/tree/v2/cmd/hwinfo) displays a table of hardware information about your platform;
* [cmd/i2cdetect](https://github.com/djthorpe/gopi/tree/v2/cmd/i2cdetect) displays a table of detected I2C devices connected on a bus;
* [cmd/fonts](https://github.com/djthorpe/gopi/tree/v2/cmd/fonts) loads fonts from a path and displays information about the loaded fonts;
* [cmd/discovery](https://github.com/djthorpe/gopi/tree/v2/cmd/discovery) displays service names and instances registered through multicast DNS \(mDNS\) on your local network.

### Embedding Version Information

The `Makefile` demonstrates embedding version information into your application through the use of linker flags. If you include version information then your application will automatically have a `-version` flag which prints out this information. For example,

```bash
bash% helloworld -version

buildtime  2020-01-03T20:02:16Z
tag        v2.0.1
branch     v2
hash       f7d5e2c4b02d2450376c84df8a61d7191b625b23

```

You can use the following linker flags in order to embed the information:

```bash
bash% GOPI=github.com/djthorpe/gopi/v2/config
bash% LDFLAGS=${LDFLAGS} -X ${GOPI}.GitTag=$(shell git describe --tags)
bash% LDFLAGS=${LDFLAGS} -X ${GOPI}.GitBranch=$(shell git name-rev HEAD --name-only --always)
bash% LDFLAGS=${LDFLAGS} -X ${GOPI}.GitHash=$(shell git rev-parse HEAD)
bash% LDFLAGS=${LDFLAGS} -X ${GOPI}.GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
bash% go build -o helloworld -ldflags "${LDFLAGS}" ./cmd/helloworld/...
```



