module github.com/ireuven89/auctions/bidder-service

go 1.23.9

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/go-kit/kit v0.13.0
	github.com/go-sql-driver/mysql v1.9.2
	github.com/google/uuid v1.6.0
	github.com/ireuven89/auctions/shared v0.0.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/pressly/goose/v3 v3.24.3
	github.com/sethvargo/go-retry v0.3.0
	github.com/stretchr/testify v1.10.0
	go.uber.org/zap v1.27.0
)

replace github.com/ireuven89/auctions/shared => ../shared

require (
	filippo.io/edwards25519 v1.1.0 // indirect
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
	github.com/redis/go-redis/v9 v9.11.0 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
