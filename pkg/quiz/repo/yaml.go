package repo

import (
	"os"
	"slices"

	"github.com/DavesBorges/quizzer/pkg/quiz"
	"gopkg.in/yaml.v3"
)

var _ quiz.QuizRepo = (*YamlRepo)(nil)

type YamlRepo struct{}

var file = "repo.yaml"

type Question struct {
	Prompt       string   `yaml:"prompt"`
	Answer       string   `yaml:"anser"`
	WrongAnswers []string `yaml:"wrong_answers"`
}
type Quiz struct {
	ID        int        `yaml:"id"`
	Name      string     `yaml:"string"`
	Questions []Question `yaml:"questions"`
}
type DB struct {
	Quizes            []Quiz
	SequenceGenerator int
}

// Create implements quiz.QuizRepo.
func (y *YamlRepo) Create(quiz *quiz.Quiz) (int, error) {
	db, err := loadDB()
	if err != nil {
		return 0, err
	}

	quizDataModel := convertQuizToDataModel(quiz)
	db.SequenceGenerator += 1
	quizDataModel.ID = db.SequenceGenerator

	db.Quizes = append(db.Quizes, *quizDataModel)

	if err = saveDB(db); err != nil {
		return 0, err
	}
	return db.SequenceGenerator, nil
}

// Delete implements quiz.QuizRepo.
func (y *YamlRepo) Delete(quizID int) (*quiz.Quiz, error) {

	db, err := loadDB()
	if err != nil {
		return nil, err
	}

	index := slices.IndexFunc(db.Quizes, func(q Quiz) bool {
		return q.ID == quizID
	})

	quiz := db.Quizes[index]
	db.Quizes = slices.Delete(db.Quizes, index, index+1)

	if err = saveDB(db); err != nil {
		return nil, err
	}
	return convertDataModelToQuiz(&quiz), nil
}

// GetAll implements quiz.QuizRepo.
func (y *YamlRepo) GetAll() ([]quiz.Quiz, error) {
	db, err := loadDB()
	if err != nil {
		return nil, err
	}

	return convertDataModelSliceToQuiz(db.Quizes), err
}

// GetByID implements quiz.QuizRepo.
func (y *YamlRepo) GetByID(quizID int) (*quiz.Quiz, error) {
	db, err := loadDB()
	if err != nil {
		return nil, err
	}

	for _, q := range db.Quizes {
		if q.ID == quizID {
			return convertDataModelToQuiz(&q), nil
		}
	}
	return nil, quiz.ErrQuizNotFound
}

// Update implements quiz.QuizRepo.
func (y *YamlRepo) Update(quiz *quiz.Quiz) error {
	db, err := loadDB()
	if err != nil {
		return err
	}

	idx := slices.IndexFunc(db.Quizes, func(q Quiz) bool {
		return q.ID == quiz.ID
	})

	db.Quizes[idx] = *convertQuizToDataModel(quiz)

	if err = saveDB(db); err != nil {
		return err
	}

	return nil
}

func loadDB() (*DB, error) {
	if _, err := os.Stat(file); err != nil {
		// Don't bother reading the file since it doesnt exist
		return &DB{}, nil
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var db DB
	if err = yaml.NewDecoder(f).Decode(&db); err != nil {
		return nil, err
	}

	return &db, nil
}
func saveDB(db *DB) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = yaml.NewEncoder(f).Encode(&db); err != nil {
		return err
	}

	return nil
}

func convertQuizToDataModel(q *quiz.Quiz) *Quiz {
	return &Quiz{
		ID:        q.ID,
		Name:      q.Name,
		Questions: convertQuestionsToDataModel(q.Questions),
	}
}
func convertDataModelToQuiz(q *Quiz) *quiz.Quiz {
	return &quiz.Quiz{
		ID:        q.ID,
		Name:      q.Name,
		Questions: convertDataModelToQuestions(q.Questions),
	}
}

func convertDataModelSliceToQuiz(quizzes []Quiz) []quiz.Quiz {
	var result = make([]quiz.Quiz, 0, len(quizzes))
	for _, q := range quizzes {
		result = append(result, *convertDataModelToQuiz(&q))
	}
	return result
}
func convertQuestionsToDataModel(questions []quiz.Question) []Question {
	result := make([]Question, 0, len(questions))
	for _, question := range questions {
		result = append(result, Question{
			Prompt:       question.Prompt,
			Answer:       question.Answer,
			WrongAnswers: question.WrongAnswers,
		})
	}

	return result
}

func convertDataModelToQuestions(questions []Question) []quiz.Question {
	result := make([]quiz.Question, 0, len(questions))
	for _, question := range questions {
		result = append(result, quiz.Question{
			Prompt:       question.Prompt,
			Answer:       question.Answer,
			WrongAnswers: question.WrongAnswers,
		})
	}

	return result
}
