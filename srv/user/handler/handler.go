package handler

import (
	"context"
	"errors"
	"fmt"
	"originals/conf"
	"originals/email"
	"originals/jwt"
	"originals/srv/user/model"
	"originals/srv/user/proto"
	"regexp"
	"time"

	"github.com/micro/go-log"
)

type UserSrvHandler struct {
	Model *model.UserSrvModel
}

const signUpEmailBody = `
<!doctype html>
<html>
	<body>
		<header><h3>Originals Beta v1.0<h3><hr></header>
		<article>
			<p>Hi there,</p>
			<p>Before use the <b>originals</b>, please take a few minutes to set your account. This link will take you to the page:<br>
			<a href="/user/signup?token=%s"><i>Account Setting</i></a></p>
			<p><b>Thanks</b></p>
		</article>
		<footer>
			<p>-------------------------------<br>O-P-T</p>
		<footer>
	</body>
</html>
`

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrEmailExist   = errors.New("email has been sign up")
	ErrEmptySecret  = errors.New("jwt secret key is empty")
)

// StartEmail send a sign up email to the user's email if the email is valid and has not been sign up
func (hdlr *UserSrvHandler) StartEmail(ctx context.Context, req *proto.StartEmailReq, rsp *proto.StartEmailRsp) error {
	log.Log("Received User.StartEmail request")

	// Verify email format
	if err := verifyEmail(req.Email); err != nil {
		return err
	}

	// Check if email has been sign up
	if users, err := hdlr.Model.GetUserByEmail(req.Email); len(users) > 0 {
		return ErrEmailExist
	} else if err != nil {
		return err
	}

	// Create sign up token
	secret := conf.SrvConf["user"].Extra["jwt_secret_email"]
	if secret == "" {
		return ErrEmptySecret
	}
	claims := jwt.EmailClaims{Email: req.Email}
	claims.StandardClaims.ExpiresAt = time.Now().Add(5 * time.Minute).Unix()
	tokenStr, err := jwt.CreateToken(claims, secret)
	if err != nil {
		return err
	}

	// Send sign up email
	mailBody := fmt.Sprintf(signUpEmailBody, tokenStr)
	signUpMail := &email.Email{
		Recivers: []string{req.Email},
		Subject:  "Originals-起源-Beta v1.0 注册测试邮件",
		Body:     mailBody,
	}
	if err := email.SendMail(signUpMail); err != nil {
		return err
	}

	rsp.SignUpToken = tokenStr
	return nil
}

func (hdlr *UserSrvHandler) SignUp(ctx context.Context, req *proto.SignUpReq, rsp *proto.SignUpRsp) error {
	return nil
}

// verifyEmail
func verifyEmail(email string) error {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}
