#!/bin/sh

/usr/local/etc/init.d/wifi wlan0 stop
mount /dev/mmcblk0p1
sudo mv /mnt/mmcblk0p1/used_wpa_supplicant.conf /var/run/wpa_supplicant.conf
echo "ctrl_interface=/var/run/wpa_supplicant
ctrl_interface_group=staff
update_config=1" > /opt/wpa.cfg
sudo wpa_supplicant -Dwext -iwlan0 -c/opt/wpa.cfg -B
sleep 3
sudo /usr/local/etc/init.d/wifi wlan0 start

#nohup "$(sleep 15; ./wifi-plus-startup.sh)" > /dev/null 2>&1 &
#sleep 1
#nohup "$(sleep 6; /usr/local/etc/init.d/wifi wlan0 start)" > /dev/null 2>&1 &
#sleep 1
#nohup "$(/usr/local/etc/init.d/wifi wlan0 stop; sleep 6; /usr/local/etc/init.d/wifi wlan0 start)" > /dev/null 2>&1 &
echo "{ \"beep\": \"boop\" }"