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
        echo "[wp-test.sh] DNSMASQ PID: $(pidof dnsmasq)" >> $LOG
        sudo kill -9 $pid
        if [ $? ]; then
          echo "[wp-test.sh] Killed dnsmasq process" >> $LOG
          sleep n+2
          sudo dnsmasq -C /usr/local/etc/pcp/dnsmasq.conf
          echo "[wp-test.sh] Create new process using new pcp_hosts file" >> $LOG
          sleep 2
          echo "[wp-test.sh] DNSMASQ PID: $(pidof dnsmasq)" >> $LOG
        fi
      fi
      if [ ! $(pidof dnsmasq) ]; then
        echo "[wp-test.sh] n is [$n]" >> LOG
        sleep 2
        [ $((n)) -eq 5 ] && exit 1
        nmb=n+1
        ./wp-test.sh $nmb
      fi
    fi