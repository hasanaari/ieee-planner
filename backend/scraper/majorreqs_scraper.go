package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"context"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	openai "github.com/sashabaranov/go-openai"
)

type Requirements interface {
	GetName() string
	GetNumRequirements() int
	GetOptions() []Option
	IsTheme() bool
	IsUnrestricted() bool
}

type Req interface {
	GetType() int
}

// very simple scheme thus it is redundant
type Requirement struct {
	Courses []string `json:"courses"`
}

type Option struct {
	Between []Requirement `json:"between"`
}

type GenericRequirements struct {
	RequirementType int      `json:"requirementType"`
	Name            string   `json:"name"`
	Requirements    []Option `json:"requirements"`
}

// these need better discrimination?
type ThemeRequirements struct {
	RequirementType int `json:"requirementType"`
	NumRequirements int `json:"numreqs"`
}
type UnrestrictedRequirements struct {
	RequirementType int `json:"requirementType"`
	NumRequirements int `json:"numreqs"`
}
type UnknownRequirements struct {
	RequirementType int `json:"requirementType"`
	NumRequirements int `json:"numreqs"`
}

type MajorRequirements struct {
	IsEngineering   bool   `json:"isEngineering"`
	Major           string `json:"major"`
	AllRequirements []any  `json:"allreqs"`
}

var CORE_ENGINEERING_REQUIREMENTS = MajorRequirements{
	Major:         "Core Engineering",
	IsEngineering: true,
	AllRequirements: []any{
		GenericRequirements{
			Name:            "Mathematics",
			RequirementType: 0,
			Requirements: []Option{
				{Between: []Requirement{Requirement{[]string{"MATH 220-1"}}}},
				{Between: []Requirement{Requirement{[]string{"MATH 220-2"}}}},
				{Between: []Requirement{Requirement{[]string{"MATH 228-1"}}}},
				{Between: []Requirement{Requirement{[]string{"MATH 228-2"}}}},
			},
		},
		GenericRequirements{
			Name:            "Engineering Analysis and Computer Proficiency",
			RequirementType: 0,
			Requirements: []Option{
				{Between: []Requirement{Requirement{[]string{"GEN_ENG 205-1"}}, Requirement{[]string{"GEN_ENG 206-1"}}}},
				{Between: []Requirement{Requirement{[]string{"GEN_ENG 205-2"}}}},
				{Between: []Requirement{Requirement{[]string{"GEN_ENG 205-3"}}}},
				{Between: []Requirement{Requirement{[]string{"GEN_ENG 205-4"}}, Requirement{[]string{"GEN_ENG 206-4"}}}},
			},
		},
		GenericRequirements{
			Name:            "Basic Sciences",
			RequirementType: 0,
			Requirements: []Option{
				// at least one of these options
				{Between: []Requirement{
					Requirement{[]string{"BIOL_SCI 202-0", "BIOL_SCI 232-0"}},
					Requirement{[]string{"CHEM 131-0", "CHEM 141-0"}},
					Requirement{[]string{"CHEM 151-0", "CHEM 161-0"}},
					Requirement{[]string{"CHEM 171-0", "CHEM 181-0"}},
					Requirement{[]string{"PHYSICS 135-2", "PHYSICS 136-2"}},
					Requirement{[]string{"PHYSICS 125-2", "PHYSICS 126-2"}},
					Requirement{[]string{"PHYSICS 140-2", "PHYSICS 136-2"}},
				}},

				//  one of above or below
				{Between: []Requirement{
					// all the options from above
					Requirement{[]string{"BIOL_SCI 202-0", "BIOL_SCI 232-0"}},
					Requirement{[]string{"CHEM 131-0", "CHEM 141-0"}},
					Requirement{[]string{"CHEM 151-0", "CHEM 161-0"}},
					Requirement{[]string{"CHEM 171-0", "CHEM 181-0"}},
					Requirement{[]string{"PHYSICS 135-2", "PHYSICS 136-2"}},
					Requirement{[]string{"PHYSICS 125-2", "PHYSICS 126-2"}},
					Requirement{[]string{"PHYSICS 140-2", "PHYSICS 136-2"}},

					// include all the additional options
					Requirement{[]string{"BIOL_SCI 150-0"}},
					Requirement{[]string{"BIOL_SCI 201-0"}},
					Requirement{[]string{"BIOL_SCI 203-0", "BIOL_SCI 233-0"}},
					Requirement{[]string{"BIOL_SCI 234-0"}},
					Requirement{[]string{"CHEM_ENG 275-0"}},
					Requirement{[]string{"CIV_ENV 202-0"}},
					Requirement{[]string{"CHEM 132-0", "CHEM 142-0"}},
					Requirement{[]string{"CHEM 152-0", "CHEM 162-0"}},
					Requirement{[]string{"CHEM 172-0", "CHEM 182-0"}},
					Requirement{[]string{"CHEM 215-1", "CHEM 235-1"}},
					Requirement{[]string{"CHEM 215-2", "CHEM 235-2"}},
					Requirement{[]string{"ASTRON 220-1"}},
					Requirement{[]string{"ASTRON 220-2"}},
					Requirement{[]string{"CIV_ENV 203-0"}},
					Requirement{[]string{"EARTH 201-0"}},
					Requirement{[]string{"EARTH 202-0"}},
					Requirement{[]string{"EARTH 203-0"}},
					Requirement{[]string{"PHYSICS 135-3", "PHYSICS 136-3"}},
					Requirement{[]string{"PHYSICS 125-3", "PHYSICS 126-3"}},
					Requirement{[]string{"PHYSICS 140-3", "PHYSICS 136-3"}},
					Requirement{[]string{"PHYSICS 239-0"}},
					Requirement{[]string{"COG_SCI 210-0"}},
					Requirement{[]string{"CSD 202-0"}},
					Requirement{[]string{"CSD 303-0"}},
					Requirement{[]string{"PSYCH 221-0"}},
				}},
			},
		},
		GenericRequirements{
			Name:            "Design and Communication",
			RequirementType: 0,
			Requirements: []Option{
				{Between: []Requirement{
					Requirement{[]string{"DSGN 106-1", "DSGN 106-2"}},
				}},

				{Between: []Requirement{
					Requirement{[]string{"ENGLISH 106-1", "ENGLISH 106-2"}},
				}},

				{Between: []Requirement{
					Requirement{[]string{"COMM_ST 102-0"}},
					Requirement{[]string{"PERF_ST 103-0"}},
					Requirement{[]string{"PERF_ST 203-0"}},
					Requirement{[]string{"BMD_ENG 390-2"}},
				}},
			},
		},
		ThemeRequirements{
			RequirementType: 1,
			NumRequirements: 7,
		},
		UnrestrictedRequirements{
			RequirementType: 2,
			NumRequirements: 5,
		},
	},
}

