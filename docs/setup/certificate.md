---
title: Certificate
description: Generate self-signed SSL certificates
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: certificate-searchlight
    name: Generate Certificate
    parent: setup
    weight: 20
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: setup
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Certificate

This article shows you how to generate SSL certificates using **`openssl`** or **`easyrsa`**.

## openssl

`openssl` can also be use to manually generate certificates for your cluster.

1. Set `HOST_IP` ENV to host IP

2. Generate a ca.key with 2048bit:
```sh
openssl genrsa -out ca.key 2048
```

3. According to the ca.key generate a ca.crt
```sh
openssl req -x509 -new -nodes -key ca.key -subj "/CN=${HOST_IP}" -days 10000 -out ca.crt
```

4. Generate a server.key with 2048bit
```sh
openssl genrsa -out server.key 2048
```

5. According to the server.key generate a server.csr:
```sh
openssl req -new -key server.key -subj "/CN=${HOST_IP}" -out server.csr
```

6. According to the ca.key, ca.crt and server.csr generate the server.crt:
```sh
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 10000
```

## easyrsa

`easyrsa` can be used to manually generate certificates for your cluster.

1. Download, unpack, and initialize the patched version of easyrsa3.
```sh
curl -L -O https://storage.googleapis.com/kubernetes-release/easy-rsa/easy-rsa.tar.gz
tar xzf easy-rsa.tar.gz
cd easy-rsa-master/easyrsa3
./easyrsa init-pki
```

2. Set `HOST_IP` ENV to Kubernetes host IP

3. Generate a CA. (--batch set automatic mode. --req-cn default CN to use.)
```sh
./easyrsa --batch "--req-cn=${HOST_IP}@`date +%s`" build-ca nopass
```

4. Generate server certificate and key
```sh
./easyrsa --subject-alt-name="IP:${HOST_IP}" build-server-full kubernetes-master nopass
```

## Acknowledgement

This documentation is adapted from [kubernetes.io]((https://kubernetes.io/docs/admin/authentication/#appendix)). 