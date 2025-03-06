package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/theQRL/zond-beaconchain-explorer/services"
	"github.com/theQRL/zond-beaconchain-explorer/templates"
	"github.com/theQRL/zond-beaconchain-explorer/types"
	"github.com/theQRL/zond-beaconchain-explorer/utils"
)

// Index will return the main "index" page using a go template
func Index(w http.ResponseWriter, r *http.Request) {
	var indexTemplateFiles = append(layoutTemplateFiles,
		"index/index.html",
		"index/depositChart.html",
		"index/networkStats.html",
		"index/postGenesis.html",
		"index/recentBlocks.html",
		"index/recentEpochs.html",
		"index/depositDistribution.html",
		"svg/bricks.html",
		"svg/professor.html",
		"svg/timeline.html",
		"components/rocket.html",
		"slotViz.html",
	)

	var indexTemplate = templates.GetTemplate(indexTemplateFiles...)

	w.Header().Set("Content-Type", "text/html")
	data := InitPageData(w, r, "index", "", "", indexTemplateFiles)
	pageData := services.LatestIndexPageData()

	// data.Data.(*types.IndexPageData).ShowSyncingMessage = data.ShowSyncingMessage

	pageData.SlotVizData = getSlotVizData(data.CurrentEpoch)

	calculateChurn(pageData)

	data.Data = pageData

	if handleTemplateError(w, r, "index.go", "Index", "", indexTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return // an error has occurred and was processed
	}
}

// IndexPageData will show the main "index" page in json format
func IndexPageData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", utils.Config.Chain.ClConfig.SecondsPerSlot)) // set local cache to the seconds per slot interval

	err := json.NewEncoder(w).Encode(services.LatestIndexPageData())

	if err != nil {
		logger.Errorf("error sending latest index page data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func getSlotVizData(currentEpoch uint64) *types.SlotVizPageData {
	// var visiblFrom uint64
	// var visibleTo uint64
	configuration, err := services.GetExplorerConfigurationsWithDefaults()
	if err != nil {
		utils.LogError(err, "Could not load SlotViz configuration for index page", 0)
		return nil
	}
	// visiblFrom, err = configuration.GetUInt64Value(services.ConfigurationCategorySlotViz, services.ConfigurationKeyVisibleFromEpoch)
	// if err != nil {
	// 	utils.LogError(err, "Could not get visbleFrom for SlotViz on index page", 0)
	// 	return nil
	// }
	// visibleTo, err = configuration.GetUInt64Value(services.ConfigurationCategorySlotViz, services.ConfigurationKeyVisibleToEpoch)
	// if err != nil {
	// 	utils.LogError(err, "Could not get visibleTo for SlotViz on index page", 0)
	// 	return nil
	// }
	return &types.SlotVizPageData{
		Epochs:   services.LatestSlotVizMetrics(),
		Selector: "slotsViz",
		Config:   configuration[services.ConfigurationCategorySlotViz]}
	/*
		if visiblFrom <= currentEpoch && visibleTo >= currentEpoch {
			return &types.SlotVizPageData{
				Epochs:   services.LatestSlotVizMetrics(),
				Selector: "slotsViz",
				Config:   configuration[services.ConfigurationCategorySlotViz]}

		}
		return nil
	*/
}

func calculateChurn(page *types.IndexPageData) {
	limit := services.GetLatestStats().ValidatorActivationChurnLimit
	pending_validators := services.GetLatestStats().PendingValidatorCount
	// calculate daily new validators
	limit_per_day := *limit * uint64(225)
	// calculate how long it will take for a new deposit to be processed
	time := float64(*pending_validators) / float64((limit_per_day))
	const hoursPerDay = 24
	wholeDays, fractionalDays := math.Modf(time)

	hours := int(fractionalDays * hoursPerDay)

	time_as_days := fmt.Sprintf("%d days and %d hours", int(wholeDays), hours)
	page.NewDepositProcessAfter = time_as_days
	page.ValidatorsPerEpoch = *limit
	page.ValidatorsPerDay = limit_per_day
}
