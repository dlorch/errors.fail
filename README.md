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
