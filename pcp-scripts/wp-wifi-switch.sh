#!/bin/sh

#-----------------------------------------------------------------------------#
# wp-wifi-switch.sh                                                           #
#                                                                             #
#                                                                             #
#                                                                             #
#                                                                             #
#-----------------------------------------------------------------------------#

set -a
. /var/www/.env
set +a

LOG=$LOGFILE

ssid=$1
pass=$2

if [ $DBUG -eq 1 ]; then

  if [ ! -f $LOG ]; then
    sudo touch $LOG
  fi
  echo "[wp-wifi-switch.sh] ------------------------------" >> $LOG
  echo "[wp-wifi-switch.sh] SSID is [$ssid]" >> $LOG
  echo "[wp-wifi-switch.sh] Pass is [$pass]" >> $LOG
  echo '{ "status": 200, "message": "have loggy" }'
  #echo -n "[wp-wifi-switch.shh] " >> $LOG
  #echo $(/usr/local/etc/init.d/wifi wlan0 stop) >> $LOG
  #echo -n "[wp-wifi-switch.sh] " >> $LOG
  #echo $(mount /dev/mmcblk0p1) >> $LOG
  #echo -n "[wp-wifi-switch.sh] " >> $LOG
  #echo $(sudo mv /mnt/mmcblk0p1/used_wpa_supplicant.conf /var/run/wpa_supplicant.conf)
  #echo "ctrl_interface=/var/run/wpa_supplicant
  #ctrl_interface_group=staff
  #update_config=1" > /opt/wpa.cfg
  #if [ -f /opt/wpa.cfg ]; then
  #  echo -n "[wp-wifi-switch.sh] " >> $LOG
  #  echo "[wp-wifi-switch.sh] opt.cfg created" >> $LOG
  #else
  #  echo -n "[wp-wifi-switch.sh] Failed to create opt.cfg file!" >> $LOG
  #fi
  #echo -n "[wp-wifi-switch.sh] " >> $LOG
  #echo $(sudo wpa_supplicant -Dwext -iwlan0 -c/opt/wpa.cfg -B) >> $LOG
  #sleep 3
  #echo -n "[wp-wifi-switch.sh] " >> $LOG
  #echo $(sudo /usr/local/etc/init.d/wifi wlan0 stop) >> $LOG
  #echo -n "[wp-wifi-switch.sh] " >> $LOG
  #echo $(sudo /usr/local/etc/init.d/wifi wlan0 start) >> $LOG
  #cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts
  #./wifi-plus-startup.sh
  #if [ $(pidof wifiplus) != "" ]; then
  #  echo "[wp-wifi-refresh.sh] wifiplus exe running" >> $LOG
  #fi

else
  echo '{ "status": 404, "message": "no loggy" }'
  #/usr/local/etc/init.d/wifi wlan0 stop
  #mount /dev/mmcblk0p1
  #sudo cp /mnt/mmcblk0p1/used_wpa_supplicant.conf /var/run/wpa_supplicant.conf
  #echo "ctrl_interface=/var/run/wpa_supplicant
  #ctrl_interface_group=staff
  #update_config=1" > /opt/wpa.cfg
  #sudo wpa_supplicant -Dwext -iwlan0 -c/opt/wpa.cfg -B
  #sleep 3
  #sudo /usr/local/etc/init.d/wifi wlan0 stop
  #sudo /usr/local/etc/init.d/wifi wlan0 start
  #cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts
  #./wifi-plus-startup.sh

fi

#echo "{ \"beep\": \"boop\" }"