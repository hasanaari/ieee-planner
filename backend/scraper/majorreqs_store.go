package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

type MajorRequirementsStore struct {
	RequirementsByMajor map[string]*MajorRequirements
	DataPath            string
}

func NewMajorRequirementsStore(dataPath string) (*MajorRequirementsStore, error) {
	store := &MajorRequirementsStore{
		RequirementsByMajor: make(map[string]*MajorRequirements),
		DataPath:            dataPath,
	}

	err := store.LoadAllMajorRequirements()
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (mrs *MajorRequirementsStore) LoadAllMajorRequirements() error {
	files, err := filepath.Glob(filepath.Join(mrs.DataPath, "*.json"))
	if err != nil {
		return err
	}

	for _, file := range files {
		reqs, err := ReadMajorreqsFromJSON(file)
		if err != nil {
			return fmt.Errorf("failed to load %s: %w", file, err)
		}

		major := strings.ToLower(reqs.Major)
		mrs.RequirementsByMajor[major] = reqs
	}

	return nil
}

func (mrs *MajorRequirementsStore) GetRequirements(major string) (*MajorRequirements, bool) {
	reqs, ok := mrs.RequirementsByMajor[strings.ToLower(major)]
	return reqs, ok
}

// NOTE: this one returns a complicated reqs list and this may need merging and/or filtering on the client
func GetMajorRequirementsHandler(store *MajorRequirementsStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		major := r.URL.Query().Get("major")
		if major == "" {
			http.Error(w, "Major parameter is required", http.StatusBadRequest)
			return
		}

		reqs, found := store.GetRequirements(major)
		if !found {
			http.Error(w, "Major not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reqs)
	}
}

func GetAvailableMajorsHandler(store *MajorRequirementsStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		majors := make([]string, 0, len(store.RequirementsByMajor))
		for major := range store.RequirementsByMajor {
			majors = append(majors, major)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(majors)
	}
}
