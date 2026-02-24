package system

type Location struct {
	ID    string
	Name  string
	Order int

	Description string

	// ID действий из ActionRegistry
	AvailableActions []string
}

var LocationRegistry = map[string]*Location{}

func InitLocations(gs *GameState) {
	gs.Locations["dust_town"] = &Location{
		ID:          "dust_town",
		Name:        "Пыльный город",
		Order:       1,
		Description: "Старый пыльный городок.",
		AvailableActions: []string{
			"ranch",
			"saloon",
			"train_accuracy",
			"buy_ammo",
			"sell_ammo",
			"first_duel",
		},
	}

	gs.Locations["frontier"] = &Location{
		ID:               "frontier",
		Name:             "Фронтир",
		Order:            2,
		Description:      "Линия фронта на гражданской войне.",
		AvailableActions: []string{},
	}

	// Устанавливаем текущую локацию, если ещё нет
	if gs.CurrentLocation == nil {
		gs.CurrentLocation = gs.Locations["dust_town"]
	}
}
