terraform {
  required_providers {
    rediscloud = {
      source  = "RedisLabs/rediscloud"
      version = "1.8.1"
    }
  }
}

provider "rediscloud" {
}

locals {
  redis_database_password = "See_H0w_Fast_Fee1s"
}

data "rediscloud_essentials_plan" "rediscloud_essentials_plan_30mb" {
  name           = "30MB"
  cloud_provider = "AWS"
  region         = "us-east-1"
}

resource "rediscloud_essentials_subscription" "rediscloud_essentials_plan_30mb" {
  name    = "rediscloud_essentials_plan_30mb"
  plan_id = data.rediscloud_essentials_plan.rediscloud_essentials_plan_30mb.id
}

resource "rediscloud_essentials_database" "redis_database" {
  subscription_id     = rediscloud_essentials_subscription.rediscloud_essentials_plan_30mb.id
  name                = "redisjson-for-gophers"
  enable_default_user = true
  password            = local.redis_database_password
  data_persistence    = "none"
  replication         = false
}

output "redis_database_url" {
    value = "redis://default:${local.redis_database_password}@${rediscloud_essentials_database.redis_database.public_endpoint}"
}
