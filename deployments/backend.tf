terraform {
  backend "gcs" {
    bucket = "errors-fail-terraform-state"
    prefix = "errors.fail"
  }
  required_version = "~> 0.13.5"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 3.51.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0.0"
    }
    docker = {
      source = "kreuzwerker/docker"
      version = "~> 2.9.0"
    }
  }
}
