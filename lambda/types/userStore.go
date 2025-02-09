package types

import "context"

// This is a single struct for operating over all the types of users available in the system
type UserStore interface {
	//DoesUserExist(username string) (bool, error)

	// Registering different types of users
	// !IMPORTANT: REGISTER METHODS ALSO UPDATES IF THE VALUE ALREADY EXISTS
	RegisterClient(context.Context, Client) error
	RegisterEnterprise(context.Context, Enterprise) error
	RegisterAdministrator(context.Context, Administrator) error
	RegisterEmployee(context.Context, Employee) error

	// Getting different types of users
	GetClient(context.Context, string) (Client, error)
	GetEnterprise(context.Context, string) (Enterprise, error)
	GetAdministrator(context.Context, string) (Administrator, error)
	GetEmployee(context.Context, string) (Employee, error)
}
