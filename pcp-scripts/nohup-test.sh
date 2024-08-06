#!/bin/sh

nohup "$(sleep 10; ./wifi-plus-startup.sh)" > /dev/null 2>&1 &
nohup "$(sleep 5; touch /mnt/UserData/wifi-plus/tfile2.txt; ll /mnt/UserData/wifi-plus/tfile2.txt)" > /dev/null 2>&1 &
nohup "$(sleep 1; touch /mnt/UserData/wifi-plus/tfile1.txt; ll /mnt/UserData/wifi-plus/tfile1.txt)" > /dev/null 2>&1 &
echo "{ \"beep\": \"boop\" }"