package main

import (
	"log"
	"os"
	"strings"

	"masjidku_backend/internals/configs"
	"masjidku_backend/internals/seeds"

	categories "masjidku_backend/internals/seeds/lessons/categories/categories"
	categories_news "masjidku_backend/internals/seeds/lessons/categories/categories_news"
	difficulties "masjidku_backend/internals/seeds/lessons/difficulties/difficulties"
	difficulties_news "masjidku_backend/internals/seeds/lessons/difficulties/difficulties_news"
	subcategories "masjidku_backend/internals/seeds/lessons/subcategories/subcategories"
	subcategories_news "masjidku_backend/internals/seeds/lessons/subcategories/subcategories_news"
	themes_or_levels "masjidku_backend/internals/seeds/lessons/themes_or_levels/themes_or_levels"
	themes_or_levels_news "masjidku_backend/internals/seeds/lessons/themes_or_levels/themes_or_levels_news"
	units "masjidku_backend/internals/seeds/lessons/units/units"
	units_news "masjidku_backend/internals/seeds/lessons/units/units_news"
	level "masjidku_backend/internals/seeds/progress/levels"
	rank "masjidku_backend/internals/seeds/progress/ranks"
	evaluations "masjidku_backend/internals/seeds/quizzes/evaluations"
	exams "masjidku_backend/internals/seeds/quizzes/exams"
	questions "masjidku_backend/internals/seeds/quizzes/questions"
	quizzes "masjidku_backend/internals/seeds/quizzes/quizzes"
	reading "masjidku_backend/internals/seeds/quizzes/readings"
	section_quizzes "masjidku_backend/internals/seeds/quizzes/section_quizzes"
	users "masjidku_backend/internals/seeds/users/auth"
	userProfiles "masjidku_backend/internals/seeds/users/users"
	tooltips "masjidku_backend/internals/seeds/utils/tooltips"
	survey "masjidku_backend/internals/seeds/users/surveys/survey_questions"
	user_survey "masjidku_backend/internals/seeds/users/surveys/user_surveys"
	test_exam "masjidku_backend/internals/seeds/users/test-exams/test_exams"
	user_test_exam "masjidku_backend/internals/seeds/users/test-exams/user_test_exams"

)

func main() {
	configs.LoadEnv()
	db := configs.InitSeederDB()

	log.Println("🚀 Menjalankan seeder...")
	if len(os.Args) < 2 {
		log.Fatalln("❌ Mohon masukkan argumen seperti: all | users | users_profile | lessons | quizzes | progress")
	}

	switch strings.ToLower(os.Args[1]) {
	case "all":
		seeds.RunAllSeeds(db)
	case "users":
		users.SeedUsersFromJSON(db, "internals/seeds/users/auth/data_users_2.json")
	case "users_profile":
		userProfiles.SeedUsersProfileFromJSON(db, "internals/seeds/users/users/data_users_profiles.json")
	case "lessons":
		difficulties.SeedDifficultiesFromJSON(db, "internals/seeds/lessons/difficulties/difficulties/data_difficulties.json")
		difficulties_news.SeedDifficultyNewsFromJSON(db, "internals/seeds/lessons/difficulties/difficulties_news/data_difficulties_news.json")
		categories.SeedCategoriesFromJSON(db, "internals/seeds/lessons/categories/categories/data_categories.json")
		categories_news.SeedCategoriesNewsFromJSON(db, "internals/seeds/lessons/categories/categories_news/data_categories_news.json")
		subcategories.SeedSubcategoriesFromJSON(db, "internals/seeds/lessons/subcategories/subcategories/data_subcategories.json")
		subcategories_news.SeedSubcategoryNewsFromJSON(db, "internals/seeds/lessons/subcategories/subcategories_news/data_subcategories_news.json")
		themes_or_levels.SeedThemesOrLevelsFromJSON(db, "internals/seeds/lessons/themes_or_levels/themes_or_levels/data_themes_or_levels.json")
		themes_or_levels_news.SeedThemesOrLevelsNewsFromJSON(db, "internals/seeds/lessons/themes_or_levels/themes_or_levels_news/data_themes_or_levels_news.json")
		units.SeedUnitsFromJSON(db, "internals/seeds/lessons/units/units/data_units.json")
		units_news.SeedUnitsNewsFromJSON(db, "internals/seeds/lessons/units/units_news/data_units_news.json")
	case "quizzes":
		evaluations.SeedEvaluationsFromJSON(db, "internals/seeds/quizzes/evaluations/data_evaluations.json")
		exams.SeedExamsFromJSON(db, "internals/seeds/quizzes/exams/data_exams.json")
		section_quizzes.SeedSectionQuizzesFromJSON(db, "internals/seeds/quizzes/section_quizzes/data_section_quizzes.json")
		questions.SeedQuestionsFromJSON(db, "internals/seeds/quizzes/questions/data_questions.json")
		quizzes.SeedQuizzesFromJSON(db, "internals/seeds/quizzes/quizzes/data_quizzes.json")
		reading.SeedReadingsFromJSON(db, "internals/seeds/quizzes/readings/data_readings.json")
	case "progress":
		level.SeedLevelRequirementsFromJSON(db, "internals/seeds/progress/levels/data_levels_requirements.json")
		rank.SeedRanksRequirementsFromJSON(db, "internals/seeds/progress/ranks/data_ranks_requirements.json")
	case "utils":
		tooltips.SeedTooltipsFromJSON(db, "internals/seeds/utils/tooltips/data_tooltips.json")

	case "survey_test_exam":
		survey.SeedSurveyQuestionsFromJSON(db, "internals/seeds/users/surveys/survey_questions/data_survey_questions.json")
		user_survey.SeedUserSurveysFromJSON(db, "internals/seeds/users/surveys/user_surveys/data_user_surveys.json")
		test_exam.SeedTestExamsFromJSON(db, "internals/seeds/users/test-exams/test_exams/data_test_exams.json")
		user_test_exam.SeedUserTestExamsFromJSON(db, "internals/seeds/users/test-exams/user_test_exams/data_user_test_exams.json")

	default:
		log.Fatalf("❌ Argumen '%s' tidak dikenali", os.Args[1])
	}
}
