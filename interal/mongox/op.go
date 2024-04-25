package mongox

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Insert(ctx context.Context, collection string, record []any) error {
	cl := instance.collection(collection)
	_, err := cl.InsertMany(ctx, record)
	return err
}

func UpdateOne(ctx context.Context, collection string, filter bson.M, update interface{}, upsert bool) error {
	cl := instance.collection(collection)
	opts := options.Update().SetUpsert(upsert)
	_, err := cl.UpdateOne(
		ctx,
		filter,
		bson.M{"$set": update},
		opts,
	)
	return err
}

func UpdateMany(ctx context.Context, collection string, filter bson.M, update interface{}, upsert bool) error {
	cl := instance.collection(collection)
	opts := options.Update().SetUpsert(upsert)
	_, err := cl.UpdateMany(
		ctx,
		filter,
		bson.M{"$set": update},
		opts,
	)
	return err
}

func FindOne(ctx context.Context, collection string, filter bson.M, record interface{}) error {
	cl := instance.collection(collection)
	err := cl.FindOne(ctx, filter).Decode(record)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	return nil
}

func FindOneWithOpts(ctx context.Context, collection string, filter bson.M, record interface{}, opts ...*options.FindOneOptions) error {
	cl := instance.collection(collection)
	err := cl.FindOne(ctx, filter, opts...).Decode(record)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	return nil
}

func Find(ctx context.Context, collection string, filter bson.M, record interface{}, opts ...*options.FindOptions) error {
	cl := instance.collection(collection)
	cur, err := cl.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	err = cur.All(ctx, record)
	if err != nil {
		return err
	}
	return nil
}

func DeleteOne(ctx context.Context, collection string, filter bson.M) (int64, error) {
	cl := instance.collection(collection)
	res, err := cl.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func DeleteMany(ctx context.Context, collection string, filter bson.M) (int64, error) {
	cl := instance.collection(collection)
	res, err := cl.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func Count(ctx context.Context, collection string, filter bson.M, opts ...*options.CountOptions) (n int64, err error) {
	cl := instance.collection(collection)
	n, err = cl.CountDocuments(ctx, filter, opts...)
	return
}
