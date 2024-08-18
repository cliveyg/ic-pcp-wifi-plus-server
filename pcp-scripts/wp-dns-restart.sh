#!/bin/sh

#-----------------------------------------------------------------------------#
# wp-dns-restart.sh                                                           #
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

cp /mnt/UserData/ic-pcp-wifi-plus-server/confs/pcp_hosts /usr/local/etc/pcp/pcp_hosts
sudo chown root:root /usr/local/etc/pcp/pcp_hosts
sudo chmod 644 /usr/local/etc/pcp/pcp_hosts

echo "[wp-dns-restart.sh] DNSMASQ before PID: $(pidof dnsmasq)" >> $LOG

if [ $(pidof dnsmasq) ]; then
  pid=$(pidof dnsmasq)
  echo "[wp-dns-restart.sh] DNSMASQ PID: $(pidof dnsmasq)" >> $LOG
  sudo kill -9 $pid
  if [ $? ]; then
    echo "[wp-dns-restart.sh] Killed dnsmasq process" >> $LOG
    sleep 4
    sudo dnsmasq -C /usr/local/etc/pcp/dnsmasq.conf
    echo "[wp-dns-restart.sh] Create new process using new pcp_hosts file" >> $LOG
    sleep 2
    echo "[wp-dns-restart.sh] DNSMASQ PID: $(pidof dnsmasq)" >> $LOG
  fi
fi
sleep 4
echo "[wp-dns-restart.sh] DNSMASQ after PID: $(pidof dnsmasq)" >> $LOG
echo "[wp-dns-restart.sh] pinging icplayer.local: " >> $LOG
ping -c 1 icplayer.local; echo $? >> $LOG
if [ ! $(ping -c 1 icplayer.local) ]; then
  echo '{ "status": 404, "message": "unable to ping icplayer.local" }'
else
  echo '{ "status": 404, "message": "pinged icplayer.local succssfully" }'
fi
echo "[wp-dns-restart.sh] ----------------------------" >> $LOG