package curriculum

type presetSubject struct {
	Code string
	Name string
	Type string
}

type presetGroup struct {
	Label        string
	MinSelect    int32
	MaxSelect    int32
	SubjectCodes []string
}

type presetLevel struct {
	GradeNumber int
	LabelSuffix string // appended to "Grade N" — used for the 5 A/L streams
	Groups      []presetGroup
}

// PresetSummary reports what the run did (or, for a dry run, would do) — a per-grade preview before anything is actually written.
type PresetSummary struct {
	DryRun          bool                 `json:"dry_run"`
	SubjectsCreated int                  `json:"subjects_created"`
	LevelsCreated   int                  `json:"levels_created"`
	GroupsCreated   int                  `json:"groups_created"`
	LinksCreated    int                  `json:"links_created"`
	GradesCovered   []int                `json:"grades_covered"` // grade numbers this preset has content for, present in this school
	GradesSkipped   []int                `json:"grades_skipped"` // grade numbers with no matching Grade row
	Levels          []PresetLevelPreview `json:"levels"`
}

// PresetLevelPreview is one row of the per-level breakdown — what a level
// is called, which grade it's under, and how much of it is new.
type PresetLevelPreview struct {
	Label          string `json:"label"`
	GradeNumber    int    `json:"grade_number"`
	AlreadyExists  bool   `json:"already_exists"`
	GroupsToAdd    int    `json:"groups_to_add"`
	SubjectsToLink int    `json:"subjects_to_link"`
}

