//+build ignore
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/cli/util"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: unable to determine working directory - %s\n", err)
		os.Exit(1)
	}

	var projectPath = filepath.Dir(currentDir)

	project := api.NewAppProject(projectPath)

	azPath, err := project.GetPath("github.com/Azure/azure-functions-go")

	if azPath == "" {
		fmt.Println("Not able to find Azure functions lib")
		os.Exit(1)
	}

	triggerPath, err := project.GetPath("github.com/skothari-tibco/flogoaztrigger")

	if err != nil {
		fmt.Println("Error in determining Trigger Path")
		os.Exit(1)
	}

	//Copy AzGoFunc in the App Dir
	azPath, err = copyAzGo(azPath, currentDir)

	//Set up main.go and function.Json
	err = copyMain(triggerPath, azPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Set up shim folder.
	err = copyShim(currentDir, azPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//Temporary Hack, Copy azfunc into the app Dir. Need to make azfunc as seperate repo.

	azFuncPath, err := project.GetPath("github.com/Azure/azure-functions-go/azfunc")

	err = copyAzFunc(azPath, azFuncPath)
	buildImage(azPath)

}

/*
	Copy Main.go, function.json into the TempHttp Dir.
*/
func copyMain(triggerPath string, azPath string) error {
	currPath, err := os.Getwd()

	if err != nil {
		return err
	}

	defer os.Chdir(currPath)

	err = os.Chdir(filepath.Join(azPath, "sample"))

	if err != nil {
		return err
	}

	err = os.Mkdir("TempHttp", os.ModePerm)

	if err != nil {
		return err
	}

	destPath := filepath.Join(azPath, "sample", "TempHttp")

	srcPath := filepath.Join(triggerPath, "shim")

	climd, err := exec.Command("cp", filepath.Join(srcPath, "shim.go"), destPath).Output()
	if err != nil {
		fmt.Println(string(climd))
		return err
	}
	climd, err = exec.Command("cp", filepath.Join(srcPath, "function.json"), destPath).Output()
	if err != nil {
		fmt.Println(string(climd))
		return err
	}
	err = os.Rename(filepath.Join(destPath, "shim.go"), filepath.Join(destPath, "main.go"))
	if err != nil {
		return err
	}
	return nil
}

/*
	Populate the shim folder with shim_support and app.
*/
func copyShim(projectPath string, azPath string) error {
	currPath, err := os.Getwd()

	if err != nil {
		return err
	}

	defer os.Chdir(currPath)

	err = os.Chdir(filepath.Join(azPath, "sample", "TempHttp"))
	if err != nil {
		return err
	}

	err = os.Mkdir("shim", os.ModePerm)
	if err != nil {
		return err
	}

	destPath := filepath.Join(azPath, "sample", "TempHttp", "shim")

	err = renameAndCopyFile(projectPath, "shim_support.go", destPath)

	if err != nil {

		return err
	}

	err = renameAndCopyFile(projectPath, "embeddedapp.go", destPath)

	if err != nil {

		return err
	}

	err = renameAndCopyFile(projectPath, "imports.go", destPath)

	if err != nil {

		return err
	}

	return nil
}

/*
	Rename the package to shim and copy the files.
*/
func renameAndCopyFile(projectPath string, fileName string, destPath string) error {

	data, err := ioutil.ReadFile(filepath.Join(projectPath, fileName))

	if err != nil {
		return err
	}
	text := strings.Replace(string(data), "package main", "package shim", 1)

	err = ioutil.WriteFile(filepath.Join(destPath, fileName), []byte(text), 0755)
	if err != nil {
		return err
	}
	return nil
}

func buildImage(azPath string) error {
	currPath, err := os.Getwd()

	if err != nil {
		return err
	}

	defer os.Chdir(currPath)

	err = util.ExecCmd(exec.Command("sh", "build_container.sh"), filepath.Join(azPath, "test"))

	if err != nil {

		return err
	}

	return nil
}

/*
	Copy the Azure-functions-go lib from pkg/mod into the App Dir so that you can add and build the Go functions.
*/

func copyAzGo(azPath, currentDir string) (string, error) {

	_, err := os.Stat(filepath.Join(azPath))

	if err != nil {
		fmt.Println("Error in determining Azure Path")
		return "", err
	}

	climd, err := exec.Command("cp", "-r", filepath.Join(azPath), currentDir).Output()
	if err != nil {
		fmt.Println(string(climd))
		return "", err
	}

	arr := strings.Split(azPath, "/")
	azPath = filepath.Join(currentDir, arr[len(arr)-1])

	err = os.Chmod(azPath, 0777)
	if err != nil {

		return "", err
	}

	err = os.Chmod(filepath.Join(azPath, "sample"), 0777)
	if err != nil {

		return "", err
	}
	return azPath, nil
}

/*
	Copy the Azure-functions-go/azfunc lib from pkg/mod into the App Dir so that you can add and build the Go functions.
*/
func copyAzFunc(azPath, azFuncPath string) error {
	currPath, err := os.Getwd()

	if err != nil {
		return err
	}

	defer os.Chdir(currPath)

	err = os.Chdir(filepath.Join(azPath))
	if err != nil {
		return err
	}

	err = os.Mkdir("azfunc", 0777)
	if err != nil {
		return err
	}

	climd, err := exec.Command("cp", filepath.Join(azFuncPath, "bindings.go"), filepath.Join(azPath, "azfunc", "bindings.go")).Output()

	if err != nil {
		fmt.Println(string(climd))
		return err
	}

	climd, err = exec.Command("cp", filepath.Join(azFuncPath, "go.mod"), filepath.Join(azPath, "azfunc", "go.mod")).Output()

	if err != nil {
		fmt.Println(string(climd))
		return err
	}
	return nil
}
