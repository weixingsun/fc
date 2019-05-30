# This project may be consideredly dangerous
# If you do not have special domain skills for computer science

## Overview

This golang program exercises `Autonomous Flight Controller`. 

## Prerequisites
```
for iNav
map AERT
set receiver_type = MSP

# for Cleanflight
map AERT1234
set receiver_type = MSP *
```

## Building

 ```
 git clone
 go build
 ```


## Usage

```
$ ./msp_set_rx --help
Usage of msp_set_rx [options]
  -a	Arm (take care now)
  -b int
    	Baud rate (default 115200)
  -d string
    	Serial Device
```

Sets random (but safe) values:

```
$ ./msp_set_rx -d /dev/ttyACM0 115200
```

### Arm / Disarm test

The application can also test arm / disarm, with the `-a` option. In this mode, the application:

* Stay quiescent state for 10 seconds
* Arms using the customary stick command
* Maintains min-throttle for 5 seconds
* Disarms (stick command)

The vehicle must be in a state that will allow arming: [iNav wiki article](https://github.com/iNavFlight/inav/wiki/%22Something%22-is-disabled----Reasons).

Summary of output (`##` indicates a comment, repeated lines removed).

```
$ ./msp_set_rx -d /dev/ttyACM0 -a
```

While this attempts to arm at a safe throttle value, removing props or using a current limiter is recommended.


## Licence

None
