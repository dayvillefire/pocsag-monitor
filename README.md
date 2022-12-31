# POCSAG-MONITOR

Golang POCSAG512 monitor, designed to be run on a Linux machine.

This code has been run on a [Raspberry Pi Zero W](https://amzn.to/3bflIyP) (the slower first-generation model) if anyone wants to run this on a standalone single-chip machine.

POCSAG monitor is currently keyed to publish to a number of different pluggable outputs, but is primarily designed to push to Discord channels.

## OUTPUT PLUGINS

* discord - Discord channels
* stdout - Print to standard output

## PREREQUISITES

* Linux-compatible SDR
* `rtl_fm` binary
* `multimon-ng` binary

## RECOMMENDED HARDWARE

* [Raspberry Pi Zero W](https://amzn.to/3bflIyP) - 79$ US
* [NESDR Mini SDR](https://amzn.to/3TXecta) - 26$ US (the antenna that comes with this unit is not great, use the replacement one)
* [VHF Antenna](https://amzn.to/3ssavjt) - 9$ US
* [High Endurance SD Card](https://amzn.to/3Szn8Uj) - 10$ US

You'll also need a power supply. If the Pi Zero W is a little too expensive for you due to the parts shortage, consider "Le Potato", instead of the Pi Zero W:

* [Le Potato](https://amzn.to/3N57YF4) - 45$ US
* [Edimax Wifi Dongle](https://amzn.to/3SDyiaF) - 9$ US
* [Case and Power Supply](https://amzn.to/3gKwPlP) - 15$ US

