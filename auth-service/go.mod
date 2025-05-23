module github.com/ireuven89/auctions/auth-service

go 1.23.0

require (
	github.com/aws/aws-sdk-go-v2/config v1.29.14
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.35.4
	github.com/go-kit/kit v0.13.0
	github.com/go-sql-driver/mysql v1.9.2
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/google/uuid v1.6.0
	github.com/ireuven89/auctions/shared v0.0.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/pressly/goose/v3 v3.24.3
	github.com/redis/go-redis/v9 v9.8.0
	github.com/sethvargo/go-retry v0.3.0
	github.com/stretchr/testify v1.10.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.38.0
)

replace github.com/ireuven89/auctions/shared => ../shared

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.36.3 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.67 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.19 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
