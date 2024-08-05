#!/bin/sh

# --------------------------------------------------------------- #
# wifi-plus.sh - called by the wifi-plus binary to interface with #
#                picoreplayer subroutines.                        #
# --------------------------------------------------------------- #

. pcp-functions
. pcp-wifi-functions

subroutine=$1
arg1=$2

#arg4=$5

# ---------------------- subroutines ---------------------- #

wp_general_hup() {
  for i in $(seq 2 4);
  do
    hup $i > /dev/null 2>&1 &
  done
  echo "hupped"
  #cmmnd1=$arg1
  #nohup cmmnd1 > /dev/null 2>&1 &
  #cmmnd2=$arg2
  #nohup cmmnd2 > /dev/null 2>&1 &
}

wp_picore_details() {
  printf "\"picore_version\": \"%s\", " $(pcp_picore_version)
  printf "\"picoreplayer_version\": \"%s\", " $(pcp_picoreplayer_version)
  printf "\"squeezelite_version\": \"%s\", " $(pcp_squeezelite_version)
  printf "\"linux_release\": \"%s\"" $(pcp_linux_release)
}

wp_status() {
  echo $arg1
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