package pikacloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a schema.Provider for Pikacloud.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PIKACLOUD_TOKEN", nil),
				Description: "The token key for API operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"pikacloud_instance": resourcePikacloudInstance(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Token: d.Get("token").(string),
	}

	return config.Client()
}
