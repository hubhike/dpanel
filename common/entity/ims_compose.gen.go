// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package entity

const TableNameCompose = "ims_compose"

// Compose mapped from table <ims_compose>
type Compose struct {
	ID    int32  `gorm:"column:id;primaryKey" json:"id"`
	Title string `gorm:"column:title" json:"title"`
	Yaml  string `gorm:"column:yaml" json:"yaml"`
	Name  string `gorm:"column:name" json:"name"`
}

// TableName Compose's table name
func (*Compose) TableName() string {
	return TableNameCompose
}
