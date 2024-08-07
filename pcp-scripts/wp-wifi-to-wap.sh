#!/bin/sh

#-----------------------------------------------------------------------------#
# wp-wifi-to-wap.sh                                                           #
#                                                                             #
#                                                                             #
#                                                                             #
#                                                                             #
#-----------------------------------------------------------------------------#

set -a
. /var/www/.env
set +a

LOG=$LOGFILE

if [ $DBUG -eq 1 ]; then

  if [ -f $LOG ]; then
    echo "[wp-wifi-to-wap.sh] ------------------------------" >> $LOG
    echo "[wp-wifi-to-wap.sh] DBUG IS [$DBUG]" >> $LOG
  else
    sudo touch $LOG
    echo "[wp-wifi-to-wap.sh] ------------------------------" >> $LOG
    echo "[wp-wifi-to-wap.sh] DBUG IS [$DBUG]" >> $LOG
  fi

  pcp_read_config
  echo "[wp-wifi-to-wap.sh] PCPCFG IS \n$PCPCFG" >> $LOG

else
  # turning wifi off
 # export WIFI="off"
 # /usr/local/etc/init.d/wifi wlan0 stop
	#pcp_wifi_unload_wifi_extns "text"
	#pcp_wifi_unload_wifi_firmware_extns "text"
  #pcp_save_to_config
  #pcp_backup "text"
  # turning wap on
	#if [ ! -x /usr/local/etc/init.d/pcp-apmode ]; then
	#	sudo -u tc pcp-load -i pcp-apmode.tcz
	#fi
	#sudo /usr/local/etc/init.d/pcp-apmode start
fi

echo '{ "beep": "boop", "yarp": "narp" }'