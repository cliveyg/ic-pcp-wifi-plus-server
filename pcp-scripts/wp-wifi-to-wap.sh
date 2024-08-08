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

. /var/www/cgi-bin/pcp-functions
. /var/www/cgi-bin/pcp-wifi-functions

# TODO:
# for some reason this particular script doesn't get the log location
# from the .env vars. all other envs appear without a problem...
# other scripts that use the same mechanism and have the same permissions
# are able to see the log location - aaargh. will take a deeper look after
# the backend is feature complete
LOGGY=/var/log/wifiplus.log

if [ $DBUG -eq 1 ]; then

  if [ -f $LOGGY ]; then
    echo "[wp-wifi-to-wap.sh] ------------ woop ------------------" >> $LOGGY
    #echo "[wp-wifi-to-wap.sh] ENVs are [$(printenv)]" >> $LOGGY
    pcp_config_file
    pcp_read_config
    echo -n"[wp-wifi-to-wap.sh] PCPCFG IS [$($PCPCFG)]" >> $LOGGY
    $PCPCFG
  else
    sudo touch $LOGGY
    echo "[wp-wifi-to-wap.sh] ------------------------------" >> $LOGGY
    echo "[wp-wifi-to-wap.sh] DBUG IS [$DBUG]" >> $LOGGY
  fi


#else
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

#echo "{ \"beep\": \"boop\", \"yarp\": \"narp\" }"
echo '{ "beep": "boop", "yarp": "narp" }'