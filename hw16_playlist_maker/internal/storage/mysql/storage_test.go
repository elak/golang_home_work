package sqlstorage

import (
	"context"
	"errors"
	"os"
	"testing"

	domain "github.com/elak/golang_home_work/hw16_playlist_maker/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategories(t *testing.T) {
	ctx := context.TODO()

	stor, err := getStorage(t)
	require.NoError(t, err)

	defer stor.Close(ctx)

	cat := domain.Category{
		ID:    "1",
		Title: "1st",
	}

	categories := stor.Categories()
	err = categories.Create(ctx, cat)
	require.NoError(t, err)

	cat2, err := categories.Read(ctx, cat.ID)
	require.NoError(t, err)
	assert.Equal(t, cat, *cat2)
	require.NotNil(t, cat2)

	cat2.Title += ", 2nd"
	err = categories.Update(ctx, *cat2)
	require.NoError(t, err)

	cat3, err := categories.Read(ctx, cat.ID)
	require.NoError(t, err)
	assert.Equal(t, cat2, cat3)
	assert.NotEqual(t, cat, *cat3)

	cat2.ID += "-1"
	err = categories.Create(ctx, *cat2)
	require.NoError(t, err)

	cts, err := categories.List(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, 2, len(cts))
	assert.Equal(t, cts[0], cat3)

	err = categories.Delete(ctx, cat3.ID)
	require.NoError(t, err)

	cts, err = categories.List(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, len(cts))
	assert.NotEqual(t, cat3.ID, cts[0].ID)
}

func TestGroups(t *testing.T) {
	ctx := context.TODO()

	stor, err := getStorage(t)

	require.NoError(t, err)

	defer stor.Close(ctx)

	group := domain.Group{
		ID:         "1",
		Title:      "1st",
		ParentID:   "",
		Order:      0,
		CategoryID: "1-1",
	}

	groups := stor.Groups()
	err = groups.Create(ctx, group)
	require.NoError(t, err)

	group2, err := groups.Read(ctx, group.ID)
	require.NoError(t, err)
	assert.Equal(t, group, *group2)
	require.NotNil(t, group2)

	group2.Title += ", 2nd"
	err = groups.Update(ctx, *group2)
	require.NoError(t, err)

	group3, err := groups.Read(ctx, group.ID)
	require.NoError(t, err)
	assert.Equal(t, group2, group3)
	assert.NotEqual(t, group, *group3)

	group2.ID += "-1"
	err = groups.Create(ctx, *group2)
	require.NoError(t, err)

	cts, err := groups.List(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, 2, len(cts))
	assert.Equal(t, cts[0], group3)

	err = groups.Delete(ctx, group3.ID)
	require.NoError(t, err)

	cts, err = groups.List(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, len(cts))
	assert.NotEqual(t, group3.ID, cts[0].ID)
}

func getStorage(t *testing.T) (*Storage, error) {
	storageURI := os.Getenv("MYSQL_URI")
	if storageURI == "" {
		t.Skipf("MYSQL_URI environment variable not set")
		return nil, errors.New("MYSQL_URI environment variable not set")
	}

	stor := New()
	err := stor.Connect(context.TODO(), storageURI)
	return stor, err
}

func TestVideos(t *testing.T) {
	ctx := context.TODO()

	stor, err := getStorage(t)
	require.NoError(t, err)

	defer stor.Close(ctx)

	Video := domain.Video{
		ID:         "1",
		Title:      "1st",
		ParentID:   "",
		Order:      0,
		CategoryID: "1-1",
		Duration:   317,
	}

	Videos := stor.Videos()
	err = Videos.Create(ctx, Video)
	require.NoError(t, err)

	Video2, err := Videos.Read(ctx, Video.ID)
	require.NoError(t, err)
	assert.Equal(t, Video, *Video2)
	require.NotNil(t, Video2)

	Video2.Title += ", 2nd"
	err = Videos.Update(ctx, *Video2)
	require.NoError(t, err)

	Video3, err := Videos.Read(ctx, Video.ID)
	require.NoError(t, err)
	assert.Equal(t, Video2, Video3)
	assert.NotEqual(t, Video, *Video3)

	Video2.ID += "-1"
	err = Videos.Create(ctx, *Video2)
	require.NoError(t, err)

	cts, err := Videos.List(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, 2, len(cts))
	assert.Equal(t, cts[0], Video3)

	err = Videos.Delete(ctx, Video3.ID)
	require.NoError(t, err)

	cts, err = Videos.List(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, len(cts))
	assert.NotEqual(t, Video3.ID, cts[0].ID)
}

func TestTemplates(t *testing.T) {
	ctx := context.TODO()

	stor, err := getStorage(t)
	require.NoError(t, err)

	defer stor.Close(ctx)

	Template := domain.Template{
		ID:           "1",
		Title:        "1st",
		StartItems:   nil,
		Items:        nil,
		EndItems:     nil,
		Restrictions: nil,
	}

	Templates := stor.Templates()
	err = Templates.Create(ctx, Template)
	require.NoError(t, err)

	Template2, err := Templates.Read(ctx, Template.ID)
	require.NoError(t, err)
	require.NotNil(t, Template2)
	assert.Equal(t, Template.ID, Template2.ID)
	assert.Equal(t, Template.Title, Template2.Title)

	Template2.Title += ", 2nd"
	err = Templates.Update(ctx, *Template2)
	require.NoError(t, err)

	Template3, err := Templates.Read(ctx, Template.ID)
	require.NoError(t, err)
	assert.Equal(t, Template2, Template3)
	assert.NotEqual(t, Template, *Template3)

	Template2.ID += "-1"
	err = Templates.Create(ctx, *Template2)
	require.NoError(t, err)

	cts, err := Templates.List(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, 2, len(cts))
	assert.Equal(t, cts[0], Template3)

	err = Templates.Delete(ctx, Template3.ID)
	require.NoError(t, err)

	cts, err = Templates.List(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, len(cts))
	assert.NotEqual(t, Template3.ID, cts[0].ID)
}
