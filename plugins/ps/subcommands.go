package ps

import (
	"errors"
	"strings"

	"github.com/dokku/dokku/plugins/common"
)

// TODO: implement me
func CommandInspect(appName string) error {
	if appName == "" {
		return errors.New("Please specify an app to run the command on")
	}

	if err := common.VerifyAppName(appName); err != nil {
		return err
	}

	return nil
}

func CommandRebuild(appName string, allApps bool, runInSerial bool, parallelCount int) error {
	if allApps {
		return RunCommandAgainstAllApps(Rebuild, "rebuild", runInSerial, parallelCount)
	}

	if appName == "" {
		return errors.New("Please specify an app to run the command on")
	}

	if err := common.VerifyAppName(appName); err != nil {
		return err
	}

	return Rebuild(appName)
}

// CommandReport displays a ps report for one or more apps
func CommandReport(appName string, infoFlag string) error {
	if strings.HasPrefix(appName, "--") {
		infoFlag = appName
		appName = ""
	}

	if len(appName) == 0 {
		apps, err := common.DokkuApps()
		if err != nil {
			return err
		}
		for _, appName := range apps {
			if err := ReportSingleApp(appName, infoFlag); err != nil {
				return err
			}
		}
		return nil
	}

	return ReportSingleApp(appName, infoFlag)
}

func CommandRestart(appName string, allApps bool, runInSerial bool, parallelCount int) error {
	if allApps {
		return RunCommandAgainstAllApps(Restart, "restart", runInSerial, parallelCount)
	}

	if appName == "" {
		return errors.New("Please specify an app to run the command on")
	}

	if err := common.VerifyAppName(appName); err != nil {
		return err
	}

	return Restart(appName)
}

// TODO: implement me
func CommandRestore(appName string) error {
	return nil
}

// TODO: implement me
func CommandRetire() error {
	return nil
}

// TODO: implement me
func CommandScale(appName string, processTuples []string) error {
	return nil
}

// TODO: implement me
func CommandSet(appName string, property string, value string) error {
	return nil
}

func CommandStart(appName string, allApps bool, runInSerial bool, parallelCount int) error {
	if allApps {
		return RunCommandAgainstAllApps(Start, "start", runInSerial, parallelCount)
	}

	if appName == "" {
		return errors.New("Please specify an app to run the command on")
	}

	if err := common.VerifyAppName(appName); err != nil {
		return err
	}

	return Start(appName)
}

func CommandStop(appName string, allApps bool, runInSerial bool, parallelCount int) error {
	if allApps {
		return RunCommandAgainstAllApps(Stop, "stop", runInSerial, parallelCount)
	}

	if appName == "" {
		return errors.New("Please specify an app to run the command on")
	}

	if err := common.VerifyAppName(appName); err != nil {
		return err
	}

	return Stop(appName)
}
