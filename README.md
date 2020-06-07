errors.fail
===========

Probing endpoint for blackbox_exporter with configurable errors: https://errors.fail/

Configuration
-------------

Environment:

```
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/secret.json
export PROJECT_ID="dlorch-bd021"
export COOKIE_DOMAIN="errors.fail"
```

First time manual setup:
* Create Cloud Storage bucket for Terraform state ```errors-fail-terraform-state```
* Setup Cloud Build triggers for each project and sub-project (see below)
  * Assign roles "Compute Instance Admin (v1)", "DNS Administrator" and "Cloud Run Admin" to Cloud Build service account
* Add Cloud Build service account to verified owners of "errors.fail" domain in https://www.google.com/webmasters/verification/home

Sub-Projects
------------

* https://github.com/dlorch/probe.errors.fail HTTPS endpoint for probing
* https://github.com/dlorch/expired.errors.fail Expired TLS/SSL certificate and ICMP endpoint packetloss.errors.fail

Running
-------

errors.fail runs on the Google Cloud Platform. Following tools were used:
* Cloud Run
* Cloud DNS
* Compute
* Secrets Manager
* Cloud Build
* Source Repositories
* Cloud Firestore
* Cloud Storage
