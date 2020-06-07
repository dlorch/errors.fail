terraform {
  backend "gcs" {
    bucket = "errors-fail-terraform-state"
    prefix = "errors.fail"
  }
}
