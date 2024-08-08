package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// the identifier can be a name OR an ID, but ID is unique
func lookupGroupByIdentifier(ctx context.Context, client *APIClient, identifier string) (*groupDataSourceModel, error) {
	readData, _ := json.Marshal(map[string]interface{}{
		"group_identifier": identifier,
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
		ID:          types.StringValue(groupReturned.ID),
		Description: types.StringValue(groupReturned.Description),
	}, nil
}
