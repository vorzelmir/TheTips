#!/bin/sh

default() {
    if [ -z "$NODE" ]; then NODE="alpine"; fi
    if [ -z "$MASK" ]; then MASK=24; fi
}
help() {
    echo "-n is a hostname"
    echo "-p is a subnet prefix"
    echo "-m is a subnet mask"
    echo "-g is a gateway"
}

if [ -z "$1" ]
then
    echo "missing args!"
    help
    exit 1
fi

while getopts "n:p:m:g:h" opt
do
case "$opt" in
   n) NODE="$OPTARG";;
   p) PREFIX="$OPTARG";;
   m) MASK="$OPTARG";;
   g) GATEWAY="$OPTARG";;
   \?) echo "Invalid arg: $OPTARG"
       help
       exit 1;;
   :) echo "Argument is required"
      help
      exit 1;;
   h) help;;
esac
done

default

#set hostname
echo "$NODE" > /etc/hostname
hostname -F /etc/hostname

#confg networking
if ! command -v ifquery 2>&1 >/dev/null
then
    apk add ifupdown-ng
fi

NETFILE=/etc/network/interfaces
cp -f $NETFILE $NETFILE.backup


cat << EOF > $NETFILE
auto lo
iface lo inet loopback

auto eth0
iface eth0
address $PREFIX/$MASK
gateway $GATEWAY
EOF

rc-service networking restart

