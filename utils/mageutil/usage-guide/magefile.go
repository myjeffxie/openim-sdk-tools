//go:build mage
// +build mage

package main

import (
	"os"
	"strings"

	"github.com/myjeffxie/openim-sdk-tools/utils/mageutil"
)

var Default = Build

func Build() {
	platforms := os.Getenv("PLATFORMS")
	if platforms == "" {
		platforms = mageutil.DetectPlatform()
	}

	for _, platform := range strings.Split(platforms, " ") {
		mageutil.CompileForPlatform(platform)
	}

	mageutil.PrintGreen("All binaries under cmd and tools were successfully compiled.")
}

func Start() {
	setMaxOpenFiles()
	mageutil.StartToolsAndServices()
}

func Stop() {
	mageutil.StopAndCheckBinaries()
}

func Check() {
	mageutil.CheckAndReportBinariesStatus()
}
