# industrialcool-pcp-wifi
Project to add currently unsupported wifi functionality to piCorePlayer.

### The background:
I make boomboxes out of wood, and I'm currently developing a new sound system that uses a RaspberryPi 4 with an DigiAmp+ sound card. The RPi4 controls a touchscreen and runs piCorePlayer, Squeezelite and Lyrion Music Server.

The boomboxes are intended to be portable and by default LMS and PiCorePlayer prefer to have an external wifi to connect to. It is possible to switch to using the RPi as a wireless access point (WAP) but as piCorePlayer currently stands you need an ethernet connection or keyboard connected to the RPi4 to manage the switch. For a portable sound system this could sometimes prove to be a bit tricky ;)

### The problem:
Controlling pcp and LMS is easiest via a phone using the local wifi network. For example your home router/modem is used to communicate between the pcp/LMS system and the phone in the most typical use case.

However, there will be occasions when there is a need to use the sound system when there is no available wifi network. piCorePlayer, and LMS in particular, expect and need a network to function. Luckily pcp provides the facility to configure the RPi as a WAP. I know it is possible to use AirPlay or Bluetooth but this leaves users unable to access parts of the complete pcp system. I will implement both but also anticipate users wanting to use wifi connectivity.

The issue is that as pcp currently stands a wired ethernet connection or keyboard is needed as the wifi needs to be turned off for the WAP page to be accessible. There are also issues around making sure to follow quite a specific set of actions to complete the switchover successfully. There are similar issues going the other way; switching off the WAP and switching on using the local wifi.

I would also like to be able to change wifi networks more easily. piCorePlayer is built using TinyCore Linux and as it currently stands changing wifi networks easily is not possible. I would also like the system to remember known wifi networks and on start up connect automatically to one if available.

### The plan:
To create an additional, cut-down and more user-friendly controller of piCorePlayer functions to supplement the current piCorePlayer whilst offering some additional functionality. After studying the pcp code I think I can plug in to existing shell code to achieve my aims. 

The new system must provide the following:

* On boot, scan for wifi networks and connect to known wifi.
* Provide a method for connecting to new wifi networks that doesn't involve the user having to fiddle with the wpa_supplicant file :)
* If neither of the above then create a wireless access point.
* Make this interface more modern looking and simpler for the end user.
* Provide access to the pcp ui if needed for power users.

### The tech stack:
As my background involves lots of REST APIs I thought it would be good to utilise an api between the RPI and an HTML5 app running on a phone. I'm currently programming a lot in Golang and would like to use this on the RPi as the backend. I also quite like React and the Material UI to provide the frontend.

### Backend:

Golang API that can leverage parts of the original code (written largely in shell script).

Positives:
* Compiled binary is small and fast.
* Can hook into existing pcp functionality and codebase.

Negatives:
* Changes to pcp could break the backend.

### Frontend

React single page application utilising Material UI.

Positives:
* Quick-ish to develop.
* Keeps me up to date on React ;D
* HTML5 app is cross-platform (iOS and Android)

Negatives:
* Performance?

### Communication Protocol

~~**gRPC & Protobuf**~~

~~Use a combination of gRPC and Protocol Buffers.~~

~~Positives:~~
~~* Faster than JSON and HTTP1.1~~
~~* More lightweight and robust than JSON/HTTP1.1~~
~~* Uses HTTP2 as communication protocol~~
~~* gRPC good with Golang~~

~~Negatives:~~
~~* Harder to grasp than more traditional Restful APIs.~~
~~* Less supported than JSON~~

**JSON & HTTP**

I was going to use gRPC but as it currently stands browsersdo not support gRPC natively. 
Browser based gRPC needs a proxy to function as a go between HTTP1 and HTTP2 and apps running in browsers.

As a result of this need for a proxy and the lightweight nature of this API I have decided to use standard JSON and HTTP1.1.

I may, time permitting, try and use websockets as I've used these with React in the past.


### Assumptions
* Only tested on RaspberryPi 4 8Gb Ram
* OS and other pcp partition on 32Gb SD card split into two equal partitions.
* Only tested with piCorePlayer version 9.2.0
* ARM binary built with Golang version 1.21.0
* Binary built on M2 MacBook Air running macOS version 14.5 (23F79)
* Music for LMS and UserData stored on separate 64Gb USB pen drive.

### HowTo

#### Building the binary
I am building the Golang binary on my M2 MacBook Air and the target is a RaspberryPi.  The go build settings I am using are:

```go
GOARCH=arm64
GOOS=linux
```

And the go build command is:

```go
CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -a -ldflags '-w' -o wifiplus
```

I am using these flags to avoid dynamically linked libs in the binary.
