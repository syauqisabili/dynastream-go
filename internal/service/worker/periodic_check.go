package worker

import (
	"fmt"
	"os"
	"strconv"
	"stream-session-api/dto"
	"stream-session-api/internal/conf/network"
	"stream-session-api/internal/repository"
	"stream-session-api/pkg"
	"time"

	"github.com/go-resty/resty/v2"
)

func inactiveSessionHandler() error {
	// Get config instance
	config := network.Get()

	// Get WebRTC Session
	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(&dto.StreamSessionList{}).
		Get(fmt.Sprintf("http://%s:%d/v3/webrtcsessions/list",
			config.MediaMtx.Http.Ip,
			config.MediaMtx.Http.Port))

	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return err
	}
	sessions := resp.Result().(*dto.StreamSessionList)

	// Get all stream
	repo := repository.NewStream()
	defer repo.Close()

	streams, err := repo.GetAll()
	if err != nil || streams == nil {
		return pkg.ErrorStatus(pkg.ErrCodeProcessFail, "Stream list not found")
	}

	// Cleanup inactive session
	for _, stream := range streams {
		match := false
		for _, session := range sessions.Items {
			if stream.Uuid == session.Path {
				pkg.LogInfo(fmt.Sprintf("%v Active", *stream))
				match = true
				break
			}
		}

		if !match {
			pkg.LogInfo(fmt.Sprintf("%v Inactive", *stream))
			// Stop stream path
			client := resty.New()
			resp, err := client.R().
				Delete(fmt.Sprintf("http://%s:%d/v3/config/paths/delete/%s",
					config.MediaMtx.Http.Ip,
					config.MediaMtx.Http.Port,
					stream.Uuid))

			if err != nil {
				return pkg.ErrorStatus(pkg.ErrCodeProcessFail, "Fail to stop stream")
			}
			if resp.StatusCode() != 200 {
				return pkg.ErrorStatus(pkg.ErrCodeProcessFail, fmt.Sprintf("Fail to delete path stream: %d", resp.StatusCode()))
			}

			// Delete stream redis log
			if err := repo.Delete(stream.Uuid); err != nil {
				return pkg.ErrorStatus(pkg.ErrCodeProcessFail, "Fail to close stream")
			}
		}
	}

	return nil
}

func PeriodicStreamSessionCheck() {
	go func() {
		val, _ := strconv.ParseInt(os.Getenv("PERIODIC_STREAM_SESSION_CHECK"), 10, 16)
		ticker := time.NewTicker(time.Second * time.Duration(val))
		defer ticker.Stop()

		// Schedule on
		for range ticker.C {
			currentTime := time.Now()
			pkg.LogInfo(fmt.Sprintf("STREAM_SESSION_CHECK: %d/%02d/%02d %d:%d:%d",
				currentTime.Year(), int(currentTime.Month()), currentTime.Day(),
				currentTime.Hour(), currentTime.Minute(), currentTime.Second()))

			// Check inactive Stream Session
			err := inactiveSessionHandler()
			if err != nil {
				pkg.LogWarn(fmt.Sprintf("Failed to check stream session: %s", err))
			}
		}
	}()

	// Prevent the main routine from exiting
	select {}

}
