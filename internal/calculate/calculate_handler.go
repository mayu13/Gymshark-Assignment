package calculate

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mayu13/gymshark-assignment/api"
	"github.com/mayu13/gymshark-assignment/internal/packs"
	"github.com/sirupsen/logrus"
)

func CalculatePacksHandler(pm packs.PacksManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logrus.NewEntry(logrus.StandardLogger())
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("Failed to parse body")
			renderStatusCode(w, http.StatusBadRequest)
			return
		}

		var req api.CalculatePacksRequest
		if err := json.Unmarshal(body, &req); err != nil {
			log.WithError(err).Error("Failed to parse request")
			renderStatusCode(w, http.StatusBadRequest)
			return
		}

		packs, err := pm.CalculatePacks(req.ItemsCount)
		if err != nil {
			log.WithError(err).Error("Failed to calculate packs")
			renderStatusCode(w, http.StatusInternalServerError)
			return
		}

		packsReponse := make([]api.Pack, len(packs))
		for i, p := range packs {
			packsReponse[i] = api.Pack{
				Size:  p.Size,
				Count: p.Quantity,
			}
		}

		response := api.CalculatePacksResponse{
			Packs: packsReponse,
		}

		data, err := json.Marshal(response)
		if err != nil {
			log.WithError(err).Error("Failed to marshal response")
			renderStatusCode(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(data); err != nil {
			log.WithError(err).Error("Failed to marshal response")
		}

		log.Info("CalculatePacks Done")
	}
}

func SetPackSizesHandler(pm packs.PacksManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logrus.NewEntry(logrus.StandardLogger())
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("Failed to parse body")
			renderStatusCode(w, http.StatusBadRequest)
			return
		}

		var req api.SetPackSizes
		if err := json.Unmarshal(body, &req); err != nil {
			log.WithError(err).Error("Failed to parse request")
			renderStatusCode(w, http.StatusBadRequest)
			return
		}

		if err := pm.SetPackSizes(req.Sizes); err != nil {
			log.WithError(err).Error("Failed to set packs")
			renderStatusCode(w, http.StatusInternalServerError)
			return
		}

		renderStatusCode(w, http.StatusOK)

		log.Info("setSize Done")
	}
}

func HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logrus.NewEntry(logrus.StandardLogger())

		renderStatusCode(w, http.StatusOK)

		log.Info("check done")
	}
}
