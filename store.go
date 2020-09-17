package main

import (
	"context"
	"errors"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DirectoryDocument is a MongoDB document struct for a directory document
type DirectoryDocument struct {
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}

var collection *mongo.Collection

func connectToDatabase() {
	logger("Connecting to MongoDB")

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logger(err)
	}

	collection = client.Database("diff-hackerone").Collection("directory")
}

func getStoredDirectoryCount() int {
	count, err := collection.CountDocuments(context.TODO(), bson.M{}, nil)
	if err != nil {
		logger(err)
	}

	logger("Number of stored programs: " + strconv.FormatInt(count, 10))
	return int(count)
}

func insertFullDirectory(directory map[string][]Asset) {
	logger("Inserting full directory into diff-hackerone.directory")

	for name, assets := range directory {
		directoryDocument := DirectoryDocument{
			Name:   name,
			Assets: assets,
		}

		_, err := collection.InsertOne(context.TODO(), directoryDocument)
		if err != nil {
			logger(err)
		}
	}
}

func updateDirectory(directory map[string][]Asset) {
	logger("Updating directory in database...")

	// Get full existing directory
	var existingDirectoryList []DirectoryDocument
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		logger(err)
	}
	err = cursor.All(context.TODO(), &existingDirectoryList)
	if err != nil {
		logger(err)
	}

	existingDirectory := make(map[string][]Asset)
	for _, existingDirectoryDocument := range existingDirectoryList {
		existingDirectory[existingDirectoryDocument.Name] = existingDirectoryDocument.Assets
	}

	// Search for changes
	for name, assets := range directory {

		// New program
		if existingDirectory[name] == nil {
			insertNewProgram(name, assets)
			continue
		}

		// Existing program
		newAssets := []string{}
		changedAssets := []string{}
		isProgramUpdated := false
		if len(assets) != len(existingDirectory[name]) {
			isProgramUpdated = true
		}

		for _, asset := range assets {
			existingAsset, err := findAsset(asset.Name, asset.Type, existingDirectory[name])

			if err != nil {
				// New asset
				if err.Error() == "Asset not found" {
					newAssets = append(newAssets, stringifyAsset(asset))
					isProgramUpdated = true
					continue
				} else {
					logger(err)
				}
			}

			// Existing asset
			if asset.Type != existingAsset.Type || asset.Severity != existingAsset.Severity || asset.Bounty != existingAsset.Bounty {
				changedAssets = append(changedAssets, stringifyAsset(existingAsset)+" -> "+stringifyAsset(asset))
				isProgramUpdated = true
			}
		}

		// Update program
		if isProgramUpdated {
			if len(newAssets) > 0 {
				logger("New asset(s) for program \"" + name + "\" found:")
				for _, newAsset := range newAssets {
					logger("\t" + newAsset)
				}
			}
			if len(changedAssets) > 0 {
				logger("Changed asset(s) for program \"" + name + "\" found:")
				for _, changedAsset := range changedAssets {
					logger("\t" + changedAsset)
				}
			}
			if len(assets)-len(newAssets) < len(existingDirectory[name]) {
				logger("Deleting dead asset(s) from program \"" + name + "\"")
			}
			updateProgram(name, assets)
		}

		// Remove existing program from list to remove
		delete(existingDirectory, name)
	}

	// Delete dead programs
	for name := range existingDirectory {
		deleteDeadProgram(name)
	}

	logger("Updated program directory")
}

func insertNewProgram(name string, assets []Asset) {
	logger("New program \"" + name + "\" found with the following assets:")
	for _, asset := range assets {
		logger("\t" + stringifyAsset(asset))
	}

	directoryDocument := DirectoryDocument{
		Name:   name,
		Assets: assets,
	}

	_, err := collection.InsertOne(context.TODO(), directoryDocument)
	if err != nil {
		logger(err)
	}
}

func updateProgram(name string, assets []Asset) {
	filter := bson.M{"name": name}
	update := bson.D{{"$set", bson.D{{"assets", assets}}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger(err)
	}
}

func deleteDeadProgram(name string) {
	logger("Deleting dead program \"" + name + "\"")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"name": name})
	if err != nil {
		logger(err)
	}
}

func stringifyAsset(asset Asset) string {
	str := "[ " + asset.Name + " | " + asset.Type + " | " + asset.Severity + " | "
	if asset.Bounty {
		str += "paid"
	} else {
		str += "free"
	}
	return str + " ]"
}

func findAsset(name string, assetType string, assets []Asset) (Asset, error) {
	for _, asset := range assets {
		if asset.Name == name && asset.Type == assetType {
			return asset, nil
		}
	}
	return Asset{}, errors.New("Asset not found")
}
