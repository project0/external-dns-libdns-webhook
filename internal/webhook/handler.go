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

// NegotiateHandler returns the domain filter, external-dns requires the exact media type as content type.
func (s *Server) NegotiateHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, http.StatusOK, s.provider.GetDomainFilter(r.Context()))
}

// RecordsHandler returns all current records of the managed zones.
func (s *Server) RecordsHandler(w http.ResponseWriter, r *http.Request) {
	records, err := s.provider.Records(r.Context())
	if err != nil {
		log.Ctx(r.Context()).Err(err).Msg("Failed to get records")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	writeJSON(w, r, http.StatusOK, records)
}

// ApplyChangesHandler applies the changes calculated by external-dns, it must respond with 204 on success.
func (s *Server) ApplyChangesHandler(w http.ResponseWriter, r *http.Request) {
	var changes externaldns.Changes
	if err := json.NewDecoder(r.Body).Decode(&changes); err != nil {
		log.Ctx(r.Context()).Err(err).Msg("Failed to decode changes")
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := s.provider.ApplyChanges(r.Context(), &changes); err != nil {
		log.Ctx(r.Context()).Err(err).Msg("Failed to apply changes")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AdjustEndpointsHandler canonicalizes the desired endpoints to match the records returned by RecordsHandler.
func (s *Server) AdjustEndpointsHandler(w http.ResponseWriter, r *http.Request) {
	endpoints := externaldns.Endpoints{}
	if err := json.NewDecoder(r.Body).Decode(&endpoints); err != nil {
		log.Ctx(r.Context()).Err(err).Msg("Failed to decode endpoints")
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	endpoints, err := s.provider.AdjustEndpoints(r.Context(), endpoints)
	if err != nil {
		log.Ctx(r.Context()).Err(err).Msg("Failed to adjust endpoints")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	writeJSON(w, r, http.StatusOK, endpoints)
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Ctx(r.Context()).Err(err).Msg("Failed to encode response")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set(contentTypeHeader, mediaTypeFormatAndVersion)
	w.WriteHeader(status)

	if _, err := w.Write(body); err != nil {
		log.Ctx(r.Context()).Err(err).Msg("Failed to write response")
	}
}
