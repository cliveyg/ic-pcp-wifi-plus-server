#!/bin/sh

nohup "$(sleep 12; ./wifi-plus-startup.sh)" > /dev/null 2>&1 &
sleep 1
nohup "$(sleep 6; /usr/local/etc/init.d/wifi wlan0 start)" > /dev/null 2>&1 &
sleep 1
nohup "$(/usr/local/etc/init.d/wifi wlan0 stop)" > /dev/null 2>&1 &
echo "{ \"beep\": \"boop\" }"