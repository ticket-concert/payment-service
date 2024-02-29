package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"payment-service/internal/pkg/errors"
	wrapper "payment-service/internal/pkg/helpers"
	"payment-service/internal/pkg/log"
)

type MongoDBLogger struct {
	mongoClient *mongo.Client
	dbName      string
	logger      log.Logger
}

func NewMongoDBLogger(mongoClient *mongo.Client, dbName string, log log.Logger) Collections {
	return MongoDBLogger{
		mongoClient: mongoClient,
		dbName:      dbName,
		logger:      log,
	}
}

const (
	SortAscending  = `asc`
	SortDescending = `desc`
)

type Sort struct {
	FieldName string
	By        string
}

func (s Sort) buildSortBy() int {
	if s.By == SortDescending {
		return -1
	}

	return 1
}

type FindAllData struct {
	Result         interface{}
	CountData      *int64
	CollectionName string
	Filter         interface{}
	Sort           *Sort
	Page           int64
	Size           int64
}

func (f FindAllData) generateOptionSkip() *int64 {
	skipNumber := f.Size * (f.Page - 1)
	return &skipNumber
}

func (m MongoDBLogger) FindAllData(payload FindAllData, ctx context.Context) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		defer close(output)

		start := time.Now()

		collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)

		findOption := options.Find()

		if payload.Sort != nil {
			findOption.SetSort(bson.D{{Key: payload.Sort.FieldName, Value: payload.Sort.buildSortBy()}})
		}

		findOption.Limit = &payload.Size
		findOption.Skip = payload.generateOptionSkip()

		cursor, err := collection.Find(ctx, payload.Filter, findOption)

		if err != nil {
			msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError(msg),
			}
		}

		defer cursor.Close(ctx)

		err = cursor.All(ctx, payload.Result)

		if err != nil {
			msg := "cannot unmarshal result"
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError(msg),
			}
		}

		finish := time.Now()
		if finish.Sub(start).Seconds() > 10 {
			j, _ := json.Marshal(payload.Filter)
			msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		}

		// handle countdata
		if payload.CountData != nil {
			resp := <-m.CountData(CountData{
				CollectionName: payload.CollectionName,
				Result:         payload.CountData,
				Filter:         payload.Filter,
			}, ctx)

			if resp.Error != nil {
				output <- wrapper.Result{
					Error: errors.InternalServerError("Error Mongodb Connection"),
				}
			}
			output <- wrapper.Result{
				Data:  payload.Result,
				Count: resp.Count,
			}
		} else {
			output <- wrapper.Result{
				Data: payload.Result,
			}
		}

	}()
	return output
}

type FindOne struct {
	Result         interface{}
	CollectionName string
	Filter         interface{}
}

func (m MongoDBLogger) FindOne(payload FindOne, ctx context.Context) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		defer close(output)
		start := time.Now()

		collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)
		documentReturned := collection.FindOne(ctx, payload.Filter)

		if documentReturned.Err() != nil {
			if documentReturned.Err() == mongo.ErrNoDocuments {
				m.logger.Error(ctx, fmt.Sprintf("%v %v", "mongo-query-noDocuments", mongo.ErrNoDocuments.Error()), fmt.Sprintf("%+v", payload))
				output <- wrapper.Result{
					Data: nil,
				}
			}

			msg := fmt.Sprintf("Error Mongodb Connection %s", documentReturned.Err())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError(msg),
			}
		}

		if err := documentReturned.Decode(payload.Result); err != nil {
			msg := "cannot unmarshal result"
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError(msg),
			}
		}
		output <- wrapper.Result{
			Data: payload.Result,
		}
		finish := time.Now()

		if finish.Sub(start).Seconds() > 10 {
			j, _ := json.Marshal(payload.Filter)
			msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		}
	}()

	return output
}

type FindMany struct {
	Result         interface{}
	CollectionName string
	Filter         interface{}
	Sort           *Sort
}

