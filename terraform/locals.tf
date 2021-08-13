locals {
  service_name       = "audiy-api"
  image_fullname_api = "gcr.io/${var.project}/${local.service_name}"

  env = {
    SLACK_BOT_TOKEN        = "xoxb-12376436915"
    GCS_INPUT_AUDIO_BUCKET = "audiy-input"
    TOPIC_NAME             = "audioHandleTopic"
  }

  cloudbuild_email = "${data.google_project.project.number}@cloudbuild.gserviceaccount.com"
}
