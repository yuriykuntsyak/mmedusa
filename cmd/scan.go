/*
Copyright © 2023 Yuriy Kuntsyak
*/
package cmd

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Define a struct to represent the media file entry
type MediaFile struct {
	ID   uint   `gorm:"primaryKey"`
	Path string `gorm:"unique"`
	Hash string
}

// var db *gorm.DB
var path string
var pattern string

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans given path for media files.",
	Long:  `Scans given path for media files. It will also scan subdirectories recursively.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if path == "" {
			path = "."
		}
		log.Println("scan called with path:", path)

		errChan := make(chan error)
		fileChan := make(chan string)

		// Initialize the SQLite database connection
		db, err := gorm.Open(sqlite.Open("media.db"), &gorm.Config{})
		if err != nil {
			return err
		}
		// Auto-migrate the MediaFile struct to create the corresponding table in the database
		err = db.AutoMigrate(&MediaFile{})
		if err != nil {
			return err
		}

		go func() {
			errChan <- filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
				if err != nil {
					errChan <- err
					return nil
				}

				if !fileInfo.IsDir() && (pattern == "" || strings.Contains(filePath, pattern)) {
					fileChan <- filePath
				}

				return nil
			})
			close(fileChan)
		}()

		var wg sync.WaitGroup
		for i := 0; i < runtime.NumCPU(); i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for filePath := range fileChan {
					var existingFile MediaFile
					result := db.First(&existingFile, "path = ?", filePath)
					if result.Error == nil {
						log.Debug("File already exists in the database ", filePath)
						continue
					}

					file, err := os.Open(filePath)
					if err != nil {
						errChan <- err
						continue
					}
					defer file.Close()

					hash := sha1.New()
					if _, err := io.Copy(hash, file); err != nil {
						errChan <- err
						continue
					}

					log.Printf("%s: %x\n", filePath, hash.Sum(nil))

					// Create a new MediaFile entry
					mediaFile := MediaFile{
						Path: filePath,
						Hash: fmt.Sprintf("%x", hash.Sum(nil)),
					}

					// Save the MediaFile entry to the database
					err = db.Create(&mediaFile).Error
					if err != nil {
						errChan <- err
						continue
					}

				}
			}()
		}

		go func() {
			wg.Wait()
			close(errChan)
		}()

		for err := range errChan {
			if err != nil {
				if os.IsPermission(err) {
					log.Errorf("Permission denied: %v", err)
				} else if os.IsNotExist(err) {
					log.Errorf("File not found: %v", err)
				} else if os.IsTimeout(err) {
					log.Errorf("Timeout: %v", err)
				} else {
					log.Errorf("Unexpected error: %v", err)
				}
			}
		}

		// Close the database connection
		// db.Close()
		return nil
	},
}

func init() {

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(&path, "path", "p", "", "Path to scan")
	scanCmd.Flags().StringVarP(&pattern, "pattern", "", "", "Pattern to match")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