func WriteMajorreqsToJSON(mr MajorRequirements, filePath string) error {
	jsonData, err := json.MarshalIndent(mr, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling majorreqs to json: %w", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing json to file: %w", err)
	}

	fmt.Printf("wrote %s major requirements to %s\n", mr.Major, filePath)
	return nil
}

func ReadMajorreqsFromJSON(filePath string) (*MajorRequirements, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading json file: %w", err)
	}

	var mr *MajorRequirements
	err = json.Unmarshal(jsonData, &mr)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json data: %w", err)
	}

	fmt.Printf("read %s major requirements from %s\n", mr.Major, filePath)
	return mr, nil
}

func ReadMajorreqsFromJSONString(jsonString string) (*MajorRequirements, error) {
	var mr *MajorRequirements
	err := json.Unmarshal([]byte(jsonString), &mr)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json data: %w", err)
	}
	return mr, nil
}

var majorURLs = map[string]string{
	"Computer Engineering": "https://catalogs.northwestern.edu/undergraduate/engineering-applied-science/electrical-computer-engineering/computer-engineering-degree/",
	"Computer Science":     "https://catalogs.northwestern.edu/undergraduate/engineering-applied-science/computer-science/computer-science-degree/",
    "Electrical Engineering": "https://catalogs.northwestern.edu/undergraduate/engineering-applied-science/electrical-computer-engineering/electrical-engineering-degree/",
	"Theatre": "https://catalogs.northwestern.edu/undergraduate/communication/theatre/theatre-major/",
	"Economics": "https://catalogs.northwestern.edu/undergraduate/arts-sciences/economics/economics-major/",
	"Psychology": "https://catalogs.northwestern.edu/undergraduate/arts-sciences/psychology/psychology-major/",
	"Philosophy": "https://catalogs.northwestern.edu/undergraduate/arts-sciences/philosophy/philosophy-major/",
	"Industrial Engineering": "https://catalogs.northwestern.edu/undergraduate/engineering-applied-science/industrial-engineering-management-sciences/industrial-engineering-degree/",
	"Biomedical Engineering": "https://catalogs.northwestern.edu/undergraduate/engineering-applied-science/biomedical-engineering/biomedical-engineering-degree/",
	"Environmental Engineering": "https://catalogs.northwestern.edu/undergraduate/engineering-applied-science/civil-environmental-engineering/environmental-engineering-degree/",
	"Mechanical Engineering": "https://catalogs.northwestern.edu/undergraduate/engineering-applied-science/mechanical-engineering/mechanical-engineering-degree/",
	"Manufacturing and Design Engineering": "https://catalogs.northwestern.edu/undergraduate/engineering-applied-science/segal-design-institute/manufacturing-design-engineering-degree/",
}

