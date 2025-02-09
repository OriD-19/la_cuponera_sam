package domain

// Domain layer implementation for the User store

import (
	"OriD19/webdev2/types"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Users struct {
	store types.UserStore
}

func NewUsersDomain(s types.UserStore) *Users {
	return &Users{
		store: s,
	}
}

func (u *Users) RegisterClient(ctx context.Context, body []byte) (*types.Client, error) {
	var client types.Client

	err := json.Unmarshal(body, &client)

	if err != nil {
		return &types.Client{}, err
	}

	validate := validator.New()

	err = validate.Struct(client)

	if err != nil {
		return &types.Client{}, err
	}

	// check if email is not taken
	_, err = u.store.GetClient(ctx, client.Email)

	if err == nil {
		return &types.Client{}, fmt.Errorf("email %s is already taken", client.Email)
	}

	// hash password before storing the user
	err = client.HashPassword()

	if err != nil {
		return &types.Client{}, err
	}

	err = u.store.RegisterClient(ctx, client)

	if err != nil {
		return &types.Client{}, err
	}

	return &client, nil
}

func (u *Users) RegisterEnterprise(ctx context.Context, body []byte) (*types.Enterprise, error) {

	var enterprise types.Enterprise

	err := json.Unmarshal(body, &enterprise)

	if err != nil {
		return &types.Enterprise{}, err
	}

	validate := validator.New()

	err = validate.Struct(enterprise)

	if err != nil {
		return &types.Enterprise{}, err
	}

	_, err = u.store.GetEnterprise(ctx, enterprise.Email)

	if err == nil {
		return &types.Enterprise{}, fmt.Errorf("email %s is already taken", enterprise.Email)
	}

	// hash password before storing the user
	err = enterprise.HashPassword()

	if err != nil {
		return &types.Enterprise{}, err
	}

	err = u.store.RegisterEnterprise(ctx, enterprise)

	if err != nil {
		return &types.Enterprise{}, err
	}

	return &enterprise, nil
}

func (u *Users) RegisterAdministrator(ctx context.Context, body []byte) (*types.Administrator, error) {

	var administrator types.Administrator

	err := json.Unmarshal(body, &administrator)

	if err != nil {
		return &types.Administrator{}, err
	}

	validate := validator.New()

	err = validate.Struct(administrator)

	if err != nil {
		return &types.Administrator{}, err
	}

	_, err = u.store.GetAdministrator(ctx, administrator.Email)

	if err == nil {
		return &types.Administrator{}, fmt.Errorf("email %s is already taken", administrator.Email)
	}

	// hash password before storing the user
	err = administrator.HashPassword()

	if err != nil {
		return &types.Administrator{}, err
	}

	err = u.store.RegisterAdministrator(ctx, administrator)

	if err != nil {
		return &types.Administrator{}, err
	}

	return &administrator, nil
}

func (u *Users) RegisterEmployee(ctx context.Context, body []byte) (*types.Employee, error) {

	var employee types.Employee

	err := json.Unmarshal(body, &employee)

	if err != nil {
		return &types.Employee{}, err
	}

	validate := validator.New()

	err = validate.Struct(employee)

	if err != nil {
		return &types.Employee{}, err
	}

	_, err = u.store.GetEmployee(ctx, employee.Email)

	if err == nil {
		return &types.Employee{}, fmt.Errorf("email %s is already taken", employee.Email)
	}

	// hash password before storing the user
	err = employee.HashPassword()

	if err != nil {
		return &types.Employee{}, err
	}

	err = u.store.RegisterEmployee(ctx, employee)

	if err != nil {
		return &types.Employee{}, err
	}

	return &employee, nil
}

func (u *Users) GetClient(ctx context.Context, username string) (*types.Client, error) {
	client, err := u.store.GetClient(ctx, username)

	if err != nil {
		return &types.Client{}, err
	}

	return &client, nil
}

func (u *Users) GetEnterprise(ctx context.Context, username string) (*types.Enterprise, error) {
	enterprise, err := u.store.GetEnterprise(ctx, username)

	if err != nil {
		return &types.Enterprise{}, err
	}

	return &enterprise, nil
}

func (u *Users) GetAdministrator(ctx context.Context, username string) (*types.Administrator, error) {
	administrator, err := u.store.GetAdministrator(ctx, username)

	if err != nil {
		return &types.Administrator{}, err
	}

	return &administrator, nil
}

func (u *Users) GetEmployee(ctx context.Context, username string) (*types.Employee, error) {
	employee, err := u.store.GetEmployee(ctx, username)

	if err != nil {
		return &types.Employee{}, err
	}

	return &employee, nil
}
