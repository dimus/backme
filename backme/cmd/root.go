package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/dimus/backme"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, buildDate, buildVersion string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "backme",
	Short: "Backup files organizer.",

	Long: `Quite often big files (like database dumps) are backed up on a
 regular basis. Such backups become a disk space hogs and when time comes to
 restoring data it is cumbersome to find the right file for the job. This
 script checks files of a certain pattern and sorts them by date into daily
 monthly, yearly directories, progressively deleting more and more of aged
 files. For example it keeps all files from the last 24 hours, but keeps only
 one file per day in the last 30 days, only one file per month for the first
 3 years, and after that only 1 file per year. As a result backup takes
 significantly less space and it is easier to find a file from a specific
 period of time.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		versionFlag(cmd)

		conf := getConfig()
		err := backme.Organize(conf)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string, d string) {
	buildVersion = v
	buildDate = d
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.backme.yaml)")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Show build version and date")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("version", "v", false, "Show build version and date")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".backme" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".backme")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Println(err)
	}
}

func versionFlag(cmd *cobra.Command) {
	version, err := cmd.Flags().GetBool("version")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if version {
		fmt.Printf("\nversion: %s\n\ndate:    %s\n\n", buildVersion, buildDate)
		os.Exit(0)
	}
}

func getConfig() *backme.Config {
	conf := backme.NewConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = backme.CheckConfig(conf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("The following backup directories are in config:")
	for i, v := range conf.InputDirs {
		log.Printf("Backup dir %d: %s", i, v.Path)
		if _, err := os.Stat(v.Path); os.IsNotExist(err) {
			log.Printf("Directory %s does not exist, exiting...", v.Path)
			os.Exit(1)
		}
	}
	log.Println(backme.LogSep)
	return conf
}
