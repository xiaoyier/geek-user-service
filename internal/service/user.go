
package service

import (
	"context"
	pb "geek-user-service/api/user/v1"
	"geek-user-service/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type UserService struct {
	pb.UnimplementedUserServer

	uc  *biz.UserUsecase
	logger *log.Helper
}

func NewUserService(uc *biz.UserUsecase, logger log.Logger ) *UserService {
	return &UserService{
		uc: uc,
		logger: log.NewHelper(log.With(logger, "modulr", "user.service")),
	}
}

func (s *UserService) UserSignUp(ctx context.Context, req *pb.UserSignUpRequest) (*pb.UserSignUpReply, error) {
	// DTO -> DO
	user := &biz.User{
		UserName: req.UserName,
		Password: req.Password,
	}
	err := s.uc.SignUp(ctx, user)
	if err != nil {
		s.logger.Errorf("UserSignUp: %v", err)
	}
	if errors.Is(err, biz.ErrRecordAlreadyExist) {
		return &pb.UserSignUpReply{
			Code: int32(ErrCodeUserExisted),
			Message: ErrCodeUserExisted.Message(),
		}, nil
	}else if err != nil {
		return nil, InternalServerWithCode(ErrCodeInternalError)
	}
	// DO -> DTO
	return &pb.UserSignUpReply{
		Code: int32(CodeSucc),
		Message: CodeSucc.Message(),
		Data: &pb.UserSignUpReply_UserSignupInfo{
			UserId: user.UserID,
			UserName: user.UserName,
		},
	}, nil
}
func (s *UserService) UserLogin(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginReply, error) {
	// DTO -> DO
	user := &biz.User{
		UserName: req.UserName,
		Password: req.Password,
	}
	result, err := s.uc.Login(ctx, user)
	if err != nil {
		s.logger.Errorf("UserLogin: %v", err)
	}
	if errors.Is(err, biz.ErrRecordNotFound) {
		return nil, NotFoundWithCode(ErrCodeUserNotFound)
	} else if errors.Is(err, biz.ErrPasswordWrong) {
		return &pb.UserLoginReply{
			Code: int32(ErrCodePasswordWrong),
			Message: ErrCodePasswordWrong.Message(),
		}, nil
	} else if err != nil {
		return nil, InternalServerWithCode(ErrCodeInternalError)
	}

	// DO -> DTO
	return &pb.UserLoginReply{
		Code: int32(CodeSucc),
		Message: CodeSucc.Message(),
		Data: &pb.UserLoginReply_UserLoginInfo{
			UserId: result.UserID,
			UserName: result.UserName,
			LoginToken: result.LoginToken,
			RefreshToken: result.RefreshToken,
		},
	}, nil
}
func (s *UserService) UserUpdate(ctx context.Context, req *pb.UserUpdateRequest) (*pb.UserUpdateReply, error) {
	// DTO -> DO
	user := &biz.User{
		UserName: req.UserName,
		Password: req.Password,
	}

	result,  err := s.uc.Update(ctx, user, req.NewName)
	if err != nil {
		s.logger.Errorf("UserUpdate: %v", err)
	}
	if errors.Is(err, biz.ErrRecordNotFound) {
		return nil, NotFoundWithCode(ErrCodeUserExisted)
	} else if errors.Is(err, biz.ErrPasswordWrong) {
		return &pb.UserUpdateReply{
			Code: int32(ErrCodePasswordWrong),
			Message: ErrCodePasswordWrong.Message(),
		}, nil
	} else if err != nil {
		return nil, InternalServerWithCode(ErrCodeInternalError)
	}

	// DO -> DTO
	return &pb.UserUpdateReply{
		Code: int32(CodeSucc),
		Message: CodeSucc.Message(),
		Data: &pb.UserLoginReply_UserLoginInfo{
			UserId: result.UserID,
			UserName: result.UserName,
			LoginToken: result.LoginToken,
			RefreshToken: result.RefreshToken,
		},
	}, nil
}
