nohup etcd --name infra0 \
--data-dir=data/infra0 \
--listen-peer-urls https://127.0.0.1:3380 \
--listen-client-urls https://127.0.0.1:3379 \
--advertise-client-urls https://127.0.0.1:3379 \
--client-cert-auth --trusted-ca-file=certs/ca.pem \
--cert-file=certs/127.0.0.1.pem \
--key-file=certs/127.0.0.1-key.pem \
--peer-client-cert-auth --peer-trusted-ca-file=certs/ca.pem \
--peer-cert-file=certs/127.0.0.1.pem \
--peer-key-file=certs/127.0.0.1-key.pem > log/infra0.log 2>&1 &

nohup etcd --name infra1 \
--data-dir=data/infra1 \
--listen-peer-urls https://127.0.0.1:4380 \
--listen-client-urls https://127.0.0.1:4379 \
--advertise-client-urls https://127.0.0.1:4379 \
--client-cert-auth --trusted-ca-file=certs/ca.pem \
--cert-file=certs/127.0.0.1.pem \
--key-file=certs/127.0.0.1-key.pem \
--peer-client-cert-auth --peer-trusted-ca-file=certs/ca.pem \
--peer-cert-file=certs/127.0.0.1.pem \
--peer-key-file=certs/127.0.0.1-key.pem > log/infra0.log 2>&1 &

nohup etcd --name infra2 \
--data-dir=data/infra2 \
--listen-peer-urls https://127.0.0.1:5380 \
--listen-client-urls https://127.0.0.1:5379 \
--advertise-client-urls https://127.0.0.1:5379 \
--client-cert-auth --trusted-ca-file=certs/ca.pem \
--cert-file=certs/127.0.0.1.pem \
--key-file=certs/127.0.0.1-key.pem \
--peer-client-cert-auth --peer-trusted-ca-file=certs/ca.pem \
--peer-cert-file=certs/127.0.0.1.pem \
--peer-key-file=certs/127.0.0.1-key.pem > log/infra0.log 2>&1 &
