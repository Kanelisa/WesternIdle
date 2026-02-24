package system

import "math/rand"

// ---------------- Active Action ----------------

type ActiveAction struct {
	ID       string
	Progress float64
	Duration float64
}

// ---------------- Action Info ----------------

type ActionInfo struct {
	ID          string
	Name        string
	Description string
	Category    string

	Duration float64
	Cost     map[Resource]float64
	Reward   map[Resource]float64
	Instant  bool

	// Кастомная логика (опционально)
	Resolver func(*GameState)
}

// ---------------- Action Registry ----------------

var ActionRegistry = map[string]ActionInfo{

	// -------- Долгие --------

	"ranch": {
		ID:          "ranch",
		Name:        "Работа на ранчо",
		Description: "Тяжёлая работа, но стабильный доход.",
		Category:    "long",
		Duration:    3,
		Reward: map[Resource]float64{
			ResourceMoney: 10,
		},
	},

	"saloon": {
		ID:          "saloon",
		Name:        "Играть в покер",
		Description: "Попробуйте удачу. Можно выиграть или проиграть.",
		Category:    "long",
		Duration:    4,
		Cost: map[Resource]float64{
			ResourceMoney: 10,
		},
		Resolver: func(s *GameState) {
			if rand.Float64() < 0.5 {
				s.AddResource(ResourceMoney, 20)
				s.LogEvent("Вы выиграли $20!")
			} else {
				s.LogEvent("Вы проиграли $10...")
			}
		},
	},

	"train_accuracy": {
		ID:          "train_accuracy",
		Name:        "Тренировка меткости",
		Description: "Повышает меткость стрелка.",
		Category:    "long",
		Duration:    5,
		Reward: map[Resource]float64{
			ResourceAccuracy: 1,
		},
	},

	//-----------Мгновенные-------------

	"buy_ammo": {
		ID:          "buy_ammo",
		Name:        "Купить патроны",
		Description: "За $10 получите 6 патронов.",
		Category:    "instant",
		Instant:     true,
		Cost: map[Resource]float64{
			ResourceMoney: 10,
		},
		Reward: map[Resource]float64{
			ResourceAmmo: 6,
		},
	},

	"sell_ammo": {
		ID:          "sell_ammo",
		Name:        "Продать патроны",
		Description: "Продать 6 патронов за $8.",
		Category:    "instant",
		Instant:     true,
		Cost: map[Resource]float64{
			ResourceAmmo: 6,
		},
		Reward: map[Resource]float64{
			ResourceMoney: 8,
		},
	},

	//----------Дуэли-------------

	"first_duel": {
		ID:          "first_duel",
		Name:        "Первая дуэль",
		Description: "dull",
		Category:    "duel",
	},
}

// ---------------- Получение информации ----------------

func GetActionInfo(id string) (ActionInfo, bool) {
	info, ok := ActionRegistry[id]
	return info, ok
}
