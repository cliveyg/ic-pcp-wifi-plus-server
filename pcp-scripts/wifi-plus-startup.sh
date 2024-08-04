#!/bin/sh

echo "-----------------------------------------------------------------------"
echo "Starting wifi-plus-startup..."
echo "Copying go binary and script files to web folders..."
sudo chmod 777 /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts/wifi-plus.sh
if sudo cp /mnt/UserData/industrialcool-pcp-wifi-plus/pcp-scripts/wifi-plus.sh /var/www/cgi-bin/wifi-plus.sh; then
  echo "Successfully copied wifi-plus shell file to cgi-bin"
else
  echo "Unable to copy shell file to cgi-bin."
  echo "Exiting..."
  exit 1
fi

if sudo cp /mnt/UserData/industrialcool-pcp-wifi-plus/wifiplus /var/www/wifiplus &&
   sudo cp /mnt/UserData/industrialcool-pcp-wifi-plus/.env /var/www/.env; then

  echo "Successfully copied wifi-plus binary to webroot"
  echo "Attempting to start binary..."

  cd /var/www/ && nohup ./wifiplus > /dev/null 2>&1 &

  printf "Binary started successfully.\nListening on port 8020...\n"
  echo "Testing connection..."
  sleep 3
  rc=$(curl -s -o /dev/null -w "%{http_code}" http://icplayer.local:8020/status)
  if [ $rc = "200" ]; then
    echo "[$rc OK] API up and running."
    exit 0
  else
    echo "Unable to connect to API successfully."
    echo "Status code is [$rc]"
    echo "Exiting..."
  fi

else
  echo "Unable to copy binary to webroot."
  echo "Exiting..."
fi

exit 1