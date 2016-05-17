resource "aws_key_pair" "dev-deployer" {
  key_name     = "dev-r31-deployer-key"
  public_key = "${file(\"ssh/insecure-r31-deployer.pub\")}"
}
