/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"math/rand/v2"
	"red-db-test/utils"
	"sync/atomic"
	"time"
)

var (
	okLoginRPS     = 100
	okLoginCounter atomic.Int64

	badLoginRPS     = 100
	badLoginCounter atomic.Int64

	commitTaskRPS     = 100
	commitTaskCounter atomic.Int64

	clickInviteRPS     = 100
	clickInviteCounter atomic.Int64
)

// loaddbCmd represents the loaddb command
var loaddbCmd = &cobra.Command{
	Use:   "loaddb",
	Short: "workload",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		go loadOkLogins()
		go loadBadLogins()
		go loadCommitTasks()
		go loadClickInvite()

		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			log.Printf("ok logins rps:  %v", okLoginCounter.Swap(0))
			log.Printf("bad logins rps: %v", badLoginCounter.Swap(0))
			log.Printf("commit task rps: %v", commitTaskCounter.Swap(0))
			log.Printf("click invite rps:  %v", clickInviteCounter.Swap(0))
			log.Println()
		}
	},
}

func loadOkLogins() {
	if okLoginRPS == 0 {
		return
	}

	db := NewDB()
	users := utils.GenerateUsers(usersSeed, initUsersCount)

	ticker := time.NewTicker(time.Second / time.Duration(okLoginRPS))
	rnd := rand.New(rand.NewPCG(1, 2))
	for {
		<-ticker.C
		go func(index int) {
			defer okLoginCounter.Add(1)

			user := &users[index]

			if err := db.Login(user.ID, user.Token); err != nil {
				log.Printf("failed to login for user id %v, err: %+v", user.ID, err)
			}
		}(rnd.IntN(initUsersCount))
	}
}
func loadBadLogins() {
	if badLoginRPS == 0 {
		return
	}

	db := NewDB()
	users := utils.GenerateUsers(usersSeed, initUsersCount)

	ticker := time.NewTicker(time.Second / time.Duration(badLoginRPS))
	rnd := rand.New(rand.NewPCG(1, 2))
	for {
		<-ticker.C
		go func(index int) {
			defer badLoginCounter.Add(1)

			user := &users[index]

			if err := db.Login(user.ID, "asd"); err == nil {
				log.Printf("failed to bad login (no error) for user id %v, err: %+v", user.ID, err)
			}
		}(rnd.IntN(initUsersCount))
	}
}

func loadCommitTasks() {
	if commitTaskRPS == 0 {
		return
	}

	db := NewDB()
	users := utils.GenerateUsers(usersSeed, initUsersCount)
	tasks := utils.GenerateTasks(tasksSeed, initTaskCount)

	ticker := time.NewTicker(time.Second / time.Duration(badLoginRPS))
	rnd := rand.New(rand.NewPCG(1, 2))
	for {
		<-ticker.C
		go func(userIndex, taskIndex int) {
			defer commitTaskCounter.Add(1)

			user := &users[userIndex]
			task := &tasks[taskIndex]

			if err := db.CompleteTask(user.ID, task.ID); err != nil {
				log.Printf("failed to commit task for user id %v, task id %v, err: %+v", user.ID, task.ID, err)
			}
		}(rnd.IntN(initUsersCount), rnd.IntN(initTaskCount))
	}
}

func loadClickInvite() {
	if clickInviteRPS == 0 {
		return
	}

	db := NewDB()
	users := utils.GenerateUsers(usersSeed, initUsersCount)

	ticker := time.NewTicker(time.Second / time.Duration(okLoginRPS))
	rnd := rand.New(rand.NewPCG(1, 2))
	for {
		<-ticker.C
		go func(index int) {
			defer clickInviteCounter.Add(1)

			user := &users[index]

			if err := db.Login(user.ID, user.Token); err != nil {
				log.Printf("failed to click invite for user id %v, err: %+v", user.ID, err)
			}
		}(rnd.IntN(initUsersCount))
	}
}

func init() {
	rootCmd.AddCommand(loaddbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loaddbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loaddbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
