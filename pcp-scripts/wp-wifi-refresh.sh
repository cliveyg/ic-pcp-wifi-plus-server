#!/bin/sh

#-----------------------------------------------------------------------------#
# wp-wifi-refresh.sh                                                       #
# for some reason the only way i have been able to run these as nohup is by   #
# putting them in separate script. at least part of the issue is that         #
# stopping and starting wlan0 wifi exhibits different behaviours in the shell #
# and via a script <shrugs>                                                   #
#-----------------------------------------------------------------------------#

set -a
. /var/www/.env
set +a

LOG=/var/log/wifiplus.log

if [ $DBUG -eq 1 ]; then

  if [ ! -f $LOG ]; then
    sudo touch $LOG
  fi
  echo "[wp-wifi-refresh.sh] ------------------------------" >> $LOG
  echo -n "[wp-wifi-refresh.sh] " >> $LOG
  echo $(/usr/local/etc/init.d/wifi wlan0 stop) >> $LOG
  echo -n "[wp-wifi-refresh.sh] " >> $LOG
  echo $(mount /dev/mmcblk0p1) >> $LOG
  echo -n "[wp-wifi-refresh.sh] " >> $LOG
  echo $(sudo cp /mnt/mmcblk0p1/used_wpa_supplicant.conf /var/run/wpa_supplicant.conf)
  echo "ctrl_interface=/var/run/wpa_supplicant
  ctrl_interface_group=staff
  update_config=1" > /opt/wpa.cfg
  if [ -f /opt/wpa.cfg ]; then
    echo -n "[wp-wifi-refresh.sh] " >> $LOG
    echo "[wp-wifi-refresh.sh] opt.cfg created" >> $LOG
  else
    echo -n "[wp-wifi-refresh.sh] Failed to create opt.cfg file!" >> $LOG
  fi
  echo -n "[wp-wifi-refresh.sh] Starting wpa_supplicant... " >> $LOG
  echo $(sudo wpa_supplicant -Dwext -iwlan0 -c/opt/wpa.cfg -B) >> $LOG
  sleep 3
  echo -n "[wp-wifi-refresh.sh] Stopping wifi " >> $LOG
  echo $(sudo /usr/local/etc/init.d/wifi wlan0 stop) >> $LOG
  echo -n "[wp-wifi-refresh.sh] Starting wifi " >> $LOG
  echo $(sudo /usr/local/etc/init.d/wifi wlan0 start) >> $LOG
  echo -n "[wp-wifi-refresh.sh] Wifi started " >> $LOG
  cd /mnt/UserData/ic-pcp-wifi-plus-server/pcp-scripts
  echo -n "[wp-wifi-refresh.sh] Starting wifi-plus-startup.sh " >> $LOG
  ./wifi-plus-startup.sh
  if [ $(pidof wifiplus) != "" ]; then
    echo "[wp-wifi-refresh.sh] wifiplus exe running" >> $LOG
  else
    echo "[wp-wifi-refresh.sh] wifiplus exe not running!" >> $LOG
    echo "[wp-wifi-refresh.sh] Attempting to start again..." >> $LOG
    ./wifi-plus-startup.sh
    sleep 3
    if [ $(pidof wifiplus) != "" ]; then
        echo "[wp-wifi-refresh.sh] wifiplus exe running!" >> $LOG
    else
      echo "[wp-wifi-refresh.sh] wifiplus still exe not running :(" >> $LOG
    fi
  fi
else

  /usr/local/etc/init.d/wifi wlan0 stop
  mount /dev/mmcblk0p1
  sudo cp /mnt/mmcblk0p1/used_wpa_supplicant.conf /var/run/wpa_supplicant.conf
  echo "ctrl_interface=/var/run/wpa_supplicant
  ctrl_interface_group=staff
  update_config=1" > /opt/wpa.cfg
  sudo wpa_supplicant -Dwext -iwlan0 -c/opt/wpa.cfg -B
  sleep 3
  sudo /usr/local/etc/init.d/wifi wlan0 stop
  sudo /usr/local/etc/init.d/wifi wlan0 start
  sleep 2
  sudo /usr/local/etc/init.d/wifi wlan0 restart
  sleep 3
  cd /mnt/UserData/ic-pcp-wifi-plus-server/pcp-scripts
  ./wifi-plus-startup.sh

fi

echo "{ \"beep\": \"boop\" }"