locals {
  project_id   = "activityhub-408810"
  project_name = "activityhub"
  region       = "europe-west3"
  zone         = "${local.region}-a"
}

provider "google" {
  project = local.project_id
  region  = local.region
  zone    = local.region
}

provider "google-beta" {
  project = local.project_id
  region  = local.region
  zone    = local.region
}







resource "google_compute_network" "private_network" {
  provider = google-beta
  name     = "private-network"
  project  = local.project_id
}

resource "google_compute_global_address" "private_ip_address" {
  provider = google-beta

  name          = "private-ip-address"
  project       = local.project_id
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.private_network.id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  provider = google-beta

  network                 = google_compute_network.private_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "activityhub-db" {
  provider = google-beta

  name             = "${local.project_name}-db-${random_id.db_name_suffix.hex}"
  project          = local.project_id
  region           = local.region
  database_version = "POSTGRES_15"

  depends_on = [google_service_networking_connection.private_vpc_connection]

  deletion_protection = false

  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled                                  = false
      private_network                               = google_compute_network.private_network.id
      enable_private_path_for_google_cloud_services = true
    }
  }
}

resource "google_sql_database" "db" {
  name      = "activityhub"
  instance  = google_sql_database_instance.activityhub-db.name
  charset   = "UTF8"
  collation = "en_US.UTF8"
}

resource "random_password" "password" {
  length           = 16
  special          = true
  override_special = "_%@"
}


resource "google_sql_user" "activityhub-user" {
  name     = "activityhub-user"
  instance = google_sql_database_instance.activityhub-db.name
  password = random_password.password.result
}



resource "google_secret_manager_secret" "secret-db-postgres-user" {
  secret_id = "db-${google_sql_database_instance.activityhub-db.name}-${google_sql_user.activityhub-user.name}"
  replication {
    auto {}

  }
}

resource "google_secret_manager_secret_version" "secret-db-postgres-user-1" {
  provider = google-beta

  secret      = google_secret_manager_secret.secret-db-postgres-user.id
  secret_data = random_password.password.result
}




