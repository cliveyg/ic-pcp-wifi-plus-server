#!/bin/sh

# --------------------------------------------------------------- #
# wifi-plus.sh - called by the wifi-plus binary to interface with #
#                picoreplayer subroutines.                        #
# --------------------------------------------------------------- #

# getting env settings from .env
set -a
. /var/www/.env
set +a

. /var/www/cgi-bin/pcp-functions
. /var/www/cgi-bin/pcp-wifi-functions

LOG=/var/log/wifiplus.log

subroutine=$1
arg1=$2

#-------------------------------- subroutines --------------------------------#

wp_pcp_config() {
  if [ "$arg1" = "read" ]; then
    echo "$(cat $PCPCFG)"
  elif [ "$arg1" = "edit" ]; then
    echo '{ "error": "not yet implemented" }'
  else
    echo '{ "error": "unknown verb" }'
  fi

}

#-----------------------------------------------------------------------------#

wp_picore_details() {
  printf "{\"picore_version\": \"%s\", " $(pcp_picore_version)
  printf "\"picoreplayer_version\": \"%s\", " $(pcp_picoreplayer_version)
  printf "\"squeezelite_version\": \"%s\", " $(pcp_squeezelite_version)
  printf "\"linux_release\": \"%s\"}" $(pcp_linux_release)
}

#-----------------------------------------------------------------------------#

wp_status() {
  if [ $DBUG -eq 1 ]; then
    echo "[wifi-plus.sh] wp_status : Debug is on. Successful write to logfile" >> $LOG
  fi
  echo "$arg1"
}

#-----------------------------------------------------------------------------#

wp_test() {

  if [ $(whoami) = "root" ]; then
    sudo -u tc echo "root sudoing echo as user tc"
  else
    sudo echo "tc sudoing echo as normal"
  fi

}

#-----------------------------------------------------------------------------#

wp_edit_wap_config() {
  if [ $DBUG -eq 1 ]; then
    echo "[wifi-plus.sh] wp_edit_wap_config" >> $LOG
    printf "arg1 is [%s]" $arg1 >> $LOG
  fi
  echo '{ "status": 200, "message": "success"}'
}

#-----------------------------------------------------------------------------#

wp_fetch_wap_config() {
  if [ $DBUG -eq 1 ]; then
    echo "[wifi-plus.sh] wp_fetch_config" >> $LOG
  fi

  if [ -f /usr/local/etc/pcp/hostapd.conf ]; then
    filename="/usr/local/etc/pcp/hostapd.conf"
  elif [ -f /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts/hostapd.conf ]; then
    filename=/mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts/hostapd.conf
  else
    echo '{ "status": 404, "error": "hostapd.conf file not found" }'
    exit 0
  fi

  ssid_line=$(grep ssid $filename | head -1)
  pass_line=$(grep wpa_passphrase $filename)
  country_code_line=$(grep country_code $filename)
  channel_line=$(grep channel $filename)

  ssid=$(echo $ssid_line | sed 's/ssid=//g')
  pass=$(echo $pass_line | sed 's/wpa_passphrase=//g')
  country_code=$(echo $country_code_line | sed 's/country_code=//g')
  channel=$(echo $channel_line | sed 's/channel=//g')

  echo '{ "ssid": "'$ssid'", "ap_ip_address": "'$AP_IP'", "password": "''"country_code": "'$country_code'", "channel": '$channel'}'
}

#-----------------------------------------------------------------------------#

