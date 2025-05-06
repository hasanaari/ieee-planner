package scraper

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
)

type CoursesStore struct {
	CoursesByQuarter map[int][]*Course
	Quarters         []int
	DataPath         string
}

func NewCoursesStore(dataPath string) (*CoursesStore, error) {
	store := &CoursesStore{
		CoursesByQuarter: make(map[int][]*Course),
		DataPath:         dataPath,
	}

	err := store.LoadAllCourseFiles()
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (cs *CoursesStore) LoadAllCourseFiles() error {
	files, err := filepath.Glob(filepath.Join(cs.DataPath, "*.json"))
	if err != nil {
		return err
	}

	for _, file := range files {
		courses, err := ReadCourseDataFromJSON(file)
		if err != nil {
			return err
		}

		for _, course := range courses {
			if course.Quarter > 0 {
				cs.CoursesByQuarter[course.Quarter] = append(cs.CoursesByQuarter[course.Quarter], course)
			}
		}
	}

	cs.UpdateQuartersList()
	return nil
}

func (cs *CoursesStore) UpdateQuartersList() {
	cs.Quarters = make([]int, 0, len(cs.CoursesByQuarter))
	for quarter := range cs.CoursesByQuarter {
		cs.Quarters = append(cs.Quarters, quarter)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(cs.Quarters)))
}

func (cs *CoursesStore) GetCoursesByQuarter(quarter int) []*Course {
	return cs.CoursesByQuarter[quarter]
}

func (cs *CoursesStore) GetAvailableQuarters() []int {
	return cs.Quarters
}

func GetCoursesByQuarterHandler(store *CoursesStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		quarterStr := r.URL.Query().Get("quarter")
		if quarterStr == "" {
			http.Error(w, "Quarter parameter is required", http.StatusBadRequest)
			return
		}

		quarter, err := strconv.Atoi(quarterStr)
		if err != nil {
			http.Error(w, "Invalid quarter format", http.StatusBadRequest)
			return
		}

		courses := store.GetCoursesByQuarter(quarter)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(courses)
	}
}

func GetAvailableQuartersHandler(store *CoursesStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		quarters := store.GetAvailableQuarters()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(quarters)
	}
}
