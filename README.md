# qaBlink

Get the current build status from Jenkins and Quality Gate Status from Sonar for multiple projects and report them on a blink1 device.

All status can be placed to 'slots'. Slots are gone through iteratively with a fixed duration per slot. Status from jobs can be assigned to an available LED of the blink1 devices that are attached.

## Requirements

This project depends on https://github.com/hink/go-blink1/ which makes use of libusb to access the blink1 device. libusb is available for many Platforms including Linux, Mac and Windows.

You need to install following dependencies first:

 * libusb
 * Go >= 1.9
 
### Linux

For Linux, you will need to add a udev rule if you want to execute the binary without root permissions. On Arch Linux, it would look like this:
```
cat /etc/udev/rules.d/10.blink1.rules 
SUBSYSTEMS=="usb", ATTRS{idVendor}=="27b8", ATTRS{idProduct}=="01ed", SYMLINK+="blink1", GROUP="wheel"
```

### Windows
In theory, it should also work with Windows, if libusb is installed. I have never tested this, though...

## Installation

Download and install to [GOPATH](https://github.com/golang/go/wiki/GOPATH):

`go get github.com/g3force/qaBlink`

It will produce a `qaBlink` executable in $GOPATH/bin. The executable assumes a `config.json` file in the current working directory or a `.qaBlink.json` in your HOME dir.

## Configuration

Copy example.config.json to config.json and adapt it to your needs. You can define multiple Jenkins and Sonar connections and jobs and assign them to slots.

