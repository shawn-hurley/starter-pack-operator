package main

import (
	"context"
	"flag"
	"runtime"

	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	stub "github.com/shawn-hurley/starter-pack-operator/pkg/stub"

	"github.com/sirupsen/logrus"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
	logrus.Infof("watching in namespace: %v", namespace)
}

var namespace string

func init() {
	flag.StringVar(&namespace, "namespace", "default", "specify the namespace for the operator to watch in")
}

func main() {
	flag.Parse()
	printVersion()
	sdk.Watch("starterpack.osbkit.com/v1alpha1", "Broker", namespace, 5)
	sdk.Handle(stub.NewHandler())
	sdk.Run(context.TODO())
}
