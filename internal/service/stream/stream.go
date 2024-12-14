package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"stream-session-api/domain"
	"stream-session-api/internal/conf/network"
	"stream-session-api/internal/repository"
	pb "stream-session-api/internal/service/stream/proto"
	"stream-session-api/pkg"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.StreamServiceServer
}

func (*Server) StartStream(ctx context.Context, in *pb.StartStreamRequest) (*pb.StartStreamResponse, error) {
	// Check value pb.StartStreamRequest
	if in == nil {
		pkg.LogError("Invalid message request")
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Invalid message request")
	}

	// Get the peer information from the context
	client, _ := peer.FromContext(ctx)
	pkg.LogInfo(fmt.Sprintf("%s requested to start stream for %s from %s", in.Username, in.GetStreamId(), client.Addr))

	// Stream request for specific id
	stream := &domain.Stream{
		Id:   in.GetStreamId(),
		Uuid: uuid.New().String(),
	}

	//* you must make a mapping between stream.id and rtsp.subpath ex => id: 001 to subpath: /74630a72-2612-478b-9c49-308c720aa619
	//? Why ? Because, you should not stream a video from a server directly if stream.id==rtsp.subPath
	//! Temporary, i make it to be the same
	subPath := "/" + stream.Id

	// Set body for http post
	body, _ := json.Marshal(
		struct {
			Name   string `json:"name"`
			Source string `json:"source"`
		}{
			Name: stream.Uuid,
			Source: fmt.Sprintf("rtsp://%s:%d%s%s",
				network.Get().MediaMtx.Rtsp.Ip,
				network.Get().MediaMtx.Rtsp.Port,
				network.Get().MediaMtx.Rtsp.Path,
				subPath,
			),
		},
	)

	// Add stream session via http request
	httpClient := resty.New()
	resp, err := httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(string(body)).
		Post(fmt.Sprintf("http://%s:%d/v3/config/paths/add/%s",
			network.Get().MediaMtx.Http.Ip,
			network.Get().MediaMtx.Http.Port,
			stream.Uuid))
	if err != nil {
		pkg.LogError(err)
		return nil, status.Errorf(
			codes.Unavailable, "Failed to configure stream")
	}

	switch resp.StatusCode() {
	case 400:
		pkg.LogError(err)
		return nil, status.Errorf(
			codes.FailedPrecondition, "Failed to add stream session: %d", resp.StatusCode())
	case 500:
		pkg.LogError(err)
		return nil, status.Errorf(
			codes.Unimplemented, "Failed to add stream session: %d", resp.StatusCode())

	}

	// Insert stream url to redis
	repo := repository.NewStream()
	if err := repo.Insert(stream); err != nil {
		pkg.LogError(err)
		return nil, status.Errorf(
			codes.Unknown, "Cannot do streaming")
	}

	// Set stream url
	url := fmt.Sprintf("http://%s:%d/%s",
		network.Get().MediaMtx.WebRtc.Ip,
		network.Get().MediaMtx.WebRtc.Port,
		stream.Uuid)
	pkg.LogInfo(fmt.Sprintf("Streaming on %s", url))

	return &pb.StartStreamResponse{
		StreamUrl: url,
	}, nil
}

func (*Server) StopStream(ctx context.Context, in *pb.StopStreamRequest) (*emptypb.Empty, error) {
	// Check value pb.StartStreamRequest
	if in == nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Invalid message request")
	}

	// Get the peer information from the context
	client, _ := peer.FromContext(ctx)
	pkg.LogInfo(fmt.Sprintf("%s requested to stop stream on %s from %s", in.Username, in.GetStreamUrl(), client.Addr))

	// Get uuid
	parsedUrl, err := url.ParseRequestURI(in.GetStreamUrl())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Invalid stream url request")
	}
	uuid := strings.TrimSuffix(parsedUrl.Path, "/")
	uuid = strings.TrimPrefix(uuid, "/")

	// Stop stream
	httpClient := resty.New()
	resp, err := httpClient.R().
		SetContext(ctx).
		Delete(fmt.Sprintf("http://%s:%d/v3/config/paths/delete/%s",
			network.Get().MediaMtx.Http.Ip,
			network.Get().MediaMtx.Http.Port,
			uuid))

	if err != nil {
		pkg.LogError(err)
		return nil, status.Errorf(
			codes.Unavailable, "Failed to stop stream")
	}
	switch resp.StatusCode() {
	case 400:
		pkg.LogError(err)
		return nil, status.Errorf(
			codes.FailedPrecondition, "Failed to remove stream session: %d", resp.StatusCode())
	case 500:
		pkg.LogError(err)
		return nil, status.Errorf(
			codes.Unimplemented, "Failed to remove stream session: %d", resp.StatusCode())

	}
	// Delete uuid on redis
	repo := repository.NewStream()
	if err := repo.Delete(uuid); err != nil {
		return nil, status.Errorf(
			codes.Unknown, "Fail to close stream")
	}

	//! Temporary
	return nil, nil
}
