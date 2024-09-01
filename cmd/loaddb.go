/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"math/rand/v2"
	"red-db-test/utils"
	"time"
)

var (
	okLoginRPS     = 100
	okLoginCounter utils.LatencyMetric

	badLoginRPS     = 100
	badLoginCounter utils.LatencyMetric

	commitTaskRPS     = 100
	commitTaskCounter utils.LatencyMetric

	clickInviteRPS     = 100
	clickInviteCounter utils.LatencyMetric
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
			printStat("ok logins", &okLoginCounter)
			printStat("bad logins", &badLoginCounter)
			printStat("commit task", &commitTaskCounter)
			printStat("click invite", &clickInviteCounter)
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
	rnd := rand.New(rand.NewPCG(304, 1341))
	for {
		<-ticker.C
		go func(index int) {
			start := time.Now()
			okLoginCounter.AddSince(start)

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
	rnd := rand.New(rand.NewPCG(534, 134))
	for {
		<-ticker.C
		go func(index int) {
			start := time.Now()
			badLoginCounter.AddSince(start)

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
	rnd := rand.New(rand.NewPCG(415, 3133))
	for {
		<-ticker.C
		go func(userIndex, taskIndex int) {
			start := time.Now()
			commitTaskCounter.AddSince(start)

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
	rnd := rand.New(rand.NewPCG(332, 231))
	for {
		<-ticker.C
		go func(index int) {
			start := time.Now()
			clickInviteCounter.AddSince(start)

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

func printStat(mess string, m *utils.LatencyMetric) {
	stat := m.Stat(0.5, 0.9, 0.99, 1.0)
	log.Printf(
		mess+": 0.5 (%v), 0.9 (%v), 0.99 (%v), 1.0 (%v), Total Count: %v",
		stat.Durations[0],
		stat.Durations[1],
		stat.Durations[2],
		stat.Durations[3],
		stat.TotalCount,
	)
}
