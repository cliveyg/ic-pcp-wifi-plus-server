#!/bin/sh

#-----------------------------------------------------------------------------#
# wp-switcher.sh                                                           #
#                                                                             #
#                                                                             #
#                                                                             #
#                                                                             #
#-----------------------------------------------------------------------------#

set -a
. /var/www/.env
set +a

. /var/www/cgi-bin/pcp-functions
. /var/www/cgi-bin/pcp-wifi-functions

arg1=$1

# TODO:
# for some reason this particular script doesn't get the log location
# from the .env vars. all other envs appear without a problem...
# other scripts that use the same mechanism and have the same permissions
# are able to see the log location - aaargh. will take a deeper look after
# the backend is feature complete
LOG=/var/log/wifiplus.log


if [ $DBUG -eq 1 ]; then

  if [ ! -f $LOG ]; then
    sudo touch $LOG
  fi
    
  echo "[wp-switcher.sh] --------------- running --------------------" >> $LOG
  echo "[wp-switcher.sh] " >> $LOG

  if [ $arg1 = "towap" ]; then
    # get all wap stuff set up
    pcp_write_var_to_config APMODE "no"
    #sudo -u tc pcp-load -i pcp-apmode.tcz
    # turn wifi off

  elif [ $arg1 = "towifi" ]; then
    echo "[wp-switcher.sh] TO WAP MODE" >> $LOG
  else
    echo '{ "status": 400, "message": "action not valid" }'
  fi

else
  echo "no loggy"

fi


# turning wifi off
#echo '{ "status": 501, "message": "not implemented yet [1]" }'
#pcp_write_var_to_config WIFI "off"
#/usr/local/etc/init.d/wifi wlan0 stop
#pcp_wifi_unload_wifi_extns "text"
#pcp_wifi_unload_wifi_firmware_extns "text"
#pcp_save_to_config
#pcp_backup "text"
# turning wap on
#if [ ! -x /usr/local/etc/init.d/pcp-apmode ]; then
#  pcp-load -i pcp-apmode.tcz
#fi
#/usr/local/etc/init.d/pcp-apmode start


