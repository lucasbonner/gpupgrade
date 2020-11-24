//  Copyright (c) 2017-2020 VMware, Inc. or its affiliates
//  SPDX-License-Identifier: Apache-2.0

package upgrade

import (
	"fmt"

	"github.com/greenplum-db/gp-common-go-libs/gplog"
	"golang.org/x/xerrors"

	"github.com/greenplum-db/gpupgrade/utils"
)

func GpupgradeVersion() (string, error) {
	return getGpupgradeVersion("")
}

func GpupgradeVersionOnHost(host string) (string, error) {
	return getGpupgradeVersion(host)
}

type GpupgradeVersions struct{}

func (g *GpupgradeVersions) HubVersion() (string, error) {
	return GpupgradeVersion()
}

func (g *GpupgradeVersions) AgentVersion(host string) (string, error) {
	return GpupgradeVersionOnHost(host)
}

func getGpupgradeVersion(host string) (string, error) {
	gpupgradePath, err := utils.GetGpupgradePath()
	if err != nil {
		return "", xerrors.Errorf("getting gpupgrade binary path: %w", err)
	}

	name := gpupgradePath
	args := []string{"version", "--format", "oneline"}
	if host != "" {
		name = "ssh"
		args = []string{host, fmt.Sprintf(`bash -c "%s version --format oneline"`, gpupgradePath)}
	}

	cmd := execCommand(name, args...)
	gplog.Debug("running cmd %q", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", xerrors.Errorf("%q failed with %q: %w", cmd.String(), string(output), err)
	}

	gplog.Debug("output: %q", output)

	return string(output), nil
}
