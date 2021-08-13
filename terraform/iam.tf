resource "google_service_account" "adminapi" {
  project      = var.project
  account_id   = "audiy-adminapi-sa"
  display_name = "Verification Admin api"
}

resource "google_service_account_iam_member" "cloudbuild-deploy-adminapi" {
  service_account_id = google_service_account.adminapi.id
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${local.cloudbuild_email}"
}

resource "google_project_iam_member" "adminapi-observability" {
  for_each = local.observability_iam_roles
  project  = var.project
  role     = each.key
  member   = "serviceAccount:${google_service_account.adminapi.email}"
}

locals {
  # observability_iam_roles is the list of IAM roles required for the service to
  # participate in observability.
  observability_iam_roles = toset([
    "roles/cloudtrace.agent",
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter",
    "roles/stackdriver.resourceMetadata.writer",
  ])
}
