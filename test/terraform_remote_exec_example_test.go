package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// This test shows how to override the systems local SSH Agent, with an in-process SSH agent, whose keys can be managed
// from within your tests. This allows you to test Terraform modules which make SSH connections to the created
// instances, useful for tasks such as provisioning.
func TestTerraformRemoteExecExample(t *testing.T) {
	t.Parallel()

	terraformDirectory := "../examples/terraform-remote-exec-example"

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, terraformDirectory)
		keyPair := test_structure.LoadEc2KeyPair(t, terraformDirectory)

		// destroy terraform resources and delete ec2 key pair
		terraform.Destroy(t, terraformOptions)
		aws.DeleteEC2KeyPair(t, keyPair)

		// remove testFile, if it exists
		testFile := filepath.Join(terraformDirectory, "public-ip")
		if _, err := os.Stat(testFile); err == nil {
			os.Remove(testFile)
		}
	})

	// Deploy the example
	test_structure.RunTestStage(t, "setup", func() {

		// A unique ID we can use to namespace resources so we don't clash with anything already in the AWS account or
		// tests running in parallel
		uniqueID := random.UniqueId()

		// Give this EC2 Instance and other resources in the Terraform code a name with a unique ID so it doesn't clash
		// with anything else in the AWS account.
		instanceName := fmt.Sprintf("terratest-remote-exec-example-%s", uniqueID)

		// Pick a random AWS region to test in. This helps ensure your code works in all regions.
		awsRegion := aws.GetRandomStableRegion(t, nil, nil)

		// Create an EC2 KeyPair that we can use for SSH access
		keyPairName := fmt.Sprintf("terratest-remote-exec-example-%s", uniqueID)
		keyPair := aws.CreateAndImportEC2KeyPair(t, awsRegion, keyPairName)

		// start an SSH agent, with our key pair added
		sshAgent := ssh.SshAgentWithKeyPair(t, keyPair.KeyPair)
		defer sshAgent.Stop()

		terraformOptions := terraform.NewTerraformOptionsWithDefaultRetryableErrors(
			// The path to where our Terraform code is located
			terraformDirectory,
			// Variables to pass to our Terraform code using -var options
			map[string]interface{}{
				"aws_region":    awsRegion,
				"instance_name": instanceName,
				"key_pair_name": keyPairName,
			},
		)
		// Override local SSH agent with our new agent
		terraformOptions.SshAgent = sshAgent

		// Save the options and key pair so later test stages can use them
		test_structure.SaveTerraformOptions(t, terraformDirectory, terraformOptions)
		test_structure.SaveEc2KeyPair(t, terraformDirectory, keyPair)

		// Because of the SshAgent option above, the terraform process will be provided an `SSH_AUTH_SOCK` environment
		// variable, which will point to the socket file of our in-process `sshAgent` instance:
		terraform.InitAndApply(t, terraformOptions)

		// save the `public_instance_ip` output variable for later steps
		publicIP := terraform.Output(t, terraformOptions, "public_instance_ip")
		test_structure.SaveString(t, terraformDirectory, "publicIP", publicIP)
	})

	// Make sure we can SSH to the public Instance directly from the public Internet and the private Instance by using
	// the public Instance as a jump host
	test_structure.RunTestStage(t, "validate", func() {
		publicIP := test_structure.LoadString(t, terraformDirectory, "publicIP")

		// Confirm that the public-ip file that was generated by the provisioner was copied back from the server using
		// the `scp` command
		testFile := filepath.Join(terraformDirectory, "public-ip")
		assert.FileExists(t, testFile)

		// Check that public IP from output matches public IP generated by script on the server
		b, err := ioutil.ReadFile(testFile)
		if err != nil {
			fmt.Print(err)
		}
		assert.Equal(t, strings.TrimSpace(publicIP), strings.TrimSpace(string(b)))
	})

}
