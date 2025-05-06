package scraper

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func GetCourseKey(c Course) string {
	parts := strings.Split(c.Number, "-")
	if len(parts) > 1 {
		return strings.ToUpper(c.Subject + " " + parts[0] + "-" + parts[1])
	} else {
		return strings.ToUpper(c.Subject + " " + parts[0])
	}
}

type CourseBySubject struct {
	Title    string `json:"title"`
	Number   string `json:"number"`
	Topic    string `json:"topic"`
	Overview string `json:"overview"`
	Quarters []int  `json:"quarters"`
}

type CoursesStore struct {
	CoursesByQuarter map[int][]*Course
	CoursesBySubject map[string]map[string]*CourseBySubject
	CoursesByKey     map[string][]*Course

	Quarters []int
	DataPath string
}

func NewCoursesStore(dataPath string) (*CoursesStore, error) {
	store := &CoursesStore{
		CoursesByQuarter: make(map[int][]*Course),
		CoursesBySubject: make(map[string]map[string]*CourseBySubject),
		CoursesByKey:     make(map[string][]*Course),
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
			mkey := GetCourseKey(*course)

			_, exists := cs.CoursesByKey[mkey]
			cs.CoursesByKey[mkey] = append(cs.CoursesByKey[mkey], course)
			if exists {
				if subjectMap, ok := cs.CoursesBySubject[course.Subject]; ok {
					if existingcoursebysubject, ok := subjectMap[mkey]; ok {
						iambetter := true
						iexist := false
						for _, quarter := range existingcoursebysubject.Quarters {
							if quarter > course.Quarter {
								iambetter = false
								break
							}
							if quarter == course.Quarter {
								iexist = true
								break
							}
						}
						if !iexist {
							if iambetter {
								existingcoursebysubject.Overview = course.Overview
							}
							existingcoursebysubject.Quarters = append(existingcoursebysubject.Quarters, course.Quarter)
						}
					}
				}
			} else {
				nq := make([]int, 0)
				nq = append(nq, course.Quarter)

				if _, exists := cs.CoursesBySubject[course.Subject]; !exists {
					cs.CoursesBySubject[course.Subject] = make(map[string]*CourseBySubject)
				}

				cs.CoursesBySubject[course.Subject][mkey] = &CourseBySubject{
					Title:    course.Title,
					Number:   course.Number,
					Topic:    course.Topic,
					Overview: course.Overview,
					Quarters: nq,
				}
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

func (cs *CoursesStore) GetCoursesByKey(key string) []*Course {
	return cs.CoursesByKey[key]
}

func (cs *CoursesStore) GetCoursesByQuarter(quarter int) []*Course {
	return cs.CoursesByQuarter[quarter]
}

func (cs *CoursesStore) GetCoursesBySubject(subject string) []*CourseBySubject {
	coursesbysubject := make([]*CourseBySubject, 0)

	for _, coursebysubject := range cs.CoursesBySubject[subject] {
		coursesbysubject = append(coursesbysubject, coursebysubject)
	}

	return coursesbysubject
}

func (cs *CoursesStore) GetAvailableQuarters() []int {
	return cs.Quarters
}

func GetCoursesByKeyHandler(store *CoursesStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyStr := r.URL.Query().Get("key")
		if keyStr == "" {
			http.Error(w, "Key parameter is required", http.StatusBadRequest)
			return
		}

		courses := store.GetCoursesByKey(keyStr)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(courses)
	}
}

func GetCoursesBySubjectHandler(store *CoursesStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subjectStr := r.URL.Query().Get("subject")
		if subjectStr == "" {
			http.Error(w, "Subject parameter is required", http.StatusBadRequest)
			return
		}

		courses := store.GetCoursesBySubject(subjectStr)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(courses)
	}
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
