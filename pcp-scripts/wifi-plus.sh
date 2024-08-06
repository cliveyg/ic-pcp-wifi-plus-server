#!/bin/sh

# --------------------------------------------------------------- #
# wifi-plus.sh - called by the wifi-plus binary to interface with #
#                picoreplayer subroutines.                        #
# --------------------------------------------------------------- #

# getting env settings from .env
set -a
. /var/www/.env
set +a

. pcp-functions
. pcp-wifi-functions

subroutine=$1
arg1=$2
#arg2=$3
#arg3=$4

# ---------------------- subroutines ---------------------- #

#wp_wifi_restart() {

  #"/usr/local/etc/init.d/wifi wlan0 stop;"
  #"sleep 3; /usr/local/etc/init.d/wifi wlan0 start;"
  #"sleep 10; /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts/wifi-plus-startup.sh;"

  #echo $arg3 | base64 --decode

  #nohup $(echo $arg3 | base64 --decode) > /var/log/wifiplus.log 2>&1 &
  #nohup $(echo $arg2 | base64 --decode) > /var/log/wifiplus.log 2>&1 &
  #nohup $(echo $arg1 | base64 --decode) > /var/log/wifiplus.log 2>&1 &
  #echo $(sleep 10; /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts/wifi-plus-startup.sh;)

  #nohup "$(sleep 1; cd /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts/; ./wifi-plus-startup.sh)" > /www/log/wifiplus.log 2>&1 &
#}

wp_picore_details() {
  printf "\"picore_version\": \"%s\", " $(pcp_picore_version)
  printf "\"picoreplayer_version\": \"%s\", " $(pcp_picoreplayer_version)
  printf "\"squeezelite_version\": \"%s\", " $(pcp_squeezelite_version)
  printf "\"linux_release\": \"%s\"" $(pcp_linux_release)
}

wp_status() {
  echo "$arg1"
}

wp_test() {
  pcp_set_coloured_text
  echo "Able to call pcp functions"
}


# ---------------------- main program ---------------------- #

case $subroutine in
  wp_general_hup)
    wp_general_hup
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
  *)
    echo "$subroutine"
  ;;
esac