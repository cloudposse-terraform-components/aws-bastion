module "vpc" {
  source  = "cloudposse/stack-config/yaml//modules/remote-state"
  version = "1.8.0"

  component = var.vpc_component_name

  defaults = {
    vpc_id             = ""
    public_subnet_ids  = []
    private_subnet_ids = []
  }

  context = module.this.context
}
