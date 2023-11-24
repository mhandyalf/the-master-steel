package handlers

import (
	"net/http"
	"os"
	"the-master-steel/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Auth struct {
	db *gorm.DB
}

func NewAuth(db *gorm.DB) *Auth {
	return &Auth{db: db}
}

func (a *Auth) Register(e echo.Context) error {
	user := new(models.Employee)
	if err := e.Bind(user); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	user.Password = string(hashedPassword)
	if err := a.db.Create(user).Error; err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return e.JSON(http.StatusCreated, user)

}

func (a *Auth) Login(e echo.Context) error {
	user := new(models.Employee)
	if err := e.Bind(user); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	var foundUser models.Employee
	if err := a.db.Where("name = ?", user.Name).First(&foundUser).Error; err != nil {
		return e.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		return e.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
	}

	// Buat token JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = foundUser.Name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token berlaku selama 1 hari

	// Sign token dengan secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Gagal membuat token JWT",
		})
	}

	// Kirim token sebagai respons
	return e.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})
}

func (a *Auth) GetEmployeeInfo(e echo.Context) error {
	var products []models.Employee
	if err := a.db.Find(&products).Error; err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, products)
}
