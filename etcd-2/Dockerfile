FROM quay.io/coreos/etcd:v3.6.7

COPY ./config/etcd.yml /etc/etcd/etcd.yml

CMD ["etcd", "--config-file", "/etc/etcd/etcd.yml"]