var presetSubjects = []presetSubject{
	// Primary (Grades 1-5)
	{"PRI-01", "First Language (Sinhala/Tamil)", "compulsory"},
	{"PRI-02", "English", "compulsory"},
	{"PRI-03", "Mathematics", "compulsory"},
	{"PRI-04", "Environmental Related Activities", "compulsory"},
	{"PRI-05", "Religion", "compulsory"},
	{"PRI-06", "Aesthetic Education (Art/Music/Dance)", "compulsory"},
	{"PRI-07", "Health & Physical Education", "compulsory"},

	// Junior Secondary (Grades 6-9)
	{"GEN-L01", "First Language (Sinhala)", "compulsory"},
	{"GEN-L03", "Second National Language", "compulsory"},
	{"GEN-L05", "English", "compulsory"},
	{"GEN-M01", "Mathematics", "compulsory"},
	{"GEN-S01", "Science", "compulsory"},
	{"GEN-SS01", "History", "compulsory"},
	{"GEN-SS02", "Geography", "compulsory"},
	{"GEN-SS03", "Civic Education", "compulsory"},
	{"GEN-P01", "Health & Physical Education", "compulsory"},
	{"GEN-AE01", "Art", "compulsory"},
	{"GEN-AE02", "Music", "compulsory"},
	{"GEN-AE03", "Dancing", "compulsory"},
	{"GEN-T01", "Practical & Technical Skills", "compulsory"},
	{"GEN-T02", "ICT", "compulsory"},
	{"GEN-R01", "Religion (Buddhism)", "religion"},
	{"GEN-R02", "Religion (Hinduism)", "religion"},
	{"GEN-R03", "Religion (Islam)", "religion"},
	{"GEN-R04", "Religion (Christianity)", "religion"},

	// O/L (Grades 10-11) mandatory + religion
	{"OL-M01", "First Language", "compulsory"},
	{"OL-M03", "English", "compulsory"},
	{"OL-M04", "Mathematics", "compulsory"},
	{"OL-M05", "Science", "compulsory"},
	{"OL-M06", "History", "compulsory"},
	{"OL-M07", "Religion (Buddhism)", "religion"},
	{"OL-M08", "Religion (Hinduism)", "religion"},
	{"OL-M09", "Religion (Islam)", "religion"},
	{"OL-M10", "Religion (Christianity)", "religion"},
	{"OL-M11", "Religion (Catholicism)", "religion"},

	// O/L Basket 1: Academic & Languages
	{"OL-B1-01", "Business & Accounting Studies", "basket"},
	{"OL-B1-02", "Geography", "basket"},
	{"OL-B1-03", "Civic Education", "basket"},
	{"OL-B1-04", "Entrepreneurship Studies", "basket"},
	{"OL-B1-L01", "Sinhala Second Language", "basket"},
	{"OL-B1-L02", "Tamil Second Language", "basket"},
	{"OL-B1-L03", "French", "basket"},
	{"OL-B1-L04", "German", "basket"},
	{"OL-B1-L05", "Hindi", "basket"},
	{"OL-B1-L06", "Japanese", "basket"},
	{"OL-B1-L07", "Arabic", "basket"},
	{"OL-B1-L08", "Korean", "basket"},
	{"OL-B1-L09", "Chinese", "basket"},
	{"OL-B1-L10", "Russian", "basket"},
	{"OL-B1-L11", "Pali", "basket"},
	{"OL-B1-L12", "Sanskrit", "basket"},

	// O/L Basket 2: Arts & Literature
	{"OL-B2-01", "Music (Oriental)", "basket"},
	{"OL-B2-02", "Music (Western)", "basket"},
	{"OL-B2-03", "Music (Carnatic)", "basket"},
	{"OL-B2-04", "Dancing (Indigenous)", "basket"},
	{"OL-B2-05", "Dancing (Bharatha)", "basket"},
	{"OL-B2-06", "Art", "basket"},
	{"OL-B2-07", "Drama & Theatre", "basket"},
	{"OL-B2-08", "English Literary Texts", "basket"},
	{"OL-B2-09", "Sinhala Literary Texts", "basket"},
	{"OL-B2-10", "Tamil Literary Texts", "basket"},
	{"OL-B2-11", "Arabic Literary Texts", "basket"},

	// O/L Basket 3: Technology & Practical
	{"OL-B3-01", "ICT", "basket"},
	{"OL-B3-02", "Agriculture & Food Technology", "basket"},
	{"OL-B3-03", "Aquatic Bio-resources Technology", "basket"},
	{"OL-B3-04", "Design & Construction Technology", "basket"},
	{"OL-B3-05", "Design & Mechanical Technology", "basket"},
	{"OL-B3-06", "Electrical & Electronic Technology", "basket"},
	{"OL-B3-07", "Arts & Crafts", "basket"},
	{"OL-B3-08", "Home Economics", "basket"},
	{"OL-B3-09", "Health & Physical Education", "basket"},
	{"OL-B3-10", "Communication & Media Studies", "basket"},
	{"OL-B3-11", "Electronic Writing & Shorthand", "basket"},

	// A/L common
	{"AL-12", "Common General Test", "compulsory"},
	{"AL-13", "General English", "compulsory"},

	// A/L stream subjects
	{"AL-01", "Physics", "stream"},
	{"AL-02", "Chemistry", "stream"},
	{"AL-10", "Combined Mathematics", "stream"},
	{"AL-20", "ICT", "stream"},
	{"AL-09", "Biology", "stream"},
	{"AL-08", "Agricultural Science", "stream"},
	{"AL-65", "Engineering Technology", "stream"},
	{"AL-66", "Biosystems Technology", "stream"},
	{"AL-67", "Science for Technology", "stream"},
	{"AL-33", "Accountancy", "stream"},
	{"AL-32", "Business Studies", "stream"},
	{"AL-21", "Economics", "stream"},
	{"AL-22", "Geography", "stream"},
	{"AL-23", "Political Science", "stream"},
	{"AL-24", "Logic", "stream"},
	{"AL-25", "History", "stream"},
	{"AL-28", "Home Economics", "stream"},
	{"AL-29", "Communication & Media", "stream"},
}

