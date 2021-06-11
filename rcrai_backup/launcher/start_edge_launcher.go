package launcher

import (
	"infra/artifacts/appversions"
	"infra/servers"
	"log"
)

func StartLauncher() {
	log.Printf("[Launcher] * 启动 Edge Device Launcher")

	// * Step 1: Load configuration

	// Load launcher configuration. with location of launcher

	cfg := LoadGlobalConfig()
	GlobalLauncherConfig = cfg // todo refactor

	// * Step 2: Start Http Server

	server := servers.NewUniversalServer()

	launcher := NewFileLauncher()

	server.RegisterServer(appversions.NewAppVersionController(
		cfg.Launcher.VersionsPath, 2000, launcher,
	))

	server.Start()

	if err := server.GracefulShutdown(); err != nil {
		panic(err)
	}

}
