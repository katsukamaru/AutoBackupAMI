package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	sess := session.Must(session.NewSession())

	awsRegion := "ap-northeast-1"
	svc := ec2.New(sess, &aws.Config{Region: aws.String(awsRegion)})
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:AutoAMI"),
				Values: []*string{
					aws.String("true"),
				},
			},
		},
	}

	resp, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", awsRegion, err.Error())
		log.Fatal(err.Error())
	}

	// instance ids
	var targetMap = map[string]string{}

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			var name string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}
			targetMap[name] = *instance.InstanceId
		}
	}

	t := time.Now()

	// This is the right format. Do not edit this date.
	// See more: https://kakakakakku.hatenablog.com/entry/2016/03/28/001145
	const layout = "20060102"
	today := t.Format(layout)

	for name, instanceId := range targetMap {
		input := ec2.CreateImageInput{
			InstanceId: aws.String(instanceId),
			Name:       aws.String(name + "_AutoAMI_" + today),
			NoReboot:   aws.Bool(true),
		}
		output, e := svc.CreateImage(&input)
		if e != nil {
			fmt.Printf("failed create Image by something: %+v", e.Error())
			break
		}
		fmt.Printf("%+v", *output.ImageId)
	}

	createDate := t.AddDate(0, 0, -7).Format(layout)
	filters := []*ec2.Filter{
		{
			Name:   aws.String("name"),
			Values: []*string{aws.String("*_AutoAMI_" + createDate)},
		},
	}

	imagesInput := ec2.DescribeImagesInput{Filters: filters}
	output, e := svc.DescribeImages(&imagesInput)
	if e != nil {
		log.Printf("Failed describeImages : %+v", e.Error())
	}

	for _, image := range output.Images {
		input := ec2.DeregisterImageInput{ImageId: image.ImageId}
		imageOutput, e := svc.DeregisterImage(&input)
		if e != nil {
			log.Printf("Failed deregister Image: %+v", e.Error())
		}
		fmt.Printf("%+v", imageOutput.GoString())
	}
}
