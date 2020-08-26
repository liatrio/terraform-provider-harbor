# Harbor Provider

This Harbor provider can be used to interact with and configure [Harbor](https://goharbor.io/)

## Example Usage

```hcl
provider "harbor" {
  url = "http://localhost:8080"
  username = "admin"
  password = "Harbor12345"
}
```

## Argument Reference

`url` - (Required) The URL of Harbor instance. Defaults to the environment variable `HARBOR_URL`

`username` - (Required) The username of the user to use while connecting to the Harbor instance. Defaults to the environment variable `HARBOR_USERNAME`.

`password` - (Required) The password corrosponding to the user used for connecting to the Harbor instance. Defaults to the environment variable `HARBOR_PASSWORD`.

`tls_insecure_skip_verify` - (Optional) Allows skipping TLS certificate verification. This variable is provided for ease of use in development, but is not recommended for production use. Defaults to `false`.
