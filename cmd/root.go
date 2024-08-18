package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "quizzer <cmd> [args]",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World")
	},
}

func Execute() {
	rootCmd.AddCommand(addQuiz)
	rootCmd.AddCommand(listQuizzes)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
