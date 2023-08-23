package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"

	"github.com/irbgeo/user-balance/internal/service"
)

//	@Summary		Get balance of a user
//	@Description	Get the balance of a user by user ID.
//	@Tags			Balance
//	@Produce		plain
//	@Param			id	path		int		true	"User ID"
//	@Success		200	{string}	string	"User balance"
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		500	{string}	string	"Internal Server Error"
//
//	@Router			/{id} [get]
func getBalance(w http.ResponseWriter, r *http.Request, svc svc) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	balance, err := svc.GetBalance(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), getCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(balance.Balance))) // nolint: errcheck
}

//	@Summary		Increase user balance
//	@Description	Increase the balance of a user by a specified amount.
//	@Tags			Balance
//	@Accept			json
//	@Produce		plain
//	@Param			id		path		int				true	"User ID"
//	@Param			request	body		changeBalance	true	"Balance change amount"
//	@Success		200		{string}	string			"OK"
//	@Failure		400		{string}	string			"Bad Request"
//	@Failure		500		{string}	string			"Internal Server Error"
//
//	@Router			/{id} [post]
func upBalance(w http.ResponseWriter, r *http.Request, svc svc) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	changeBalance, err := getBalanceChanging(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bc := service.BalanceChange{
		User:     user,
		Changing: changeBalance.Amount,
	}

	if err := svc.UpBalance(r.Context(), bc); err != nil {
		http.Error(w, err.Error(), getCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

//	@Summary		Decrease user balance
//	@Description	Decrease the balance of a user by a specified amount.
//	@Tags			Balance
//	@Accept			json
//	@Produce		plain
//	@Param			id		path		int				true	"User ID"
//	@Param			request	body		changeBalance	true	"Balance change amount"
//	@Success		200		{string}	string			"OK"
//	@Failure		400		{string}	string			"Bad Request"
//	@Failure		500		{string}	string			"Internal Server Error"
//
//	@Router			/{id} [put]
func downBalance(w http.ResponseWriter, r *http.Request, svc svc) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	changeBalance, err := getBalanceChanging(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bc := service.BalanceChange{
		User:     user,
		Changing: changeBalance.Amount,
	}

	if err := svc.DownBalance(r.Context(), bc); err != nil {
		http.Error(w, err.Error(), getCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getUser(r *http.Request) (service.User, error) {
	idStr := path.Base(r.URL.Path)
	if idStr == "" {
		return service.User{}, fmt.Errorf("Invalid URL format: not found user id")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return service.User{}, fmt.Errorf("Invalid URL format: %w", err)
	}

	return service.User{ID: id}, nil
}

func getCode(err error) int {
	switch err {
	case service.ErrInvalidBalance, service.ErrInvalidNewBalance, service.ErrUserNotExist:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func getBalanceChanging(r *http.Request) (changeBalance, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return changeBalance{}, fmt.Errorf("Failed to read body")
	}
	defer r.Body.Close()

	var cb changeBalance

	if err := json.Unmarshal(body, &cb); err != nil {
		return changeBalance{}, fmt.Errorf("Failed to unmarshal body: %w", err)
	}
	return cb, nil
}
