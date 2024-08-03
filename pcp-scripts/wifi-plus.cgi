#!/bin/sh

# --------------------------------------------------------------- #
# wifi-plus.sh - called by the wifi-plus binary to interface with #
#                picoreplayer subroutines.                        #
# --------------------------------------------------------------- #

. pcp-functions
. pcp-wifi-functions

subroutine=$1
arg1=$2
arg2=$3
arg3=$4
arg4=$5

case "subroutine" in
	wp_status)
	  wp_status
	;;
	*)
	  echo "404"
	;;
esac

wp_status() {
  echo "$arg1"
}