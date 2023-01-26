// Code generated by protoc-gen-gohttp. DO NOT EDIT.
// source: contracts/service.proto

package contracts

import (
	context "context"
	gin "github.com/gin-gonic/gin"
	easyjson "github.com/mailru/easyjson"
	contracts "prototodo/pkg/domain/contracts"
)

// Tasks
type TasksHTTPServer interface {
	// - Commands
	Create(context.Context, *contracts.CreateTaskCommand) (*contracts.TaskEvent, error)
	Delete(context.Context, *contracts.DeleteTaskCommand) (*contracts.TaskEvent, error)
	Update(context.Context, *contracts.UpdateTaskCommand) (*contracts.TaskEvent, error)
	// Update existing task state to progress
	Progress(context.Context, *contracts.ProgressTaskCommand) (*contracts.TaskEvent, error)
	// Update existing task to complete
	Complete(context.Context, *contracts.CompleteTaskCommand) (*contracts.TaskEvent, error)
	// Query for existing tasks
	ListQuery(context.Context, *contracts.ListTasksQuery) (*contracts.TaskEntityList, error)
}
type tasks struct {
	app TasksHTTPServer
}

// creates a new task
func (p *tasks) create(ctx *gin.Context) {
	body := contracts.CreateTaskCommand{}
	easyjson.UnmarshalFromReader(ctx.Request.Body, &body)
	res, err := p.app.Create(
		ctx,
		&body,
	)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
	easyjson.MarshalToHTTPResponseWriter(
		res,
		ctx.Writer,
	)
}

// deletes an existing task
func (p *tasks) delete(ctx *gin.Context) {
	body := contracts.DeleteTaskCommand{}
	easyjson.UnmarshalFromReader(ctx.Request.Body, &body)
	res, err := p.app.Delete(
		ctx,
		&body,
	)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
	easyjson.MarshalToHTTPResponseWriter(
		res,
		ctx.Writer,
	)
}

// updates an existing task
func (p *tasks) update(ctx *gin.Context) {
	body := contracts.UpdateTaskCommand{}
	easyjson.UnmarshalFromReader(ctx.Request.Body, &body)
	res, err := p.app.Update(
		ctx,
		&body,
	)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
	easyjson.MarshalToHTTPResponseWriter(
		res,
		ctx.Writer,
	)
}

// update state of existing task to progress
func (p *tasks) progress(ctx *gin.Context) {
	body := contracts.ProgressTaskCommand{}
	easyjson.UnmarshalFromReader(ctx.Request.Body, &body)
	res, err := p.app.Progress(
		ctx,
		&body,
	)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
	easyjson.MarshalToHTTPResponseWriter(
		res,
		ctx.Writer,
	)
}

// update state of existing task to complete
func (p *tasks) complete(ctx *gin.Context) {
	body := contracts.CompleteTaskCommand{}
	easyjson.UnmarshalFromReader(ctx.Request.Body, &body)
	res, err := p.app.Complete(
		ctx,
		&body,
	)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
	easyjson.MarshalToHTTPResponseWriter(
		res,
		ctx.Writer,
	)
}

// query all existing tasks
func (p *tasks) listQuery(ctx *gin.Context) {
	body := contracts.ListTasksQuery{}
	easyjson.UnmarshalFromReader(ctx.Request.Body, &body)
	res, err := p.app.ListQuery(
		ctx,
		&body,
	)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
	easyjson.MarshalToHTTPResponseWriter(
		res,
		ctx.Writer,
	)
}
func RegisterTasksHTTPServer(
	grp *gin.RouterGroup,
	srv TasksHTTPServer,
) {
	ctrl := tasks{app: srv}
	grp.POST("/commands/createTask", ctrl.create)
	grp.POST("/commands/deleteTask", ctrl.delete)
	grp.POST("/commands/updateTask", ctrl.update)
	grp.POST("/commands/progressTask", ctrl.progress)
	grp.POST("/commands/completeTask", ctrl.complete)
	grp.POST("/queries/listTasks", ctrl.listQuery)
}

// Quotes
type QuotesHTTPServer interface {
	// Get a quote
	Get(context.Context, *contracts.GetQuoteQuery) (*contracts.QuoteData, error)
}
type quotes struct {
	app QuotesHTTPServer
}

// get a random quote
func (p *quotes) get(ctx *gin.Context) {
	body := contracts.GetQuoteQuery{}
	easyjson.UnmarshalFromReader(ctx.Request.Body, &body)
	res, err := p.app.Get(
		ctx,
		&body,
	)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(200)
	easyjson.MarshalToHTTPResponseWriter(
		res,
		ctx.Writer,
	)
}
func RegisterQuotesHTTPServer(
	grp *gin.RouterGroup,
	srv QuotesHTTPServer,
) {
	ctrl := quotes{app: srv}
	grp.POST("/queries/getQuote", ctrl.get)
}
