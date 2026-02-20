package user

import (
	"context"
	"errors"
	cfgUtil "github.com/StewardMcCormick/Paste_Bin/config/cfg_util"
	"github.com/StewardMcCormick/Paste_Bin/internal/domain"
	"github.com/StewardMcCormick/Paste_Bin/internal/dto"
	errs "github.com/StewardMcCormick/Paste_Bin/internal/error"
	"github.com/StewardMcCormick/Paste_Bin/internal/usecase/user/mocks"
	appctx "github.com/StewardMcCormick/Paste_Bin/internal/util/app_context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"time"
)

type UseCaseTestSuite struct {
	suite.Suite
	repo     *mocks.MockRepository
	security *mocks.MockSecurityUtil
	valid    *mocks.MockValidator
	useCase  *UseCase
}

func TestUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UseCaseTestSuite))
}

func (s *UseCaseTestSuite) SetupTest() {
	s.repo = mocks.NewMockRepository(s.T())
	s.security = mocks.NewMockSecurityUtil(s.T())
	s.valid = mocks.NewMockValidator(s.T())
	s.useCase = NewUseCase(s.repo, s.security,
		s.valid,
		Config{APIKeyExpireDuration: 162 * time.Hour},
	)
}

func (s *UseCaseTestSuite) Test_Registration_Success() {
	now := time.Now()
	expected := &domain.User{
		Id:       0,
		Username: "test",
		Password: "test_pass",
		APIKey: domain.APIKey{
			Key:       "hashed_key",
			Prefix:    "pb_test",
			ExpiresAt: now.Add(s.useCase.cfg.APIKeyExpireDuration),
		},
		CreatedAt: now,
	}

	s.valid.EXPECT().
		Validate(mock.Anything).
		Return(nil).
		Once()

	s.repo.EXPECT().
		Exists(mock.Anything, mock.Anything).
		Return(false, nil).
		Once()

	s.security.EXPECT().
		HashPassword(mock.Anything).
		Return("hashed_pass", nil).
		Once()

	s.security.EXPECT().
		GenerateAPIKey(mock.Anything).
		Return("pb_test", "pb_test_api-key", nil).
		Once()

	s.security.EXPECT().
		HashAPIKey(mock.Anything).
		Return("hashed_key").
		Once()

	s.repo.EXPECT().
		CreateUser(mock.Anything, mock.Anything).
		Return(expected, nil).
		Once()

	ctx := context.WithValue(context.Background(), appctx.EnvKey, cfgUtil.DevelopmentEnv)
	result, err := s.useCase.Registration(ctx,
		&dto.CreateUserRequest{Username: expected.Username, Password: expected.Password})

	s.NoError(err)
	s.Equal(expected.Id, result.Id)
	s.Equal(expected.Username, result.Username)
	s.Equal(expected.APIKey.Key, result.APIKey.Key)
	s.Equal(expected.APIKey.ExpiresAt, result.APIKey.ExpiresAt)
	s.True(strings.HasPrefix(result.APIKey.Key, expected.APIKey.Prefix))
	s.NotNil(result.CreatedAt)
}

func (s *UseCaseTestSuite) Test_Registration_NotValidationError() {
	cases := []struct {
		name       string
		setupMocks func()
		wantError  error
	}{
		{
			"Check user existing - Already Exists Error",
			func() {
				s.valid.EXPECT().
					Validate(mock.Anything).
					Return(nil).
					Once()

				s.repo.EXPECT().
					Exists(mock.Anything, mock.Anything).
					Return(true, nil).
					Once()
			},
			errs.UserAlreadyExists,
		},
		{
			"Check user existing - Internal Error",
			func() {
				s.valid.EXPECT().
					Validate(mock.Anything).
					Return(nil).
					Once()

				s.repo.EXPECT().
					Exists(mock.Anything, mock.Anything).
					Return(false, errors.New("db error")).
					Once()
			},
			errs.InternalError,
		},
		{
			"Hashing password error",
			func() {
				s.valid.EXPECT().
					Validate(mock.Anything).
					Return(nil).
					Once()

				s.repo.EXPECT().
					Exists(mock.Anything, mock.Anything).
					Return(false, nil).
					Once()

				s.security.EXPECT().
					HashPassword(mock.Anything).
					Return("", errors.New("hashing error")).
					Once()
			},
			errs.InternalError,
		},
		{
			"Generate API Key error",
			func() {
				s.valid.EXPECT().
					Validate(mock.Anything).
					Return(nil).
					Once()

				s.repo.EXPECT().
					Exists(mock.Anything, mock.Anything).
					Return(false, nil).
					Once()

				s.security.EXPECT().
					HashPassword(mock.Anything).
					Return("hash", nil).
					Once()

				s.security.EXPECT().
					GenerateAPIKey(mock.Anything).
					Return("", "", errors.New("generate API Key error")).
					Once()
			},
			errs.InternalError,
		},
		{
			"Registration internal error",
			func() {
				s.valid.EXPECT().
					Validate(mock.Anything).
					Return(nil).
					Once()

				s.repo.EXPECT().
					Exists(mock.Anything, mock.Anything).
					Return(false, nil).
					Once()

				s.security.EXPECT().
					HashPassword(mock.Anything).
					Return("hash", nil).
					Once()

				s.security.EXPECT().
					GenerateAPIKey(mock.Anything).
					Return("pb_test", "key", nil).
					Once()

				s.security.EXPECT().
					HashAPIKey(mock.Anything).
					Return("key_hash").
					Once()

				s.repo.EXPECT().
					CreateUser(mock.Anything, mock.Anything).
					Return(nil, errors.New("db error")).
					Once()
			},
			errs.InternalError,
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			s.SetupTest()
			tc.setupMocks()
			result, err := s.useCase.Registration(context.Background(),
				&dto.CreateUserRequest{Username: "test", Password: "test_pass"},
			)

			s.Nil(result)
			s.NotNil(err)
			s.ErrorIs(err, tc.wantError)
		})
	}
}

func (s *UseCaseTestSuite) Test_Registration_ValidationError() {
	cases := []struct {
		name      string
		value     *dto.CreateUserRequest
		setupMock func()
		wantErr   error
	}{
		{
			name:  "user with empty fields",
			value: &dto.CreateUserRequest{Username: "", Password: ""},
			setupMock: func() {
				s.valid.EXPECT().
					Validate(mock.Anything).
					Return(errors.New("user with empty fields error")).
					Once()
			},
			wantErr: errors.New("user with empty fields error"),
		},
		{
			name:  "user with too short fields",
			value: &dto.CreateUserRequest{Username: "Us", Password: "pass"},
			setupMock: func() {
				s.valid.EXPECT().
					Validate(mock.Anything).
					Return(errors.New("user with too short fields error")).
					Once()
			},
			wantErr: errors.New("user with too short fields error"),
		},
		{
			name:  "user with too long fields",
			value: &dto.CreateUserRequest{Username: "too_long_field", Password: "too_long_password"},
			setupMock: func() {
				s.valid.EXPECT().
					Validate(mock.Anything).
					Return(errors.New("user with too long fields error")).
					Once()
			},
			wantErr: errors.New("user with too long fields error"),
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			s.SetupTest()
			tc.setupMock()

			result, err := s.useCase.Registration(context.Background(), tc.value)

			s.Nil(result)
			s.NotNil(err)
			s.Equal(err, tc.wantErr)
		})
	}
}
