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
	},
}

func init() {
	rootCmd.AddCommand(initdbCmd)

}