func (m MongoDBLogger) FindMany(payload FindMany, ctx context.Context) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		defer close(output)
		start := time.Now()

		collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)
		findOption := options.Find()

		if payload.Sort != nil {
			findOption.SetSort(bson.D{{Key: payload.Sort.FieldName, Value: payload.Sort.buildSortBy()}})
		}

		cursor, err := collection.Find(ctx, payload.Filter, findOption)

		if err != nil {
			msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError(msg),
			}
		}

		defer cursor.Close(ctx)

		if err := cursor.All(ctx, payload.Result); err != nil {
			msg := "cannot unmarshal result"
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError(msg),
			}
		}
		output <- wrapper.Result{
			Data: payload.Result,
		}

		finish := time.Now()

		if finish.Sub(start).Seconds() > 10 {
			j, _ := json.Marshal(payload.Filter)
			msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		}
	}()

	return output
}

type CountData struct {
	Result         *int64
	CollectionName string
	Filter         interface{}
}

func (m MongoDBLogger) CountData(payload CountData, ctx context.Context) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		defer close(output)

		start := time.Now()

		collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)
		countDoc, err := collection.CountDocuments(ctx, payload.Filter)

		if err != nil {
			msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError(msg),
			}
		}

		if payload.Result != nil {
			output <- wrapper.Result{
				Count: countDoc,
			}
		}

		finish := time.Now()

		if finish.Sub(start).Seconds() > 10 {
			j, _ := json.Marshal(payload.Filter)
			msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		}
	}()

	return output
}

type UpdateOne struct {
	CollectionName string
	Filter         interface{}
	Document       interface{}
}

func (m MongoDBLogger) UpsertOne(payload UpdateOne, ctx context.Context) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		defer close(output)
		start := time.Now()

		wc := writeconcern.Majority()
		rc := readconcern.Snapshot()
		txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

		collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)
		pByte, err := bson.Marshal(payload.Document)
		if err != nil {
			msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb"),
			}
		}

		var update bson.M
		err = bson.Unmarshal(pByte, &update)
		if err != nil {
			msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb"),
			}
		}

		doc := bson.D{{Key: "$set", Value: update}}
		opts := options.Update().SetUpsert(true)

		callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
			// Important: You must pass sessCtx as the Context parameter to the operations for them to be executed in the
			// transaction.
			_, err = collection.UpdateOne(sessCtx, payload.Filter, doc, opts)

			if err != nil {
				msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
				m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
				return nil, errors.InternalServerError("Error mongodb connection")
			}
			return nil, nil
		}

		session, err := m.mongoClient.StartSession()
		if err != nil {
			msg := fmt.Sprintf("Error Mongodb Session : %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb session"),
			}
		}
		defer session.EndSession(context.Background())

		_, err = session.WithTransaction(ctx, callback, txnOpts)
		if err != nil {
			msg := fmt.Sprintf("Error Mongodb Transaction : %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb transaction"),
			}
		}

		finish := time.Now()

		if finish.Sub(start).Seconds() > 10 {
			j, _ := json.Marshal(payload.Filter)
			msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		}
	}()

	return output
}

func (m MongoDBLogger) Close(ctx context.Context) error {
	return m.mongoClient.Disconnect(ctx)
}

type InsertOne struct {
	CollectionName string
	Document       interface{}
}

func (m MongoDBLogger) InsertOne(payload InsertOne, ctx context.Context) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		defer close(output)
		start := time.Now()

		collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)

		_, err := collection.InsertOne(ctx, payload.Document)
		if err != nil {
			msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb connection"),
			}
		}

		finish := time.Now()

		if finish.Sub(start).Seconds() > 10 {
			j, _ := json.Marshal(payload.Document)
			msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		}

		output <- wrapper.Result{
			Data: "Success insert data",
		}
	}()

	return output
}

func (m MongoDBLogger) UpdateOne(payload UpdateOne, ctx context.Context) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		defer close(output)
		start := time.Now()

		collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)

		pByte, err := bson.Marshal(payload.Document)
		if err != nil {
			msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb"),
			}
		}

		var update bson.M
		err = bson.Unmarshal(pByte, &update)
		if err != nil {
			msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb"),
			}
		}

		doc := bson.D{{Key: "$set", Value: update}}
		_, err = collection.UpdateOne(ctx, payload.Filter, doc)

		if err != nil {
			msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb connection"),
			}
		}

		finish := time.Now()

		if finish.Sub(start).Seconds() > 10 {
			j, _ := json.Marshal(payload.Filter)
			msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		}

		output <- wrapper.Result{
			Data: "Success update data",
		}
	}()

	return output
}

type Aggregate struct {
	Result         interface{}
	CollectionName string
	Filter         interface{}
}

