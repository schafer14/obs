#!/bin/bash

curl -d '{ 
  "observationLocation": {
    "type": "Point",
    "coordinates": [
      149.12558555603027,
      -35.30694800929574
    ]
  },

  "feature": {"id": "https://example.com/banners-garden", "label": "Banners Garden", "reference": "https://example.com/banners-garden" } ,

  "featureType": {"id": "urn:example:garden", "label": "Garden" }, 

  "property": {"id": "urn:example:garden-health", "label": "Garden Health" }, 

  "propertyType": {"id": "urn:example:scale-1-5", "label": "Scale 1 - 5" }, 

  "process": {"id": "urn:example:by-eye", "label": "Measured by eye" }, 

  "result": {
    "wisteria": 5,
    "magnolia": 4,
    "citrus": 4
  }

}' \
  localhost:3000/v1/observations  | jq
