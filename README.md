errors.fail
===========

Probing endpoint for blackbox_exporter with configurable errors: https://errors.fail/

Environment
-----------

Required environment variables to run the Golang webserver:

```
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/secret.json
export PROJECT_ID="dlorch-bd021"
export COOKIE_DOMAIN="errors.fail"
```

For the Google application credentials, create a service account for the project and
assign "Owner" privileges to it, then retrieve the service account's secret as .json
file.

Run the server with ```go run main.go```.

Sub-Projects
------------

* https://github.com/dlorch/probe.errors.fail HTTPS endpoint for probing
* https://github.com/dlorch/expired.errors.fail Expired TLS/SSL certificate and ICMP endpoint packetloss.errors.fail

Infrastructure as Code - CI/CD
------------------------------

All required infrastructure is described "as code" in Terraform. errors.fail runs on
the Google Cloud Platform. Following tools were used:
* Cloud Run
* Cloud DNS
* Compute
* Secrets Manager
* Cloud Build
* Source Repositories
* Cloud Firestore
* Cloud Storage

First time manual project setup:
* Create Cloud Storage bucket for Terraform state ```errors-fail-terraform-state```
* Setup Cloud Build triggers for each project and sub-project (see below)
  * Assign roles "Compute Instance Admin (v1)", "DNS Administrator" and "Cloud Run Admin" to Cloud Build service account
* Add Cloud Build service account to verified owners of "errors.fail" domain in https://www.google.com/webmasters/verification/home
* TLS/SSL certificate for expired.errors.fail created with Let's Encrypt and Certbot as described here: https://certbot.eff.org/lets-encrypt/debianstretch-nginx, then take a copy with ```sudo tar cf letsencrypt.tar /etc/letsencrypt``` and store it as a secret ```expired-errors-fail_letsencrypt-tar``` in the Secret Manager
* Give Compute service account "Secret Manager Secret Accessor" privileges

Enjoy! Daniel Lorch, 2020
