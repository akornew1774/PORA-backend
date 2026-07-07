// entities - пакет с сущностями для БД
package entities

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Family - структура сущности семьи. Содержит
// название, ссылку на владельца, массив участников
// семьи и массив списков покупок.
type Family struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerID uuid.UUID `gorm:"type:uuid;not null"`

	InviteCode string `gorm:"size:10;uniqueIndex;not null"`

	Name string `gorm:"not null"`

	Members []FamilyMember `gorm:"foreignKey:FamilyID;constraint:OnDelete:CASCADE"`
	Lists   []List         `gorm:"foreignKey:FamilyID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"not null"`
}

// alphabet задает символы, которые могут быть в InviteCode
const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz2346789"

// GenerateInviteCode создает случайный новый InviteCode для семьи
func (f *Family) GenerateInviteCode() (string, error) {

	code := make([]byte, 8)

	for i := range code {
		n, err := rand.Int(rand.Reader,
			big.NewInt(int64(len(alphabet))),
		)
		if err != nil {
			return "", err
		}

		code[i] = alphabet[n.Int64()]
	}

	return string(code), nil
}

// BeforeCreate создает необходимые отсутствующие поля при создании сущности
func (f *Family) BeforeCreate(db *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now().UTC()
	}
	return nil
}
