package entity

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Email     string
	Password  string
	FirstName string
	LastName  string
	Role      string
}

// UserRole defines the available roles for users
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// NewUser creates a new user with the given details
func NewUser(email, password, firstName, lastName string, role UserRole) (*User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		Role:      string(role),
	}, nil
}

func (u *User) Update(firstName string, lastName string) error {
	if firstName == "" {
		return errors.New("first name cannot be empty")
	}
	if lastName == "" {
		return errors.New("last name cannot be empty")
	}

	u.FirstName = firstName
	u.LastName = lastName

	return nil
}

// ComparePassword checks if the provided password matches the stored hash
func (u *User) ComparePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if u.Password == "" {
		return errors.New("user password is not set")
	}

	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// UpdatePassword updates the user's password
func (u *User) UpdatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == string(RoleAdmin)
}
