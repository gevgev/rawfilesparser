variable "access_key" { 
  description = "AWS access key"
}

variable "secret_key" { 
  description = "AWS secret access key"
}

variable "region"     				{ default = "us-west-2" }
variable "zone" 					{ default = "us-west-2a"}
variable "vpc_id"       			{ default = "vpc-a88757cc" }
variable "subnet_id"    			{ default = "subnet-dde815ab" }
variable "vpc_security_group_id"    { default = "sg-74f8ac13" }

/* Ubuntu 14.04 amis by region */
variable "amis" {
  description = "Base AMI to launch the instances with"
  default = {
    us-west-1 = "ami-049d8641" 
    us-east-1 = "ami-a6b8e7ce"
    us-west-2 = "ami-9abea4fb"
  }
}

variable "instances" {
	description = "Number of angoss worker instances to launch. 0 = shut down everything"
}

variable "ftp-folder" {
	description = "FTP path to folder"
}

variable "ftp-file" {
	description = "File name"
}

variable "s3-bucket" {
	description = "Target S3 bucket full name"
}

variable "ftp-url" {
	description = "FTP URL"
	default = "ftp://adsawsftp:quickaccess123@ec2-52-24-138-152.us-west-2.compute.amazonaws.com"
}
