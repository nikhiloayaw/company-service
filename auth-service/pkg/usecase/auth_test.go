package usecase

import (
	"auth-service/pkg/api/service/response"
	"auth-service/pkg/domain"
	"auth-service/pkg/mock"
	"auth-service/pkg/utils"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {

	testCases := map[string]struct {
		input          domain.User
		buildStub      func(authRepo *mock.MockAuthRepo, input domain.User)
		expectedOutput domain.User
		checkSameError bool // to notify that there is a need of checking error is same
		expectedError  error
	}{
		"failed_to_check_user_already_exist_on_db_should_return_error": {
			input: domain.User{
				Email:    "user@gmail.com",
				Password: "password",
			},
			buildStub: func(authRepo *mock.MockAuthRepo, input domain.User) {

				dbErr := errors.New("error from db")
				// expecting a call to auth repo is user already exist
				authRepo.EXPECT().IsUserExist("user@gmail.com").
					Times(1).Return(false, dbErr)
			},
			expectedOutput: domain.User{},
			checkSameError: false,
			expectedError:  errors.New("expecting db error"),
		},
		"user_already_exist_should_return_error": {
			input: domain.User{
				Email:    "already_exist_user@gmail.com",
				Password: "password",
			},
			buildStub: func(authRepo *mock.MockAuthRepo, input domain.User) {

				// expecting call to auth repo for checking user already exist
				authRepo.EXPECT().IsUserExist("already_exist_user@gmail.com").
					Times(1).Return(true, nil) // return true for auth repo call
			},
			expectedOutput: domain.User{},
			checkSameError: true,
			expectedError:  ErrAlreadyExist,
		},
		"failed_save_user_on_db_should_return_error": {
			input: domain.User{
				Email:    "new_user@gmail.com",
				Password: "password",
			},
			buildStub: func(authRepo *mock.MockAuthRepo, input domain.User) {

				// expecting call to auth repo for checking user already exist
				authRepo.EXPECT().IsUserExist("new_user@gmail.com").
					Times(1).Return(false, nil) // return false for use not exist

				dbErr := errors.New("failed to save user on db")
				// expecting a call to auth repo for save user on db
				authRepo.EXPECT().SaveUser(gomock.Any()).Return(input, dbErr)
			},
			expectedOutput: domain.User{},
			checkSameError: false,
			expectedError:  errors.New("some error"),
		},
		"successful_user_sign_up": {
			input: domain.User{
				Email:    "new_user@gmail.com",
				Password: "password",
			},
			buildStub: func(authRepo *mock.MockAuthRepo, input domain.User) {
				// expecting call to auth repo for checking user already exist
				authRepo.EXPECT().IsUserExist("new_user@gmail.com").
					Times(1).Return(false, nil) // return false for use not exist

				// expecting a call to auth repo for save user on db
				authRepo.EXPECT().SaveUser(gomock.Any()).Return(input, nil)
			},
			expectedOutput: domain.User{
				Email:    "new_user@gmail.com",
				Password: "password",
			},
		},
	}

	for name, test := range testCases {

		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// create a gomock controller with testing object
			ctl := gomock.NewController(t)
			// create a mock repository for the usecase
			mockAuthRepo := mock.NewMockAuthRepo(ctl)

			// pass the mock auth repo to setup repo before using in usecase
			test.buildStub(mockAuthRepo, test.input)

			// create auth usecase with the mock deps
			authUseCase := NewAuthUseCase(mockAuthRepo, nil)

			// run the sign up with test input
			actualOutput, actualErr := authUseCase.SignUp(test.input)

			// check the actual output and expected output is equal
			assert.Equal(t, test.expectedOutput, actualOutput)

			// if the test case expect and output error should be same or no error.
			if test.checkSameError || test.expectedError == nil {
				// then error should be same even if it's nil on both side.
				assert.Equal(t, test.expectedError, actualErr)
			} else {
				// otherwise just confirm it's an error
				assert.Error(t, actualErr)
			}
		})
	}
}

