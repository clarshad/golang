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

// run installs terraform, initializes and applies terraform configuration
func Run(tfversion string, path ...string) {
	tmpDir := getTmpDir()
	defer os.RemoveAll(tmpDir)
	execPath := installTerraform(tfversion, tmpDir)
	workingDir := getworkingDir(path)
	tf := createTfInstance(workingDir, execPath)
	tfinit(tf)
	tfapply(tf)
}

// apply runs terraform apply command
func tfapply(tf *tfexec.Terraform) {
	err := tf.Apply(context.Background())
	if err != nil {
		log.Fatalf("ERROR: error applying terraform config: %s", err)
	}
	fmt.Println("INFO: Terraform configuration applied successfully")

}

// init runs terraform init command
func tfinit(tf *tfexec.Terraform) {
	err := tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("ERROR: error running Init: %s", err)
	}
	fmt.Println("INFO: Terraform successfully initialized")
}

// createTfInstance creates a terraform object to run further commands on it
func createTfInstance(workingDir string, execPath string) *tfexec.Terraform {
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("ERROR: error running NewTerraform: %s", err)
	}
	fmt.Println("INFO: Terraform instance created successfully")
	return tf
}

// installTerraform installs specific terraform version on the given path/directory
func installTerraform(tfversion string, dir string) string {
	execPath, err := tfinstall.Find(context.Background(), tfinstall.ExactVersion(tfversion, dir))
	if err != nil {
		log.Fatalf("ERROR: error locating Terraform binary: %s", err)
	}
	fmt.Printf("INFO: Terraform version %v installed successfully", tfversion)
	return execPath
}

// getTmpDir creates a temporary directory
func getTmpDir() string {
	tmpDir, err := ioutil.TempDir("", "tfinstall")
	if err != nil {
		log.Fatalf("ERROR: error creating temp dir: %s", err)
	}
	return tmpDir
}

// getworkingDir retrieve the working directory
func getworkingDir(path []string) string {
	var workingDir string
	if len(path) != 0 {
		workingDir = path[0]
	} else {
		workingDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("ERROR: error getting current working directory: %s", err)
		}
		return workingDir
	}
	return workingDir
}
