# Build Instructions

##### Build plugin

```
./hack/make.py build hyperalert
```

##### Push hyperalert plugin
```
./hack/make.py push hyperalert
```
This `push` command upload binary to cloud using gsutil
```
gsutil cp hyperalert-linux-amd64 gs://appscode-dev/binaries/hyperalert/<tag>/hyperalert-linux-amd64
gsutil acl ch -u AllUsers:R gs://appscode-dev/binaries/hyperalert/<tag>/hyperalert-linux-amd64
```
