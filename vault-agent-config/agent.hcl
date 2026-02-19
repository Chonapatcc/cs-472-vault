pid_file = "/tmp/vault-agent.pid"

vault {
  address = "http://vault:8200"
  retry {
    num_retries = 5
  }
}

auto_auth {
  method {
    type = "token_file"
    config = {
      token_file_path = "/vault/config/token"
    }
  }

  sink "file" {
    config = {
      path = "/vault/config/token"
    }
  }
}

template {
  source      = "/vault/config/mongo-root-password.tpl"
  destination = "/secrets/mongo_root_password"
}

template {
  source      = "/vault/config/mongo-root-username.tpl"
  destination = "/secrets/mongo_root_username"
} 

template {
  source      = "/vault/config/mongo.tpl"
  destination = "/secrets/mongo.env"
}

listener "tcp" {
   address     = "0.0.0.0:8100"
   tls_disable = true
}

api_proxy {
   use_auto_auth_token = true
   enforce_consistency = "always"
}

cache {
   // Requires Vault Enterprise 1.16 or later
   cache_static_secrets = true
   static_secret_token_capability_refresh_interval = "5m"
}
