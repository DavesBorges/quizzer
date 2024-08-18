package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/DavesBorges/quizzer/pkg/quiz"
	"github.com/DavesBorges/quizzer/pkg/quiz/repo"
	"github.com/spf13/cobra"
)

var addQuiz = &cobra.Command{
	Use:   "add <name> ",
	Short: "Adds a quiz",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("quiz name not provided")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		service := quiz.NewService(&repo.YamlRepo{})
		id, err := service.Add(args[0])
		if err != nil {
			fmt.Printf("Failed to create quiz: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(id)
	},
}

var listQuizzes = &cobra.Command{
	Use:   "list",
	Short: "Shows all quizzes",
	Run: func(cmd *cobra.Command, args []string) {
		service := quiz.NewService(&repo.YamlRepo{})
		quizes, err := service.List()
		if err != nil {
			fmt.Printf("Failed to create quiz: %v\n", err)
			os.Exit(1)
		}

		for _, q := range quizes {
			fmt.Printf("- %v\n", q.Name)
		}
	},
}
