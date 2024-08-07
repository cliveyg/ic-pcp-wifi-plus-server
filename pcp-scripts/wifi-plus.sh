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

wp_wap_install() {

	sudo -u tc pcp-load -r $PCP_REPO -w pcp-apmode.tcz 2>&1

	if [ -f $TCEMNT/tce/optional/pcp-apmode.tcz ]; then
		sudo -u tc pcp-load -i firmware-atheros.tcz
		sudo -u tc pcp-load -i firmware-brcmwifi.tcz
		sudo -u tc pcp-load -i firmware-mediatek.tcz

		sudo -u tc pcp-load -i firmware-ralinkwifi.tcz
		sudo -u tc pcp-load -i firmware-rtlwifi.tcz
		sudo -u tc pcp-load -i firmware-rpi-wifi.tcz
		sudo -u tc pcp-load -i pcp-apmode.tcz
		pcp_wifi_update_wifi_onbootlst
		pcp_wifi_update_onbootlst "add" "pcp-apmode.tcz"
	fi
  # echo "{\"progress\": \"$ap\"}"
}


# ---------------------- main program ---------------------- #

case $subroutine in
  wp_picore_details)
    wp_picore_details
  ;;
  wp_status)
    wp_status
  ;;
  wp_test)
    wp_test
  ;;
  wp_wap_install)
    wp_wap_install
  ;;
  *)
    echo "$subroutine"
  ;;
esac