var primaryCompulsory = []string{"PRI-01", "PRI-02", "PRI-03", "PRI-04", "PRI-05", "PRI-06", "PRI-07"}
var juniorCompulsory = []string{
	"GEN-L01", "GEN-L03", "GEN-L05", "GEN-M01", "GEN-S01", "GEN-SS01", "GEN-SS02",
	"GEN-SS03", "GEN-P01", "GEN-AE01", "GEN-AE02", "GEN-AE03", "GEN-T01", "GEN-T02",
}
var juniorReligion = []string{"GEN-R01", "GEN-R02", "GEN-R03", "GEN-R04"}
var olMandatory = []string{"OL-M01", "OL-M03", "OL-M04", "OL-M05", "OL-M06"}
var olReligion = []string{"OL-M07", "OL-M08", "OL-M09", "OL-M10", "OL-M11"}
var olBasket1 = []string{
	"OL-B1-01", "OL-B1-02", "OL-B1-03", "OL-B1-04", "OL-B1-L01", "OL-B1-L02", "OL-B1-L03",
	"OL-B1-L04", "OL-B1-L05", "OL-B1-L06", "OL-B1-L07", "OL-B1-L08", "OL-B1-L09", "OL-B1-L10",
	"OL-B1-L11", "OL-B1-L12",
}
var olBasket2 = []string{
	"OL-B2-01", "OL-B2-02", "OL-B2-03", "OL-B2-04", "OL-B2-05", "OL-B2-06", "OL-B2-07",
	"OL-B2-08", "OL-B2-09", "OL-B2-10", "OL-B2-11",
}
var olBasket3 = []string{
	"OL-B3-01", "OL-B3-02", "OL-B3-03", "OL-B3-04", "OL-B3-05", "OL-B3-06", "OL-B3-07",
	"OL-B3-08", "OL-B3-09", "OL-B3-10", "OL-B3-11",
}
var alCommon = []string{"AL-12", "AL-13"}

// alStreams: label -> pool of subjects to pick 3 from (Commerce's pool is exactly 3, so min=max=3 there correctly means "take all three", not a bug).
var alStreams = []struct {
	Label string
	Pool  []string
}{
	{"Physical Science", []string{"AL-01", "AL-02", "AL-10", "AL-20"}},
	{"Biological Science", []string{"AL-09", "AL-02", "AL-01", "AL-08"}},
	{"Technology", []string{"AL-65", "AL-66", "AL-67", "AL-20"}},
	{"Commerce", []string{"AL-33", "AL-32", "AL-21"}},
	// Arts pool uses only the concretely-named subjects from the source
	// spec; the "AL-41-49 / 51-59 / 71-87" ranges are placeholder numbering
	// with no named subjects behind them, so nothing was fabricated for them.
	{"Arts", []string{"AL-22", "AL-23", "AL-24", "AL-25", "AL-28", "AL-29"}},
}

func buildPresetLevels() []presetLevel {
	levels := make([]presetLevel, 0, 40)

	for g := 1; g <= 5; g++ {
		levels = append(levels, presetLevel{
			GradeNumber: g,
			Groups: []presetGroup{
				{"Compulsory Subjects", 7, 7, primaryCompulsory},
			},
		})
	}

	for g := 6; g <= 9; g++ {
		levels = append(levels, presetLevel{
			GradeNumber: g,
			Groups: []presetGroup{
				{"Compulsory Subjects", int32(len(juniorCompulsory)), int32(len(juniorCompulsory)), juniorCompulsory},
				{"Religion", 1, 1, juniorReligion},
			},
		})
	}

	for g := 10; g <= 11; g++ {
		levels = append(levels, presetLevel{
			GradeNumber: g,
			Groups: []presetGroup{
				{"Mandatory", int32(len(olMandatory)), int32(len(olMandatory)), olMandatory},
				{"Religion", 1, 1, olReligion},
				{"Basket 1", 1, 1, olBasket1},
				{"Basket 2", 1, 1, olBasket2},
				{"Basket 3", 1, 1, olBasket3},
			},
		})
	}

	for g := 12; g <= 13; g++ {
		for _, stream := range alStreams {
			levels = append(levels, presetLevel{
				GradeNumber: g,
				LabelSuffix: " - " + stream.Label,
				Groups: []presetGroup{
					{"Common Compulsory", int32(len(alCommon)), int32(len(alCommon)), alCommon},
					{"Stream Subjects", 3, 3, stream.Pool},
				},
			})
		}
	}

	return levels
}
