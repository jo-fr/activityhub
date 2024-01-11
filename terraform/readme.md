# terraform

The terraform folder contains everything that related with backendcloud infrastructure.

## setup

If you want to deploy the backend infrastructure using terraform make sure you take care of the follow these steps:

1. install [terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli) on your system
1. [Create a new project](https://cloud.google.com/resource-manager/docs/creating-managing-projects) on GCP
1. [Create a service account](https://cloud.google.com/iam/docs/service-accounts-create) and download the [credentials](https://cloud.google.com/iam/docs/keys-create-delete) as json
1. add the ` $GOOGLE_APPLICATION_CREDENTIALS` env variable with the path to your credentials
1. run `$ terraform init`
1. test if everything works with `$ terraform plan`
1. Deploy with `$ terrform apply`
1. Update `.github/workflows/env/.env-prod.yaml` accordingly.
