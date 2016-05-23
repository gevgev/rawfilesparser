provider "aws" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  region      = "${var.region}"
}

/* App servers */
resource "aws_instance" "worker-server" {
  count = "${var.instances}"
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
      "sudo apt-gey update",
      /* Install samba client */
      /*"sudo apt-get -y install smbclient",*/
      /* Connect to shared folder */
      /*"smbclient //52.24.138.152/share",*/
      /* Install docker */ 
      "sudo apt-get -y install lxc wget bsdtar curl",
      "sudo apt-get -y install linux-image-extra-$(uname -r)",
      "sudo modprobe aufs",
      "curl -sSL https://get.docker.com/ | sh",
      "sudo usermod -aG docker ubuntu",
      /* get the files from specified location */
      "mkdir data",
      "cd data",
      "wget --no-passive-ftp ${var.ftp-url}/${var.ftp-folder}/${var.ftp-file}",
      "export DATAFOLDER=$(pwd)",
      "cd ..",
      /* Start container */
      "sudo docker run -d -v $DATAFOLDER:/data -e 'FILE=/data/${var.ftp-file}' -e 'BUCKET=${var.s3-bucket}' -e AWS_ACCESS_KEY_ID='${var.access_key}' -e AWS_SECRET_ACCESS_KEY='${var.secret_key}' gevgev/s3uploader"
      /* "sudo docker run --volumes-from ovpn-data --rm gosuri/openvpn ovpn_genconfig -p ${var.vpc_cidr} -u udp://${aws_instance.nat.public_ip}" */
    ]
  }

}
