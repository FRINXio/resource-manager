module github.com/net-auto/resourceManager

go 1.14

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/prometheus v0.2.0
	contrib.go.opencensus.io/integrations/ocsql v0.1.6
	github.com/99designs/gqlgen v0.11.3
	github.com/AlekSi/pointer v1.1.0
	github.com/Azure/azure-amqp-common-go/v2 v2.1.0 // indirect
	github.com/NYTimes/gziphandler v1.1.1
	github.com/apaxa-go/eval v0.0.0-20171223182326-1d18b251d679 // indirect
	github.com/apaxa-go/helper v0.0.0-20180607175117-61d31b1c31c3 // indirect
	github.com/cenkalti/backoff/v4 v4.0.2
	github.com/facebookincubator/ent v0.2.8-0.20200726173043-ff6163f1a068
	// IMPORTANT!! if symphony version is updated, also update generate.go in ent/
	github.com/facebookincubator/symphony v0.0.0-20200727135522-f17846fd0b94
	github.com/fatih/color v1.7.0 // indirect
	github.com/go-sql-driver/mysql v1.5.1-0.20200311113236-681ffa848bae
	github.com/golang/protobuf v1.4.2
	github.com/google/addlicense v0.0.0-20200422172452-68a83edd47bc // indirect
	github.com/google/uuid v1.1.1
	github.com/google/wire v0.4.0
	github.com/gorilla/mux v1.7.4
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/ugorji/go/codec v1.1.7
	github.com/vektah/gqlparser/v2 v2.0.1
	go.opencensus.io v0.22.4
	go.uber.org/zap v1.15.0
	gocloud.dev v0.20.0
	gocloud.dev/pubsub/natspubsub v0.20.0
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200702021140-07506425bd67 // indirect
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)
