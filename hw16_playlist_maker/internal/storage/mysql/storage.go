package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	domain "github.com/elak/golang_home_work/hw16_playlist_maker/internal/storage"
	_ "github.com/go-sql-driver/mysql"
)

//nolint:structcheck
type CommonStorage struct {
	s           *Storage
	table       string
	fields      []string
	upsertQuery string
	selectQuery string
	listQuery   string
}

type GroupsStorage struct {
	CommonStorage
}

type VideosStorage struct {
	CommonStorage
}

type CategoriesStorage struct {
	CommonStorage
}

type HistoryStorage struct {
	CommonStorage
}

type Storage struct {
	db         *sql.DB
	groups     GroupsStorage
	videos     VideosStorage
	categories CategoriesStorage
	templates  TemplatesStorage
	history    HistoryStorage
}

func (s *Storage) Groups() domain.GroupsManager {
	return &s.groups
}

func (s *Storage) Videos() domain.VideosManager {
	return &s.videos
}

func (s *Storage) Categories() domain.CategoriesManager {
	return &s.categories
}

func (s *Storage) Templates() domain.TemplatesManager {
	return &s.templates
}

func (s *Storage) History() domain.HistoryManager {
	return &s.history
}

func New() *Storage {
	stor := &Storage{}

	stor.groups.Init(stor)
	stor.videos.Init(stor)
	stor.categories.Init(stor)
	stor.templates.Init(stor)
	stor.history.Init(stor)

	return stor
}

func (s *Storage) Connect(ctx context.Context, storageURI string) (err error) {
	s.db, err = sql.Open("mysql", storageURI)
	if err != nil {
		return err
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) deleteImpl(ctx context.Context, table string, key domain.UUID) error {
	qText := fmt.Sprintf("DELETE FROM %s WHERE ID=?", table)
	_, err := s.db.ExecContext(ctx, qText, key)
	return err
}

func (s *Storage) generatePurge(table, ownerField string, palaceHolders int) string {
	if palaceHolders == 0 {
		return fmt.Sprintf("DELETE FROM %s WHERE %s=?", table, ownerField)
	}

	allPlaceHolders := strings.Repeat(", ?", palaceHolders)[1:]
	return fmt.Sprintf("DELETE FROM %s WHERE %s=? AND ID NOT IN (%s)", table, ownerField, allPlaceHolders)
}

func (s *Storage) generatePurgeList(table, ownerField string, palaceHolders int) string {
	if palaceHolders == 0 {
		return fmt.Sprintf("SELECT FROM %s WHERE %s=?", table, ownerField)
	}
	allPlaceHolders := strings.Repeat(", ?", palaceHolders)[1:]
	return fmt.Sprintf("SELECT FROM %s WHERE %s=? AND ID NOT IN (%s)", table, ownerField, allPlaceHolders)
}

func (s *Storage) generateUpsert(table string, fields []string) string {
	allFields := strings.Join(fields, ", ")
	allPlaceHolders := strings.Repeat(", ?", len(fields))[1:]
	qText := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s) on duplicate KEY update ", table, allFields, allPlaceHolders)

	for i := 1; i < len(fields); i++ {
		if i > 1 {
			qText += ", "
		}
		qText += fields[i] + "=VALUES(" + fields[i] + ")"
	}

	return qText
}

func (s *Storage) generateSelect(table string, selectFields []string, conditionFields []string) string {
	allFields := strings.Join(selectFields, ", ")
	allConditions := ""

	if len(conditionFields) > 0 {
		allConditions = "WHERE "
		for i, field := range conditionFields {
			if i != 0 {
				allConditions += " AND "
			}
			allConditions += fmt.Sprintf("(%s=?)", field)
		}
	}

	return fmt.Sprintf("SELECT %s FROM %s %s", allFields, table, allConditions)
}

func (s *Storage) generateList(table string, fields []string) string {
	return s.generateSelect(table, fields, []string{})
}

func (g *GroupsStorage) Init(s *Storage) {
	g.s = s
	g.table = "`GROUPS`"
	g.fields = append(g.fields, "ID", "Title", "ParentID", "`Order`", "CategoryID")
	g.upsertQuery = s.generateUpsert(g.table, g.fields)
	g.selectQuery = s.generateSelect(g.table, g.fields, []string{"ID"})
	g.listQuery = s.generateList(g.table, g.fields)
}

