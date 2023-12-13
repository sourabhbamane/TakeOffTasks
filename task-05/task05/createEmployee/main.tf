provider "google" {
  project = "employee-management-403415"
  region  = "us-central1"
  zone    = "us-central1-c"
}

resource "google_storage_bucket" "bucket" {
  name = "terraform-create-emp"
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "function.zip"
}



#2nd Gen
resource "google_cloudfunctions2_function" "default" {
  name        = "dev-usc1-dispatch-post-employee"
  location    = "us-central1"
  description = "Function To Add Employee Details"

  build_config {
    runtime     = "go121"
    entry_point = "AddEmployee" # Set the entry point
    source {
      storage_source {
        bucket = google_storage_bucket.bucket.name
        object = google_storage_bucket_object.archive.name
      }
    }
  }

  service_config {
    min_instance_count = 1
    max_instance_count = 10
    available_memory   = "256M"
    timeout_seconds    = 60
    all_traffic_on_latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "member" {
  location = google_cloudfunctions2_function.default.location
  service  = google_cloudfunctions2_function.default.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}