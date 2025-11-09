package text

import (
	"io"

	"text/template"
)

const (
	PrivateController = "private-controller"
	PublicController  = "public-controller"
	DAO               = "dao"
	Service           = "service"
)

var templates = map[string]string{
	PrivateController: privateController,
	PublicController:  publicController,
	DAO:               dao,
	Service:           service,
}

func ExecuteTemplate(wr io.Writer, name string, model any) {
	t, err := template.New(name).Parse(templates[name])
	if err != nil {
		panic(err)
	}
	if err = t.ExecuteTemplate(wr, name, model); err != nil {
		panic(err)
	}
}

var privateController = `package private

import (
	gctx "context"
	"net/http"

	"github.com/carsonkrueger/main/internal/builders"
	"github.com/carsonkrueger/main/internal/context"
)

type {{ .NameLower }} struct {
	*context.AppContext
}

func New{{ .Name }}(ctx *context.AppContext) *{{ .NameLower }} {
	return &{{ .NameLower }}{
		AppContext: ctx,
	}
}

func (r {{ .NameLower }}) Path() string {
	return "/{{ .NameLower }}"
}

func (r *{{ .NameLower }}) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().Register(http.MethodGet, "/", r.{{ .NameLower }}Get).Build(ctx)
	b.NewHandler().Register(http.MethodPost, "/", r.{{ .NameLower }}Post).Build(ctx)
	b.NewHandler().Register(http.MethodPut, "/", r.{{ .NameLower }}Put).Build(ctx)
	b.NewHandler().Register(http.MethodPatch, "/", r.{{ .NameLower }}Patch).Build(ctx)
	b.NewHandler().Register(http.MethodDelete, "/", r.{{ .NameLower }}Delete).Build(ctx)
}

func (r *{{ .NameLower }}) {{ .NameLower }}Get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "{{ .NameLower }}.{{ .NameLower }}Get")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Post(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "{{ .NameLower }}.{{ .NameLower }}Post")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Put(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "{{ .NameLower }}.{{ .NameLower }}Put")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Patch(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "{{ .NameLower }}.{{ .NameLower }}Patch")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Delete(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "{{ .NameLower }}.{{ .NameLower }}Delete")
	lgr.Info("Called")
}
`

var publicController = `package public

import (
	"net/http"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/go-chi/chi/v5"
)

type {{ .NameLower }} struct {
	*context.AppContext
}

func New{{ .Name }}(ctx *context.AppContext) *{{ .NameLower }} {
	return &{{ .NameLower }}{
		AppContext: ctx,
	}
}

func (r *{{ .NameLower }}) Path() string {
	return "/{{ .NameLower }}"
}

func (r *{{ .NameLower }}) PublicRoute(router chi.Router) {
	router.Get("/", r.{{ .NameLower }}Get)
	router.Post("/", r.{{ .NameLower }}Post)
	router.Put("/", r.{{ .NameLower }}Put)
	router.Patch("/", r.{{ .NameLower }}Patch)
	router.Delete("/", r.{{ .NameLower }}Delete)
}

func (r *{{ .NameLower }}) {{ .NameLower }}Get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "{{ .NameLower }}.{{ .NameLower }}Get")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Post(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "{{ .NameLower }}.{{ .NameLower }}Post")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Put(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "{{ .NameLower }}.{{ .NameLower }}Put")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Patch(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "{{ .NameLower }}.{{ .NameLower }}Patch")
	lgr.Info("Called")
}

func (r *{{ .NameLower }}) {{ .NameLower }}Delete(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lgr := context.GetLogger(ctx, "{{ .NameLower }}.{{ .NameLower }}Delete")
	lgr.Info("Called")
}
`

var dao = `package dao

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/pkg/db/{{ .Schema }}/model"
	"github.com/carsonkrueger/main/internal/database/gen/{{ .DB }}/{{ .Schema }}/table"
	"github.com/go-jet/jet/v2/postgres"
)

type {{ .Name }}PrimaryKey int64;

type {{ .NameLower }}DAO struct {
	db *sql.DB
	context.DAOBaseQueries[{{ .Name }}PrimaryKey, model.{{ .Name }}]
}

func New{{ .Name }}DAO(db *sql.DB) *{{ .NameLower }}DAO {
	dao := &{{ .NameLower }}DAO{
		db:              db,
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[{{ .Name }}PrimaryKey, model.{{ .Name }}](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *{{ .NameLower }}DAO) Table() context.PostgresTable {
	return table.{{ .Name }}
}

func (dao *{{ .NameLower }}DAO) InsertCols() postgres.ColumnList {
	return table.{{ .Name }}.AllColumns.Except(
		table.{{ .Name }}.ID,
		table.{{ .Name }}.CreatedAt,
		table.{{ .Name }}.UpdatedAt,
	)
}

func (dao *{{ .NameLower }}DAO) UpdateCols() postgres.ColumnList {
	return table.{{ .Name }}.AllColumns.Except(
		table.{{ .Name }}.ID,
		table.{{ .Name }}.CreatedAt,
	)
}

func (dao *{{ .NameLower }}DAO) AllCols() postgres.ColumnList {
	return table.{{ .Name }}.AllColumns
}

func (dao *{{ .NameLower }}DAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{}
}

func (dao *{{ .NameLower }}DAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{}
}

func (dao *{{ .NameLower }}DAO) PKMatch(pk {{ .Name }}PrimaryKey) postgres.BoolExpression {
	return table.{{ .Name }}.ID.EQ(postgres.Int64(int64(pk)))
}

func (dao *{{ .NameLower }}DAO) GetUpdatedAt(row *model.{{ .Name }}) *time.Time {
	return row.UpdatedAt
}
`

var service = `package services

import "github.com/carsonkrueger/main/internal/context"

type {{ .NameLower }}Service struct {
	*context.AppContext
}

func New{{ .Name }}Service(ctx *context.AppContext) *{{ .NameLower }}Service {
	return &{{ .NameLower }}Service{
		ctx,
	}
}
`
