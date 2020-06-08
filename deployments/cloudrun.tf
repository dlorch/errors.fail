resource "google_cloud_run_service" "errors-fail-service" {
  name     = "errors-fail"
  location = "us-east1"

  template {
    spec {
      containers {
        // the ":latest" tag is very important, otherwise Terraform will not
        // re-deploy the Cloud Run service even if the image changes
        image = "gcr.io/dlorch-bd021/errors-fail:latest"
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

// Make sure to add service account to verified owners of domain
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
