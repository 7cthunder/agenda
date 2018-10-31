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

	"github.com/7cthunder/agenda/entity"
	"github.com/spf13/cobra"
)

// createMeetingCmd represents the createMeeting command
var createMeetingCmd = &cobra.Command{
	Use:   "cm",
	Short: "Create a meeting with title, startTime, endTime and participators",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("createMeeting called")
		title, _ := cmd.Flags().GetString("title")
		startTimeS, _ := cmd.Flags().GetString("starttime")
		endTimeS, _ := cmd.Flags().GetString("endtime")
		//ptcpt, _ := cmd.Flags().GetStringSlice("participator")
		ptcpt := cmd.Flags().Args()

		instance := entity.GetStorage()

		if instance.GetCurUser().GetName() == "" {
			fmt.Println("You have not logged in yet, please log in first!")
			return
		}
		if title == "" {
			fmt.Println("You do not enter title, please input again!")
			return
		}
		if startTimeS == "" {
			fmt.Println("You do not enter start time, please input again!")
			return
		}
		if endTimeS == "" {
			fmt.Println("You do not enter end time, please input again!")
			return
		}
		if len(ptcpt) == 0 {
			fmt.Println("No participator!")
			return
		}

		sponsor := instance.GetCurUser().GetName()
		startTime := entity.StringToDate(startTimeS)
		endTime := entity.StringToDate(endTimeS)

		if !startTime.IsValid() {
			fmt.Println("Invalid start time!")
			return
		}
		if !endTime.IsValid() {
			fmt.Println("Invalid end time!")
			return
		}
		if startTime.IsGreaterThanEqual(endTime) {
			fmt.Println("Start time cannot be later or equal than end time!")
			return
		}

		mfilter1 := func(m *entity.Meeting) bool {
			return m.GetTitle() == title
		}
		if len(instance.QueryMeeting(mfilter1)) > 0 {
			fmt.Println("Duplicate title, please change it!")
			return
		}

		for _, p := range ptcpt {
			if p == sponsor {
				fmt.Println("Sponsor cannot be participator!")
				return
			}
		}

		for i := 0; i < len(ptcpt); i++ {
			for j := i + 1; j < len(ptcpt); j++ {
				if ptcpt[i] == ptcpt[j] {
					fmt.Println("Duplicate participators!")
					return
				}
			}
		}

		// ufilter1 := func(u *entity.User) bool {
		// 	return u.GetName() == sponsor
		// }
		// if len(instance.QueryUser(ufilter1)) == 0 {
		// 	fmt.Println("Non-existent sponsor!")
		// 	return
		// }

		for _, p := range ptcpt {
			ufilter1 := func(u *entity.User) bool {
				return u.GetName() == p
			}
			if len(instance.QueryUser(ufilter1)) == 0 {
				fmt.Println("There is at least one non-existent participator!")
				return
			}
		}

		mfilter2 := func(m *entity.Meeting) bool {
			if m.GetSponsor() != sponsor && !m.IsParticipator(sponsor) {
				return false
			}
			if startTime.IsGreaterThanEqual(m.GetEndTime()) ||
				endTime.IsLessThanEqual(m.GetStartTime()) {
				return false
			}
			return true
		}
		if len(instance.QueryMeeting(mfilter2)) > 0 {
			fmt.Println("Sponsor's time conflict!")
			return
		}

		for _, p := range ptcpt {
			mfilter3 := func(m *entity.Meeting) bool {
				if m.GetSponsor() != p && !m.IsParticipator(p) {
					return false
				}
				if startTime.IsGreaterThanEqual(m.GetEndTime()) ||
					endTime.IsLessThanEqual(m.GetStartTime()) {
					return false
				}
				return true
			}
			if len(instance.QueryMeeting(mfilter3)) > 0 {
				fmt.Println("Participator's time conflict!")
				return
			}
		}
		fmt.Println("Create meeting successfully!")
		instance.CreateMeeting(*entity.NewMeeting(sponsor, title, startTime, endTime, ptcpt))
	},
}

func init() {
	rootCmd.AddCommand(createMeetingCmd)

	// Here you will define your flags and configuration settings.
	createMeetingCmd.Flags().StringP("title", "t", "", "title of meeitng")
	createMeetingCmd.Flags().StringP("starttime", "s", "", "start time of meeting")
	createMeetingCmd.Flags().StringP("endtime", "e", "", "end time of meeting")
	//createMeetingCmd.Flags().StringSliceP("participator", "p", []string{}, "participators of meeting")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createMeetingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createMeetingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
