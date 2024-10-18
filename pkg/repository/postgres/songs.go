package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/leonideliseev/songLibraryCrud/models"
)

type SongsPostgres struct {
	db *gorm.DB
}

func NewSongsPostgres(db *gorm.DB) *SongsPostgres {
	return &SongsPostgres{
		db: db,
	}
}

func (d *SongsPostgres) GetAll(limit, offset int) ([]*models.Song, error) {
	var songs []*models.Song

	if err := d.db.Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
			return nil, err
	}

	return songs, nil
}

func (d *SongsPostgres) CreateSong(s models.Song) (*models.Song, error) {
	if err := d.db.Create(&s).Error; err != nil {
		return nil, err
	}

	return &s, nil
}

func (d *SongsPostgres) GetSong(group, name string) (*models.Song, error) {
	var song models.Song

    // Найдем существующую запись по уникальной паре "group" и "name"
    if err := d.db.Where("group = ? AND name = ?", group, name).First(&song).Error; err != nil {
        return nil, err
    }

	return &song, nil
}

func (d *SongsPostgres) DeleteSong(group, name string) error {
	var song models.Song
    
    // Найдем существующую запись по уникальной паре "group" и "name"
    if err := d.db.Where("group = ? AND name = ?", group, name).First(&song).Error; err != nil {
        return err
    }

    // Удаляем найденную запись
    if err := d.db.Delete(&song).Error; err != nil {
        return err
    }

    return nil
}

func (d *SongsPostgres) UpdateSong(group, name string, updatedData *models.Song) (*models.Song, error) {
	var song models.Song

    // Найдем существующую запись по уникальной паре "group" и "name"
    if err := d.db.Where("group = ? AND name = ?", group, name).First(&song).Error; err != nil {
        return nil, err
    }

	// Создаем карту для обновляемых полей
    updates := make(map[string]interface{})

    // Проверяем каждое поле и добавляем его в карту, если оно не пустое
	if updatedData.Group != "" {
		updates["group"] = updatedData.Group
	}
	if updatedData.Name != "" {
		updates["name"] = updatedData.Name
	}
    if updatedData.ReleaseDate != "" {
        updates["release_date"] = updatedData.ReleaseDate
    }
    if updatedData.Text != "" {
        updates["text"] = updatedData.Text
    }
    if updatedData.Link != "" {
        updates["link"] = updatedData.Link
    }

    // Обновляем только указанные поля
    if err := d.db.Model(&song).Updates(updates).Error; err != nil {
        return nil, err
    }

    return &song, nil
}