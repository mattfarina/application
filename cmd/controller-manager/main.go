/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"log"

	controllerlib "github.com/najena/kubebuilder/pkg/controller"
	"github.com/najena/kubebuilder/pkg/install"
	extensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	"github.com/kubernetes-sigs/apps_application/pkg/controller"
	"github.com/kubernetes-sigs/apps_application/pkg/apis"
)

var kubeconfig = flag.String("kubeconfig", "", "path to kubeconfig")
var installCRDs = flag.Bool("install-crds", true, "install the CRDs used by the controller as part of startup")

// Controller-manager main.
func main() {
	flag.Parse()
	config, err := controllerlib.GetConfig(*kubeconfig)
	if err != nil {
		log.Fatalf("Could not create Config for talking to the apiserver: %v", err)
	}

    if *installCRDs {
        err = install.NewInstaller(config).Install(&InstallStrategy{crds: apis.APIMeta.GetCRDs()})
        if err != nil {
            log.Fatalf("Could not create CRDs: %v", err)
        }
    }

    // Start the controllers
	controllers, _ := controller.GetAllControllers(config)
	controllerlib.StartControllerManager(controllers...)

	// Blockforever
	select {}
}

type InstallStrategy struct {
	install.EmptyInstallStrategy
	crds []extensionsv1beta1.CustomResourceDefinition
}

func (s *InstallStrategy) GetCRDs() []extensionsv1beta1.CustomResourceDefinition {
	return s.crds
}
