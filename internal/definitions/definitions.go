package definitions

var Data = map[string]FeatureType{
	"people": FeatureType{
		ID:          "34edda82-0f22-4115-b5cf-406db1330436",
		Name:        "Person",
		Slug:        "people",
		Description: "A person that may be observed",
		Properties: map[string]Property{
			"optimism": Property{
				ID:          "98d1e62b-14c5-487e-bee5-81348edede77",
				Name:        "Optimism",
				Slug:        "optimism",
				Description: "A persons optimism",
				Category:    "optimism",
				PropertyTypes: map[string]PropertyType{
					"learned-optimism":     PropertyType{},
					"learned-optimism-raw": PropertyType{},
					"cave":                 PropertyType{},
				},
			},
			"depression": Property{
				ID:          "b1f141f6-d957-4006-9b3d-0cc1c883fffe",
				Name:        "Depression",
				Slug:        "depression",
				Description: "A persons depression",
				Category:    "optimism",
				PropertyTypes: map[string]PropertyType{
					"learned-optimism":     PropertyType{},
					"learned-optimism-raw": PropertyType{},
				},
			},
			"goal": Property{
				ID:          "4b46d2af-e908-4643-8060-3c85f991a8bf",
				Name:        "Goal",
				Slug:        "goal",
				Description: "A personal goal",
				Category:    "future",
				PropertyTypes: map[string]PropertyType{
					"textual":   PropertyType{},
					"goal-plan": PropertyType{},
				},
			},
			// "motivation": Property{},
			// "software-aptitude": Property{},
			"personality": Property{
				ID:          "8955d4f2-6968-4548-8ee6-6ae3501b9afe",
				Name:        "Personality",
				Slug:        "personality",
				Description: "A persons personality",
				Category:    "personality",
				PropertyTypes: map[string]PropertyType{
					"sixteen-and-me": PropertyType{},
					"myers-briggs":   PropertyType{},
				},
			},
		},
	},
	"groups": FeatureType{
		ID:          "012c7b88-d55c-4309-98b1-f009f5608f2d",
		Name:        "Group",
		Slug:        "group",
		Description: "A group of people",
		Properties: map[string]Property{
			"safety": Property{
				ID:            "efc2378a-e325-47af-b4ca-f4ffa3d52afe",
				Name:          "Safety",
				Slug:          "safety",
				Description:   "How safe members of the group feel (from the book Culture Code)",
				PropertyTypes: map[string]PropertyType{},
			},
			"belonging": Property{
				ID:            "57f62e37-77dc-4712-ab17-156f31d1ea5e",
				Name:          "Belonging",
				Slug:          "belonging",
				Description:   "The sense of belonging the group members feel to the group (from the book Culture Code)",
				PropertyTypes: map[string]PropertyType{},
			},
		},
	},
}

type FeatureType struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Description string              `json:"description"`
	Properties  map[string]Property `json:"properties"`
}

type Property struct {
	ID            string                  `json:"id"`
	Name          string                  `json:"name"`
	Slug          string                  `json:"slug"`
	Description   string                  `json:"description"`
	Category      string                  `json:"category"`
	PropertyTypes map[string]PropertyType `json:"propertyTypes"`
}

type PropertyType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     int    `json:"version"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}
