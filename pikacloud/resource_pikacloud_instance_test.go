package pikacloud

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/bjorand/gopikacloud"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPikacloudInstance_Basic(t *testing.T) {
	var instance gopikacloud.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPikacloudInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPikacloudInstanceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPikacloudInstanceExists("pikacloud_instance.foobar", &instance),
					testAccCheckPikacloudInstanceAttributes(&instance),
					resource.TestCheckResourceAttr(
						"pikacloud_instance.foobar", "region", "3"),
					resource.TestCheckResourceAttr(
						"pikacloud_instance.foobar", "hosts.0", "example.com"),
				),
			},
		},
	})
}

func TestAccPikacloudInstance_Update(t *testing.T) {
	var instance gopikacloud.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPikacloudInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPikacloudInstanceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPikacloudInstanceExists("pikacloud_instance.foobar", &instance),
					testAccCheckPikacloudInstanceAttributes(&instance),
				),
			},

			resource.TestStep{
				Config: testAccCheckPikacloudInstanceConfig_ChangeHosts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPikacloudInstanceExists("pikacloud_instance.foobar", &instance),
					testAccCheckPikacloudInstanceChangeHosts(&instance),
					resource.TestCheckResourceAttr(
						"pikacloud_instance.foobar", "region", "3"),
					resource.TestCheckResourceAttr(
						"pikacloud_instance.foobar", "hosts.0", "foobar.com"),
				),
			},
		},
	})
}

func testAccCheckPikacloudInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*gopikacloud.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pikacloud_instance" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the Instance
		_, _, err = client.Instances.Get(id)

		// Wait

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for instance (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckPikacloudInstanceAttributes(instance *gopikacloud.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Region != 3 {
			return fmt.Errorf("Bad region: %d", instance.Region)
		}
		if instance.Hosts[0] != "example.com" {
			return fmt.Errorf("Bad hostname: %d", instance.Hosts[0])
		}
		return nil
	}
}
func testAccCheckPikacloudInstanceChangeHosts(instance *gopikacloud.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Region != 3 {
			return fmt.Errorf("Bad region: %d", instance.Region)
		}

		if instance.Hosts[0] != "foobar.com" {
			return fmt.Errorf("Bad hostname: %s", instance.Hosts[0])
		}

		return nil
	}
}

func testAccCheckPikacloudInstanceExists(n string, instance *gopikacloud.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Instance ID is set")
		}

		client := testAccProvider.Meta().(*gopikacloud.Client)

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the Instance
		retrieveInstance, _, err := client.Instances.Get(id)

		if err != nil {
			return err
		}

		if strconv.Itoa(retrieveInstance.ID) != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *retrieveInstance

		return nil
	}
}

func testAccCheckPikacloudInstanceRecreated(t *testing.T,
	before, after *gopikacloud.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ID == after.ID {
			t.Fatalf("Expected change of instance IDs, but both were %v", before.ID)
		}
		return nil
	}
}

var testAccCheckPikacloudInstanceConfig_basic = `
resource "pikacloud_instance" "foobar" {
  region = 3
  hosts  = ["example.com"]
}
`

var testAccCheckPikacloudInstanceConfig_ChangeHosts = `
resource "pikacloud_instance" "foobar" {
  region      = 3
  hosts  = ["foobar.com"]
}
`
