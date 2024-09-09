#!/bin/bash

OSD_COUNT=3

# stop ceph rgw
sudo docker stop ceph-rgw

# stop ceph manager
sudo docker stop ceph-mgr

# stop ceph monitor
sudo docker stop ceph-mon

# stop ceph osd
for (( i=0; i<$OSD_COUNT; i++ )); do
    sudo docker stop ceph-osd-$i
done