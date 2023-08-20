package main

import (
	"fmt"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	chartNginx = "data/nginx"
	chartUpf   = "data/free5gc-upf"
)

func main() {
	chrt, err := loader.Load(chartUpf)
	if err != nil {
		panic(err)
	}

	vals, err := chartutil.CoalesceValues(chrt, map[string]any{})
	if err != nil {
		panic(err)
	}

	/*
		v, err := yaml.Marshal(vals)
		if err != nil {
			panic(err)
		}
		//fmt.Println(string(v))

		c, err := config.GetConfig()
		if err != nil {
			panic(err)
		}
		//fmt.Println(c)
	*/

	/*
		httpClient, err := rest.HTTPClientFor(c)
		if err != nil {
			panic(err)
		}
	*/

	/*
		restmapper, err := apiutil.NewDynamicRESTMapper(c, httpClient)
		if err != nil {
			panic(err)
		}
	*/

	namespace := "test"

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(
		&genericclioptions.ConfigFlags{
			Namespace: &namespace,
		},
		namespace,
		os.Getenv("HELM_DRIVER"),
		log.Printf,
	); err != nil {
		panic(err)
	}

	client := action.NewInstall(actionConfig)
	client.ReleaseName = chrt.Metadata.Version
	client.Namespace = namespace
	client.DryRun = true

	r, err := client.Run(chrt, vals)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(r.Manifest)
}