wp_wap_add() {

  if [ $DBUG -eq 1 ]; then
    echo "[wifi-plus.sh] wp_wap_add : Attempting to add apmode" >> $LOG
    printf "WHOAMI %s" $(whoami) >> $LOG
  fi

  # not sure why but sometimes this script runs as root and sometimes as tc
  # hence this check here. very cludgy but it works
  # TODO: work out why script runs under different users sometimes
  if [ $(whoami) = "root" ]; then
    echo "root user running pcp-load repo as tc" >> $LOG
    sudo -u tc pcp-load -r $PCP_REPO -w pcp-apmode.tcz 2>&1
  else
    echo "tc user running pcp-load repo" >> $LOG
    pcp-load -r $PCP_REPO -w pcp-apmode.tcz 2>&1
  fi

	if [ -f $TCEMNT/tce/optional/pcp-apmode.tcz ]; then

    if [ $(whoami) = "root" ]; then
      sudo -u tc pcp-load -i firmware-atheros.tcz
      sudo -u tc pcp-load -i firmware-brcmwifi.tcz
      sudo -u tc pcp-load -i firmware-mediatek.tcz
      sudo -u tc pcp-load -i firmware-ralinkwifi.tcz
      sudo -u tc pcp-load -i firmware-rtlwifi.tcz
      sudo -u tc pcp-load -i firmware-rpi-wifi.tcz
      sudo -u tc pcp-load -i pcp-apmode.tcz
		else
      pcp-load -i firmware-atheros.tcz
      pcp-load -i firmware-brcmwifi.tcz
      pcp-load -i firmware-mediatek.tcz
      pcp-load -i firmware-ralinkwifi.tcz
      pcp-load -i firmware-rtlwifi.tcz
      pcp-load -i firmware-rpi-wifi.tcz
      pcp-load -i pcp-apmode.tcz
		fi

		pcp_wifi_update_wifi_onbootlst
		pcp_wifi_update_onbootlst "add" "pcp-apmode.tcz"

    pcp_write_var_to_config APMODE "no"
    pcp_write_var_to_config AP_IP "10.10.10.1"
    pcp_save_to_config
    pcp_backup "text"

    sudo chown tc:staff /usr/local/etc/pcp/dnsmasq.conf
    cp /mnt/UserData/industrialcool-pcp-wifi-plus/confs/hostapd.conf /usr/local/etc/pcp/hostapd.conf
    sudo chown tc:staff /usr/local/etc/pcp/hostapd.conf
  else
    echo '{"status": "500", "message": "Failed to download ap mode file."}'
	fi

	[ $DBUG -eq 1 ] && echo "[wifi-plus.sh] wp_wap_add: Added apmode" >> $LOG

  echo "{ \"boop\": \"soup\" }"

}

#-----------------------------------------------------------------------------#

wp_wap_remove() {

  if [ $DBUG -eq 1 ]; then
      if [ ! -f $LOG ]; then
        sudo touch $LOG
      fi
    echo "[wifi-plus.sh] wp_wap_remove: Attempting to remove apmode" >> $LOG
  fi

	pcp_write_var_to_config APMODE "no"
	pcp_save_to_config

	sudo /usr/local/etc/init.d/pcp-apmode stop >/dev/null 2>&1

  if [ $(whoami) = "root" ]; then
	  sudo -u tc tce-audit builddb
	  sudo -u tc tce-audit delete pcp-apmode.tcz
  else
	  tce-audit builddb
	  tce-audit delete pcp-apmode.tcz
	fi

	sed -i '/pcp-apmode.tcz/d' $ONBOOTLST

	rm -f $APMODECONF >/dev/null 2>&1
	rm -f $HOSTAPDCONF >/dev/null 2>&1
	rm -f $DNSMASQCONF >/dev/null 2>&1
	rm -f /usr/local/etc/pcp/pcp_hosts >/dev/null 2>&1

	pcp_backup "text"

	[ $DBUG -eq 1 ] && echo "[wifi-plus.sh] wp_wap_remove: Removed apmode" >> $LOG

  echo "{ \"soup\": \"boop\" }"

}

#------------------------------- main program --------------------------------#

case $subroutine in
  wp_edit_wap_config)
    wp_edit_wap_config
  ;;
  wp_fetch_wap_config)
    wp_fetch_wap_config
  ;;
  wp_pcp_config)
    wp_pcp_config
  ;;
  wp_picore_details)
    wp_picore_details
  ;;
  wp_status)
    wp_status
  ;;
  wp_test)
    wp_test
  ;;
  wp_wap_add)
    wp_wap_add
  ;;
  wp_wap_remove)
    wp_wap_remove
  ;;
  *)
    echo "$subroutine not found"
  ;;
esac