func (m MongoDBLogger) Aggregate(payload Aggregate, ctx context.Context) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		defer close(output)
		start := time.Now()

		collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName)
		cursor, err := collection.Aggregate(ctx, payload.Filter)

		if err != nil {
			msg := fmt.Sprintf("Error Mongodb Connection : %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb connection"),
			}
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, payload.Result); err != nil {
			msg := "cannot unmarshal result"
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError(msg),
			}
		}
		output <- wrapper.Result{
			Data: payload.Result,
		}

		finish := time.Now()

		if finish.Sub(start).Seconds() > 10 {
			j, _ := json.Marshal(payload.Filter)
			msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		}
	}()
	return output
}

type FindOneAndUpdate struct {
	CollectionName string
	Filter         interface{}
	Update         interface{}
	Result         interface{}
	Upsert         bool
}

// FindOneAndUpdate executes a findAndModify command to update at most one document in the collection and returns the document BEFORE or AFTER updating.
func (m MongoDBLogger) FindOneAndUpdate(payload FindOneAndUpdate, rd options.ReturnDocument, ctx context.Context) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		start := time.Now()

		wc := writeconcern.Majority()
		rc := readconcern.Snapshot()
		txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

		collection := m.mongoClient.Database(m.dbName).Collection(payload.CollectionName, options.Collection().SetReadPreference(readpref.Primary()))
		pByte, err := bson.Marshal(payload.Update)
		if err != nil {
			msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb connection"),
			}
		}

		var update bson.M
		err = bson.Unmarshal(pByte, &update)
		if err != nil {
			msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb connection"),
			}
		}

		callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
			// Important: You must pass sessCtx as the Context parameter to the operations for them to be executed in the
			// transaction.
			opts := options.FindOneAndUpdate().SetUpsert(payload.Upsert).SetReturnDocument(rd)
			res := collection.FindOneAndUpdate(ctx, payload.Filter, update, opts)
			if res.Err() != nil {
				msg := fmt.Sprintf("Error Mongodb: %s", err.Error())
				m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
				return nil, errors.InternalServerError("Error mongodb connection")
			}

			if err := res.Decode(payload.Result); err != nil {
				msg := "cannot unmarshal result: " + err.Error()
				m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
				return nil, errors.InternalServerError("Error mongodb connection")
			}
			return payload.Result, nil
		}

		session, err := m.mongoClient.StartSession()
		if err != nil {
			msg := fmt.Sprintf("Error Mongodb Session : %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb session"),
			}
		}
		defer session.EndSession(context.Background())

		result, err := session.WithTransaction(ctx, callback, txnOpts)
		if err != nil {
			msg := fmt.Sprintf("Error Mongodb Transaction : %s", err.Error())
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			output <- wrapper.Result{
				Error: errors.InternalServerError("Error mongodb transaction"),
			}
		}
		rs, _ := json.Marshal(result)
		fmt.Println("-----Result FindOneAndUpdate------")
		fmt.Println(string(rs))

		output <- wrapper.Result{
			Data: payload.Result,
		}

		finish := time.Now()

		if finish.Sub(start).Seconds() > 10 {
			j, _ := json.Marshal(payload.Filter)
			msg := fmt.Sprintf("slow query: %v second, query: %s", finish.Sub(start).Seconds(), string(j))
			m.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		}
	}()

	return output
}

// Collections is mongodb's collection of function
type Collections interface {
	FindAllData(payload FindAllData, ctx context.Context) <-chan wrapper.Result
	FindOne(payload FindOne, ctx context.Context) <-chan wrapper.Result
	FindMany(payload FindMany, ctx context.Context) <-chan wrapper.Result
	FindOneAndUpdate(payload FindOneAndUpdate, rd options.ReturnDocument, ctx context.Context) <-chan wrapper.Result
	CountData(payload CountData, ctx context.Context) <-chan wrapper.Result
	UpsertOne(payload UpdateOne, ctx context.Context) <-chan wrapper.Result
	InsertOne(payload InsertOne, ctx context.Context) <-chan wrapper.Result
	UpdateOne(payload UpdateOne, ctx context.Context) <-chan wrapper.Result
	Aggregate(payload Aggregate, ctx context.Context) <-chan wrapper.Result
	Close(ctx context.Context) error
}
