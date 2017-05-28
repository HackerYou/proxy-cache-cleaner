package main


import (
	"fmt"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Cache represents the structure of our cache objects
type Cache struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Date string `json:"date"`
	Response string `json:"response"`
	Endpoint string `json:"endpoint"`
}

func main() {
	fmt.Printf("Cleaning cache \n")
	cleanCount := 0
	format := "Mon Jan 02 2006 15:04:05 GMT-0700 (MST)"

	session, err := mgo.Dial("mongodb://localhost")
	defer session.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	cacheItems := make([]Cache,0,100)

	collection := session.DB("cache").C("caches")

	collection.Find(bson.M{}).All(&cacheItems)

	if err != nil {
		fmt.Println(err.Error())
	}
	itemsToDelete := make(chan string)
	go func() {
		for i := 0; i < len(cacheItems); i++ {
			singleItem := cacheItems[i]

			t, err := time.Parse(format,singleItem.Date)

			if err != nil {
				fmt.Println(err.Error())
			}
			since := time.Since(t).Hours()
			
			if since > 1.0 {
				itemsToDelete <- "deleting"
				err := collection.RemoveId(singleItem.ID)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
		close(itemsToDelete)
	}()

	for range itemsToDelete {
		cleanCount = cleanCount + 1
	}

	fmt.Printf("%d items removed from cache", cleanCount)
}
