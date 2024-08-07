package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func lookupGroupByName(ctx context.Context, client *APIClient, name string) (*groupDataSourceModel, error) {
	readData, _ := json.Marshal(map[string]interface{}{
		"group_identifier": name,
		"record_offset":    0,
		"record_size":      10,
	})
	tflog.Trace(ctx, fmt.Sprintf("My postBody will be something like %v", readData))
	itemReturned, err := client.LookupExactlyOne(ctx, "/groups/search", readData)
	if err != nil {
		return nil, fmt.Errorf("Error from APIClient:", err)
	}
	var groupReturned groupSearchResult
	err = json.Unmarshal(itemReturned, &groupReturned)
	if err != nil {
		return nil, fmt.Errorf("I couldn't unmarshall the JSON because", err)
	}
	return &groupDataSourceModel{
		Name:        types.StringValue(groupReturned.Name),
		Description: types.StringValue(groupReturned.Description),
	}, nil
}
