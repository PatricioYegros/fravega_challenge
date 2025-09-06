package models

type Counter struct {
	ID            string `bson:"_id"`
	SequenceValue int64  `bson:"sequence_value"`
}
