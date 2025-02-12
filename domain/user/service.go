// File:		service.go
// Created by:	Hoven
// Created on:	2025-02-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package user

import (
	"context"
	"io"
	"mime/multipart"

	"gitea.hoven.com/core/auth-core/pkg/dto"
	"gitea.hoven.com/core/auth-core/proto/authenticationpb"
	"gitea.hoven.com/core/auth-core/proto/userpb"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/yazl-tech/beauty-rating-server/pkg/exception"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Service interface {
	WxLogin(ctx context.Context, deviceId, code string) (*Token, error)
	GetUserInfo(ctx context.Context) (*User, error)
	UpdateUsername(ctx context.Context, username string) error
	UpdateGender(ctx context.Context, gender int) error
	UploadAvatar(ctx context.Context, fh *multipart.FileHeader) (string, error)
	GetAvatar(ctx context.Context, avatarId string, writer io.Writer) error
}

var _ Service = (*DefaultUserService)(nil)

type DefaultUserService struct {
	wxConfig   *WechatConfig
	authClient authenticationpb.AuthCoreAuthenticationHandlerClient
	userClient userpb.AuthCoreUserHandlerClient
}

func NewUserService(wxConf *WechatConfig, authCoreConn grpc.ClientConnInterface) *DefaultUserService {
	userClient := userpb.NewAuthCoreUserHandlerClient(authCoreConn)
	authClient := authenticationpb.NewAuthCoreAuthenticationHandlerClient(authCoreConn)

	return &DefaultUserService{
		wxConfig:   wxConf,
		authClient: authClient,
		userClient: userClient,
	}
}

func (us *DefaultUserService) WxLogin(ctx context.Context, deviceId, code string) (*Token, error) {
	resp, err := us.authClient.WechatLogin(ctx, &dto.WechatLoginRequest{
		AppId:    us.wxConfig.AppId,
		SecretId: us.wxConfig.SecretId,
		DeviceId: deviceId,
		Code:     code,
	})
	if err != nil {
		return nil, exception.ParseGrpcError(err)
	}

	return &Token{
		UserID:       int(resp.GetUserId()),
		AccessToken:  resp.GetToken().AccessToken,
		RefreshToken: resp.GetToken().RefreshToken,
	}, nil
}

func (us *DefaultUserService) GetUserInfo(ctx context.Context) (*User, error) {
	resp, err := us.userClient.GetUserProfile(ctx, &dto.GetUserProfileRequest{})
	if err != nil {
		return nil, exception.ParseGrpcError(err)
	}

	return &User{
		ID:     int(resp.GetUserId()),
		Name:   resp.GetName(),
		Avatar: resp.GetAvatar(),
		Gender: Gender(resp.GetGender()),
		Email:  resp.GetEmail(),
		Status: Status(resp.GetStatus()),
		Role:   Role(resp.GetRole()),
	}, nil
}

func (us *DefaultUserService) UpdateUsername(ctx context.Context, username string) error {
	_, err := us.userClient.UpdateUserName(ctx, &dto.UpdateUserNameRequest{Username: username})
	if err != nil {
		return exception.ParseGrpcError(err)
	}

	return nil
}

func (us *DefaultUserService) UpdateGender(ctx context.Context, gender int) error {
	_, err := us.userClient.UpdateUserGender(ctx, &dto.UpdateUserGenderRequest{Gender: int32(gender)})
	if err != nil {
		return exception.ParseGrpcError(err)
	}

	return nil
}

func (us *DefaultUserService) UploadAvatar(ctx context.Context, fh *multipart.FileHeader) (string, error) {
	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ctx = metadata.AppendToOutgoingContext(ctx, "filename", fh.Filename)
	stream, err := us.userClient.UploadAvatar(ctx)
	if err != nil {
		return "", exception.ParseGrpcError(err)
	}

	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			plog.Errorc(ctx, "failed to read avatar, error: %v", err)
			return "", err
		}

		err = stream.Send(&dto.AvatarByte{
			ByteData: buf[:n],
		})
		if err != nil {
			plog.Errorc(ctx, "failed to send avatar, error: %v", err)
			return "", exception.ParseGrpcError(err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", exception.ParseGrpcError(err)
	}

	return resp.GetAvatarId(), nil
}

func (us *DefaultUserService) GetAvatar(ctx context.Context, avatarId string, writer io.Writer) error {
	stream, err := us.userClient.GetAvatar(ctx, &dto.GetAvatarRequest{AvatarId: avatarId})
	if err != nil {
		return exception.ParseGrpcError(err)
	}

	reader := &streamReader{Stream: stream}
	_, err = io.Copy(writer, reader)
	if err != nil {
		return errors.Wrap(err, "ReadFromStream")
	}

	return nil
}

type streamReader struct {
	Stream grpc.ServerStreamingClient[dto.AvatarByte]
	buffer []byte
}

func (sr *streamReader) Read(p []byte) (n int, err error) {
	if len(sr.buffer) > 0 {
		n = copy(p, sr.buffer)
		sr.buffer = sr.buffer[n:]

		return n, nil
	}

	req, err := sr.Stream.Recv()
	if err == io.EOF {
		return 0, io.EOF
	} else if status.Code(err) == codes.ResourceExhausted {
		return 0, exception.ErrFileTooLarge
	} else if err != nil {
		return 0, err
	}

	sr.buffer = req.ByteData

	n = copy(p, sr.buffer)
	sr.buffer = sr.buffer[n:]
	return n, nil
}
