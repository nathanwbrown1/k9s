# kubetui integration fork

Based on upstream K9s v0.51.0 (`558caafe7ba067467de46b320cc22ef11fef9c34`).
This is a modified distribution, not an official upstream release.

Changes:

- `internal/view/actions.go` installs companion actions after user plugins load.
- `internal/view/kubetui.go` provides F6 guided workflows and F7 related logs,
  using K9s suspend/resume and literal argv. It requires the launcher's absolute
  `KUBETUI_EXECUTABLE` environment variable and rejects connection overrides.
- `internal/kubetui/` validates supported resource selections and constructs the
  explicit context/namespace/resource/name handoff with unit tests.

Build this module as `kubetui-k9s` and install beside `kubetui`. Keep upstream
module paths, branding, LICENSE and COPYING. F6/F7 are reserved for this
distribution; do not assign them in custom hotkeys. No plugin YAML is installed.

The `kubetui-integration` branch follows reviewed upstream release tags. Keep
the integration patch small and test handoff arguments, terminal restoration,
read-only mode and builds on Linux and Windows when updating upstream.
