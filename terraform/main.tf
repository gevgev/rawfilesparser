provider "aws" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  region      = "${var.region}"
}

/* App servers */
resource "aws_instance" "worker-server" {
  count = 1
  ami = "${lookup(var.amis, var.region)}"
  instance_type = "t2.micro"
  subnet_id = "${var.subnet_id}"
  vpc_security_group_ids = ["${var.vpc_security_group_id}"]
  key_name = "${aws_key_pair.dev-deployer.key_name}"
  source_dest_check = false
  tags = { 
    Name = "R31-app-server-${count.index}"
  }
  connection {
    user = "ubuntu"
    key_file = "ssh/insecure-r31-deployer"
  }
  provisioner "remote-exec" {
    inline = [
      /* Install docker */ 
      "sudo apt-gey update",
      "sudo apt-get -y install lxc wget bsdtar curl",
      "sudo apt-get -y install linux-image-extra-$(uname -r)",
      "sudo modprobe aufs",
      "curl -sSL https://get.docker.com/ | sh",
      "sudo usermod -aG docker ubuntu",
      /* Start my container */
      "sudo docker run -d -p 80:80 gevgev/contributors"
      /* "sudo docker run --volumes-from ovpn-data --rm gosuri/openvpn ovpn_genconfig -p ${var.vpc_cidr} -u udp://${aws_instance.nat.public_ip}" */
    ]
  }

}
