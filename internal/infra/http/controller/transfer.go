package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	http_middleware "github.com/marcoscoutinhodev/pp_chlg/internal/infra/http/middleware"
	"github.com/marcoscoutinhodev/pp_chlg/internal/usecase"
)

type Transfer struct {
	transferUseCase usecase.Transfer
}

func NewTransfer(tuc usecase.Transfer) *Transfer {
	return &Transfer{
		transferUseCase: tuc,
	}
}

func (t Transfer) Transfer(w http.ResponseWriter, r *http.Request) {
	var input usecase.TransferInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "failed to parse request body",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	payerID := r.Context().Value(http_middleware.UserIDKeyContext{})
	if payerID == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "internal error, try again in a few minutes",
		})
		fmt.Println("[transfer controller]: userID not found in context of protected route")
		return
	}

	input.Payer = payerID.(string)

	output, err := t.transferUseCase.Transfer(ctx, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "internal error, try again in a few minutes",
		})
		fmt.Printf("internal error: %v\n", err)
		return
	}

	w.WriteHeader(output.StatusCode)
	json.NewEncoder(w).Encode(output)
}

func (t Transfer) List(w http.ResponseWriter, r *http.Request) {
	pageAsString := r.Header.Get("page")
	limitAsString := r.Header.Get("limit")
	page, err := strconv.Atoi(pageAsString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "page must be valid",
		})
		return
	}

	limit, err := strconv.Atoi(limitAsString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "limit must be valid",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	userID := r.Context().Value(http_middleware.UserIDKeyContext{})
	if userID == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "internal error, try again in a few minutes",
		})
		fmt.Println("[transfer controller]: userID not found in context of protected route")
		return
	}

	input := usecase.TransferListInputDTO{
		UserID: userID.(string),
		Page:   int64(page),
		Limit:  int64(limit),
	}

	output, err := t.transferUseCase.List(ctx, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "internal error, try again in a few minutes",
		})
		fmt.Printf("internal error: %v\n", err)
		return
	}

	w.WriteHeader(output.StatusCode)
	json.NewEncoder(w).Encode(output)
}
