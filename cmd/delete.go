/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

var AccessKeyID string
var SecretAccessKey string
var MyRegion string
var bucket string

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}
}

func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete called")
		LoadEnv()
		AccessKeyID = GetEnvWithKey("AWS_ACCESS_KEY_ID")
		SecretAccessKey = GetEnvWithKey("AWS_SECRET_ACCESS_KEY")
		MyRegion = GetEnvWithKey("AWS_REGION")
		if GetEnvWithKey("ENVIRONMENT") ==  "production" {
			bucket = GetEnvWithKey("AWS_BUCKET_PUBLIC")
		}
		else {
			bucket = GetEnvWithKey("AWS_BUCKET_PUBLIC_TEST")
		}
		

		sess, err := session.NewSession(
			&aws.Config{
				Region: aws.String(MyRegion),
				Credentials: credentials.NewStaticCredentials(
					AccessKeyID,
					SecretAccessKey,
					"", // a token will be created when the session it's used.
				),
			})

		if err != nil {
			panic(err)
		}
		
		// Create S3 service client
		svc := s3.New(sess)
		obj := os.args[1]
		_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(obj)})
		if err != nil {
			exitErrorf("Unable to delete object %q from bucket %q, %v", obj, bucket, err)
		}

		err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(obj),
		})

		fmt.Printf("Object %q successfully deleted\n", obj)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}