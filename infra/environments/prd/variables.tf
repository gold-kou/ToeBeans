variable rds_password {
  type        = string
  default     = "VeryStrongPassword!"
}

variable codepipeline_webhook_secret_token {
  type        = string
  default     = "VeryRandomStringMoreThan20Byte!"
}

variable github_webhook_secret {
  type        = string
  default     = "VeryRandomStringMoreThan20Byte!"
}