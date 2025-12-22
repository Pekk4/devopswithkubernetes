terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.51.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

resource "google_storage_bucket" "postgres_backups" {
  name          = "dwk-postgres-backups-${var.project_id}"
  location      = var.region
  force_destroy = true
  storage_class = "STANDARD"
  uniform_bucket_level_access = true

  lifecycle_rule {
    condition {
      age = 3
    }
    action {
      type = "Delete"
    }
  }
  versioning {
    enabled = false
  }
}

resource "google_service_account" "backup_sa" {
  account_id   = "postgres-backup-sa"
  display_name = "Todo App Postgres Backup Service Account"
}

resource "google_storage_bucket_iam_binding" "backup_uploader" {
  bucket = google_storage_bucket.postgres_backups.name
  role   = "roles/storage.objectAdmin"

  members = [
    "serviceAccount:${google_service_account.backup_sa.email}"
  ]
}

resource "google_service_account_iam_binding" "workload_identity_binding" {
  service_account_id = google_service_account.backup_sa.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:${var.project_id}.svc.id.goog[project/postgres-backup-job-sa]"
  ]
}
