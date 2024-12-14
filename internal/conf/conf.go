package config

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"stream-session-api/internal/conf/network"
	"stream-session-api/pkg"

	"github.com/charmbracelet/log"
	"gopkg.in/ini.v1"
)

func getDir() (string, error) {
	// Select os runtime
	runtime := runtime.GOOS
	switch runtime {
	case "windows":
		return os.Getenv("DATA_DIR_WIN"), nil
	case "linux":
		dir, err := os.UserHomeDir()
		if err != nil {
			return "", pkg.ErrorStatus(pkg.ErrCodeDirNotFound, "Fail to find home dir")
		}
		applicationDir := dir + "/"
		return applicationDir, nil

	default:
		return "", pkg.ErrorStatus(pkg.ErrCodeUnsupportedOs, "Unsupported OS")
	}

}

func setDefault() {
	log.Info("Set default conf...")

	conf := network.Get()

	// MediaMtx: http, rtsp, webrtc server
	conf.MediaMtx.Http.Ip = os.Getenv("DEFAULT_MEDIAMTX_HTTP_SERVER_URI")
	port, _ := strconv.ParseInt(os.Getenv("DEFAULT_MEDIAMTX_HTTP_SERVER_PORT"), 10, 16)
	conf.MediaMtx.Http.Port = uint16(port)

	conf.MediaMtx.Rtsp.Path = os.Getenv("DEFAULT_MEDIAMTX_RTSP_SERVER_PATH")
	conf.MediaMtx.Rtsp.Ip = os.Getenv("DEFAULT_MEDIAMTX_RTSP_SERVER_URI")
	port, _ = strconv.ParseInt(os.Getenv("DEFAULT_MEDIAMTX_RTSP_SERVER_PORT"), 10, 16)
	conf.MediaMtx.Rtsp.Port = uint16(port)

	conf.MediaMtx.WebRtc.Ip = os.Getenv("DEFAULT_MEDIAMTX_WEBRTC_SERVER_URI")
	port, _ = strconv.ParseInt(os.Getenv("DEFAULT_MEDIAMTX_WEBRTC_SERVER_PORT"), 10, 16)
	conf.MediaMtx.WebRtc.Port = uint16(port)

	// Grpc
	conf.Grpc.Ip = os.Getenv("DEFAULT_GRPC_SERVER_URI")
	port, _ = strconv.ParseInt(os.Getenv("DEFAULT_GRPC_SERVER_PORT"), 10, 16)
	conf.Grpc.Port = uint16(port)

	// Redis
	conf.Redis.Ip = os.Getenv("DEFAULT_REDIS_SERVER_URI")
	port, _ = strconv.ParseInt(os.Getenv("DEFAULT_REDIS_SERVER_PORT"), 10, 16)
	conf.Redis.Port = uint16(port)
	conf.Redis.Password = os.Getenv("DEFAULT_REDIS_SERVER_PASSWORD")
	idx, _ := strconv.ParseUint(os.Getenv("DEFAULT_REDIS_SERVER_DB_INDEX"), 10, 8)
	conf.Redis.DatabaseIndex = uint8(idx)

	// Set net conf
	network.Set(conf)
}

func show() error {
	log.Info("Show conf...")

	// Get application dir
	applicationDir, err := getDir()
	if err != nil {
		return err
	}

	// Check config.ini if it doesn't exist
	pathFile := fmt.Sprintf("%s%s/%s", applicationDir, os.Getenv("APPLICATION_NAME"), os.Getenv("FILENAME_CONFIG"))
	if _, err := os.Stat(pathFile); err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeFileNotFound, fmt.Sprintf("%s does not exist!", pathFile))
	}

	// Read config.ini
	config, err := ini.Load(pathFile)
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeReadFile, fmt.Sprintf("Fail to read %s ", pathFile))
	}

	// Iterate over all sections
	for _, section := range config.Sections() {

		if section.Name() == "DEFAULT" {
			continue
		}
		log.Info(fmt.Sprintf("[%s]", section.Name()))

		// Iterate over all keys in the current section
		for _, key := range section.Keys() {
			log.Info(fmt.Sprintf(" %s = %s", key.Name(), key.Value()))
		}
	}

	return nil
}

