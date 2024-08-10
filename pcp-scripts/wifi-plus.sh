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
#arg2=$3
#arg3=$4

# ---------------------- subroutines ---------------------- #

wp_pcp_config() {
  if [ "$arg1" = "read" ]; then
    echo "$(cat $PCPCFG)"
  elif [ "$arg1" = "edit" ]; then
    echo '{ "error": "not yet implemented" }'
  else
    echo '{ "error": "unknown verb" }'
  fi

}

wp_picore_details() {
  printf "{\"picore_version\": \"%s\", " $(pcp_picore_version)
  printf "\"picoreplayer_version\": \"%s\", " $(pcp_picoreplayer_version)
  printf "\"squeezelite_version\": \"%s\", " $(pcp_squeezelite_version)
  printf "\"linux_release\": \"%s\"}" $(pcp_linux_release)
}

wp_status() {
  if [ $DBUG -eq 1 ]; then
    echo "[wifi-plus.sh] wp_status : Debug is on. Successful write to logfile" >> $LOG
  fi
  echo "$arg1"
}

wp_test() {
  #sudo -u tc printf '{ "message": "I am testy [%s]"}' $(whoami)
  sudo -u tc echo "sudoing echo"
  #if [ $(whoami) = "root" ]; then
  #  sudo -u tc printf '{ "message": "I am root [%s]"}' $(whoami)
  #else
  #  printf '{ "message": "I am tc [%s]"}' $(whoami)
  #fi
}

wp_wap_add() {

  if [ $DBUG -eq 1 ]; then
    echo "[wifi-plus.sh] wp_wap_add : Attempting to add apmode" >> $LOG
    printf "WHOAMI %s" $(whoami) >> $LOG
  fi

	pcp-load -r $PCP_REPO -w pcp-apmode.tcz 2>&1

	if [ -f $TCEMNT/tce/optional/pcp-apmode.tcz ]; then

		pcp-load -i firmware-atheros.tcz
		pcp-load -i firmware-brcmwifi.tcz
		pcp-load -i firmware-mediatek.tcz
		pcp-load -i firmware-ralinkwifi.tcz
		pcp-load -i firmware-rtlwifi.tcz
		pcp-load -i firmware-rpi-wifi.tcz
		pcp-load -i pcp-apmode.tcz

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

	tce-audit builddb
	tce-audit delete pcp-apmode.tcz

	#sed -i '/firmware-atheros.tcz/d' $ONBOOTLST
	#sed -i '/firmware-brcmwifi.tcz/d' $ONBOOTLST
	#sed -i '/firmware-mediatek.tcz/d' $ONBOOTLST
	#sed -i '/firmware-rpi-wifi.tcz/d' $ONBOOTLST
	#sed -i '/firmware-ralinkwifi.tcz/d' $ONBOOTLST
	#sed -i '/firmware-rtlwifi.tcz/d' $ONBOOTLST
	sed -i '/pcp-apmode.tcz/d' $ONBOOTLST

	rm -f $APMODECONF >/dev/null 2>&1
	rm -f $HOSTAPDCONF >/dev/null 2>&1
	rm -f $DNSMASQCONF >/dev/null 2>&1
	rm -f /usr/local/etc/pcp/pcp_hosts >/dev/null 2>&1

	pcp_backup "text"

	[ $DBUG -eq 1 ] && echo "[wifi-plus.sh] wp_wap_remove: Removed apmode" >> $LOG

  echo "{ \"soup\": \"boop\" }"

}

# ---------------------- main program ---------------------- #

case $subroutine in
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