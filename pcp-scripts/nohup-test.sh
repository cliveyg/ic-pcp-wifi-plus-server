#!/bin/sh

nohup "$(sleep 10; touch /mnt/UserData/wifi-plus/tfile3.txt; ll /mnt/UserData/wifi-plus/tfile3.txt)" > /www/log/wifiplus.log 2>&1 &
nohup "$(sleep 5; touch /mnt/UserData/wifi-plus/tfile2.txt; ll /mnt/UserData/wifi-plus/tfile2.txt)" > /www/log/wifiplus.log 2>&1 &
nohup "$(sleep 1; touch /mnt/UserData/wifi-plus/tfile1.txt; ll /mnt/UserData/wifi-plus/tfile1.txt)" > /www/log/wifiplus.log 2>&1 &
echo "{ \"beep\": \"boop\" }"