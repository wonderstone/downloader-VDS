module vds

go 1.21.5

replace (
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/data => ../data
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/hiscenter => ../hiscenter
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/rtcenter => ../rtcenter
)

require (
	github.com/go-gota/gota v0.12.0
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/data v0.0.0-00010101000000-000000000000
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/hiscenter v0.0.0-00010101000000-000000000000
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/rtcenter v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.1
)

require (
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gonum.org/v1/gonum v0.9.1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)
