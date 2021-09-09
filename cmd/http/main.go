package main

type Config struct {
}

func getCfg() Config {
	return Config{}
}

func main() {
	_ = getCfg()
}
