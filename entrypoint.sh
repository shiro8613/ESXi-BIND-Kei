#!/bin/bash

MODIFIED_PORT=`echo $(echo ${BIND_PORT} | sed -e 's/{{/${/g' -e 's/}}/}/g')`
MODIFIED_IP=`echo $(echo ${ESXI_IP} | sed -e 's/{{/${/g' -e 's/}}/}/g')`
MODIFIED_DOMAIN=`echo $(echo ${DOMAIN_NAME} | sed -e 's/{{/${/g' -e 's/}}/}/g')`

cat << EOF > /data/config.yml
port: $MODIFIED_PORT
esxi_addr: $MODIFIED_IP
domain_name: $MODIFIED_DOMAIN
EOF

./data/dns
