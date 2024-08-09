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

if [ -f $TCEMNT/tce/optional/pcp-apmode.tcz ]; then

  if [ $DBUG -eq 1 ]; then

    if [ -f $LOGGY ]; then
      echo "[wp-wifi-to-wap.sh] ------------------------------------" >> $LOGGY
      #pcp_write_var_to_config USER_COMMAND_3 "echo 'boop'"
      #pcp_save_to_config
      #pcp_backup "text"
      if [ $(pcp_write_var_to_config USER_COMMAND_3 "echo 'boop'") ] && [ $(pcp_save_to_config) ] && [ $(pcp_backup "text") ]; then
        echo '{ "status": 200, "message": "var written to config successfully" }'
      else
        echo '{ "status": 500, "message": "failed to save to pcp config and backup" }'
      fi

      #echo "[wp-wifi-to-wap.sh] ENVs are [$(printenv)]" >> $LOGGY
      #$(pcp_config_file)
      #$(pcp_read_config)
      #[ cd /var/www/cgi-bin ] &&  pcp_picore_version
      #echo -n"[wp-wifi-to-wap.sh] PCPCFG IS [$($PCPCFG)]" >> $LOGGY
      #echo "$(pcp_picore_version)" >> $LOGGY
      #echo "$PCPCFG" >> $LOGGY
      #echo "$(pcp_config_file)"
      #echo "$(. $PCPCFG)"
      #echo "$(pcp_variables)"
      #echo "$(cat /etc/httpd.conf)" >> $LOGGY
    else
      sudo touch $LOGGY
      echo "[wp-wifi-to-wap.sh] ------------------------------" >> $LOGGY
      echo "[wp-wifi-to-wap.sh] DBUG IS [$DBUG]" >> $LOGGY
    fi


  else
    # turning wifi off
    echo '{ "status": 501, "message": "not implemented yet [1]" }'
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
  fi

else
  echo '{ "status": 404, "message": "ap mode not installed" }'
fi
