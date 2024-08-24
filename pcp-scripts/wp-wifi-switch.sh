#!/bin/sh

#-----------------------------------------------------------------------------#
# wp-wifi-switch.sh                                                           #
#                                                                             #
#                                                                             #
#-----------------------------------------------------------------------------#

set -a
. /var/www/.env
set +a

ssid=$1
pass=$2

#-------------------------------- subroutines --------------------------------#

wp_backup() {
  # basically a copy of pcp_backup() without pcp bits.

  # delete any previous backup_done file
  [ -e /tmp/backup_done ] && sudo rm -f /tmp/backup_done >/dev/null 2>&1

  # do a backup - filetool.sh backs up files in .filetool.lst
  sudo filetool.sh -b >/dev/null 2>&1
  sync > /dev/null 2>&1 >/dev/null 2>&1

  # if backup_status file exists and is non-zero in size, then an error has occurred
  if [ -s /tmp/backup_status ]; then
    return 1
  fi

  # if backup_done exists, then the backup was successful
  if [ -f /tmp/backup_done ]; then
    return 0
  else
    return 1
  fi
}

#------------------------------- main program --------------------------------#

if [ $DBUG -eq 1 ]; then
  LOG=/var/log/wifiplus.log
  if [ ! -f $LOG ]; then
    sudo touch $LOG
  fi
else
  LOG=/dev/null
fi

echo "[wp-wifi-switch.sh] ------------------------------" >> $LOG
echo -n "[wp-wifi-switch.sh] whoami: " >> $LOG
echo $(whoami) >> $LOG
echo "[wp-wifi-switch.sh] SSID is [$ssid]" >> $LOG
echo "[wp-wifi-switch.sh] Pass is [$pass]" >> $LOG
sudo cp /usr/local/etc/pcp/wpa_supplicant.conf /usr/local/etc/pcp/wpa_supplicant.conf~
sudo cp /mnt/UserData/ic-pcp-wifi-plus-server/confs/wpa_supplicant.conf /usr/local/etc/pcp/wpa_supplicant.conf
sudo sed -i "s/90909090909090909090c909090909/$ssid/g" /usr/local/etc/pcp/wpa_supplicant.conf
sudo sed -i "s/\+\+/$pass/g" /usr/local/etc/pcp/wpa_supplicant.conf
sudo chown root:root /usr/local/etc/pcp/wpa_supplicant.conf

echo "[wp-wifi-switch.sh] Attempting to switch wifi" >> $LOG
echo "[wp-wifi-switch.sh] Stopping current wifi" >> $LOG

sudo /usr/local/etc/init.d/wifi wlan0 stop
if [ $? -eq 0 ]; then
  echo "[wp-wifi-switch.sh] Current wifi stopped" >> $LOG
  sleep 1
  wpa_cli -i wlan0 reconfigure
  sleep 3
  iwconfig wlan0 | grep "Frequency"
  if [ $? -ne 0 ]; then
    echo "[wp-wifi-switch.sh] Failed to switch wifi networks " >> $LOG
    echo "[wp-wifi-switch.sh] Switching back... " >> $LOG
    sudo cp /usr/local/etc/pcp/wpa_supplicant.conf~ /usr/local/etc/pcp/wpa_supplicant.conf
    sudo chown root:root /usr/local/etc/pcp/wpa_supplicant.conf

    wpa_cli -i wlan0 reconfigure
    sleep 3
    iwconfig wlan0 | grep "Frequency"
    if [ $? -ne 0 ]; then
      echo "[wp-wifi-switch.sh] Failed to switch back!" >> $LOG
      echo '{ "status": 500, "message": "Failed to switch back" }'
    else
      echo "[wp-wifi-switch.sh] Switched back to old wifi settings" >> $LOG
      echo '{ "status": 400, "message": "Switched back to old wifi settings" }'
    fi
  else
    # backup stuff
    echo -n "[wp-wifi-switch.sh] backup status: " >> $LOG
    if wp_backup; then
      echo "success!" >> $LOG
      echo '{ "status": 202, "message": "the good ting" }'
    else
      echo "fail :(" >> $LOG
      echo '{ "status": 500, "message": "bad stuff" }'
    fi
  fi
else
       echo "[wp-wifi-switch.sh] Unable to stop running wifi :(" >> $LOG
       echo '{ "status": 500, "message": "Unable to stop running wifi" }'
fi



#echo "{ \"beep\": \"boop\" }"