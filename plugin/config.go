package plugin

import (
	"github.com/doubletrey/crawlab-core/config"
	"path"
)

const DefaultPluginFsPathBase = "plugins"
const DefaultPluginDirName = "plugins"
const DefaultPluginBinName = "plugin"
const DefaultPluginInstallCmd = "go build -o ./build/start"

var DefaultPluginDirPath = path.Join(config.DefaultConfigDirPath, DefaultPluginDirName)
