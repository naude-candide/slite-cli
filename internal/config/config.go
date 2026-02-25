package config

import "os"

func APIKey() string {
	return os.Getenv("SLITE_API_KEY")
}
