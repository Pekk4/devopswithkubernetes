variable "project_id" {
  description = "The Google Cloud Project ID"
  type        = string
}

variable "region" {
  description = "The Google Cloud region"
  type        = string
  default     = "europe-north1"
}

variable "zone" {
  description = "The Google Cloud zone"
  type        = string
  default     = "europe-north1-b"
}
