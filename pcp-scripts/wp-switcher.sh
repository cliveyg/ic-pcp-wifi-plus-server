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
    echo "[wp-switcher.sh] TO WAP MODE" >> $LOG
    echo '{ "status": 202, "message": "Attempting to switch to wap" }'
    # turn wifi off
    echo "[wp-switcher.sh] TO WRITING WIFI off TO CONFIG" >> $LOG
    pcp_write_var_to_config WIFI "off"
    echo "[wp-switcher.sh] STOPPING WIFI" >> $LOG
    /usr/local/etc/init.d/wifi wlan0 stop
    # get all wap stuff set up
    echo "[wp-switcher.sh] TO WRITING APMODE yes TO CONFIG" >> $LOG
    pcp_write_var_to_config APMODE "yes"

    echo "[wp-switcher.sh] LOADING pcp-apmode.tcz" >> $LOG
    sudo -u tc pcp-load -i pcp-apmode.tcz
    echo "[wp-switcher.sh] starting /usr/local/etc/init.d/pcp-apmode" >> $LOG
    sudo /usr/local/etc/init.d/pcp-apmode start
    sleep 2

    echo "[wp-switcher.sh] copying and permissions for pcp_hosts" >> $LOG

    #[ ! $(sudo ./wp-test.sh) ] && $(echo "[wp-switcher.sh] Bad JuJu" >> $LOG) && exit 1

    cp /mnt/UserData/industrialcool-pcp-wifi-plus/confs/pcp_hosts /usr/local/etc/pcp/pcp_hosts
    sudo chown root:root /usr/local/etc/pcp/pcp_hosts
    sudo chmod 644 /usr/local/etc/pcp/pcp_hosts
    echo "[wp-switcher.sh]  cat of my pcp_hosts $(/usr/local/etc/pcp/pcp_hosts)" >> LOG
    if [ ! $(sudo dnsmasq -C /usr/local/etc/pcp/dnsmasq.conf) ]; then
      pid=$(pidof dnsmasq)
      echo "[wp-switcher.sh] DNSMASQ PID: $(pidof dnsmasq)" >> $LOG
      sudo /usr/local/etc/init.d/pcp-apmode restart
      sleep 4
      echo "[wp-switcher.sh] DNSMASQ PID: $(pidof dnsmasq)" >> $LOG
      echo "[wp-switcher.sh]  cat after of pcp_hosts $(/usr/local/etc/pcp/pcp_hosts)" >> LOG
    fi

    pcp_backup "text"
    cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts
    echo "[wp-switcher.sh] Run  wifi-plus-startup.sh" >> $LOG
    ./wifi-plus-startup.sh
    echo "[wp-switcher.sh] wifiplus PID: $(pidof wifiplus)" >> $LOG
    #./wp-test.sh

    #if [ $(whoami) = "root" ]; then
    #  sudo -u tc echo "root sudoing echo as user tc"
    #else
    #  sudo echo "tc sudoing echo as normal"
    #fi
    #echo '{ "status": 202, "message": "Attempting to switch to wap" }'

  elif [ $arg1 = "towifi" ]; then
    echo "[wp-switcher.sh] TO WIFI MODE" >> $LOG
    [  ! -f "/usr/local/etc/pcp/wpa_supplicant.conf" ] && exit 1
    # before we can do this we need to check we have a wpa_supp file
    echo '{ "status": 202, "message": "Attempting to switch to wifi" }'
    # turn wifi on in config
    pcp_write_var_to_config WIFI "on"
    pcp_backup "text"
    # stop wap stuff
    pcp_write_var_to_config APMODE "no"
    pcp_backup "text"
    sudo /usr/local/etc/init.d/pcp-apmode stop
    sleep 2
    # start wifi
    pcp_wifi_load_wifi_extns
    pcp_wifi_load_wifi_firmware_extns

    pcp_backup "text"
    ./wp-wifi-refresh.sh
    #./wifi-plus-startup.sh
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


