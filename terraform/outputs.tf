output "app-server.0.ip" {
  value = "${aws_instance.worker-server.0.private_ip}"
}

/*
output "app-server.1.ip" {
  value = "${aws_instance.angoss-server.1.private_ip}"
}
*/