func (g *GroupsStorage) Create(ctx context.Context, item domain.Group) error {
	_, err := g.s.db.ExecContext(ctx, g.upsertQuery, item.ID, item.Title, item.ParentID, item.Order, item.CategoryID)

	return err
}

func (g *GroupsStorage) Update(ctx context.Context, item domain.Group) error {
	_, err := g.s.db.ExecContext(ctx, g.upsertQuery, item.ID, item.Title, item.ParentID, item.Order, item.CategoryID)

	return err
}

func (g *GroupsStorage) Read(ctx context.Context, id domain.UUID) (*domain.Group, error) {
	rows, err := g.s.db.QueryContext(ctx, g.selectQuery, id)
	if err != nil {
		return nil, err
	}

	res, err := g.walkRows(rows, nil)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

func (g *GroupsStorage) Delete(ctx context.Context, id domain.UUID) error {
	return g.s.deleteImpl(ctx, g.table, id)
}

func (g *GroupsStorage) List(ctx context.Context, filter func(*domain.Group) bool) (res []*domain.Group, err error) {
	rows, err := g.s.db.QueryContext(ctx, g.listQuery)
	if err != nil {
		return nil, err
	}
	return g.walkRows(rows, filter)
}

func (g *GroupsStorage) walkRows(rows *sql.Rows, filter func(*domain.Group) bool) (res []*domain.Group, err error) {
	defer func() {
		err = rows.Close()
		_ = rows.Err()
	}()

	res = make([]*domain.Group, 0)

	for rows.Next() {
		var item domain.Group
		if err := rows.Scan(&item.ID, &item.Title, &item.ParentID, &item.Order, &item.CategoryID); err != nil {
			return nil, err
		}

		if filter == nil || filter(&item) {
			res = append(res, &item)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (v *VideosStorage) Init(s *Storage) {
	v.s = s
	v.table = "VIDEOS"
	v.fields = append(v.fields, "ID", "Title", "ParentID", "`Order`", "CategoryID", "Duration")
	v.upsertQuery = s.generateUpsert(v.table, v.fields)
	v.selectQuery = s.generateSelect(v.table, v.fields, []string{"ID"})
	v.listQuery = s.generateList(v.table, v.fields)
}

func (v *VideosStorage) Create(ctx context.Context, item domain.Video) error {
	_, err := v.s.db.ExecContext(ctx, v.upsertQuery, item.ID, item.Title, item.ParentID, item.Order, item.CategoryID, item.Duration)

	return err
}

func (v *VideosStorage) Update(ctx context.Context, item domain.Video) error {
	_, err := v.s.db.ExecContext(ctx, v.upsertQuery, item.ID, item.Title, item.ParentID, item.Order, item.CategoryID, item.Duration)

	return err
}

func (v *VideosStorage) Read(ctx context.Context, id domain.UUID) (*domain.Video, error) {
	rows, err := v.s.db.QueryContext(ctx, v.selectQuery, id)
	if err != nil {
		return nil, err
	}

	res, err := v.walkRows(rows, nil)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

func (v *VideosStorage) Delete(ctx context.Context, id domain.UUID) error {
	return v.s.deleteImpl(ctx, v.table, id)
}

func (v *VideosStorage) List(ctx context.Context, filter func(*domain.Video) bool) ([]*domain.Video, error) {
	rows, err := v.s.db.QueryContext(ctx, v.listQuery)
	if err != nil {
		return nil, err
	}
	return v.walkRows(rows, filter)
}

func (v *VideosStorage) walkRows(rows *sql.Rows, filter func(*domain.Video) bool) (res []*domain.Video, err error) {
	defer func() {
		err = rows.Close()
		_ = rows.Err()
	}()

	res = make([]*domain.Video, 0)

	for rows.Next() {
		var item domain.Video
		if err := rows.Scan(&item.ID, &item.Title, &item.ParentID, &item.Order, &item.CategoryID, &item.Duration); err != nil {
			return nil, err
		}

		if filter == nil || filter(&item) {
			res = append(res, &item)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *CategoriesStorage) Init(s *Storage) {
	c.s = s
	c.table = "CATEGORIES"
	c.fields = append(c.fields, "ID", "Title")
	c.upsertQuery = s.generateUpsert(c.table, c.fields)
	c.selectQuery = s.generateSelect(c.table, c.fields, []string{"ID"})
	c.listQuery = s.generateList(c.table, c.fields)
}

func (c *CategoriesStorage) Create(ctx context.Context, item domain.Category) error {
	_, err := c.s.db.ExecContext(ctx, c.upsertQuery, item.ID, item.Title)

	return err
}

func (c *CategoriesStorage) Update(ctx context.Context, item domain.Category) error {
	_, err := c.s.db.ExecContext(ctx, c.upsertQuery, item.ID, item.Title)

	return err
}

func (c *CategoriesStorage) Read(ctx context.Context, id domain.UUID) (*domain.Category, error) {
	rows, err := c.s.db.QueryContext(ctx, c.selectQuery, id)
	if err != nil {
		return nil, err
	}

	res, err := c.walkRows(rows, nil)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

func (c *CategoriesStorage) Delete(ctx context.Context, id domain.UUID) error {
	return c.s.deleteImpl(ctx, c.table, id)
}

func (c *CategoriesStorage) List(ctx context.Context, filter func(*domain.Category) bool) ([]*domain.Category, error) {
	rows, err := c.s.db.QueryContext(ctx, c.listQuery)
	if err != nil {
		return nil, err
	}
	return c.walkRows(rows, filter)
}

func (c *CategoriesStorage) walkRows(rows *sql.Rows, filter func(*domain.Category) bool) (res []*domain.Category, err error) {
	defer func() {
		err = rows.Close()
		_ = rows.Err()
	}()

	res = make([]*domain.Category, 0)

	for rows.Next() {
		var item domain.Category
		if err := rows.Scan(&item.ID, &item.Title); err != nil {
			return nil, err
		}

		if filter == nil || filter(&item) {
			res = append(res, &item)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (v *HistoryStorage) Init(s *Storage) {
	v.s = s
	v.table = "HISTORY"
	v.fields = append(v.fields, "VideoID", "LastSeen")
	v.upsertQuery = s.generateUpsert(v.table, v.fields)
	v.selectQuery = s.generateSelect(v.table, v.fields, []string{"VideoID"})
	v.listQuery = s.generateList(v.table, v.fields)
}

func (v *HistoryStorage) Create(ctx context.Context, item domain.History) error {
	_, err := v.s.db.ExecContext(ctx, v.upsertQuery, item.VideoID, item.LastSeen)

	return err
}

func (v *HistoryStorage) Update(ctx context.Context, item domain.History) error {
	_, err := v.s.db.ExecContext(ctx, v.upsertQuery, item.VideoID, item.LastSeen)

	return err
}

func (v *HistoryStorage) Read(ctx context.Context, id domain.UUID) (*domain.History, error) {
	rows, err := v.s.db.QueryContext(ctx, v.selectQuery, id)
	if err != nil {
		return nil, err
	}

	res, err := v.walkRows(rows, nil)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

func (v *HistoryStorage) Delete(ctx context.Context, id domain.UUID) error {
	return v.s.deleteImpl(ctx, v.table, id)
}

func (v *HistoryStorage) List(ctx context.Context, filter func(*domain.History) bool) ([]*domain.History, error) {
	rows, err := v.s.db.QueryContext(ctx, v.listQuery)
	if err != nil {
		return nil, err
	}
	return v.walkRows(rows, filter)
}

func (v *HistoryStorage) walkRows(rows *sql.Rows, filter func(*domain.History) bool) (res []*domain.History, err error) {
	defer func() {
		err = rows.Close()
		_ = rows.Err()
	}()

	res = make([]*domain.History, 0)

	for rows.Next() {
		var item domain.History
		if err := rows.Scan(&item.VideoID, &item.LastSeen); err != nil {
			return nil, err
		}

		if filter == nil || filter(&item) {
			res = append(res, &item)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}
