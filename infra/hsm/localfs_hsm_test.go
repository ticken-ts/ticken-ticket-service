package hsm_test

import (
	"testing"
	"ticken-event-service/infra/hsm"
	"ticken-event-service/utils"
)

const testEncryptionKey = "@NcRfUjXnZr4u7x!A%D*G-KaPdSgVkYp"

func TestStoreAndRetrieveData(t *testing.T) {
	data := "how you doing?"

	path, _ := utils.GetServiceRootPath()

	localFileSystemHSM, err := hsm.NewLocalFileSystemHSM(testEncryptionKey, path)
	if err != nil {
		t.Errorf("fail to create localFileSystemHSM: %s", err.Error())
	}

	key, err := localFileSystemHSM.Store([]byte(data))
	if err != nil {
		t.Errorf("fail to store data: %s", err.Error())
	}

	retrieved, err := localFileSystemHSM.Retrieve(key)
	if err != nil {
		t.Errorf("fail to retrieve data: %s", err.Error())
	}

	if string(retrieved) != data {
		t.Errorf("retrived data doesnt match stored data: %s != %s", data, string(retrieved))
	}
}