func GetMajorreqs(major string) (MajorRequirements, error) {
	if major == "Core Engineering" {
		return CORE_ENGINEERING_REQUIREMENTS, nil
	} else {
		url, ok := majorURLs[major]
		if ok {
			return *ScrapeMajorRequirements(url), nil
		} else {
			return MajorRequirements{}, fmt.Errorf("error finding major url %v", major)
		}
	}
}

func getScrapePrompt(html string) string {
	prompt := `
## instructions
task: parse this HTML course catalog page and convert it into a structured JSON representation following a specific format.
The output should be a single JSON object representing all major requirements. Respond with just the json.

## response format
The JSON structure should represent three types of requirements:
- Type 0: GenericRequirements (for standard course requirements)
- Type 1: ThemeRequirements (for theme electives)
- Type 2: UnrestrictedRequirements (for unrestricted electives)
- Type 3: UnknownRequirements (for unknown requirements/ when you can only parse number of requirements/ has just the number of requirements like 1, 2)

1. The top-level object is a major requirements object with "major", "allreqs", and "isEngineering" fields

2. If the course is an engineering course you will not have detailed information about its core courses just record true for isEngineering and proceed with just parsing major requirements without including theme or unrestricted

3. Each requirement in "allreqs:" must have a "type" field:
   - 0 for GenericRequirements (needs "name" and "requirements" fields)
   - 1 for ThemeRequirements (needs only "numreqs" field)
   - 2 for UnrestrictedRequirements (needs only "numreqs" field)

4. For type 0 (GenericRequirements):
   - Each item in "requirements" is an Option with a "between" field
   - Each item in "between" is a Requirement with a "courses" field
   - "courses" contains an array of course identifiers as strings

5. Handle these specific relationships correctly:
   - Paired courses (with "&") should be in the same "courses" array
   - Alternative courses (with "or") should be in separate objects within the "between" array
   - When all courses must be taken, create separate Option objects for each
   - For "one of above or below" requirements, include all options in the same "between" array

6. For Theme and Unrestricted requirements:
   - Simply include the number of required courses in "numreqs"

7. Be thorough try to place every course you see in a requirements block! Do not mark anything unknown easily!

8. If the major is an engineering major skip theme and unrestricted requirements as those will be in core engineering requirements!

9. Worst comes to worst for courses you cannot place mark remainder of count of courses in a section(still skip engineering requirements if this is an engineering course) as unknown

10. If there is a "Requirement" in the name of a requirement get rid of it

## response example
{
  "major": "Core Engineering", // name of major
  "isEngineering": true,
  "allreqs": [
    {
      "requirementType": 0, // type of requirement(this is generic)
      "name": "Engineering Analysis and Computer Proficiency", // name the requirement with its title
      "requirements": [ // this is an array of all requirements
        {
          "between": [ // this denotes an one of relationship
            {
              "courses": [ // can be a number of courses that need to be taken together
                "GEN_ENG 205-1" // courses are represented like so
              ]
            },
            {
              "courses": [
                "GEN_ENG 206-1"
              ]
            }
          ]
        },
        {
          "between": [
            {
              "courses": [
                "GEN_ENG 205-2"
              ]
            }
          ]
        },
        {
          "between": [
            {
              "courses": [
                "GEN_ENG 205-3"
              ]
            }
          ]
        },
        {
          "between": [
            {
              "courses": [
                "GEN_ENG 205-4"
              ]
            },
            {
              "courses": [
                "GEN_ENG 206-4"
              ]
            }
          ]
        }
      ]
    },
    // More requirements...
    {
      "type": 1,
      "numreqs": 7
    },
    {
      "type": 2,
      "numreqs": 5
    }
  ]
}

##### HTML:  ` + html
	return prompt
}

func ScrapeMajorRequirements(url string) *MajorRequirements {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var mr *MajorRequirements

	doc.Find("div#textcontainer.page_content").Each(func(i int, s *goquery.Selection) {
		html, err := s.Html()
		if err == nil {
			key := os.Getenv("OPENAI_API_KEY")
			client := openai.NewClient(key)
			resp, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: openai.GPT4oLatest, // use versioned
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: getScrapePrompt(html),
						},
					},
				},
			)

			if err != nil {
				log.Fatalf("ChatCompletion error: %v\n", err)
			}

			mr, err = ReadMajorreqsFromJSONString(cleanJSONblock(resp.Choices[0].Message.Content))
			if err != nil {
				log.Fatalf("Failed to parse ChatCompletion output %v\n", err)
			}
		} else {
			log.Fatalf("Error getting HTML: %v", err)
		}
	})

	return mr
}

func cleanJSONblock(s string) string {
	s = strings.TrimPrefix(s, "```json\n")
	s = strings.TrimPrefix(s, "```\n")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
