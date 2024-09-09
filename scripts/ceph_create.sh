#!/bin/bash

# set network 
NETWORK_NAME="ceph-network"
NETWORK_SUBNET="172.20.0.0/16"
MON_IP="172.20.0.10"
MGR_IP="172.20.0.14"
# OSD_IPS=("172.20.0.20" "172.20.0.21" "172.20.0.22")
OSD_COUNT=3
RGW_IP="172.20.0.15"

# create docker network
sudo docker network create --driver bridge --subnet $NETWORK_SUBNET ceph-network

# start ceph monitor(mon)
sudo docker run -d \
	--name ceph-mon \
	--hostname ceph-mon \
	--network $NETWORK_NAME \
	--ip $MON_IP \
	-e CLUSTER=ceph \
	-e WEIGHT=1.0 \
	-e MON_IP=$MON_IP \
	-e MON_NAME=ceph-mon \
	-e CEPH_PUBLIC_NETWORK=$NETWORK_SUBNET \
	-v /etc/ceph:/etc/ceph \
	-v /var/lib/ceph/:/var/lib/ceph/ \
	-v /var/log/ceph/:/var/log/ceph/ \
	ceph/daemon:latest mon

# start ceph manager(mgr)
sudo docker run -d \
	--privileged=true \
	--name ceph-mgr \
	--hostname ceph-mgr \
	--network $NETWORK_NAME \
	--ip $MGR_IP \
	-e CLUSTER=ceph \
	-p 28080:8080 \
	-p 28443:8443 \
	--pid=container:ceph-mon \
	-v /etc/ceph:/etc/ceph \
	-v /var/lib/ceph/:/var/lib/ceph/ \
	ceph/daemon:latest mgr

# create ceph osd keyring
sudo docker exec ceph-mon ceph auth get client.bootstrap-osd -o /var/lib/ceph/bootstrap-osd/ceph.keyring

# start ceph osd(3 nodes)
for (( i=0; i<$OSD_COUNT; i++ )); do
    OSD_IP="172.20.0.$((20 + i))"
    sudo docker run -d \
        --privileged=true \
        --name ceph-osd-$i \
        --hostname ceph-osd-$i \
        --network $NETWORK_NAME \
        --ip $OSD_IP \
        -e CLUSTER=ceph \
        -e WEIGHT=1.0 \
        -e MON_NAME=ceph-mon \
        -e MON_IP=$MON_IP \
        -e OSD_TYPE=directory \
        -v /etc/ceph:/etc/ceph \
        -v /var/lib/ceph/:/var/lib/ceph/ \
        -v /var/lib/ceph/osd/$i:/var/lib/ceph/osd \
        -v /etc/localtime:/etc/localtime:ro \
        ceph/daemon:latest osd
done

# create ceph rgw keyring
sudo docker exec ceph-mon ceph auth get client.bootstrap-rgw -o /var/lib/ceph/bootstrap-rgw/ceph.keyring

# start ceph rgw
sudo docker run -d \
    --privileged=true \
    --name ceph-rgw \
    --hostname ceph-rgw \
    --network $NETWORK_NAME \
    --ip $RGW_IP \
    -e CLUSTER=ceph \
    -e RGW_NAME=ceph-rgw \
    -p 27480:7480 \
    -v /var/lib/ceph/:/var/lib/ceph/ \
    -v /etc/ceph:/etc/ceph \
    -v /etc/localtime:/etc/localtime:ro \
    ceph/daemon:latest rgw

echo "Ceph containers create successfully."