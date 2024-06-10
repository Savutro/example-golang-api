package auth

import (
	"fmt"
	"net/http"
	"time"

	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/config"
	"git.gibb.ch/faf141769/infw-22a-m152-teamsigma/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
	"golang.org/x/crypto/bcrypt"
)

var (
	cfg         config.Config
	Store       *sessions.CookieStore
	SessionName string
	db          *gorm.DB
)

func Init(database *gorm.DB, config config.Config) {
	db = database
	cfg = config
	SessionName = cfg.SessionName
	Store = sessions.NewCookieStore([]byte(cfg.CookiePassword))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

func RegisterNewUser(username string, password string, secret string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Correct implementation
	user := &models.User{
		UUID:            uuid.New().String(),
		Username:        username,
		Password:        string(passwordHash),
		TwoFactorSecret: string(secret),
		Role:            "standard",
	}
	db.Create(&user)

	// NOTE: Vulnerable code: raw SQL with direct user input (example', 'password', 'secret', 'admin'); -- )
	// uuid := uuid.New().String()
	// query := fmt.Sprintf(`
	// INSERT INTO users (uuid, username, password, two_factor_secret, role)
	// VALUES ('%s', '%s', '%s', '%s', 'standard')`,
	// 	uuid, username, string(passwordHash), secret)

	// db.Exec(query)

	return nil
}

func LoginUser(username, password string) error {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return fmt.Errorf("username couldn't be found: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return fmt.Errorf("password doesn't match the hash: %v", err)
	}

	return nil
}

func CreateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte(cfg.JwtPassword))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateTOTPWithSecret(randomSecret string, username string) string {
	uri := gotp.NewDefaultTOTP(randomSecret).ProvisioningUri(username, "Golang-API")

	qrcode.WriteFile(uri, qrcode.Medium, 256, "qr.png")

	// NOTE: to print qr code to terminal
	// qrterminal.GenerateWithConfig(uri, qrterminal.Config{
	// 	Level:     qrterminal.L,
	// 	Writer:    os.Stdout,
	// 	BlackChar: qrterminal.BLACK,
	// 	WhiteChar: qrterminal.WHITE,
	// })

	return uri
}

func VerifyOTP(randomSecret string, code string) bool {
	totp := gotp.NewDefaultTOTP(randomSecret) //TODO Read from Database?

	ok := totp.Verify(code, time.Now().Unix())

	return ok
}

func GetSecretFromDB(username string) (string, error) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", fmt.Errorf("username couldn't be found: %v", err)
	}

	return user.TwoFactorSecret, nil
}
