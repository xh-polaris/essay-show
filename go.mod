module github.com/xh-polaris/essay-show

go 1.22.0

toolchain go1.22.12

replace go.opentelemetry.io/otel v1.32.0 => go.opentelemetry.io/otel v1.24.0

require (
	github.com/bytedance/gopkg v0.1.1
	github.com/cloudwego/hertz v0.9.3
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.6.0
	github.com/hertz-contrib/monitor-prometheus v0.1.2
	github.com/hertz-contrib/obs-opentelemetry/tracing v0.4.1
	github.com/jinzhu/copier v0.3.5
	github.com/xh-polaris/gopkg v0.0.0-20250312141711-7327267f4ea6
	github.com/xh-polaris/service-idl-gen-go v0.0.0-20241008065911-312508da7b2f
	github.com/zeromicro/go-zero v1.7.3
	go.mongodb.org/mongo-driver v1.17.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.60.0
	go.opentelemetry.io/contrib/propagators/b3 v1.20.0
	go.opentelemetry.io/otel v1.35.0
	golang.org/x/net v0.30.0
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.35.1
)

require (
	github.com/andeya/ameda v1.5.3 // indirect
	github.com/andeya/goutil v1.0.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bufbuild/protocompile v0.14.1 // indirect
	github.com/bytedance/go-tagexpr/v2 v2.9.11 // indirect
	github.com/bytedance/sonic v1.12.10 // indirect
	github.com/bytedance/sonic/loader v0.2.2 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/cloudwego/configmanager v0.2.2 // indirect
	github.com/cloudwego/dynamicgo v0.5.2 // indirect
	github.com/cloudwego/fastpb v0.0.5 // indirect
	github.com/cloudwego/frugal v0.2.3 // indirect
	github.com/cloudwego/gopkg v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/cloudwego/kitex v0.12.3 // indirect
	github.com/cloudwego/localsession v0.1.2 // indirect
	github.com/cloudwego/netpoll v0.6.5 // indirect
	github.com/cloudwego/runtimex v0.1.1 // indirect
	github.com/cloudwego/thriftgo v0.3.18 // indirect
	github.com/coze-dev/coze-go v0.0.0-20250701071254-8082a2c3da76 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/pprof v0.0.0-20240727154555-813a5fbdbec8 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/jhump/protoreflect v1.17.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kitex-contrib/monitor-prometheus v0.2.0 // indirect
	github.com/kitex-contrib/obs-opentelemetry v0.2.3 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nyaruka/phonenumbers v1.3.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.20.5 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/tidwall/gjson v1.17.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	golang.org/x/arch v0.6.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240711142825-46eb208f015d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240701130421-f6361c86f094 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/utils v0.0.0-20240711033017-18e509b52bc8 // indirect
)
