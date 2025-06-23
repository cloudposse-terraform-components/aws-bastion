module "vpc" {
  source  = "cloudposse/stack-config/yaml//modules/remote-state"
  version = "1.8.0"

  component = var.vpc_component_name

  context = module.this.context
}
