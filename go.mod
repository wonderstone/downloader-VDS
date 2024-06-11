module github.com/wonderstone/downloader-vds

go 1.21.5

require (
	github.com/go-gota/gota v0.12.0
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/vds v0.0.0-00010101000000-000000000000
)

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

require (
	github.com/emirpasic/gods v1.18.1
	github.com/spf13/viper v1.19.0
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/data v0.0.0-00010101000000-000000000000 // indirect
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/hiscenter v0.0.0-00010101000000-000000000000 // indirect
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/rtcenter v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gonum.org/v1/gonum v0.9.1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/data => ./data
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/hiscenter => ./hiscenter
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/rtcenter => ./rtcenter
	gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/vds => ./vds

)
