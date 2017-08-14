//     Digota <http://digota.com> - eCommerce microservice
//     Copyright (C) 2017  Yaron Sumel <yaron@digota.com>. All Rights Reserved.
//
//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published
//     by the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.
//
//     You should have received a copy of the GNU Affero General Public License
//     along with this program.  If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"path/filepath"
	"testing"
)

func TestNewConfig(t *testing.T) {
	absValidPath, _ := filepath.Abs("testFiles/valid.yaml")
	if _, err := LoadConfig(absValidPath); err != nil {
		t.Fatal(err)
	}
	absNotFoundPath, _ := filepath.Abs("testFiles/notfound.yaml")
	if _, err := LoadConfig(absNotFoundPath); err == nil {
		t.Fatal(err)
	}
	absNotValidPath, _ := filepath.Abs("testFiles/invalid.yaml")
	if _, err := LoadConfig(absNotValidPath); err == nil {
		t.Fatal(err)
	}
}
