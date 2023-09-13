package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type JwtCustomClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {

	var resp generated.HelloResponse
	fmt.Println(resp)
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func validatePassword(password string) bool {
	// Length validation
	if len(password) < 6 || len(password) > 64 {
		return false
	}

	// Uppercase letter validation
	uppercaseRegex := regexp.MustCompile("[A-Z]")
	if !uppercaseRegex.MatchString(password) {
		return false
	}

	// Numeric digit validation
	digitRegex := regexp.MustCompile("[0-9]")
	if !digitRegex.MatchString(password) {
		return false
	}

	// Special character validation
	specialRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	if !specialRegex.MatchString(password) {
		return false
	}

	return true
}

func (s *Server) Register(ctx echo.Context) error {
	var resQuery repository.UsersEntity
	var req generated.RegisterRequest
	var isErr bool = true
	var errMap generated.RegisterErrorResponse = make(generated.RegisterErrorResponse)
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	if len(req.PhoneNumber) < 10 || len(req.PhoneNumber) > 13 {
		isErr = false
		errMap[`phone_number`] = `Phone numbers must be at minimum 10 characters and maximum 13 characters`
	}
	if req.PhoneNumber[0:3] != `+62` {
		isErr = false
		errMap[`phone_number`] = `Phone numbers must start with the Indonesia country code “+62”`
	}
	if len(req.Fullname) < 3 || len(req.Fullname) > 60 {
		isErr = false
		errMap[`fullname`] = `Full name must be at minimum 3 characters and maximum 60 characters`
	}
	if !validatePassword(req.Password) {
		isErr = false
		errMap[`password`] = `Passwords must be minimum 6 characters and maximum 64 characters,
		containing at least 1 capital characters AND 1 number AND 1 special (nonalpha-numeric) characters`
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := repository.UsersEntity{
		FullName:          req.Fullname,
		PhoneNumber:       req.PhoneNumber,
		Password:          string(pwd),
		SuccessLoginCount: 0,
		CreatedAt:         time.Now().UTC(),
	}
	resQuery, err = s.Repository.InsertNewUser(ctx.Request().Context(), newUser)
	if err != nil {
		return err
	}
	var resp generated.RegisterSuccessResponse
	fmt.Println(resp)
	resp.Id = resQuery.Id
	// resp.
	if !isErr {
		return ctx.JSON(http.StatusOK, errMap)
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Login(ctx echo.Context) error {
	var resQuery repository.UsersEntity
	var req generated.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	loginUser := repository.UsersEntity{
		PhoneNumber: req.PhoneNumber,
	}
	resQuery, err := s.Repository.Login(ctx.Request().Context(), loginUser)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(resQuery.Password)
	err = bcrypt.CompareHashAndPassword([]byte(resQuery.Password), []byte(req.Password))
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		resQuery.Id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(utils.SecretKey))
	if err != nil {
		return err
	}

	s.Repository.UpdateLoginCount(ctx.Request().Context(), loginUser)

	var resp generated.LoginSuccessResponse
	resp.AccessToken = t

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) MyProfile(ctx echo.Context) error {
	var resQuery repository.UsersEntity
	user := repository.UsersEntity{
		Id: 1,
	}
	resQuery, err := s.Repository.GetMyMyProfile(ctx.Request().Context(), user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var resp generated.MyProfileSuccessResponse
	resp.Name = resQuery.FullName
	resp.PhoneNumber = resQuery.PhoneNumber
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateProfile(ctx echo.Context) error {
	var req generated.UpdateProfileRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	updateUser := repository.UsersEntity{
		PhoneNumber: req.PhoneNumber,
		FullName:    req.Fullname,
		Id:          1,
	}
	_, err := s.Repository.UpdateUser(ctx.Request().Context(), updateUser)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return ctx.JSON(http.StatusOK, `{}`)
}
