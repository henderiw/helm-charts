package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/kustomize/kyaml/kio"
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

	namespace := ""
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
	client.Devel = true
	client.DryRun = true
	client.ReleaseName = "release-name"
	client.Replace = true // Skip the name check
	client.ClientOnly = true
	client.APIVersions = chartutil.VersionSet([]string{})
	client.IncludeCRDs = true

	rel, err := client.Run(chrt, vals)
	if err != nil {
		fmt.Println(err)
	}

	if rel != nil {
		r := &kio.ByteReader{Reader: bytes.NewBufferString(string(rel.Manifest)), OmitReaderAnnotations: true}
		nodes, err := r.Read()
		if err != nil {
			panic(err)
		}

		for i := range nodes {
			o, err := fn.ParseKubeObject([]byte(nodes[i].MustString()))
			if err != nil {
				if strings.Contains(err.Error(), "expected exactly one object, got 0") {
					// sometimes helm produces some messages in between resources, we can safely
					// ignore these
					continue
				}
				err = fmt.Errorf("failed to parse %s: %s", nodes[i].MustString(), err.Error())
				panic(err)
			}
			fmt.Println("---- kubeObject start ----")
			fmt.Println(o.String())
			fmt.Println("---- kubeObject end   ----")
		}
	}
}
