package dbclient

import (
		"github.com/Ankurk99/compile_microsevice/cmd/model"
)

type IBoltClient interface {
		OpenBoltDb()
		QueryAccount(accountId string) (model.Account, error)
		Seed()
}

type BoltClient struct {
		boltDB *bolt.DB
}

func (bc *BoltClient) OpenBoltDb() {
		var err error
		bc.boltDB, err = bolt.Open("accounts.db", 0600, nil)
		if err !=nil {
				log.Fatal(err)
		}
}

func (bc *BoltClient) Seed() {
		initializeBucket()
		seedAccounts()
}

func (bc *BoltClient) initializeBucket() {
		bc.boltDB.Update(func(tx *bolt.Tx) error {
				_, err := tx.CreateBucket([]byte("AccountBucket"))
				if err != nil {
						return fmt.Errorf("create bucket failed %s", err)
				}
				return nil
		})
}

func (bc *BoltClient) seedAccounts() {

        total := 100
        for i := 0; i < total; i++ {

                // Generate a key 10000 or larger
                key := strconv.Itoa(10000 + i)

                // Create an instance of our Account struct
                acc := model.Account{
                        Id: key,
                        Name: "Person_" + strconv.Itoa(i),
                }

                // Serialize the struct to JSON
                jsonBytes, _ := json.Marshal(acc)

                // Write the data to the AccountBucket
                bc.boltDB.Update(func(tx *bolt.Tx) error {
                        b := tx.Bucket([]byte("AccountBucket"))
                        err := b.Put([]byte(key), jsonBytes)
                        return err
                })
        }
        fmt.Printf("Seeded %v fake accounts...\n", total)
}


