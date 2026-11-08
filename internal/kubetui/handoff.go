// SPDX-License-Identifier: Apache-2.0
// Copyright kubetui contributors

// Package kubetui defines the selected-resource boundary with the companion UI.
package kubetui

import (
	"fmt"
	"strings"
)

// Resources with relationship/workflow support in kubetui. Other resources remain
// fully browsable in K9s; a handoff must not silently widen to the whole cluster.
var resources = map[string]bool{
	"pods": true, "services": true, "secrets": true, "configmaps": true,
	"persistentvolumeclaims": true, "nodes": false,
	"deployments.apps": true, "replicasets.apps": true,
	"statefulsets.apps": true, "daemonsets.apps": true,
	"jobs.batch": true, "cronjobs.batch": true,
	"ingresses.networking.k8s.io":  true,
	"certificates.cert-manager.io": true, "certificaterequests.cert-manager.io": true,
	"issuers.cert-manager.io": true, "clusterissuers.cert-manager.io": false,
	"orders.acme.cert-manager.io": true, "challenges.acme.cert-manager.io": true,
}

func Supported(resource, group string) bool {
	_, ok := resources[qualified(resource, group)]
	return ok
}

func qualified(resource, group string) string {
	if group == "" {
		return resource
	}
	return resource + "." + group
}

// Args uses literal argv elements, never shell interpolation or plugin expansion.
func Args(env map[string]string, mode string) ([]string, error) {
	kind := qualified(env["RESOURCE_NAME"], env["RESOURCE_GROUP"])
	namespaced, ok := resources[kind]
	if !ok {
		return nil, fmt.Errorf("kubetui investigations do not yet support %s", kind)
	}
	if mode != "workflows" && mode != "logs" {
		return nil, fmt.Errorf("unknown investigation %q", mode)
	}
	context, name, ns := env["CONTEXT"], env["NAME"], env["NAMESPACE"]
	if context == "" || name == "" || strings.HasPrefix(name, "-") || strings.ContainsAny(name, "/\\\r\n") {
		return nil, fmt.Errorf("select a named Kubernetes resource first")
	}
	if namespaced && (ns == "" || ns == "-" || ns == "all" || ns == "*") {
		return nil, fmt.Errorf("selected resource has no concrete namespace")
	}
	if !namespaced {
		ns = ""
	}
	return []string{"--classic", "--handoff", mode, "--context", context,
		"--namespace", ns, "--resource", kind, "--name", name}, nil
}
