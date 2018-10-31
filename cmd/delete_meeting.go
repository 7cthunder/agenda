// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	entity "github.com/7cthunder/agenda/entity"
	"github.com/spf13/cobra"
)

// deleteMeetingCmd represents the deleteMeeting command
var deleteMeetingCmd = &cobra.Command{
	Use:   "delm",
	Short: "Delete a meeting which current user created",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		instance := entity.GetStorage()
		curU := instance.GetCurUser()
		title, _ := cmd.Flags().GetString("title")

		if curU.GetName() == "" {
			fmt.Println("You have not logged in yet, please log in first!")
		}

		if title == "" {
			fmt.Println("Please input the title of the meeting you want to delete")
		}

		mfilter := func(m *entity.Meeting) bool {
			return curU.GetName() == m.GetSponsor() && title == m.GetTitle()
		}

		if instance.DeleteMeeting(mfilter) > 0 {
			fmt.Println("Delete successfully!")
		} else {
			fmt.Println("you don't sponsor this meeting")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteMeetingCmd)

	// Here you will define your flags and configuration settings.

	deleteMeetingCmd.Flags().StringP("title", "t", "", "meeting title")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteMeetingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteMeetingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
