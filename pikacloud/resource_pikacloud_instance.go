package pikacloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bjorand/gopikacloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePikacloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourcePikacloudInstanceCreate,
		Read:   resourcePikacloudInstanceRead,
		Delete: resourcePikacloudInstanceDelete,
		Update: resourcePikacloudInstanceUpdate,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"certificate": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"hosts": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"servers": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"dns": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePikacloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gopikacloud.Client)

	// Create the new instance
	newInstance := &gopikacloud.InstanceCreateRequest{
		Region: d.Get("region").(int),
		Hosts:  d.Get("hosts").([]interface{}),
	}
	log.Printf("[DEBUG] Pikacloud Instance create configuration: %#v", newInstance)

	instance, _, err := client.Instances.Create(newInstance)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(instance.ID))

	log.Printf("[INFO] Instance ID: %s", d.Id())

	return resourcePikacloudInstanceRead(d, meta)
}

func resourcePikacloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gopikacloud.Client)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid pikacloud instance id: %v", err)
	}

	// Retrieve the instance properties for updating the state
	instance, resp, err := client.Instances.Get(id)
	if err != nil {
		// check if the droplet no longer exists.
		if resp.StatusCode == 404 {
			log.Printf("[WARN] Pikacloud Instance (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving instance: %s", err)
	}
	log.Printf("[DEBUG] Pikacloud Instance read: %#v", instance)

	d.Set("region", instance.Region)
	d.Set("hosts", instance.Hosts)
	d.Set("dns", instance.DNS)
	d.Set("certificate", instance.Certificate)
	d.Set("servers", instance.Servers)
	return nil
}

func resourcePikacloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gopikacloud.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid instance id: %v", err)
	}

	log.Printf("[INFO] Deleting instance: %s", d.Id())

	// Destroy the droplet
	_, err = client.Instances.Delete(id)

	// Handle remotely destroyed droplets
	if err != nil && strings.Contains(err.Error(), "404 Not Found") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error deleting instance: %s", err)
	}

	return nil
}

func resourcePikacloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gopikacloud.Client)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid record ID: %v", err)
	}

	updateInstance := gopikacloud.InstanceUpdateRequest{
		Hosts:  d.Get("hosts").([]interface{}),
		Region: d.Get("region").(int),
	}

	log.Printf("[DEBUG] instance update configuration: %#v", updateInstance)

	if err != nil {
		return fmt.Errorf("invalid instance id: %v", err)
	}
	if d.HasChange("hosts") {
		_, _, err = client.Instances.Update(id, &updateInstance)
	}

	return resourcePikacloudInstanceRead(d, meta)
}
