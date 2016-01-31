package pikacloud

import (
	"log"

	"github.com/bjorand/gopikacloud"
)

type Config struct {
	Token string
}

// Client() returns a new client for accessing pikacloud.
func (c *Config) Client() (*gopikacloud.Client, error) {
	client := gopikacloud.NewClient(nil)
	// if err != nil {
	// return nil, fmt.Errorf("Error setting up client: %s", err)
	// }
	log.Printf("[INFO] Pikacloud API Client configured for URL: %s", client.BaseURL.String())
	return client, nil
}
