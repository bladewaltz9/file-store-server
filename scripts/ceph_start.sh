#!/bin/bash

OSD_COUNT=3

# start ceph monitor
sudo docker start ceph-mon

# start ceph manager
sudo docker start ceph-mgr

# start ceph osd
for (( i=0; i<$OSD_COUNT; i++ )); do
    sudo docker start ceph-osd-$i
done

# start ceph rgw
sudo docker start ceph-rgw