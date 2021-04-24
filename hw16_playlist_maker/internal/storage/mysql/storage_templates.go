package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	domain "github.com/elak/golang_home_work/hw16_playlist_maker/internal/storage"
)

type TemplatesStorage struct {
	CommonStorage
	itemsStorage        TemplatesItemsStorage
	restrictionsStorage TemplateRestrictionsStorage
}

type TemplateRestrictionsStorage struct {
	CommonStorage
}

type TemplatesItemsStorage struct {
	CommonStorage
	restrictionsStorage TemplateRestrictionsStorage
	fillerStorage       TemplateFillerStorage
}

type TemplateFillerStorage struct {
	CommonStorage
}

func (t *TemplatesStorage) Init(s *Storage) {
	t.s = s
	t.table = "TEMPLATES"
	t.fields = append(t.fields, "ID", "Title")
	t.upsertQuery = s.generateUpsert(t.table, t.fields)
	t.selectQuery = s.generateSelect(t.table, t.fields, []string{"ID"})
	t.listQuery = s.generateList(t.table, t.fields)

	t.itemsStorage.Init(s)
	t.restrictionsStorage.Init(s)
}

func (t *TemplatesStorage) Create(ctx context.Context, item domain.Template) (err error) {
	return t.store(ctx, item)
}

func (t *TemplatesStorage) Update(ctx context.Context, item domain.Template) error {
	return t.store(ctx, item)
}

func (t *TemplatesStorage) Read(ctx context.Context, id domain.UUID) (*domain.Template, error) {
	tx, err := t.s.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		err = tx.Commit()
	}()

	return t.readImpl(ctx, tx, id)
}

