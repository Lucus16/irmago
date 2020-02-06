module github.com/privacybydesign/irmago

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/bwesterb/go-atum v1.0.0
	github.com/certifi/gocertifi v0.0.0-20180118203423-deb3ae2ef261 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/raven-go v0.0.0-20180121060056-563b81fc02b7
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/go-chi/cors v1.0.0
	github.com/go-errors/errors v1.0.0
	github.com/hashicorp/go-retryablehttp v0.6.2
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jasonlvhit/gocron v0.0.0-20180312192515-54194c9749d4
	github.com/magiconair/properties v1.8.0 // indirect
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/mdp/qrterminal v1.0.1
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pkg/errors v0.8.0
	github.com/privacybydesign/gabi v0.0.0-20190503104928-ce779395f4c9
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/afero v1.2.0 // indirect
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v0.0.1
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v0.0.0-00010101000000-000000000000
	github.com/spf13/viper v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.2.2
	github.com/timshannon/bolthold v0.0.0-20190812165541-a85bcc049a2e
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	go.etcd.io/bbolt v1.3.2
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	gopkg.in/antage/eventsource.v1 v1.0.0-20150318155416-803f4c5af225
)

replace github.com/spf13/pflag => github.com/sietseringers/pflag v1.0.4-0.20190111213756-a45bfec10d59

replace github.com/spf13/viper => github.com/sietseringers/viper v1.0.1-0.20190113114857-554683669b21
