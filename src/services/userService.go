/*
    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package services

import (
    "github.com/jinzhu/gorm"

    "skaioskit/models"
)

type IUserService interface {
    CreateUser(models.User) models.User
    UpdateUser(models.User) models.User
    GetUser(string) (models.User, error)
    EnsureUserTable()
}

type UserService struct {
    db *gorm.DB
}
func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}
func (p *UserService) CreateUser(user models.User) models.User {
    p.db.Create(&user)
    return user
}
func (p *UserService) UpdateUser(user models.User) models.User {
    p.db.Save(&user)
    return user
}
func (p *UserService) GetUser(name string) (models.User, error) {
    var user models.User
    err := p.db.Where(&models.User{Name: name}).First(&user).Error
    return user, err
}
func (p *UserService) EnsureUserTable() {
    p.db.AutoMigrate(&models.User{})
}
func (p *UserService) EnsureUser(user models.User) {
    existing, err := p.GetUser(user.Name)
    if err != nil {
        p.CreateUser(user)
    } else {
        existing.Name = user.Name
        p.UpdateUser(existing)
    }
}
