package config

import "os"

type Config struct {
	ADDR string;
	SHUTDOWN_TIMEOUT string;
}

func Load() Config {
	val, ok:= os.LookupEnv(key)
	if ok! {
		return 
	}
	return 
}