package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// IncidentModel is the model used to describe incidents.
// Most fields are selected from options as described in the incidentOptions collection in the database
type IncidentModel struct {
	ID                     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Detection              string             `json:"detection" bson:"detection"`
	Actor                  string             `json:"actor" bson:"actor"`
	Plan                   string             `json:"plan" bson:"plan"`
	FileSet                []string           `json:"fileset" bson:"fileset"`
	DateCreated            int64              `json:"datecreated,omitempty" bson:"datecreated,omitempty"`
	LastModified           int64              `json:"lastmodified,omitempty" bson:"lastmodified,omitempty"`
	IsStarred              bool               `json:"isstarred" bson:"isstarred"`
	Subject                string             `json:"subject" bson:"subject"`
	Description            string             `json:"description" bson:"description"`
	Severity               int32              `json:"severity" bson:"severity"`
	IsIncident             bool               `json:"isincident" bson:"isincident"`
	IsMajor                bool               `json:"ismajor" bson:"ismajor"`
	Status                 string             `json:"status" bson:"status"`
	Confidentiality        int32              `json:"confidentiality" bson:"confidentiality"`
	Category               string             `json:"category" bson:"category"`
	OpenedBy               string             `json:"openedby" bson:"openedby"`
	ConcernedBusinessLines []string           `json:"concernedbusinesslines" bson:"concernedbusinesslines"`
}
