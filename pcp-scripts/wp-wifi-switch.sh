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
  echo "[wp-wifi-switch.sh] wpa_cli -i wlan0 reconfigure: " >> LOG
  echo "$(wpa_cli -i wlan0 reconfigure)" >> LOG
  sleep 3
  echo "[wp-wifi-switch.sh] after reconfiguring" >> LOG

  case "$(iwconfig wlan0)" in
    *Frequency*)
      echo "[wp-wifi-switch.sh] New wifi running" >> LOG
      # backup stuff
      echo -n "[wp-wifi-switch.sh] backup status: " >> $LOG
      if wp_backup; then
        echo "success!" >> $LOG
        echo '{ "status": 202, "message": "the good ting" }'
        return 0
      else
        echo "fail :(" >> $LOG
        echo '{ "status": 500, "message": "Unable to back up pcp" }'
        return 0
      fi
    ;;
    *)
      echo "[wp-wifi-switch.sh] Failed to switch wifi networks " >> $LOG
      echo "[wp-wifi-switch.sh] Switching back... " >> $LOG
      echo "$(sudo cp /usr/local/etc/pcp/wpa_supplicant.conf~ /usr/local/etc/pcp/wpa_supplicant.conf)" >> LOG
      echo "$(sudo chown root:root /usr/local/etc/pcp/wpa_supplicant.conf)" >> LOG

      echo "$(wpa_cli -i wlan0 reconfigure)" >> LOG
      sleep 3
      echo "$(sudo /usr/local/etc/init.d/wifi wlan0 restart)" >> LOG

      echo "[wp-wifi-switch.sh] Attempting to kill and restart udhcpc" >> $LOG
      sudo kill `ps | grep udhcpc | grep wlan0 | awk '{print $1}'` > /dev/null 2>&1
      if [ $? -eq 0 ]; then
        rm -f /var/run/udhcpc.wlan0.pid
        if [ $? -eq 0 ]; then
          echo "[wp-wifi-switch.sh] Killed udhcpc process"
        else
          echo "[wp-wifi-switch.sh] Failed to kill udhcpc process."
          echo '{ "status": 500, "message": "Failed to kill udhcpc process" }'
          return 0
        fi
      else
        echo "[wp-wifi-switch.sh] Failed to kill udhcpc process."
        echo '{ "status": 500, "message": "Failed to kill udhcpc process" }'
        return 0
      fi

      sudo /sbin/udhcpc -b -i wlan0 -A 5 -x hostname:$(/bin/hostname) -p /var/run/udhcpc.wlan0.pid
      if [ $? -ne 0 ]; then
        echo "[wp-wifi-switch.sh] Unable to restart udhcpc." >> $LOG
        echo '{ "status": 500, "message": "Unable to restart udhcpc" }'
        return 0
      else
        echo "[wp-wifi-switch.sh] Restart udhcpc [OK]" >> $LOG
      fi

      sleep 3
      iwconfig wlan0 | grep "Frequency"
      if [ $? -ne 0 ]; then
        echo "[wp-wifi-switch.sh] Failed to switch back!" >> $LOG
        echo '{ "status": 500, "message": "Failed to switch back" }'
        return 0
      else
        echo "[wp-wifi-switch.sh] Switched back to old wifi settings" >> $LOG
        echo '{ "status": 400, "message": "Switched back to old wifi settings" }'
        return 0
      fi
    ;;
  esac

else
       echo "[wp-wifi-switch.sh] Unable to stop running wifi :(" >> $LOG
       echo '{ "status": 500, "message": "Unable to stop running wifi" }'
       return 0
fi



#echo "{ \"beep\": \"boop\" }"