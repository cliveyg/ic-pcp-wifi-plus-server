#!/bin/sh

echo "-----------------------------------------------------------------------"
echo "Starting copy-to-www..."
echo "Copying go binary and script files to web folders..."
if sudo cp /mnt/UserData/wifi-settings/wifi-plus.sh /var/www/cgi-bin/wifi-plus.sh; then
  echo "Successfully copied wifi-plus shell file to cgi-bin"
else
  echo "Unable to copy shell file to cgi-bin."
  echo "Exiting..."
  exit 1
fi

if sudo cp /mnt/UserData/wifi-settings/wifiplus /var/www/wifiplus &&
   sudo cp /mnt/UserData/wifi-settings/.env /var/www/.env; then

  echo "Successfully copied wifi-plus binary to webroot"
  echo "Attempting to start binary..."

  echo "$(sudo /var/www/wifiplus > /dev/null 2>&1 &)"
  echo "Binary started successfully.\nListening on port 8020..."
  echo "Testing connection..."
  if curl -s -o /dev/null -w "%{http_code}" http://pcp.local:8020/status; then
    echo " API up and running."
    exit 0
  else
    echo "Unable to connect to API"
    echo "Exiting..."
  fi

else
  echo "Unable to copy binary to webroot."
  echo "Exiting..."
fi

exit 1