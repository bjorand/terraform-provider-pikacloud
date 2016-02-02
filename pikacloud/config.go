package pikacloud

import (
	"log"

	"golang.org/x/oauth2"

	"github.com/bjorand/gopikacloud"
)

type Config struct {
	Token string
}

// Client() returns a new client for accessing pikacloud.
func (c *Config) Client() (*gopikacloud.Client, error) {
	// client, err := gopikacloud.NewClient(c.Token)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error setting up client: %s", err)
	// }
	// log.Printf("[INFO] Pikacloud API Client configured")
	// return client, nil
	tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: c.Token,
	})

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSrc)
	client := gopikacloud.NewClient(oauthClient)

	log.Printf("[INFO] Pikacloud Client configured for URL: %s", client.BaseURL.String())

	return client, nil
}
