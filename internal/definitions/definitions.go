package definitions

import (
	"gopkg.in/go-playground/validator.v9"
)

// validate holds the settings and caches for validating request struct values.
var validate = validator.New()

var Data = map[string]FeatureType{
	"people": FeatureType{
		ID:          "34edda82-0f22-4115-b5cf-406db1330436",
		Name:        "Person",
		Slug:        "people",
		Description: "A person that may be observed",
		Properties: map[string]Property{
			"profession": Property{
				ID:          "3eac25e8-a83b-4cf3-a3b0-f9f347c0c1c7",
				Name:        "Profession",
				Slug:        "profession",
				Description: "Observations about persons career, skills and profession.",
				Category:    "optimism",
				PropertyTypes: map[string]PropertyType{
					"role": PropertyType{},
				},
			},
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
					"textual": PropertyType{
						ID:          "717988a9-f139-4875-b7d7-ac132d7df75b",
						Name:        "Textual Goals",
						Slug:        "textual",
						Description: "Free form description of a goal.",
						Schema: map[string]interface{}{
							"$schema":              "http://json-schema.org/draft-07/schema",
							"$id":                  "https://linked-data-land.appspot.com/v1/definitions/people/goal/textual",
							"additionalProperties": false,
							"type":                 "object",
							"title":                "Text Goal",
							"description":          "An single goal written as a textual observation",
							"required": []string{
								"goal",
							},
							"properties": map[string]interface{}{
								"goal": map[string]interface{}{
									"$id":         "#/properties/goal",
									"type":        "string",
									"title":       "The goal",
									"description": "The textual goal.",
									"examples": []string{
										"I want to write build an awesome observation platform.",
									},
								},
							},
						},
					},
					"structured": PropertyType{
						ID:          "9668ea74-9a7d-4665-8e03-402fb39c1365",
						Name:        "Structured Goal",
						Slug:        "structured",
						Description: "A goal with a series of steps to achieve the goal",
						Schema: map[string]interface{}{
							"$schema":              "http://json-schema.org/draft-07/schema",
							"additionalProperties": false,
							"$id":                  "https://linked-data-land.appspot.com/v1/definitions/people/goal/structured",
							"type":                 "object",
							"title":                "Structured Goal",
							"description":          "A structured goal is a goal with a plan of the steps required to accomplish the goal.",
							"required": []string{
								"goal",
								"steps",
							},
							"properties": map[string]interface{}{
								"goal": map[string]interface{}{
									"$id":         "#/properties/goal",
									"type":        "string",
									"title":       "The goal",
									"description": "The textual goal.",
								},
								"steps": map[string]interface{}{
									"$id":  "#/properties/goal",
									"type": "array",
									"items": map[string]interface{}{
										"type":     "object",
										"required": []string{"description"},
										"properties": map[string]interface{}{
											"step": map[string]interface{}{
												"$id":         "#/properties/goal",
												"type":        "string",
												"title":       "step",
												"description": "A step of the process",
											},
											"by": map[string]interface{}{
												"type":        "string",
												"format":      "date-time",
												"title":       "Due Date",
												"description": "The date you hope to accomplish this step by.",
											},
										},
									},
									"title":       "Steps",
									"description": "The steps required to accomplish this goal.",
								},
							},
							"examples": []interface{}{
								map[string]interface{}{
									"goal": "I want to write build an awesome observation platform.",
									"steps": []interface{}{
										map[string]interface{}{
											"description": "build the platform",
											"by":          "2020-03-29T00:00:00+00:00",
										},
										map[string]interface{}{
											"description": "build a cool implementation using the platform.",
											"by":          "2020-04-10T00:00:00+00:00",
										},
										map[string]interface{}{
											"description": "help other people build implementations using the platform.",
											"by":          "2020-05-10T00:00:00+00:00",
										},
									},
								},
							},
						},
					},
					"daily-goals": PropertyType{
						ID:          "a41e13a9-d4d6-47f5-b3d6-580a0e0e44d1",
						Name:        "Daily Goals",
						Slug:        "daily-goal",
						Description: "A set of goals to accomplish for a given time.",
						Schema: map[string]interface{}{
							"$schema":              "http://json-schema.org/draft-07/schema",
							"additionalProperties": false,
							"$id":                  "https://linked-data-land.appspot.com/v1/definitions/people/goal/daily-goals",
							"type":                 "object",
							"title":                "Daily goals",
							"description":          "A list of days to accomplish in a day",
							"required": []string{
								"goals",
								"day",
							},
							"properties": map[string]interface{}{
								"goals": map[string]interface{}{
									"type": "array",
									"items": map[string]interface{}{
										"title":       "Goal",
										"description": "One of the goals to accomplish",
										"type":        "string",
									},
								},
								"day": map[string]interface{}{
									"type":        "string",
									"format":      "date",
									"title":       "day",
									"description": "the day these goals are relevant",
								},
							},
							"examples": []interface{}{
								map[string]interface{}{
									"goals": []string{
										"Read a chapter of my book.",
										"Update how validation of observations work.",
										"Vacuum and tidy the house.",
									},
									"day": "2020-03-23",
								},
							},
						},
					},
					"daily-goals-result": PropertyType{
						ID:          "71d6330c-0f02-4ee9-85d5-b27dfa45aab7",
						Name:        "Daily Goals Result",
						Slug:        "daily-goal-result",
						Description: "A result of a daily goal.",
						Schema: map[string]interface{}{
							"$schema":              "http://json-schema.org/draft-07/schema",
							"additionalProperties": false,
							"$id":                  "https://linked-data-land.appspot.com/v1/definitions/people/goal/daily-goals-result",
							"type":                 "object",
							"title":                "Daily goals result",
							"description":          "A review of the days goals and if they were accomplished",
							"required": []string{
								"goals",
								"day",
							},
							"properties": map[string]interface{}{
								"goals": map[string]interface{}{
									"type": "array",
									"items": map[string]interface{}{
										"type":     "object",
										"required": []string{"goal", "accomplished"},
										"properties": map[string]interface{}{
											"goal": map[string]interface{}{
												"title":       "goal",
												"description": "A description of the goal.",
												"type":        "string",
											},
											"accomplished": map[string]interface{}{
												"title":       "Accomplished",
												"description": "Whether teh goal was accomplished or not",
												"type":        "boolean",
											},
										},
									},
								},
								"day": map[string]interface{}{
									"type":        "string",
									"format":      "date",
									"title":       "day",
									"description": "the day these goals are relevant",
								},
							},
							"examples": []interface{}{
								map[string]interface{}{
									"goals": []interface{}{
										map[string]interface{}{
											"goal":         "Read a chapter of my book",
											"accomplished": false,
										},
										map[string]interface{}{
											"goal":         "Update how validation of observations work.",
											"accomplished": true,
										},
										map[string]interface{}{
											"goal":         "Vacuum and tidy the house.",
											"accomplished": true,
										},
									},
									"day": "2020-03-23",
								},
							},
						},
					},
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
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Version     int                    `json:"version"`
	Slug        string                 `json:"slug"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
	SchemaURL   string                 `json:"schemaUrl,omitempty"`
}
