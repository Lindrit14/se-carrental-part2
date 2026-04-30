package outbox

import "go.mongodb.org/mongo-driver/v2/bson"

// bsonUnmarshal is a thin wrapper kept here so relay.go stays free of
// driver imports — easier to swap to a different bson library later.
func bsonUnmarshal(raw []byte, out any) error {
	return bson.Unmarshal(raw, out)
}
