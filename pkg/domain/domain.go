package domain

type API struct {
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
}

type Config struct {
	API *API `yaml:"api"`
}
