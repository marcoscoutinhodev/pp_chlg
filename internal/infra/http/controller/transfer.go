package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

func (t Transfer) Handle(w http.ResponseWriter, r *http.Request) {
	var input usecase.InputTransferDTO
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

	output, err := t.transferUseCase.Execute(ctx, &input)
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
