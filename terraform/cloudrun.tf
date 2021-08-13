resource "google_cloud_run_service" "audiy_api" {
  name     = local.service_name
  location = var.region

  template {
    spec {
      containers {
        image = "${local.image_fullname_api}:latest"

        resources {
          limits = {
            "cpu" : "1000m",
            "memory" : "128Mi",
          }
        }

        dynamic "env" {
          for_each = merge(
            local.env,
          )

          content {
            name  = env.key
            value = env.value
          }
        }
      }

      timeout_seconds = 10
    }
  }

  depends_on = [
    google_project_service.services["run.googleapis.com"],
    google_project_iam_member.adminapi-observability,
  ]

  lifecycle {
    ignore_changes = [
      template[0].spec[0].containers[0].image,
    ]
  }
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.audiy_api.location
  project  = google_cloud_run_service.audiy_api.project
  service  = google_cloud_run_service.audiy_api.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

resource "google_service_account_iam_member" "admin-account-iam" {
  service_account_id = google_service_account.adminapi.id
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${local.cloudbuild_email}"
}

output "adminapi_urls" {
  value = google_cloud_run_service.audiy_api.status.0.url
}
