// SPDX-License-Identifier: Apache-2.0
// Copyright kubetui contributors
package kubetui

import (
	"reflect"
	"testing"
)

func TestSelectedScope(t *testing.T) {
	env := map[string]string{"CONTEXT": "lab $literal; name", "RESOURCE_NAME": "deployments", "RESOURCE_GROUP": "apps", "NAMESPACE": "team-a", "NAME": "api"}
	args, err := Args(env, "logs")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"--classic", "--handoff", "logs", "--context", "lab $literal; name", "--namespace", "team-a", "--resource", "deployments.apps", "--name", "api"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("%q", args)
	}
	// Cluster-scoped rows have no namespace, regardless of the current table scope.
	env["RESOURCE_NAME"], env["RESOURCE_GROUP"] = "nodes", ""
	args, err = Args(env, "workflows")
	if err != nil || args[6] != "" {
		t.Fatalf("%q %v", args, err)
	}
}

func TestRejectAmbiguousSelection(t *testing.T) {
	for _, change := range []map[string]string{
		{"NAMESPACE": ""}, {"NAMESPACE": "all"}, {"CONTEXT": ""},
		{"NAME": ""}, {"NAME": "--help"}, {"NAME": "pod/container"},
		{"RESOURCE_GROUP": "unrelated.example"}, {"RESOURCE_NAME": "containers"},
	} {
		env := map[string]string{"CONTEXT": "lab", "RESOURCE_NAME": "pods", "NAMESPACE": "team-a", "NAME": "api"}
		for k, v := range change {
			env[k] = v
		}
		if args, err := Args(env, "logs"); err == nil {
			t.Fatalf("accepted %v: %v", change, args)
		}
	}
}
