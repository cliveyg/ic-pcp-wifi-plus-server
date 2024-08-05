#!/bin/sh

# --------------------------------------------------------------- #
# wifi-plus.sh - called by the wifi-plus binary to interface with #
#                picoreplayer subroutines.                        #
# --------------------------------------------------------------- #

# getting env settings from .env
set -a
. /mnt/UserData/industrialcool-pcp-wifi-plus/.env
set +a

. pcp-functions
. pcp-wifi-functions

subroutine=$1
arg1=$2
arg2=$3
arg3=$4

# ---------------------- subroutines ---------------------- #

wp_general_hup() {
  #echo "--=-=-= arg 1 =-=-=--"
  echo $arg1
  #echo "--=-=-= arg 2 =-=-=--"
  echo $arg2
  #echo "--=-=-= arg 3 =-=-=--"
  echo $arg3

  echo $arg1 | base64 --decode
  echo $arg2 | base64 --decode
  echo $arg3 | base64 --decode
  #nohup $(echo -n $arg1 | base64 --decode) > /var/log/wifiplus.log 2>&1 &
  #nohup $(echo -n $arg2 | base64 --decode) > /var/log/wifiplus.log 2>&1 &
  #nohup $(echo -n $arg3 | base64 --decode) > /var/log/wifiplus.log 2>&1 &

  #for i in $(seq 2 4);
  #do
  #  if [ "$DBUG" -eq 1 ]; then
  #    nohup $(echo -n "$i" | base64 --decode) > /var/log/wifiplus.log 2>&1 &
  #  else
  #    nohup $(echo -n "$i" | base64 --decode) > /dev/null 2>&1 &
  #  fi
  #done
  #echo "nohupped"
}

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