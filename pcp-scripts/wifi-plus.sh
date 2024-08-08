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

#LOG=$LOGFILE
#LOG="wifi-plus.sh.log"

subroutine=$1
arg1=$2
#arg2=$3
#arg3=$4

# ---------------------- subroutines ---------------------- #

wp_picore_details() {
  printf "{\"picore_version\": \"%s\", " $(pcp_picore_version)
  printf "\"picoreplayer_version\": \"%s\", " $(pcp_picoreplayer_version)
  printf "\"squeezelite_version\": \"%s\", " $(pcp_squeezelite_version)
  printf "\"linux_release\": \"%s\"}" $(pcp_linux_release)
}

wp_status() {
  echo "$arg1"
}

wp_test() {
  pcp_set_coloured_text
  echo "Able to call pcp functions"
}

wp_wap_add() {

  #echo "[wifi-plus.sh] wp_wap_add : ------------------------------" >> $LOG

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
    APMODE="no"
    AP_IP="10.10.10.1"
    pcp_save_to_config
    pcp_backup "text"
  else
    echo '{"status": "500", "message": "Failed to download ap mode file."}'
	fi
  echo "{ \"boop\": \"soup\" }"
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