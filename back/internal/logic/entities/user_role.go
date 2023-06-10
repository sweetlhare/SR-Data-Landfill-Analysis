package logicentities

import "strings"

type UserRole string

const (
	UserRoleUser    UserRole = "user"
	UserRoleAdmin   UserRole = "admin"
	UserRoleInvalid UserRole = "invalid"
)

// validRoles ...
var validRoles map[string]UserRole
var validRolesString string

// init ...
func init() {
	addValidRoles(
		UserRoleUser,
		UserRoleAdmin,
	)
}

// addValidRoles ...
func addValidRoles(urs ...UserRole) {
	rolesStringArr := make([]string, len(urs))
	m := make(map[string]UserRole, len(urs))
	for i, ur := range urs {
		roleString := ur.ToString()
		m[roleString] = ur
		rolesStringArr[i] = roleString
	}
	validRoles = m
	validRolesString = strings.Join(rolesStringArr, ", ")
}

// GetValidRoles ...
func GetValidRoles() string {
	return validRolesString
}

// ToString ...
func (r UserRole) ToString() string {
	return string(r)
}

// IsValid ...
func (r UserRole) IsValid() bool {
	return r != UserRoleInvalid
}

// UserRoleFromString ...
func UserRoleFromString(s string) UserRole {
	val, ok := validRoles[s]
	if ok {
		return val
	}
	return UserRoleInvalid
}

// IsValid ...
func (r UserRole) Equel(userRoles ...UserRole) bool {
	for _, nextR := range userRoles {
		if r.ToString() == nextR.ToString() {
			return true
		}
	}
	return false
}
