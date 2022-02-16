package pass

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

// Used as interface to write gopass secrets
type secretContainer struct {
	value string
}

func (c secretContainer) Bytes() []byte {
	return []byte(c.value)
}

// generic function to be used in resource and data source
func populateResourceData(d *schema.ResourceData, provider *passProvider, path string, readData bool) error {
	st := provider.store
	log.Printf("reading %s from gopass", path)

	sec, err := st.Get(context.Background(), path, "") //TODO: support getting a revision via terraform?
	if err != nil {
		return fmt.Errorf("failed to read password at %s: %w", path, err)
	}

	// Retrieve all data items if keys exist
	if readData {
		var keys = sec.Keys()
		if len(keys) != 0 {
			log.Printf("populating data with keys")
			var data = make(map[string]interface{})
			for _, key := range keys {
				data[key], _ = sec.Get(key)
			}
			err := d.Set("data", data)
			if err != nil {
				return err
			}
		}
	}

	if err := d.Set("password", sec.Password()); err != nil {
		log.Printf("error when setting password: %v", err)
		return err
	}

	if err := d.Set("body", sec.Body()); err != nil {
		log.Printf("error when setting body: %v", err)
		return err
	}

	if err := d.Set("full", string(sec.Bytes())); err != nil {
		log.Printf("error when setting full: %v", err)
		return err
	}

	return nil

}
