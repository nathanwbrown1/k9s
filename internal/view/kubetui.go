// SPDX-License-Identifier: Apache-2.0
// Copyright kubetui contributors

package view

import (
	"os"
	"path/filepath"

	"github.com/derailed/k9s/internal/kubetui"
	"github.com/derailed/k9s/internal/ui"
	"github.com/derailed/tcell/v2"
)

// Built-in companion actions do not edit or depend on the user's plugins.yaml.
func kubetuiActions(r Runner, aa *ui.KeyActions) {
	binary := os.Getenv("KUBETUI_EXECUTABLE")
	if binary == "" || !filepath.IsAbs(binary) || r.EnvFn() == nil {
		return
	}
	env := r.EnvFn()()
	if !kubetui.Supported(env["RESOURCE_NAME"], env["RESOURCE_GROUP"]) {
		return
	}
	for _, action := range []struct {
		key         tcell.Key
		label, mode string
	}{
		{tcell.KeyF6, "Kubetui workflows", "workflows"},
		{tcell.KeyF7, "Kubetui related logs", "logs"},
	} {
		mode := action.mode
		// These two keys are reserved by the kubetui distribution.
		aa.Add(action.key, ui.NewKeyActionWithOpts(action.label, func(event *tcell.EventKey) *tcell.EventKey {
			if r.GetSelectedItem() == "" {
				return event
			}
			args, err := kubetui.Args(r.EnvFn()(), mode)
			if err != nil {
				r.App().Flash().Err(err)
				return nil
			}
			flags := r.App().Conn().Config().Flags()
			// The launcher exposes only context/namespace/kubeconfig. Refuse
			// direct launches with connection overrides rather than changing identity.
			for _, value := range []*string{flags.Impersonate, flags.AuthInfoName, flags.ClusterName,
				flags.APIServer, flags.BearerToken, flags.CAFile, flags.CertFile, flags.KeyFile} {
				if value != nil && *value != "" {
					r.App().Flash().Warn("Use a named kubeconfig context for kubetui investigations; connection overrides are unsupported")
					return nil
				}
			}
			if (flags.Insecure != nil && *flags.Insecure) || (flags.ImpersonateGroup != nil && len(*flags.ImpersonateGroup) > 0) {
				r.App().Flash().Warn("Use a named kubeconfig context for kubetui investigations; connection overrides are unsupported")
				return nil
			}
			if flags.KubeConfig != nil && *flags.KubeConfig != "" {
				args = append(args, "--kubeconfig", *flags.KubeConfig)
			}
			ok, errors, _ := run(r.App(), &shellOpts{binary: binary, args: args})
			if !ok {
				r.App().Flash().Warn("Could not suspend explorer for kubetui")
				return nil
			}
			for err := range errors {
				r.App().Flash().Err(err)
			}
			return nil
		}, ui.ActionOpts{Visible: true, Plugin: true}))
	}
}
