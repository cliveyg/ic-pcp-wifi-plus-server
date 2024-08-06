#!/bin/sh

LOG=/mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts/testy.log

if [ -f $LOG ]; then
  echo "------------------------------" >> $LOG
  echo "$DBUG" >> $LOG
else
  touch testy.log
fi
echo $(/usr/local/etc/init.d/wifi wlan0 stop) >> $LOG
echo $(mount /dev/mmcblk0p1) >> $LOG
echo $(sudo mv /mnt/mmcblk0p1/used_wpa_supplicant.conf /var/run/wpa_supplicant.conf)
echo "ctrl_interface=/var/run/wpa_supplicant
ctrl_interface_group=staff
update_config=1" > /opt/wpa.cfg
if [ -f /opt/wpa.cfg ]; then
  echo "opt.cfg created" >> $LOG
fi
echo $(sudo wpa_supplicant -Dwext -iwlan0 -c/opt/wpa.cfg -B) >> $LOG
sleep 3
echo $(sudo /usr/local/etc/init.d/wifi wlan0 stop) >> $LOG
echo $(sudo /usr/local/etc/init.d/wifi wlan0 start) >> $LOG
cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts
./wifi-plus-startup.sh
if [ $(pidof wifiplus) != "" ]; then
  echo "wifiplus running" >> $LOG
fi

#nohup "$(sleep 15; ./wifi-plus-startup.sh)" > /dev/null 2>&1 &
#sleep 1
#nohup "$(sleep 6; /usr/local/etc/init.d/wifi wlan0 start)" > /dev/null 2>&1 &
#sleep 1
#nohup "$(/usr/local/etc/init.d/wifi wlan0 stop; sleep 6; /usr/local/etc/init.d/wifi wlan0 start)" > /dev/null 2>&1 &
echo "{ \"beep\": \"boop\" }"