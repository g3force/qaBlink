# qaBlink

Get the current build status from Jenkins and Quality Gate Status from Sonar for multiple projects and report them on a blink1 device.

All stati can be placed to 'slots'. Slots are gone through iteratively with a fixed duration per slot. Stati from jobs can be assigned to an available LED of the blink1 devices that are attached.

## Installation

This project is go gettable:

`go get github.com/g3force/qaBlink`

It will produce a `qaBlink` executable. The executable assumes a `config.json` file in the current working directory.

## Requirements

This project depends on https://github.com/hink/go-blink1/ which makes use of libusb to access the blink1 device. libusb is available for many Platforms including Linux, Mac and Windows.

## Configuration

Copy example.config.json to config.json and adapt it to your needs. You can define multiple Jenkins and Sonar connections and jobs and assign them to slots.

