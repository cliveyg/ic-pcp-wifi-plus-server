#!/bin/sh

#-----------------------------------------------------------------------------#
# wp-stop-start-wifi.sh                                                       #
# for some reason the only way i have been able to run these as nohup is by   #
# putting them in separate script. at least part of the issue is that         #
# stopping and starting wlan0 wifi exhibits different behaviours in the shell #
# and via a script <shrugs>                                                   #
#-----------------------------------------------------------------------------#

set -a
. /var/www/.env
set +a

. /var/www/cgi-bin/pcp-functions

LOG=/mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts/testy.log

if [ $DBUG -eq 1 ]; then

  if [ -f $LOG ]; then
    echo "------------------------------" >> $LOG
    echo "DBUG IS [$DBUG]" >> $LOG
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

else

  /usr/local/etc/init.d/wifi wlan0 stop
  mount /dev/mmcblk0p1
  sudo mv /mnt/mmcblk0p1/used_wpa_supplicant.conf /var/run/wpa_supplicant.conf
  echo "ctrl_interface=/var/run/wpa_supplicant
  ctrl_interface_group=staff
  update_config=1" > /opt/wpa.cfg
  sudo wpa_supplicant -Dwext -iwlan0 -c/opt/wpa.cfg -B
  sleep 3
  sudo /usr/local/etc/init.d/wifi wlan0 stop
  sudo /usr/local/etc/init.d/wifi wlan0 start
  cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts
  ./wifi-plus-startup.sh

fi

echo "{ \"beep\": \"boop\" }"