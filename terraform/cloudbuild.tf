resource "google_project_service" "cloudbuild" {
  service = "cloudbuild.googleapis.com"
}

resource "google_cloudbuild_trigger" "github_audiy-importer" {
  name = "audiy-importer-trigger"

  trigger_template {
    repo_name   = "github_hayashiki_audiy-importer"
    branch_name = "master"
  }

  filename = "cloudbuild.yaml"
}

locals {
  cloudbuild_roles = [
    "roles/run.admin",
  ]
}

// cloudbuildサービスアカウントにroleを付与する
resource "google_project_iam_binding" "cloudbuild" {
  for_each = toset(local.cloudbuild_roles)
  role     = each.value

  members = ["serviceAccount:${data.google_project.project.number}@cloudbuild.gserviceaccount.com"]
}
