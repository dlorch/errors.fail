resource "google_cloud_run_service" "errors-fail-service" {
  autogenerate_revision_name = "true"
  name                       = "errors-fail"
  location                   = "us-east1"

  template {
    spec {
      containers {
        image = data.google_container_registry_image.errors_fail_latest.image_url
        env {
          name  = "PROJECT_ID"
          value = "dlorch-bd021"
        }
        env {
          name  = "COOKIE_DOMAIN"
          value = "errors.fail"
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "errors-fail-iam-member" {
  location = google_cloud_run_service.errors-fail-service.location
  project  = google_cloud_run_service.errors-fail-service.project
  service  = google_cloud_run_service.errors-fail-service.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

// Make sure to add the Cloud Build service account to verified owners of domain
// in https://www.google.com/webmasters/verification/home
resource "google_cloud_run_domain_mapping" "errors-fail-domain-mapping" {
  location = "us-east1"
  name     = "errors.fail"

  metadata {
    namespace = "dlorch-bd021"
  }

  spec {
    route_name = google_cloud_run_service.errors-fail-service.name
  }
}

// in order to update the Cloud Run deployment when the underlying :latest
// image changes, we will retrieve the sha256_digest of the image through
// the docker provider, since the google container registry does not seem
// to have an API. Thanks @c2thorn for the pointers
// https://github.com/terraform-providers/terraform-provider-google/issues/6706#issuecomment-652009984
// https://github.com/terraform-providers/terraform-provider-google/issues/6635#issuecomment-647858867
data "google_client_config" "default" {}

provider "docker" {
  registry_auth {
    address  = "gcr.io"
    username = "oauth2accesstoken"
    password = data.google_client_config.default.access_token
  }
}

data "docker_registry_image" "errors_fail_image" {
  name = "gcr.io/dlorch-bd021/errors-fail"
}

data "google_container_registry_image" "errors_fail_latest" {
  name    = "errors-fail"
  project = "dlorch-bd021"
  digest  = data.docker_registry_image.errors_fail_image.sha256_digest
}
