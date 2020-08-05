Terraform Provider Harbor
==================

A Terraform provider for [Harbor](https://goharbor.io/)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.12

Using the provider
----------------------

Developing the Provider
---------------------

This project uses two primary packages, `harbor` implements the interface for interacting
with the Harbor API, `provider` implements Terraform specific provider configuration.

To run the provider locally for development, you'll require a local instance of
Harbor. This project includes scripts for setting up a local environment using
docker-for-desktop's Kubernetes cluster, and Helm. If you'd prefer to use another instance
of Harbor, that will also work, simply don't use the local environment scripts.

If you would like to use the local environment, run the command `make local`

Your Harbor instance will be running on a random port that can be obtained by running
the command `kubectl get service harbor`
The Harbor helm chart default admin username and password will be used. `admin:Harbor12345`

To build and install the provider locally, run `make install`, this will install
the provider globally, under the development version `v0.0.0`.

Examples of how to use the provider can be found under the `examples/` directory.

To run automated linting on the codebase, run the command `make lint`
To run the automated acceptance testing on the provider, run the command `make acctest`
and provide the following data to the provided script.
- HarborURL: The base url of the Harbor instance to test against. e.g. `http://localhost:8080`
- Username: The username of the user to run acceptance tests as.
- Password: The password of the user to run acceptance tests as.
*Note:* Acceptance tests create real resources, and often cost money to run.

