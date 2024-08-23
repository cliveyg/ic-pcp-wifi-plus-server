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

if [ $DBUG -eq 1 ]; then
  LOG=/var/log/wifiplus.log
  if [ ! -f $LOG ]; then
    sudo touch $LOG
  fi
else
  LOG=/dev/null
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
#echo -n "[wp-wifi-refresh.sh] Stopping wifi " >> $LOG
#echo $(sudo /usr/local/etc/init.d/wifi wlan0 stop) >> $LOG
#echo -n "[wp-wifi-refresh.sh] Starting wifi " >> $LOG
#echo $(sudo /usr/local/etc/init.d/wifi wlan0 start) >> $LOG
#echo -n "[wp-wifi-refresh.sh] Wifi started " >> $LOG
echo -n "[wp-wifi-refresh.sh] Restarting wifi with /usr/local/etc/init.d/wifi restart " >> $LOG
echo $(sudo /usr/local/etc/init.d/wifi wlan0 restart) >> $LOG

sleep 2
echo -n "[wp-wifi-refresh.sh] Attempting to kill and restart udhcpc" >> $LOG
sudo kill `ps | grep udhcpc | grep wlan0 | awk '{print $1}'` > /dev/null 2>&1
if [ $? -eq 0 ]; then
  rm -f /var/run/udhcpc.wlan0.pid
  if [ $? -eq 0 ]; then
    echo "[wp-wifi-refresh.sh] Killed udhcpc process"
  else
    echo "[wp-wifi-refresh.sh] Failed to kill udhcpc process. Exiting..."
    return 2
else
  echo "[wp-wifi-refresh.sh] Failed to kill udhcpc process. Exiting..."
  return 2
fi

sudo /sbin/udhcpc -b -i wlan0 -A 5 -x hostname:$(/bin/hostname) -p /var/run/udhcpc.wlan0.pid
if [ $? -ne 0 ]; then
  echo "[wp-wifi-refresh.sh] unable to restart udhcpc. Exiting..." >> $LOG
  return 2
else
  echo "[wp-wifi-refresh.sh] Restart udhcpc [OK]" >> $LOG
fi

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


#  echo "[wp-wifi-refresh.sh] wifiplus exe running" >> $LOG
#  /usr/local/etc/init.d/wifi wlan0 stop
#  mount /dev/mmcblk0p1
#  sudo cp /mnt/mmcblk0p1/used_wpa_supplicant.conf /var/run/wpa_supplicant.conf
#  echo "ctrl_interface=/var/run/wpa_supplicant
#  ctrl_interface_group=staff
#  update_config=1" > /opt/wpa.cfg
#  sudo wpa_supplicant -Dnl80211,wext -iwlan0 -c/opt/wpa.cfg -B > /dev/null 2>&1
#  sleep 3
#  sudo /usr/local/etc/init.d/wifi wlan0 restart
#  sleep 2
#
#  # kill udhcpc and restart
#  sudo kill `ps | grep udhcpc | grep wlan0 | awk '{print $1}'` > /dev/null 2>&1
#  [ $? -eq 0 ] || return 2
#  rm -f /var/run/udhcpc.wlan0.pid
#  [ $? -eq 0 ] || return 2

#  sudo /sbin/udhcpc -b -i wlan0 -A 5 -x hostname:$(/bin/hostname) -p /var/run/udhcpc.wlan0.pid
#  [ $? -eq 0 ] || return 99

#  cd /mnt/UserData/ic-pcp-wifi-plus-server/pcp-scripts
#  ./wifi-plus-startup.sh


echo "{ \"beep\": \"boop\" }"