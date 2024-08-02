#!/bin/sh

#------------------------------ procs -----------------------------#

pcp_apmode_status() {
        sudo /usr/local/etc/init.d/pcp-apmode status >/dev/null 2>&1
        echo $?
}

pcp_wifi_status() {
        echo "$(sudo /usr/local/etc/init.d/wifi wlan0 status)"
}

#------------------------------ main ------------------------------#
sleep 5
echo '[INFO] Checking wifi status... '
output=$(pcp_wifi_status)
echo '[INFO] ----------'
echo "$output"
echo '[INFO] ----- boop -----'
case "$output" in
        *"wpa_supplicant running"*)
                # Do stuff
                echo 'wpa_supplicant running!'
        ;;
        *)
                # everything else
                echo 'meep'
        ;;
esac
echo '[INFO] ----------'
myssid=$(iwgetid -r)
echo '[INFO] myssid'
echo "$myssid"
echo '[PARTY!] MEEEEEEEPPPP'