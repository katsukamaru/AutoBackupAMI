# Auto Backup AMI

## What is this?
AutoBackupAMI creates Automatically EC2 AMI and deletes the AMI 7days after.

## Requirement

IAM:
- ec2:CreateImage
- ec2:DeregisterImage
- ec2:DescribeImages
- ec2:DescribeInstances

## Run
Add EC2 tag `AutoAMI : true` to your instance to backup.
Creating AMI runs with NoReboot.
 
```
## how to run
$ go main.go
```