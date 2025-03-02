package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/rs/zerolog/log"
)

const (
	mediaTypeFormatAndVersion = "application/external.dns.webhook+json;version=1"
	contentTypeHeader         = "Content-Type"
)

type Handler struct {
	Provider externaldns.Provider
}

func (h *Handler) NegotiateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, mediaTypeFormatAndVersion)

	if err := json.NewEncoder(w).Encode(h.Provider.GetInitialization(r.Context())); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) RecordsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		records, err := h.Provider.Records(r.Context())
		if err != nil {
			log.Err(err).Msg("Failed to get Records")
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Set(contentTypeHeader, mediaTypeFormatAndVersion)
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(records); err != nil {
			log.Err(err).Msg("Failed to encode records")
		}

		return
	case http.MethodPost:
		var changes externaldns.Changes
		if err := json.NewDecoder(r.Body).Decode(&changes); err != nil {
			log.Err(err).Msg("Failed to decode changes")
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		err := h.Provider.ApplyChanges(r.Context(), &changes)
		if err != nil {
			log.Err(err).Msg("Failed to apply changes")
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusNoContent)

		return
	default:
		log.Error().Str("method", r.Method).Msg("Unsupported method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) AdjustendpointsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Error().Str("method", r.Method).Msg("Unsupported method")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	pve := make(externaldns.Endpoints, 0)
	if err := json.NewDecoder(r.Body).Decode(&pve); err != nil {
		log.Err(err).Msg("Failed to decode in adjustEndpointsHandler")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	w.Header().Set(contentTypeHeader, mediaTypeFormatAndVersion)

	pve, err := h.Provider.AdjustEndpoints(r.Context(), pve)
	if err != nil {
		log.Err(err).Msg("Failed to call adjust endpoints")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(&pve); err != nil {
		log.Err(err).Msg("Failed to encode in adjustEndpointsHandler")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
