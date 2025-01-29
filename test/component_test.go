package test

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/aws-component-helper"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/assert"
)

func TestComponent(t *testing.T) {
	t.Parallel()
	// Define the AWS region to use for the tests
	awsRegion := "us-east-2"

	// Initialize the test fixture
	fixture := helper.NewFixture(t, "../", awsRegion, "test/fixtures")

	// Ensure teardown is executed after the test
	defer fixture.TearDown()
	fixture.SetUp(&atmos.Options{})

	// Define the test suite
	fixture.Suite("default", func(t *testing.T, suite *helper.Suite) {
		t.Parallel()
		suite.AddDependency("vpc", "default-test")

		// Test phase: Validate the functionality of the bastion component
		suite.Test(t, "basic", func(t *testing.T, atm *helper.Atmos) {

			defer atm.GetAndDestroy("bastion/basic", "default-test", map[string]interface{}{})
			component := atm.GetAndDeploy("bastion/basic", "default-test", map[string]interface{}{})
			assert.NotNil(t, component)

			iamInstanceProfile := atm.Output(component, "iam_instance_profile")
			assert.True(t, strings.HasPrefix(iamInstanceProfile, "eg-default-ue2-test-bastion"))

			autoscalingGroupId := atm.Output(component, "autoscaling_group_id")
			assert.True(t, strings.HasPrefix(autoscalingGroupId, iamInstanceProfile))

			securityGroupId := atm.Output(component, "security_group_id")
			assert.True(t, strings.HasPrefix(securityGroupId, "sg-"))

			instanceIds := aws.GetInstanceIdsForAsg(t, autoscalingGroupId, awsRegion)
			assert.Equal(t, 1, len(instanceIds))

			instance := GetEc2Instances(t, instanceIds[0], awsRegion)
			assert.EqualValues(t, "t2.micro", instance.InstanceType)
			assert.EqualValues(t, "running", *&instance.State.Name)
		})
	})
}

// GetPrivateIpsOfEc2InstancesE gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
func GetEc2Instances(t *testing.T, instanceID string, awsRegion string) types.Instance {
	ec2Client := aws.NewEc2Client(t, awsRegion)
	// TODO: implement pagination for cases that extend beyond limit (1000 instances)
	input := ec2.DescribeInstancesInput{InstanceIds: []string{instanceID}}
	output, err := ec2Client.DescribeInstances(context.Background(), &input)
	assert.NoError(t, err)

	return output.Reservations[0].Instances[0]
}
