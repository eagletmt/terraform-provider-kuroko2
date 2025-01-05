provider "kuroko2" {
  endpoint = "http://localhost:3000/v1"
  username = "dev"
  apikey   = "devkey"
}

import {
  id = "6"
  to = kuroko2_job_definition.a
}
resource "kuroko2_job_definition" "a" {
  name = "a"
  description = <<-EOT
    Greet hourly.
    No workaround available.
  EOT
  script = <<-EOS
    execute: echo 'hello'
  EOS
  admins = [1]
  cron = [
    "0 * * * *",
  ]
  tags = ["x"]
}
