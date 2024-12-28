package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

var (
	tpl *template.Template
	fetchError bool
)

const (
	apiURL            = "https://groupietrackers.herokuapp.com/api"
	artistsEndpoint   = apiURL + "/artists"
	locationsEndpoint = apiURL + "/locations"
	datesEndpoint     = apiURL + "/dates"
	relationsEndpoint = apiURL + "/relation"
)


// Artist data structure remains the same
type Artist struct {
	ID            int      `json:"id"`
	Image         string   `json:"image"`
	Name          string   `json:"name"`
	Members       []string `json:"members"`
	CreationDate  int      `json:"creationDate"`
	FirstAlbum    string   `json:"firstAlbum"`
	LocationCount int
}

// Artist data structure remains the same
type ArtistPageData struct {
	ID            int      `json:"id"`
	Image         string   `json:"image"`
	Name          string   `json:"name"`
	Members       []string `json:"members"`
	CreationDate  int      `json:"creationDate"`
	FirstAlbum    string   `json:"firstAlbum"`
	Locations     []string
	Dates         []string
	Relations     []string
	LocationCount int
}

type LocationAPIResponse struct {
	Index []Location `json:"index"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type DatesAPIResponse struct {
	Index []Dates `json:"index"`
}

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

/* type Relation struct {
	Index []struct {
		ID        int      `json:"id"`
		Dates     []string `json:"dates"`
		Locations []string `json:"locations"`
	} `json:"index"`
} */

type Relation struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

type ErrorResponse struct {
	Code    int
	Message string
}

var (
	artists    []Artist
	locations  []Location
	dates      []Dates
	relations  Relation
)

func LocationCount(locations []string) int {
	LocationList := make(map[string]string)
	for _, loc := range locations {
		split := strings.Split(loc, "-")
		if len(split) == 2 {
			LocationList[split[1]] = split[0]
		}
	}
	return len(LocationList)
}

func ToUpper(s string) string {
	str := strings.Fields(s)
	for i:= range str {
		str[i] = strings.ToUpper(string(str[i][0])) + string(str[i][1:])
	}
	return strings.Join(str, " ")
}

func FetchArtists() error {
	resp, err := http.Get(artistsEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &artists); err != nil {
		return err
	}
	log.Println("Artists fetched successfully")
	return nil
}

func FetchLocations() error {
	resp, err := http.Get(locationsEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiResponse LocationAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return err
	}
	locations = apiResponse.Index
	log.Println("Locations fetched successfully")
	return nil
}

func FetchDates() error {
	resp, err := http.Get(datesEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Use the wrapper struct to unmarshal the data
	var apiResponse DatesAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return err
	}

	// Assign the fetched dates
	dates = apiResponse.Index
	log.Println("Dates fetched successfully")
	return nil
}

func FetchRelations() error {
	resp, err := http.Get(relationsEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &relations); err != nil {
		return err
	}
	log.Println("Relations fetched successfully")
	return nil
}

// Fetch All Data at Once
func FetchAllData() {
	if len(artists) == 0 {
		if err := FetchArtists(); err != nil {
			log.Fatal("Error fetching artists:", err)
		}
		if err := FetchLocations(); err != nil {
			log.Fatal("Error fetching locations:", err)
		}
		if err := FetchDates(); err != nil {
			log.Fatal("Error fetching dates:", err)
		}
		if err := FetchRelations(); err != nil {
			log.Fatal("Error fetching relations:", err)
		}
	}
}

func FetchArtistData(id int) (Artist, []string, []string, map[string][]string, error) {
	// Fetch the artist by ID
	var selectedArtist Artist
	for _, artist := range artists {
		if artist.ID == id {
			selectedArtist = artist
			break
		}
	}

	// Fetch associated data (locations, dates, and relations)
	var associatedLocations []string
	var associatedDates []string
	var associatedRelations map[string][]string
	var locationCount int

	for _, location := range locations {
		if location.ID == id {
			for _, loc := range location.Locations {
				cleanLoc := strings.ReplaceAll(loc, "_", " ")
				cleanLoc = strings.ReplaceAll(cleanLoc, "-", " ")
				associatedLocations = append(associatedLocations, cleanLoc)
			}
			locationCount = LocationCount(associatedLocations)
		}
	}
	selectedArtist.LocationCount = locationCount

	for _, date := range dates {
		if date.ID == id {
			for _, d := range date.Dates {
				cleanDate := strings.ReplaceAll(d, "*", "")
				associatedDates = append(associatedDates, cleanDate)
			}
		}
	}

	for _, relation := range relations.Index {
		if relation.ID == id {
			associatedRelations = make(map[string][]string)
			for loc, relDates := range relation.DatesLocations {
				cleanLoc := strings.ReplaceAll(loc, "_", " ")
				cleanLoc = strings.ReplaceAll(cleanLoc, "-", " ")
				cleanRelDates := []string{}
				for _, d := range relDates {
					cleanDate := strings.ReplaceAll(d, "*", "")
					cleanRelDates = append(cleanRelDates, cleanDate)
				}
				associatedRelations[cleanLoc] = cleanRelDates
			}
		}
	}

	return selectedArtist, associatedLocations, associatedDates, associatedRelations, nil
}
