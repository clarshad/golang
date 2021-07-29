package terraform

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

// run installs terraform with provided version, initializes and applies terraform configuration
func Run(tfversion string, path ...string) error {
	tmpDir, err := getTmpDir("", "tfinstall")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmpDir)

	execPath, err := installTerraform(tfversion, tmpDir)
	if err != nil {
		return err
	}

	workingDir, err := getworkingDir(path)
	if err != nil {
		return err
	}

	tf, err := createTfInstance(workingDir, execPath)
	if err != nil {
		return err
	}

	err = tfinit(tf)
	if err != nil {
		return err
	}

	err = tfapply(tf)
	if err != nil {
		return err
	}
	return nil
}

// getTmpDir creates a temporary directory
func getTmpDir(dir string, pattern string) (string, error) {
	tmpDir, err := ioutil.TempDir(dir, pattern)
	if err != nil {
		fmt.Printf("ERROR: Terraform: Unable to create temporary directory: %s\n", err)
		return "", err
	}

	fmt.Printf("INFO: Terraform: Temporary directory %v created for terraform installation\n", tmpDir)
	return tmpDir, nil
}

// installTerraform installs specific terraform version on the given path/directory
func installTerraform(tfversion string, dir string) (string, error) {
	execPath, err := tfinstall.Find(context.Background(), tfinstall.ExactVersion(tfversion, dir))
	if err != nil {
		fmt.Printf("ERROR: Terraform: Unable to install and locate Terraform binary: %s\n", err)
		return "", err
	}

	fmt.Printf("INFO: Terraform: Version %v installed successfully\n", tfversion)
	return execPath, nil
}

// getworkingDir retrieve the working directory
func getworkingDir(path []string) (string, error) {
	var wd string
	if path[0] != "" {
		wd = path[0]
	} else {
		d, err := os.Getwd()
		if err != nil {
			fmt.Printf("ERROR: Terraform: Unable to get current working directory: %s\n", err)
			return "", err
		}
		wd = d + "/scripts/terraform-config"
	}

	fmt.Printf("INFO: Terraform: Running terraform configuration from directory %v\n", wd)
	return wd, nil
}

// createTfInstance creates a terraform object to run further commands on it
func createTfInstance(workingDir string, execPath string) (*tfexec.Terraform, error) {
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		fmt.Printf("ERROR: Terraform: Unable to run NewTerraform instance: %s\n", err)
		return nil, err
	}
	fmt.Println("INFO: Terraform: Instance for terraform object created successfully")
	return tf, nil
}

// init runs terraform init command
func tfinit(tf *tfexec.Terraform) error {
	err := tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Printf("\nERROR: Terraform: Unable to run terraform initialization: %s\n", err)
		return err
	}
	fmt.Println("INFO: Terraform: Successfully initialized, 'terraform init' command equivalent")
	return nil
}

// apply runs terraform apply command
func tfapply(tf *tfexec.Terraform) error {
	err := tf.Apply(context.Background())
	if err != nil {
		fmt.Printf("ERROR: Terraform: Unable to apply terraform configuration: %s\n", err)
		return err
	}
	fmt.Println("INFO: Terraform: Configuration applied successfully, 'terraform apply' command equivalent")
	return nil
}
