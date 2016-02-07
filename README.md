# Terraform provider for Pikacloud

This is a terraform provider that lets you provision
servers on Pikacloud via [Terraform](https://terraform.io/).

Build status: [![Build Status](https://travis-ci.org/bjorand/terraform-provider-pikacloud.svg?branch=master)](https://travis-ci.org/bjorand/terraform-provider-pikacloud)


## Installing

Binaries are published on Bintray: [ ![Download](https://api.bintray.com/packages/bjorand/terraform-provider-pikacloud/terraform-provider-pikacloud/images/download.svg) ](https://bintray.com/bjorand/terraform-provider-pikacloud/terraform-provider-pikacloud/_latestVersion)

[Copied from the Terraform documentation](https://www.terraform.io/docs/plugins/basics.html):
> To install a plugin, put the binary somewhere on your filesystem, then configure Terraform to be able to find it. The configuration where plugins are defined is ~/.terraformrc for Unix-like systems and %APPDATA%/terraform.rc for Windows.

The binary should be renamed to terraform-provider-pikacloud.

You should update your .terraformrc and refer to the binary:

```hcl
providers {
  pikacloud = "/path/to/terraform-provider-pikacloud"
}
```

## Using the provider

Here is an example that will setup the following:
+ An SSH key resource.
+ A virtual server resource that uses an existing SSH key.
+ A virtual server resource using an existing SSH key and a Terraform managed SSH key (created as "test_key_1" in the example below).

(create this as pikacloud.tf and run terraform commands from this directory):
```hcl
provider "pikacloud" {
    token = "XXXX"
}

resource "pikacloud_instance" "mylb" {
    region = 3
    hosts = ["example.com", "foobar.com"]
    server {
        droplet_id = 3434443
        port = 8000
    }
    server {
        droplet_id = ${digitalocean_droplet.myserver.id}
        port = 8000
    }
}

output "mylb_dns" {
  value = "${pikacloud_instance.mylb.dns}"
}

You'll need to provide your Pikacloud API token,
so that Terraform can connect. If you don't want to put
credentials in your configuration file, you can leave them
out:

```
provider "pikacloud" {}
```

...and instead set these environment variables:

- **PIKACLOUD_TOKEN**: Your Pikacloud API token

### Pikacloud regions

| ID | Region name  |
|---|---------------|
| 0 | London        |
| 1 | San Francisco |
| 2 | New York      |
| 3 | Singapore     |
| 4 | Toronto       |

## Building from source

1.  [Install Go](https://golang.org/doc/install) on your machine
2.  [Set up Gopath](https://golang.org/doc/code.html)
3.  `git clone` this repository into `$GOPATH/src/github.com/bjorand/terraform-provider-pikacloud`
4.  Run `go get` to get dependencies
5.  Run `go install` to build the binary. You will now find the
    binary at `$GOPATH/bin/terraform-provider-pikacloud`.

## Releasing

1. Update `.bintray.json` with the new version to be released (e.g. 0.1337), and the release date.
2. Tag the release: `git tag 0.1337`
3. Push the tag: `git push --tags`
4. Travis CI will build and release the binaries.

## Running
1.  create the example file pikacloud.tf in your working directory
2.  terraform plan
3.  terraform apply
