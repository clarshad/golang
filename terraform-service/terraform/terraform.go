package terraform

import (
	"context"
	"fmt"
	"os"

	"github.com/clarshad/golang/terraform-service/utils"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"gopkg.in/src-d/go-git.v4"
)

// run installs terraform with provided version, initializes and applies terraform configuration
func Run(tfversion string, action string, path string) error {
	wd, _ := os.Getwd()

	dstDir := wd + "/repo"
	tfconfigDir, err := getConfigDir(path, dstDir)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dstDir)

	tfinstallDir := wd + "/tfinstall"
	os.Mkdir(tfinstallDir, 0755)
	tf, err := installTerraformAndCreateInstance(tfinstallDir, tfconfigDir, tfversion)
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

// getworkingDir retrieve the working directory
func getConfigDir(srcpath string, dstpath string) (string, error) {
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

	tfcd := dstpath + "/" + srcpath
	utils.Log(fmt.Sprintf("INFO: Terraform: Running terraform configuration from directory %v", tfcd))
	return tfcd, nil
}

func installTerraformAndCreateInstance(installDir string, configDir string, tfversion string) (*tfexec.Terraform, error) {
	execPath := installDir + "/terraform"
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		utils.Log(fmt.Sprintf("WARNING: Terraform: Installation binary %v not found", execPath))
		_, err := installTerraform(tfversion, installDir)
		if err != nil {
			return nil, err
		}
	} else {
		utils.Log(fmt.Sprintf("INFO: Terraform: Installation binary %v found, skipping Terraform Install", execPath))
	}

	tf, err := createTfInstance(configDir, execPath)
	if err != nil {
		return nil, err
	}

	v, _, err := tf.Version(context.Background(), true)
	if err != nil {
		utils.Log(fmt.Sprintf("ERROR: Terraform: Unable to retrieve terraform version: %s", err))
		return nil, err
	}
	if v.String() == tfversion {
		return tf, nil
	}

	utils.Log(fmt.Sprintf("WARNING: Terraform: Version mistmatch, expected %v, found %v", tfversion, v.String()))
	os.Remove(execPath)
	tf, err = installTerraformAndCreateInstance(installDir, configDir, tfversion)
	return tf, err
}

// installTerraform installs specific terraform version on the given path/directory
func installTerraform(tfversion string, dir string) (string, error) {
	execPath, err := tfinstall.Find(context.Background(), tfinstall.ExactVersion(tfversion, dir))
	if err != nil {
		utils.Log(fmt.Sprintf("ERROR: Terraform: Unable to locate and install Terraform binary: %s", err))
		return "", err
	}

	utils.Log(fmt.Sprintf("INFO: Terraform: Version %v installed successfully", tfversion))
	return execPath, nil
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
