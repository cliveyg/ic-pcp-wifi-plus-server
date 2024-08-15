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
  # basically a copy of pcp_backup() without pcp bits

  # delete any previous backup_done file
  [ -e /tmp/backup_done ] && sudo rm -f /tmp/backup_done

  # do a backup - filetool.sh backs up files in .filetool.lst
  sudo filetool.sh -b
  sync > /dev/null 2>&1

  # if backup_status file exists and is non-zero in size, then an error has occurred
  if [ -s /tmp/backup_status ]; then
    echo "BAD JUJU" >> $LOG
    return 1
  fi

  # if backup_done exists, then the backup was successful
  if [ -f /tmp/backup_done ]; then
    echo "Yarp" >> $LOG
    return 0
  else
    echo "Narp" >> $LOG
    return 1
  fi
}

#------------------------------- main program --------------------------------#


if [ $DBUG -eq 1 ]; then

  if [ ! -f $LOG ]; then
    sudo touch $LOG
  fi
  echo "[wp-wifi-switch.sh] ------------------------------" >> $LOG
  echo -n "[wp-wifi-switch.sh] WHOAMI: " >> $LOG
  echo $(whoami) >> $LOG
  echo "[wp-wifi-switch.sh] SSID is [$ssid]" >> $LOG
  echo "[wp-wifi-switch.sh] Pass is [$pass]" >> $LOG
  #sudo cp /usr/local/etc/pcp/wpa_supplicant.conf /usr/local/etc/pcp/wpa_supplicant.conf~
  #sudo cp /mnt/UserData/industrialcool-pcp-wifi-plus/confs/wpa_supplicant.conf /usr/local/etc/pcp/wpa_supplicant.conf
  #sudo sed -i "s/90909090909090909090909090909/$ssid/g" /usr/local/etc/pcp/wpa_supplicant.conf
  #sudo sed -i "s/\+\+/$pass/g" /usr/local/etc/pcp/wpa_supplicant.conf
  #sudo chown root:root /usr/local/etc/pcp/wpa_supplicant.conf

  # backup stuff
  echo -n "[wp-wifi-switch.sh] backup status: " >> $LOG
  if wp_backup; then
    echo "success!" >> $LOG
  else
    echo "fail :(" >> $LOG
  fi

  #sudo wpa_cli -i wlan0 reconfigure
  sleep 3
  echo '{ "status": 200, "message": "have loggy" }'

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