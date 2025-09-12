module "vpc" {
  source  = "cloudposse/stack-config/yaml//modules/remote-state"
  version = "1.5.0"

  component = "vpc"

  defaults = {
    vpc_id             = ""
    public_subnet_ids  = []
    private_subnet_ids = []
  }

  context = module.this.context
}
