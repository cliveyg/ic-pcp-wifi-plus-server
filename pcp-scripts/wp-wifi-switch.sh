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

  if [ ! -f $LOG ]; then
    sudo touch $LOG
  fi
  echo "[wp-wifi-switch.sh] ------------------------------" >> $LOG
  echo -n "[wp-wifi-switch.sh] whoami: " >> $LOG
  echo $(whoami) >> $LOG
  echo "[wp-wifi-switch.sh] SSID is [$ssid]" >> $LOG
  echo "[wp-wifi-switch.sh] Pass is [$pass]" >> $LOG
  #sudo cp /usr/local/etc/pcp/wpa_supplicant.conf /usr/local/etc/pcp/wpa_supplicant.conf~
  #sudo cp /mnt/UserData/industrialcool-pcp-wifi-plus/confs/wpa_supplicant.conf /usr/local/etc/pcp/wpa_supplicant.conf
  #sudo sed -i "s/90909090909090909090909090909/$ssid/g" /usr/local/etc/pcp/wpa_supplicant.conf
  #sudo sed -i "s/\+\+/$pass/g" /usr/local/etc/pcp/wpa_supplicant.conf
  #sudo chown root:root /usr/local/etc/pcp/wpa_supplicant.conf

  sleep 2
  # backup stuff
  echo -n "[wp-wifi-switch.sh] backup status: " >> $LOG
  if wp_backup; then
    echo "success!" >> $LOG
    echo '{ "status": 200, "message": "the good ting" }'
  else
    echo "fail :(" >> $LOG
    echo '{ "status": 500, "message": "bad stuff" }'
  fi

else
  echo '{ "status": 404, "message": "no loggy" }'

fi



#echo "{ \"beep\": \"boop\" }"