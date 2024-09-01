/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"red-db-test/utils"
	"slices"

	"github.com/spf13/cobra"
)

var (
	initUsersCount = 1000000
	initTaskCount  = 100
)

// initdbCmd represents the initdb command
var initdbCmd = &cobra.Command{
	Use:   "initdb",
	Short: "initialize selected database for test",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Initializing db %q", dbName)
		db := NewDB()
		utils.Must0(db.CreateTables())

		log.Println("Generating users")
		users := utils.GenerateUsers(usersSeed, initUsersCount)
		usersParts := slices.Collect(slices.Chunk(users, 1000))
		for index, part := range usersParts {
			if err := db.UploadUsers(part); err == nil {
				log.Printf("uploaded part %v/%v", index+1, len(usersParts))
			} else {
				log.Fatalf("Failed upload part %v/%v: %+v", index+1, len(usersParts), err)
			}
		}

		log.Println("Generate tasks")
		tasks := utils.GenerateTasks(tasksSeed, initTaskCount)
		for i, task := range tasks {
			if err := db.CreateTask(task); err == nil {
				log.Printf("created task %v/%v", i+1, initTaskCount)
			}
		}

		log.Println("set some tasks checked")
		maxUserIndex := min(len(users)-1, 100)
		maxTaskIndex := min(len(tasks)-1, 100)
		for userIndex := 0; userIndex <= maxUserIndex; userIndex++ {
			log.Printf("Setting task completed for user index: %v (%v)", userIndex, users[userIndex].ID)
			for taskIndex := 0; taskIndex <= userIndex && taskIndex <= maxTaskIndex; taskIndex++ {
				if err := db.CompleteTask(users[userIndex].ID, tasks[taskIndex].ID); err != nil {
					log.Fatalf("failed set task completed for user id '%v', task id '%v', err: %+v",
						users[userIndex].ID, tasks[taskIndex].ID, err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initdbCmd)

}