func (t *TemplatesStorage) readImpl(ctx context.Context, tx *sql.Tx, id domain.UUID) (*domain.Template, error) {
	rows, err := tx.QueryContext(ctx, t.selectQuery, id)
	if err != nil {
		return nil, err
	}

	res, err := t.loadRows(ctx, tx, rows, nil)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

func (t *TemplatesStorage) Delete(ctx context.Context, id domain.UUID) error {
	tx, err := t.s.db.Begin()
	if err != nil {
		return err
	}

	commit := false
	defer func() {
		if commit {
			err = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	template, err := t.readImpl(ctx, tx, id)
	if err != nil {
		return err
	}

	for _, item := range template.StartItems {
		item.Restrictions = nil
		item.Fillers = nil
	}

	for _, item := range template.Items {
		item.Restrictions = nil
		item.Fillers = nil
	}

	for _, item := range template.EndItems {
		item.Restrictions = nil
		item.Fillers = nil
	}

	template.Restrictions = nil

	err = t.storeImpl(ctx, tx, *template)
	if err != nil {
		return err
	}

	qText := t.s.generatePurge(t.itemsStorage.table, "TemplateID", 0)
	_, err = tx.ExecContext(ctx, qText, id)
	if err != nil {
		return err
	}

	qText = fmt.Sprintf("DELETE FROM %s WHERE ID=?", t.table)
	_, err = tx.ExecContext(ctx, qText, id)
	if err != nil {
		return err
	}

	commit = true
	return nil
}

func (t *TemplatesStorage) List(ctx context.Context, filter func(*domain.Template) bool) (res []*domain.Template, err error) {
	tx, err := t.s.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		err = tx.Commit()
	}()

	rows, err := t.s.db.QueryContext(ctx, t.listQuery)
	if err != nil {
		return nil, err
	}

	res, err = t.loadRows(ctx, tx, rows, filter)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *TemplatesStorage) store(ctx context.Context, template domain.Template) (err error) {
	tx, err := t.s.db.Begin()
	if err != nil {
		return err
	}

	commit := false
	defer func() {
		if commit {
			err = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	err = t.storeImpl(ctx, tx, template)
	if err != nil {
		return err
	}

	commit = true
	return nil
}

func (t *TemplatesStorage) storeImpl(ctx context.Context, tx *sql.Tx, template domain.Template) (err error) {
	keepItems := make([]interface{}, 0)

	_, err = tx.ExecContext(ctx, t.upsertQuery, template.ID, template.Title)

	if err != nil {
		return err
	}

	stored, err := t.itemsStorage.store(ctx, tx, template.StartItems, template.ID, 1)
	if err != nil {
		return err
	}
	keepItems = append(keepItems, stored...)

	stored, err = t.itemsStorage.store(ctx, tx, template.Items, template.ID, 2)
	if err != nil {
		return err
	}
	keepItems = append(keepItems, stored...)

	stored, err = t.itemsStorage.store(ctx, tx, template.EndItems, template.ID, 3)
	if err != nil {
		return err
	}
	keepItems = append(keepItems, stored...)

	err = t.restrictionsStorage.store(ctx, tx, template.Restrictions, template.ID)
	if err != nil {
		return err
	}

	qText := t.s.generatePurge(t.itemsStorage.table, "TemplateID", len(keepItems))
	purgeParams := []interface{}{template.ID}
	purgeParams = append(purgeParams, keepItems...)
	_, err = tx.ExecContext(ctx, qText, purgeParams...)

	return err
}

func (t *TemplatesStorage) loadRows(ctx context.Context, tx *sql.Tx, rows *sql.Rows, filter func(*domain.Template) bool) (res []*domain.Template, err error) {
	res, err = t.walkRows(rows, filter)
	if err != nil {
		return nil, err
	}

	err = t.loadItems(ctx, tx, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *TemplatesStorage) loadItems(ctx context.Context, tx *sql.Tx, items []*domain.Template) (err error) {
	for _, item := range items {
		item.Restrictions, err = t.restrictionsStorage.list(ctx, tx, item.ID)
		if err != nil {
			return err
		}

		item.StartItems, err = t.itemsStorage.list(ctx, tx, item.ID, 1)
		if err != nil {
			return err
		}

		item.Items, err = t.itemsStorage.list(ctx, tx, item.ID, 2)
		if err != nil {
			return err
		}

		item.EndItems, err = t.itemsStorage.list(ctx, tx, item.ID, 3)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TemplatesStorage) walkRows(rows *sql.Rows, filter func(*domain.Template) bool) (res []*domain.Template, err error) {
	defer func() {
		err = rows.Close()
		_ = rows.Err()
	}()

	res = make([]*domain.Template, 0)

	for rows.Next() {
		var item domain.Template
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

func (ti *TemplatesItemsStorage) Init(s *Storage) {
	ti.s = s
	ti.table = "TEMPLATES_ITEMS"
	ti.fields = append(ti.fields, "ID", "Title", "`Order`", "Duration", "TemplateID", "TemplateBlock")
	ti.upsertQuery = s.generateUpsert(ti.table, ti.fields)
	ti.selectQuery = s.generateSelect(ti.table, ti.fields, []string{"ID"})
	ti.listQuery = s.generateSelect(ti.table, ti.fields, []string{"TemplateID", "TemplateBlock"})

	ti.restrictionsStorage.Init(s)
	ti.fillerStorage.Init(s)
}

func (ti *TemplatesItemsStorage) list(ctx context.Context, tx *sql.Tx, ownerID domain.UUID, itemsBlock int) ([]domain.TemplateItem, error) {
	rows, err := tx.QueryContext(ctx, ti.listQuery, ownerID, itemsBlock)
	if err != nil {
		return nil, err
	}

	res, err := ti.walkRows(rows)
	if err != nil {
		return nil, err
	}

	for _, item := range res {
		item.Restrictions, err = ti.restrictionsStorage.list(ctx, tx, item.ID)
		if err != nil {
			return nil, err
		}

		item.Fillers, err = ti.fillerStorage.list(ctx, tx, item.ID)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (ti *TemplatesItemsStorage) walkRows(rows *sql.Rows) (res []domain.TemplateItem, err error) {
	defer func() {
		err = rows.Close()
		_ = rows.Err()
	}()

	res = make([]domain.TemplateItem, 0)

	for rows.Next() {
		var item domain.TemplateItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Order, &item.Duration); err != nil {
			return nil, err
		}

		res = append(res, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ti *TemplatesItemsStorage) store(ctx context.Context, tx *sql.Tx, items []domain.TemplateItem, ownerID domain.UUID, itemsBlock int) ([]interface{}, error) {
	storedItems := make([]interface{}, 0, len(items))

	for _, item := range items {
		_, err := tx.ExecContext(ctx, ti.upsertQuery, item.ID, item.Title, item.Order, item.Duration, ownerID, itemsBlock)
		if err != nil {
			return nil, err
		}

		storedItems = append(storedItems, item.ID)

		err = ti.restrictionsStorage.store(ctx, tx, item.Restrictions, item.ID)
		if err != nil {
			return nil, err
		}

		err = ti.fillerStorage.store(ctx, tx, item.Fillers, item.ID)
		if err != nil {
			return nil, err
		}
	}
	return storedItems, nil
}

func (tfs *TemplateFillerStorage) Init(s *Storage) {
	tfs.s = s
	tfs.table = "TEMPLATES_FILLERS"
	tfs.fields = []string{"ID", "Order", "CategoryID", "AllowRepeat", "GroupsPriority", "VideosPriority", "OwnerID"}
	tfs.upsertQuery = s.generateUpsert(tfs.table, tfs.fields)
	tfs.selectQuery = s.generateSelect(tfs.table, tfs.fields, []string{"ID"})
	tfs.listQuery = s.generateSelect(tfs.table, tfs.fields, []string{"OwnerID"})
}

func (tfs *TemplateFillerStorage) list(ctx context.Context, tx *sql.Tx, ownerID domain.UUID) ([]domain.TemplateFiller, error) {
	rows, err := tx.QueryContext(ctx, tfs.listQuery, ownerID)
	if err != nil {
		return nil, err
	}

	return tfs.walkRows(rows)
}

func (tfs *TemplateFillerStorage) walkRows(rows *sql.Rows) (res []domain.TemplateFiller, err error) {
	defer func() {
		err = rows.Close()
		_ = rows.Err()
	}()

	res = make([]domain.TemplateFiller, 0)

	for rows.Next() {
		var item domain.TemplateFiller
		if err := rows.Scan(&item.ID, &item.Order, &item.Order, &item.CategoryID, &item.AllowRepeat, &item.GroupsPriority, &item.VideosPriority); err != nil {
			return nil, err
		}

		res = append(res, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (tfs *TemplateFillerStorage) store(ctx context.Context, tx *sql.Tx, items []domain.TemplateFiller, ownerID domain.UUID) error {
	storedItems := make([]interface{}, 0, len(items))

	for _, item := range items {
		_, err := tx.ExecContext(ctx, tfs.upsertQuery, item.ID, item.Order, item.CategoryID, item.AllowRepeat, item.GroupsPriority, item.VideosPriority, ownerID)
		if err != nil {
			return err
		}

		storedItems = append(storedItems, item.ID)
	}

	// очистим записи, больше не входящие в коллекцию
	qText := tfs.s.generatePurge(tfs.table, "ownerID", len(storedItems))
	purgeParams := []interface{}{ownerID}
	purgeParams = append(purgeParams, storedItems...)
	_, err := tx.ExecContext(ctx, qText, purgeParams...)
	if err != nil {
		return err
	}

	return nil
}

func (trs *TemplateRestrictionsStorage) Init(s *Storage) {
	trs.s = s
	trs.table = "RESTRICTIONS"
	trs.fields = append(trs.fields, "ID", "Title", "Scope", "CategoryID", "GroupID", "Duration", "Amount", "OwnerID")
	trs.upsertQuery = s.generateUpsert(trs.table, trs.fields)
	trs.selectQuery = s.generateSelect(trs.table, trs.fields, []string{"ID"})
	trs.listQuery = s.generateSelect(trs.table, trs.fields, []string{"OwnerID"})
}

func (trs *TemplateRestrictionsStorage) list(ctx context.Context, tx *sql.Tx, ownerID domain.UUID) ([]domain.TemplateRestriction, error) {
	rows, err := tx.QueryContext(ctx, trs.listQuery, ownerID)
	if err != nil {
		return nil, err
	}

	return trs.walkRows(rows)
}

func (trs *TemplateRestrictionsStorage) walkRows(rows *sql.Rows) (res []domain.TemplateRestriction, err error) {
	defer func() {
		err = rows.Close()
		_ = rows.Err()
	}()

	res = make([]domain.TemplateRestriction, 0)

	for rows.Next() {
		var item domain.TemplateRestriction
		if err := rows.Scan(&item.ID, &item.Title, &item.Scope, &item.CategoryID, &item.GroupID, &item.Duration, &item.Amount); err != nil {
			return nil, err
		}

		res = append(res, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (trs *TemplateRestrictionsStorage) store(ctx context.Context, tx *sql.Tx, items []domain.TemplateRestriction, ownerID domain.UUID) error {
	storedItems := make([]interface{}, 0, len(items))
	for _, item := range items {
		_, err := tx.ExecContext(ctx, trs.upsertQuery, item.ID, item.Title, item.Scope, item.CategoryID, item.GroupID, item.Duration, item.Amount, ownerID)
		if err != nil {
			return err
		}

		storedItems = append(storedItems, item.ID)
	}

	// очистим записи, больше не входящие в коллекцию
	qText := trs.s.generatePurge(trs.table, "ownerID", len(storedItems))
	purgeParams := []interface{}{ownerID}
	purgeParams = append(purgeParams, storedItems...)
	_, err := tx.ExecContext(ctx, qText, purgeParams...)
	if err != nil {
		return err
	}

	return nil
}
