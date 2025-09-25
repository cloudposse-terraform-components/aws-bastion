package test

import (
	"context"
	"strings"
	"testing"

	"github.com/cloudposse/test-helpers/pkg/atmos"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	awshelper "github.com/cloudposse/test-helpers/pkg/aws"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/assert"
)

type ComponentSuite struct {
	helper.TestSuite
}

func (s *ComponentSuite) TestBasic() {
	const component = "bastion/basic"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	defer s.DestroyAtmosComponent(s.T(), component, stack, nil)
	options, _ := s.DeployAtmosComponent(s.T(), component, stack, nil)
	assert.NotNil(s.T(), options)

	iamInstanceProfile := atmos.Output(s.T(), options, "iam_instance_profile")
	assert.True(s.T(), strings.HasPrefix(iamInstanceProfile, "eg-default-ue2-test-bastion"))

	autoscalingGroupId := atmos.Output(s.T(), options, "autoscaling_group_id")
	assert.True(s.T(), strings.HasPrefix(autoscalingGroupId, iamInstanceProfile))

	securityGroupId := atmos.Output(s.T(), options, "security_group_id")
	assert.True(s.T(), strings.HasPrefix(securityGroupId, "sg-"))

	instanceIds := aws.GetInstanceIdsForAsg(s.T(), autoscalingGroupId, awsRegion)
	assert.Equal(s.T(), 1, len(instanceIds))

	instance := awshelper.GetEc2Instances(s.T(), context.Background(), instanceIds[0], awsRegion)
	assert.EqualValues(s.T(), "t2.micro", instance.InstanceType)
	assert.EqualValues(s.T(), "running", *&instance.State.Name)

	s.DriftTest(component, stack, nil)
}

func (s *ComponentSuite) TestEnabledFlag() {
	const component = "bastion/disabled"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	s.VerifyEnabledFlag(component, stack, nil)
}


func TestRunSuite(t *testing.T) {
	suite := new(ComponentSuite)

	suite.AddDependency(t, "vpc", "default-test", nil)
	helper.Run(t, suite)
}
