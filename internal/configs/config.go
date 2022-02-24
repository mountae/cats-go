// Package configs init env variables by caarlos0
package configs

// Config contains all env variables
type Config struct {
	Port string `env:"PORT" envDefault:":8000"`

	PgUser     string `env:"POSTGRES_USER" envDefault:"postgres"`
	PgPassword string `env:"POSTGRES_PASSWORD" envDefault:"root"`
	PgHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PgPort     string `env:"POSTGRES_PORT" envDefault:"5432"`
	PgDBName   string `env:"POSTGRES_DATABASE" envDefault:"postgres"`

	Mongo           string `env:"MONGO_CONNSTRING" envDefault:"mongodb://userm:testpassw@localhost:27017/"`
	MongoUser       string `env:"MONGO_USERNAME" envDefault:"userm"`
	MongoPassword   string `env:"MONGO_PASSWORD" envDefault:"testpassw"`
	MongoHost       string `env:"MONGO_HOST" envDefault:"localhost"`
	MongoPort       string `env:"MONGO_PORT" envDefault:"27017"`
	MongoDBName     string `env:"MONGO_DBNAME" envDefault:"mongodb"`
	MongoCollection string `env:"MONGO_COLLECTION" envDefault:"mongocl"`

	KeyForSignatureJwt string `env:"KEY_FOR_SIGNATURE_JWT" envDefault:"mySecret"`
	Salt               string `env:"SALT_FOR_GENERATE_PASSWORD" envDefault:"l337c0d3"`
}
