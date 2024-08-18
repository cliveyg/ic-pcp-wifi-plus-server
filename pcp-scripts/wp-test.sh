#!/bin/sh

#-----------------------------------------------------------------------------#
# wp-test.sh                                                           #
#                                                                             #
#                                                                             #
#                                                                             #
#                                                                             #
#-----------------------------------------------------------------------------#

set -a
. /var/www/.env
set +a

n=$1
LOG=/var/log/wifiplus.log

#-------------------------------- subroutines --------------------------------#

    cp /mnt/UserData/industrialcool-pcp-wifi-plus/confs/pcp_hosts /usr/local/etc/pcp/pcp_hosts
    sudo chown root:root /usr/local/etc/pcp/pcp_hosts
    sudo chmod 644 /usr/local/etc/pcp/pcp_hosts
    if [ ! $(sudo dnsmasq -C /usr/local/etc/pcp/dnsmasq.conf) ]; then

      if [ $(pidof dnsmasq) ]; then
        pid=$(pidof dnsmasq)
        echo "[1] DNSMASQ PID: $(pidof dnsmasq)" >> $LOG
        if [ $(sudo kill -9 $pid) ]; then
          echo "Killed dnsmasq process" >> $LOG
          sleep n+2
          sudo dnsmasq -C /usr/local/etc/pcp/dnsmasq.conf
          echo "Create new process using new pcp_hosts file" >> $LOG
          sleep 2
          echo "[2] DNSMASQ PID: $(pidof dnsmasq)" >> $LOG
        fi
      fi
      if [ ! $(pidof dnsmasq) ]; then
        sleep 2
        [ $((n)) -eq 5 ] && exit 1
        nmb=n+1
        ./wp-test.sh $nmb
      fi
    fi