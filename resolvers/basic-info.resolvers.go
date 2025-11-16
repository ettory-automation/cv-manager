package resolvers

import (
	"context"
	"fmt"

	"github.com/ettory-automation/cv-manager/generated"
	"github.com/ettory-automation/cv-manager/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionName = "basic-info"

// CreateBasicInfo creates a new BasicInfo document
func (r *mutationResolver) CreateBasicInfo(ctx context.Context, basicInfo model.BasicInfoInput) (*model.BasicInfo, error) {
	collection := r.DB.Collection(collectionName)

	newID := primitive.NewObjectID()
	newBasicInfo := model.BasicInfo{
		ID:             newID.Hex(), // store ID as string
		FirstName:      basicInfo.FirstName,
		LastName:       basicInfo.LastName,
		AdditionalName: *basicInfo.AdditionalName,
		Pronouns:       basicInfo.Pronouns,
		Headline:       basicInfo.Headline,
	}

	_, err := collection.InsertOne(ctx, newBasicInfo)
	if err != nil {
		return nil, err
	}

	return &newBasicInfo, nil
}

// UpdateBasicInfo updates an existing BasicInfo document by string ID
func (r *mutationResolver) UpdateBasicInfo(ctx context.Context, id string, basicInfo *model.BasicInfoInput) (*model.BasicInfo, error) {
	collection := r.DB.Collection(collectionName)

	update := bson.M{
		"$set": bson.M{
			"firstName":      basicInfo.FirstName,
			"lastName":       basicInfo.LastName,
			"additionalName": basicInfo.AdditionalName,
			"pronouns":       basicInfo.Pronouns,
			"headline":       basicInfo.Headline,
		},
	}

	var updated model.BasicInfo
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(ctx, bson.M{"id": id}, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("BasicInfo with ID %s not found", id)
		}
		return nil, err
	}

	return &updated, nil
}

// DeleteBasicInfo deletes a BasicInfo document by string ID
func (r *mutationResolver) DeleteBasicInfo(ctx context.Context, id string) (bool, error) {
	collection := r.DB.Collection(collectionName)

	res, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return false, err
	}

	return res.DeletedCount > 0, nil
}

// BasicInfo fetches a single BasicInfo document by string ID
func (r *queryResolver) BasicInfo(ctx context.Context, id string) (*model.BasicInfo, error) {
	collection := r.DB.Collection(collectionName)

	var basicInfo model.BasicInfo
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&basicInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("BasicInfo with ID %s not found", id)
		}
		return nil, err
	}

	return &basicInfo, nil
}

// BasicInfos fetches all BasicInfo documents
func (r *queryResolver) BasicInfos(ctx context.Context) ([]*model.BasicInfo, error) {
	collection := r.DB.Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var basicInfos []*model.BasicInfo
	for cursor.Next(ctx) {
		var basicInfo model.BasicInfo
		if err := cursor.Decode(&basicInfo); err != nil {
			return nil, err
		}
		basicInfos = append(basicInfos, &basicInfo)
	}

	return basicInfos, nil
}

// Mutation returns the MutationResolver implementation
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns the QueryResolver implementation
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