func write() error {
	log.Info("Write conf...")

	// Get application dir
	applicationDir, err := getDir()
	if err != nil {
		return err
	}

	// Check config.ini if it doesn't exist
	pathFile := fmt.Sprintf("%s%s/%s", applicationDir, os.Getenv("APPLICATION_NAME"), os.Getenv("FILENAME_CONFIG"))
	if _, err := os.Stat(pathFile); err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeFileNotFound, fmt.Sprintf("%s does not exist!", pathFile))
	}

	// Read .ini file
	settings, err := ini.Load(pathFile)
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeReadFile, fmt.Sprintf("Fail to read %s ", pathFile))
	}

	// Get config instance
	conf := network.Get()

	// Grpc server section
	sec, err := settings.NewSection("grpc")
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("ip", conf.Grpc.Ip)
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("port", strconv.FormatUint(uint64(conf.Grpc.Port), 10))
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}

	// Mediamtx http server section
	sec, err = settings.NewSection("mediamtx.http")
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("ip", conf.MediaMtx.Http.Ip)
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("port", strconv.FormatUint(uint64(conf.MediaMtx.Http.Port), 10))
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}

	// Mediamtx rtsp server section
	sec, err = settings.NewSection("mediamtx.rtsp")
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("ip", conf.MediaMtx.Rtsp.Ip)
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("port", strconv.FormatUint(uint64(conf.MediaMtx.Rtsp.Port), 10))
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("path", conf.MediaMtx.Rtsp.Path)
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}

	// Mediammtx webrtc server section
	sec, err = settings.NewSection("mediamtx.webrtc")
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("ip", conf.MediaMtx.WebRtc.Ip)
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("port", strconv.FormatUint(uint64(conf.MediaMtx.WebRtc.Port), 10))
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}

	// Redis
	sec, err = settings.NewSection("redis")
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("ip", conf.Redis.Ip)
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("port", strconv.FormatUint(uint64(conf.Redis.Port), 10))
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("password", conf.Redis.Password)
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}
	_, err = sec.NewKey("database_index", strconv.FormatUint(uint64(conf.Redis.DatabaseIndex), 10))
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeWriteFile, fmt.Sprintf("Fail to write to %s ", pathFile))
	}

	// Save to file
	err = settings.SaveTo(pathFile)
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeSaveFile, fmt.Sprintf("Fail to save %s ", pathFile))
	}

	return nil
}

func read() error {
	log.Info("Read conf...")

	// Get application dir
	applicationDir, err := getDir()
	if err != nil {
		return err
	}

	// Check .ini file if it doesn't exist
	pathFile := fmt.Sprintf("%s%s/%s", applicationDir, os.Getenv("APPLICATION_NAME"), os.Getenv("FILENAME_CONFIG"))
	if _, err := os.Stat(pathFile); err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeFileNotFound, fmt.Sprintf("%s does not exist!", pathFile))
	}

	// Read config.ini
	settings, err := ini.Load(pathFile)
	if err != nil {
		return pkg.ErrorStatus(pkg.ErrCodeReadFile, fmt.Sprintf("Fail to read %s ", pathFile))
	}

	// Get config instance
	conf := network.Get()

	// Grpc server section
	section := settings.Section("grpc")
	conf.Grpc.Ip = section.Key("ip").String()
	port, _ := section.Key("port").Uint64()
	conf.Grpc.Port = uint16(port)

	// Mediamtx http server section
	section = settings.Section("mediamtx.http")
	conf.MediaMtx.Http.Ip = section.Key("ip").String()
	port, _ = section.Key("port").Uint64()
	conf.MediaMtx.Http.Port = uint16(port)
	// Mediamtx rtsp server section
	section = settings.Section("mediamtx.rtsp")
	conf.MediaMtx.Rtsp.Ip = section.Key("ip").String()
	port, _ = section.Key("port").Uint64()
	conf.MediaMtx.Rtsp.Port = uint16(port)
	conf.MediaMtx.Rtsp.Path = section.Key("path").String()
	// Mediamtx webrtc server section
	section = settings.Section("mediamtx.webrtc")
	conf.MediaMtx.WebRtc.Ip = section.Key("ip").String()
	port, _ = section.Key("port").Uint64()
	conf.MediaMtx.WebRtc.Port = uint16(port)

	// Redis
	section = settings.Section("redis")
	conf.Redis.Ip = section.Key("ip").String()
	port, _ = section.Key("port").Uint64()
	conf.Redis.Port = uint16(port)
	conf.Redis.Password = section.Key("password").String()
	idx, _ := section.Key("database_index").Uint64()
	conf.Redis.DatabaseIndex = uint8(idx)

	// Set net conf
	network.Set(conf)

	return nil
}

func Get() error {
	log.Info("Get conf...")

	// Get application dir
	applicationDir, err := getDir()
	if err != nil {
		return err
	}

	// Create config.ini if it doesn't exist
	pathFile := fmt.Sprintf("%s%s/%s", applicationDir, os.Getenv("APPLICATION_NAME"), os.Getenv("FILENAME_CONFIG"))
	if _, err := os.Stat(pathFile); err != nil {
		// File does not exist, so create it.
		file, err := os.Create(pathFile)
		if err != nil {
			return pkg.ErrorStatus(pkg.ErrCodeCreateFile, fmt.Sprintf("Fail to create %s", pathFile))
		}
		defer file.Close()

		// Set default value param
		setDefault()

		// Write config to file
		if err := write(); err != nil {
			return err
		}
	}

	// Read config file
	if err := read(); err != nil {
		return err
	}

	// Show config to console
	if err := show(); err != nil {
		return err
	}

	return nil
}
