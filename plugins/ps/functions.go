package ps

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dokku/dokku/plugins/common"
	dockeroptions "github.com/dokku/dokku/plugins/docker-options"
)

func extractProcfile(appName, image string, procfilePath string) error {
	if err := removeProcfile(appName); err != nil {
		return err
	}

	destination := filepath.Join(common.MustGetEnv("DOKKU_ROOT"), appName, "DOKKU_PROCFILE")
	common.CopyFromImage(appName, image, "Procfile", destination)
	if !common.FileExists(destination) {
		common.LogInfo1Quiet("No Procfile found in app image")
		return nil
	}

	common.LogInfo1Quiet("App Procfile file found")

	checkCmd := common.NewShellCmd(strings.Join([]string{
		"procfile-util",
		"check",
		"--procfile",
		destination,
	}, " "))
	var stderr bytes.Buffer
	checkCmd.ShowOutput = false
	checkCmd.Command.Stderr = &stderr
	_, err := checkCmd.Output()

	if err != nil {
		return fmt.Errorf(strings.TrimSpace(stderr.String()))
	}

	return nil
}

func extractOrGenerateScaleFile(appName string, image string) error {
	destination := filepath.Join(common.MustGetEnv("DOKKU_ROOT"), appName, "DOKKU_SCALE")
	extracted := fmt.Sprintf("%s.extracted", destination)
	if err := common.CopyFromImage(appName, image, "DOKKU_SCALE", destination); err != nil {
		os.Remove(extracted)
		os.Remove(destination)
	} else if err := common.CopyFile(destination, extracted); err != nil {
		return err
	}

	if common.FileExists(destination) {
		return nil
	}

	common.LogInfo1Quiet("DOKKU_SCALE file not found in app image. Generating one based on Procfile...")
	return generateScaleFile(appName, destination)
}

func generateScaleFile(appName string, destination string) error {
	procfilePath := filepath.Join(common.MustGetEnv("DOKKU_ROOT"), appName, "DOKKU_PROCFILE")
	content := []string{"web=1"}
	if !common.FileExists(procfilePath) {
		return common.WriteSliceToFile(destination, content)
	}

	lines, err := common.FileToSlice(procfilePath)
	if err != nil {
		return common.WriteSliceToFile(destination, content)
	}

	content = []string{""}
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		count := 0
		if parts[0] == "web" {
			count = 1
		}

		content = append(content, fmt.Sprintf("%s=%d", parts[0], count))
	}

	return common.WriteSliceToFile(destination, content)
}

func getProcfileCommand(procfilePath string, processType string, port int) (string, error) {

	shellCmd := common.NewShellCmd(strings.Join([]string{
		"procfile-util",
		"show",
		"--procfile",
		procfilePath,
		"--process-type",
		processType,
		"--default-port",
		strconv.Itoa(port),
	}, " "))
	var stderr bytes.Buffer
	shellCmd.ShowOutput = false
	shellCmd.Command.Stderr = &stderr
	b, err := shellCmd.Output()

	if err != nil {
		return "", fmt.Errorf(strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(string(b[:])), nil
}

func getRestartPolicies(appName string) ([]string, error) {
	policies := []string{}

	options, err := dockeroptions.GetDockerOptionsForPhase(appName, "deploy")
	if err != nil {
		return policies, err
	}
	for _, option := range options {
		if strings.HasPrefix(option, "--restart=") {
			policies = append(policies, option)
		}
	}

	return policies, nil
}

func getRunningState(appName string) string {
	scheduler := common.GetAppScheduler(appName)
	b, _ := common.PlugnTriggerOutput("scheduler-app-status", []string{scheduler, appName}...)
	return strings.Split(strings.TrimSpace(string(b[:])), " ")[1]
}

// removeProcfile removes the DOKKU_PROCFILE file from the repo root
func removeProcfile(appName string) error {
	procfile := filepath.Join(common.MustGetEnv("DOKKU_ROOT"), appName, "DOKKU_PROCFILE")
	if !common.FileExists(procfile) {
		return nil
	}

	return os.Remove(procfile)
}
