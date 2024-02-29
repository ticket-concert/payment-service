package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"payment-service/configs"
	"payment-service/internal/pkg/constants"
	"regexp"
)

var blackListEmail []string

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func InitReadBlackListEmail() bool {
	fd, error := os.Open("blacklistedEmail.csv")

	if error != nil {
		return false
	}

	defer fd.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(fd)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	CreateBlackListEmail(data)

	return true
}

func CreateBlackListEmail(data [][]string) []string {
	for i, line := range data {
		if i > 0 { // omit header line
			// for _, field := range line {
			// 	blackListEmail = append(blackListEmail, field)
			// }
			blackListEmail = append(blackListEmail, line...)
		}
	}
	return blackListEmail
}

func IsBlacklistedEmail(searchterm string) bool {
	for _, value := range blackListEmail {
		matched, _ := regexp.MatchString(value, searchterm)
		if matched {
			return true
		}
	}
	return false
}

func IsValidPassword(password string) bool {
	// Check if password length is greater than 8 characters
	if len(password) < 8 {
		return false
	}

	// Check if password contains an alphanumeric character
	matched, err := regexp.Match("[a-zA-Z0-9]", []byte(password))
	if err != nil || !matched {
		return false
	}

	// Check if password contains a capital letter
	matched, err = regexp.Match("[A-Z]", []byte(password))
	if err != nil || !matched {
		return false
	}

	// Check if password contains a special character
	matched, err = regexp.Match(".*[!@#$%^&*()_+].*", []byte(password))
	if err != nil || !matched {
		return false
	}

	return true
}

func VerifyPhoneNumber62(phoneNumber string) string {
	mobileNumber := ""
	if phoneNumber != "" {
		// remove spaces and dashes
		phoneNumber = regexp.MustCompile(`[\s-]+`).ReplaceAllString(phoneNumber, "")
		// convert +628 or 628 or 8 to +628
		mobileNumber = regexp.MustCompile(`^\+628|^628|^08`).ReplaceAllString(phoneNumber, "+628")
	}
	return mobileNumber
}

// Generated Password
func GeneratePassword(password string) string {
	//generate signed content
	secret := configs.GetConfig().SecretHashPass
	id := configs.GetConfig().IdHash
	msgHashSig := hmac.New(sha256.New, []byte(secret))
	_, err := msgHashSig.Write([]byte(fmt.Sprintf("%v&%v", id, password)))
	if err != nil {
		return "error hash signature"
	}

	msgHashSumSig := msgHashSig.Sum(nil)

	// sign
	signatureSig := base64.StdEncoding.EncodeToString(msgHashSumSig)
	return signatureSig
}

func GenerateRandomOtp() string {
	// Generate random otp
	letterRunes := []rune("0123456789")
	otp := make([]rune, 6)
	for i := range otp {
		otp[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(otp)
}

func MaskEmail(email string) string {
	atIndex := -1
	for i, char := range email {
		if char == '@' {
			atIndex = i
			break
		}
	}

	if atIndex == -1 || atIndex <= 1 {
		return email
	}

	maskedEmail := email[:1] + "***" + email[atIndex-1:]
	return maskedEmail
}

func GenerateRandomPassword(length int) (string, error) {
	const lowercaseCharset = "abcdefghijklmnopqrstuvwxyz"
	const uppercaseCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const numericCharset = "0123456789"
	const specialCharset = "!@#$%^&*()_+"

	// Ensure there's enough length for each character type
	if length < 8 {
		return "", fmt.Errorf("password length must be at least 8 characters")
	}

	// Generate one character from each character set
	lowercase, err := getRandomChar(lowercaseCharset)
	if err != nil {
		return "", err
	}

	uppercase, err := getRandomChar(uppercaseCharset)
	if err != nil {
		return "", err
	}

	numeric, err := getRandomChar(numericCharset)
	if err != nil {
		return "", err
	}

	special, err := getRandomChar(specialCharset)
	if err != nil {
		return "", err
	}

	// Generate the remaining characters randomly
	remainingLength := length - 4
	randomPart, err := generateRandomPart(lowercaseCharset+uppercaseCharset+numericCharset+specialCharset, remainingLength)
	if err != nil {
		return "", err
	}

	// Concatenate the characters
	password := lowercase + uppercase + numeric + special + randomPart

	return password, nil
}

func getRandomChar(charset string) (string, error) {
	result := charset[rand.Intn(len(charset))]
	return string(result), nil
}

func generateRandomPart(charset string, length int) (string, error) {
	randomPart := make([]byte, length)
	for i := 0; i < length; i++ {
		randomPart[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomPart), nil
}

func CustomIfEmpty(param1 string, param2 string) string {
	if param1 != "" {
		return param1
	}
	return param2
}

func GenerateMetaData(totalData, count int64, page, limit int64) constants.MetaData {
	metaData := constants.MetaData{
		Page:      page,
		Count:     count,
		TotalPage: int64(math.Ceil(float64(totalData) / float64(limit))),
		TotalData: totalData,
	}

	return metaData
}
