package terraform

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/clarshad/golang/terraform-service/utils"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"gopkg.in/src-d/go-git.v4"
)

// run installs terraform with provided version, initializes and applies terraform configuration
func Run(tfversion string, action string, path ...string) error {
	tmpDir, err := getTmpDir("", "tfinstall")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	execPath, err := installTerraform(tfversion, tmpDir)
	if err != nil {
		return err
	}

	wd, _ := os.Getwd()
	dstDir := wd + "/repo"
	configDir, err := getConfigDir(path, dstDir)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dstDir)

	tf, err := createTfInstance(configDir, execPath)
	if err != nil {
		return err
	}

	err = tfInit(tf)
	if err != nil {
		return err
	}

	err = tfApply(tf)
	if err != nil {
		return err
	}

	if action == "destroy" {
		err = tfDestroy(tf)
		if err != nil {
			return err
		}
	}

	return nil
}

// getTmpDir creates a temporary directory
func getTmpDir(dir string, pattern string) (string, error) {
	tmpDir, err := ioutil.TempDir(dir, pattern)
	if err != nil {
		utils.Log(fmt.Sprintf("ERROR: Terraform: Unable to create temporary directory: %s", err))
		return "", err
	}

	utils.Log(fmt.Sprintf("INFO: Terraform: Temporary directory %v created for terraform installation", tmpDir))
	return tmpDir, nil
}

// installTerraform installs specific terraform version on the given path/directory
func installTerraform(tfversion string, dir string) (string, error) {
	execPath, err := tfinstall.Find(context.Background(), tfinstall.ExactVersion(tfversion, dir))
	if err != nil {
		utils.Log(fmt.Sprintf("ERROR: Terraform: Unable to install and locate Terraform binary: %s", err))
		return "", err
	}

	utils.Log(fmt.Sprintf("INFO: Terraform: Version %v installed successfully", tfversion))
	return execPath, nil
}

// getworkingDir retrieve the working directory
func getConfigDir(srcpath []string, dstpath string) (string, error) {
	username := os.Getenv("GIT_USERNAME")
	password := os.Getenv("GIT_PASSWORD")
	repo := os.Getenv("GIT_REPOSITORY")

	url := fmt.Sprintf("https://%s:%s@%s", username, password, repo)
	_, err := git.PlainClone(dstpath, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		utils.Log(fmt.Sprintf("ERROR: Terraform: Unable to git clone respository %v, error: %v", repo, err))
		return "", err
	}

	tfcd := dstpath + "/" + srcpath[0]
	utils.Log(fmt.Sprintf("INFO: Terraform: Running terraform configuration from directory %v", tfcd))
	return tfcd, nil
}

// createTfInstance creates a terraform object to run further commands on it
func createTfInstance(workingDir string, execPath string) (*tfexec.Terraform, error) {
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		utils.Log(fmt.Sprintf("ERROR: Terraform: Unable to run NewTerraform instance: %s", err))
		return nil, err
	}

	utils.Log("INFO: Terraform: Instance for terraform object created successfully")
	return tf, nil
}

// init runs terraform init command
func tfInit(tf *tfexec.Terraform) error {
	err := tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		utils.Log(fmt.Sprintf("ERROR: Terraform: Unable to run terraform initialization: %s", err))
		return err
	}

	utils.Log("INFO: Terraform: Successfully initialized, equivalent to 'terraform init' command")
	return nil
}

// apply runs terraform apply command
func tfApply(tf *tfexec.Terraform) error {
	err := tf.Apply(context.Background())
	if err != nil {
		utils.Log(fmt.Sprintf("ERROR: Terraform: Unable to apply terraform configuration: %s", err))
		return err
	}

	utils.Log("INFO: Terraform: Configuration applied successfully, equivalent to 'terraform apply' command")
	return nil
}

func tfDestroy(tf *tfexec.Terraform) error {
	err := tf.Destroy(context.Background())
	if err != nil {
		utils.Log(fmt.Sprintf("ERROR: Terraform: Unable to destroy terraform configuration: %s", err))
		return err
	}

	utils.Log("INFO: Terraform: Destroyed configuration successfully, equivalent to 'terraform destroy' command")
	return nil
}
