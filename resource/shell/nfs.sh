nfs server
yum install nfs-utils -y
mkdir /nfs && chmod 0777 /nfs

# /nfs is the mount point for the NFS share
# * all servers are accessible
# rw is the access mode for the NFS share
cat >/etc/exports<<EOF
/nfs *(rw,all_squash,anonuid=500,anongid=500)
EOF

systemctl restart rpcbind.service
systemctl restart nfs.service
chkconfig nfs on

# client-config

showmount -e "10.2.103.113"
Export list for 10.2.103.113:
/nfs113 *

mkdir /nfs
mount -o rw -t nfs 10.2.103.113:/nfs /nfs

touch /nfs/tidb-01
check files on 127