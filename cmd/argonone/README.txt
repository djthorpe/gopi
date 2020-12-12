
# Argon One controller

This command allows you to control the fan, power and IR interface of the
Argon One case. In order to install, you will need to:

  1. Ensure the LIRC and I2C kernel modules are installed;
  2. Set the daemon to run as a systemctl service.

The following lines can be added to your `/boot/config.txt` file in order
to enable the kernel modules:

```
dtoverlay=gpio-ir,gpio_pin=23
dtparam=i2c_arm=on
```

## Making the debian package

Make the Debian `.deb` package using the following commands:

```bash
bash# go get github.com/djthorpe/gopi/v3
bash# cd gopi
bash# make debian
bash# sudo dpkg -i build/argonone_VERSION_amd64.deb
```

