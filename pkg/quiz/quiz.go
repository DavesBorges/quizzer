package quiz

import "errors"

type Question struct {
	ID           int
	Answer       string
	WrongAnswers []string
	Prompt       string
}

type PublicQuestion struct {
	ID      int
	Prompt  string
	Options []string
}

type Quiz struct {
	ID        int
	Questions []Question
	Name      string
}
type PublicQuiz struct {
	Questions []PublicQuestion
	Name      string
}

type Service struct {
	repo QuizRepo
}

func NewService(repo QuizRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(quizName string) (int, error) {
	return s.repo.Create(&Quiz{Name: quizName})
}

var ErrDuplicatedQuestion = errors.New("duplicated question")
var ErrDuplicatedAnswer = errors.New("duplicated answer")
var ErrQuizNotFound = errors.New("quiz not found")

func (s *Service) AddQuestionToQuiz(quizID int, question *Question) error {
	quiz, err := s.repo.GetByID(quizID)
	if err != nil {
		return err
	}

	if err = checkForDuplicatedAnswer(question); err != nil {
		return err
	}

	if err = checkForDuplicatedQuestion(question, quiz.Questions); err != nil {
		return err
	}

	quiz.Questions = append(quiz.Questions, *question)
	if err = s.repo.Update(quiz); err != nil {
		return err
	}

	return nil
}

func (s *Service) List() ([]Quiz, error) {

	return s.repo.GetAll()
}
func (s *Service) Delete(quizID int) (*Quiz, error) {
	return s.repo.Delete(quizID)
}

type QuizRepo interface {
	Create(quiz *Quiz) (int, error)
	Update(quiz *Quiz) error
	GetAll() ([]Quiz, error)
	GetByID(id int) (*Quiz, error)
	Delete(quizID int) (*Quiz, error)
}

func checkForDuplicatedAnswer(question *Question) error {
	for i := 0; i < len(question.WrongAnswers); i++ {
		if question.WrongAnswers[i] == question.Answer {
			return ErrDuplicatedAnswer
		}

		for j := i + 1; j < len(question.WrongAnswers); j++ {
			if question.WrongAnswers[i] == question.WrongAnswers[j] {
				return ErrDuplicatedAnswer
			}
		}
	}

	return nil
}

func checkForDuplicatedQuestion(newQuestion *Question, existingQuestions []Question) error {
	for _, q := range existingQuestions {
		if q.Prompt == newQuestion.Prompt {
			return ErrDuplicatedQuestion
		}
	}

	return nil
}