func TestSignIn(t *testing.T) {

	// this hash pass to return from mock db while build stub and expected output's password to be same
	hashPass, err := utils.GenerateHashFromPassword("password")
	assert.NoError(t, err)

	testCases := map[string]struct {
		input     domain.User
		buildStub func(authRepo *mock.MockAuthRepo,
			tokenAuth *mock.MockTokenAuth, input domain.User)
		expectedOutput domain.User
		checkSameError bool // to notify that there is a need of checking error is same
		expectedError  error
	}{
		"failed_to_find_user_from_db_should_return_error": {
			input: domain.User{
				Email:    "user@gmail.com",
				Password: "password",
			},
			buildStub: func(authRepo *mock.MockAuthRepo, tokenAuth *mock.MockTokenAuth, input domain.User) {

				// expecting a call to find user by email and return an error
				dbErr := errors.New("error from db")
				authRepo.EXPECT().FindUserByEmail(input.Email).Times(1).
					Return(input, dbErr)
			},
			expectedOutput: domain.User{},
			checkSameError: false,
			expectedError:  errors.New("some error the notify failed to find user from db"),
		},
		"user_not_exist_should_return_not_exist_error": {
			input: domain.User{
				Email:    "not_existing_user@gmail.com",
				Password: "password",
			},
			buildStub: func(authRepo *mock.MockAuthRepo, tokenAuth *mock.MockTokenAuth, input domain.User) {

				notExistUser := domain.User{} // an empty user
				authRepo.EXPECT().FindUserByEmail(input.Email).Times(1).
					Return(notExistUser, nil)
			},
			expectedOutput: domain.User{},
			checkSameError: true, // expecting the error as 'ErrNotExist'
			expectedError:  ErrNotExist,
		},
		"wrong_password_should_return_error_wrong_password": {
			input: domain.User{
				Email:    "user@gmail.com",
				Password: "wrong_password",
			},
			buildStub: func(authRepo *mock.MockAuthRepo, tokenAuth *mock.MockTokenAuth, input domain.User) {

				// update the input with id and a different password hash
				input.ID = 12
				input.ID26 = "random-uuid"
				// hash with different password
				hashPass, err := utils.GenerateHashFromPassword("diff_password")
				assert.NoError(t, err)
				input.Password = hashPass

				authRepo.EXPECT().FindUserByEmail(input.Email).Times(1).
					Return(input, nil)
			},
			expectedOutput: domain.User{},
			checkSameError: true, // expecting an error as 'ErrWrongPassword' and it should be same
			expectedError:  ErrWrongPassword,
		},
		"successful_user_sign_up": {
			input: domain.User{
				Email:    "existing_user@gmail.com",
				Password: "password",
			},
			buildStub: func(authRepo *mock.MockAuthRepo, tokenAuth *mock.MockTokenAuth, input domain.User) {

				// update the input with correct password hash and id
				input.ID = 12
				input.ID26 = "random-uuid"

				// set the already created hash password to input
				input.Password = hashPass

				authRepo.EXPECT().FindUserByEmail(input.Email).Times(1).
					Return(input, nil)
			},
			expectedOutput: domain.User{
				ID:       12,
				ID26:     "random-uuid",
				Email:    "existing_user@gmail.com",
				Password: hashPass, // use the already created hash pass to make to output same
			},
		},
	}

	for name, test := range testCases {

		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// create a gomock controller with testing object
			ctl := gomock.NewController(t)
			// create a mock repository for the usecase
			mockAuthRepo := mock.NewMockAuthRepo(ctl)
			// create a mock token auth
			mockTokenAuth := mock.NewMockTokenAuth(ctl)

			// pass the mock auth repo and token auth to setup repo before using in usecase
			test.buildStub(mockAuthRepo, mockTokenAuth, test.input)

			// create auth usecase with the mock deps
			authUseCase := NewAuthUseCase(mockAuthRepo, mockTokenAuth)

			// run the sign in with test input
			actualOutput, actualErr := authUseCase.SignIn(test.input)

			// check the actual output and expected output is equal
			assert.Equal(t, test.expectedOutput, actualOutput)

			// if the test case expect and output error should be same or no error.
			if test.checkSameError || test.expectedError == nil {
				// then error should be same even if it's nil on both side.
				assert.Equal(t, test.expectedError, actualErr)
			} else {
				// otherwise just confirm it's an error
				assert.Error(t, actualErr)
			}
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {

	testCases := map[string]struct {
		input          domain.User
		buildStub      func(tokenAuth *mock.MockTokenAuth)
		expectedOutput response.Token
		checkSameError bool // to notify that there is a need of checking error is same
		expectedError  error
	}{
		"failed_to_generate_token_should_return_error": {
			input: domain.User{
				ID:       12,
				ID26:     "uuid",
				Email:    "user@gmail.com",
				Password: "password",
			},
			buildStub: func(tokenAuth *mock.MockTokenAuth) {

				tokenErr := errors.New("failed to generate token")
				tokenAuth.EXPECT().GenerateToken(gomock.Any()).Times(1).
					Return("", tokenErr)
			},
			expectedOutput: response.Token{},
			checkSameError: false,
			expectedError:  errors.New("some error"),
		},
		"successful_access_token_generate": {
			input: domain.User{
				ID:       12,
				ID26:     "uuid",
				Email:    "user@gmail.com",
				Password: "password",
			},
			buildStub: func(tokenAuth *mock.MockTokenAuth) {

				// return a successful access token
				tokenAuth.EXPECT().GenerateToken(gomock.Any()).Times(1).
					Return("access_token", nil)
			},
			expectedOutput: response.Token{AccessToken: "access_token"},
		},
	}

	for name, test := range testCases {

		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// create a gomock controller with testing object
			ctl := gomock.NewController(t)
			// create a mock token auth
			mockTokenAuth := mock.NewMockTokenAuth(ctl)

			// pass the mock token auth to setup repo before using in usecase
			test.buildStub(mockTokenAuth)

			// create auth usecase with the mock deps
			authUseCase := NewAuthUseCase(nil, mockTokenAuth)

			// run the sign in with test input
			actualOutput, actualErr := authUseCase.GenerateAccessToken("role", test.input)

			// check the output token string only,no need of expire time
			assert.Equal(t, test.expectedOutput.AccessToken, actualOutput.AccessToken)

			// if the test case expect and output error should be same or no error.
			if test.checkSameError || test.expectedError == nil {
				// then error should be same even if it's nil on both side.
				assert.Equal(t, test.expectedError, actualErr)
			} else {
				// otherwise just confirm it's an error
				assert.Error(t, actualErr)
			}
		})
	}
}
