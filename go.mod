module github.com/appscode/searchlight

go 1.12

require (
	cloud.google.com/go v0.39.0 // indirect
	contrib.go.opencensus.io/exporter/ocagent v0.5.0 // indirect
	github.com/Azure/go-autorest v12.0.0+incompatible // indirect
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/StackExchange/wmi v0.0.0-20181212234831-e0a55b97c705 // indirect
	github.com/Unknwon/com v0.0.0-20190321035513-0fed4efef755 // indirect
	github.com/appscode/go v0.0.0-20190523031839-1468ee3a76e8
	github.com/codeskyblue/go-sh v0.0.0-20190412065543-76bd3d59ff27
	github.com/emicklei/go-restful v2.9.5+incompatible // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8 // indirect
	github.com/go-macaron/auth v0.0.0-20161228062157-884c0e6c9b92
	github.com/go-macaron/inject v0.0.0-20160627170012-d8a0b8677191 // indirect
	github.com/go-macaron/toolbox v0.0.0-20180818072302-a77f45a7ce90
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-openapi/analysis v0.17.0 // indirect
	github.com/go-openapi/errors v0.17.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.0 // indirect
	github.com/go-openapi/jsonreference v0.19.0 // indirect
	github.com/go-openapi/loads v0.19.0
	github.com/go-openapi/spec v0.19.0
	github.com/go-openapi/swag v0.19.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/gophercloud/gophercloud v0.0.0-20190515011819-1992d5238d78 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190430165422-3e4dfb77656c // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.0 // indirect
	github.com/mailru/easyjson v0.0.0-20190403194419-1ea4449da983 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.3
	github.com/prometheus/common v0.4.1 // indirect
	github.com/shirou/gopsutil v0.0.0-20180227225847-5776ff9c7c5d
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/smartystreets/assertions v0.0.0-20190401211740-f487f9de1cd3 // indirect
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.4
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20190513172903-22d7a77e9e5f // indirect
	golang.org/x/net v0.0.0-20190514140710-3ec191127204 // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58 // indirect
	golang.org/x/sys v0.0.0-20190514135907-3a4b5fb9f71f // indirect
	gomodules.xyz/cert v1.0.0
	gomodules.xyz/envconfig v1.3.1-0.20190308184047-426f31af0d45
	gomodules.xyz/notify v0.0.0-20190424183923-af47cb5a07a4
	google.golang.org/appengine v1.6.0 // indirect
	google.golang.org/genproto v0.0.0-20190513181449-d00d292a067c // indirect
	gopkg.in/ini.v1 v1.42.0
	gopkg.in/macaron.v1 v1.3.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	k8s.io/api v0.0.0-20190515023547-db5a9d1c40eb
	k8s.io/apiextensions-apiserver v0.0.0-20190515024537-2fd0e9006049
	k8s.io/apimachinery v0.0.0-20190515023456-b74e4c97951f
	k8s.io/apiserver v0.0.0-20190515064100-fc28ef5782df
	k8s.io/cli-runtime v0.0.0-20190515024640-178667528169 // indirect
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/component-base v0.0.0-20190515024022-2354f2393ad4 // indirect
	k8s.io/klog v0.3.1 // indirect
	k8s.io/kube-aggregator v0.0.0-20190515024249-81a6edcf70be // indirect
	k8s.io/kube-openapi v0.0.0-20190510232812-a01b7d5d6c22
	kmodules.xyz/client-go v0.0.0-20190527113919-eb165836b702
	kmodules.xyz/webhook-runtime v0.0.0-20190508093950-b721b4eba5e5
)

replace (
	github.com/graymeta/stow => github.com/appscode/stow v0.0.0-20190506085026-ca5baa008ea3
	gopkg.in/robfig/cron.v2 => github.com/appscode/cron v0.0.0-20170717094345-ca60c6d796d4
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190315093550-53c4693659ed
	k8s.io/apimachinery => github.com/kmodules/apimachinery v0.0.0-20190508045248-a52a97a7a2bf
	k8s.io/apiserver => github.com/kmodules/apiserver v0.0.0-20190508082252-8397d761d4b5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190314001948-2899ed30580f
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190314002645-c892ea32361a
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190314000054-4a91899592f4
	k8s.io/klog => k8s.io/klog v0.3.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190314000639-da8327669ac5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190314001731-1bd6a4002213
	k8s.io/utils => k8s.io/utils v0.0.0-20190221042446-c2654d5206da
)
