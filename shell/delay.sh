# $1 delay time, unit ms
# $2 delay port

tc qdisc add dev eth0 root handle 1: prio bands 4
tc qdisc add dev eth0 parent 1:4 handle 40: netem delay $1ms
tc filter add dev eth0 protocol ip parent 1:0 prio 4 u32 match ip dport $2 0xffff flowid 1:4
