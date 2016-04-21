#!/bin/bash

c=$(cat /sys/class/net/eth0/carrier)
echo "Initial cable state $c"
while :;do
	echo "waiting, cable state: $c"
	while [[ $c != $(cat /sys/class/net/eth0/carrier) ]];do
		echo "$c != $(cat /sys/class/net/eth0/carrier)"
		c=$(cat /sys/class/net/eth0/carrier)
		/home/pi/bin/showIp -topic "eth change to $c" -target example@gmail.com -pass passw0rd
	done
	sleep 10
done
