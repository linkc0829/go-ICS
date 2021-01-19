package utils

//ServerConfig defines the configuration of server
type ServerConfig struct {
	Host          string
	Port          string
	URISchema     string //http:// or https://
	ApiVer        string
	StaticPath    string
	SessionSecret string
	JWT           JWTConfig
	GraphQL       GQLConfig
	MongoDB       MGDBConfig
	Redis         RedisConfig
	AuthProviders []AuthProvider
}

//JWTConfig defines the options for JWT tokens
type JWTConfig struct {
	Secret             string
	Algorithm          string
	AccessTokenExpire  string
	RefreshTokenExpire string
}

// GQLConfig defines the configuration for the GQL Server
type GQLConfig struct {
	Path                string
	PlaygroundPath      string
	IsPlaygroundEnabled bool
}

// MGDBConfig defines the configuration for the MongoDB config
type MGDBConfig struct {
	DSN string
}

type RedisConfig struct {
	EndPoint string
	PWD      string
}

// AuthProvider defines the configuration for the Goth config
type AuthProvider struct {
	Provider  string
	ClientKey string
	Secret    string
	Domain    string // If needed, like with auth0
	Scopes    []string
}

//ListenEndpoint returns the endpoint string
func (s *ServerConfig) ListenEndpoint() string {
	if s.Port == "80" || s.Port == "443" {
		return s.Host
	}
	return s.Host + ":" + s.Port
}

//VersioningEndpoint retruns the versioning api path, path should have "/" as prefix
func (s *ServerConfig) VersioningEndpoint(path string) string {
	return "/api/" + s.ApiVer + "/" + path
}

//SchemaVersioningEndpoint return the complete URI path
func (s *ServerConfig) SchemaVersioningEndpoint(path string) string {
	if s.Port == "80" {
		return s.URISchema + s.Host + "/" + s.ApiVer + path
	}
	return s.URISchema + s.Host + ":" + s.Port + "/api/" + s.ApiVer + path
}
