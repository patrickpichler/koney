// Copyright (c) 2025 Dynatrace LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package v1alpha1

import (
	"fmt"
	"path/filepath"
)

// FilesystemHoneytoken defines the configuration for a filesystem honeytoken trap.
type FilesystemHoneytoken struct {
	// FilePath is the path of the file to be created.
	FilePath string `json:"filePath" yaml:"filePath"`

	// FileContent is the content of the file to be created.
	// +optional
	// +kubebuilder:default=""
	FileContent string `json:"fileContent" yaml:"fileContent"`

	// ReadOnly is a flag to make the file read-only.
	// +optional
	// +kubebuilder:default=true
	ReadOnly bool `json:"readOnly" yaml:"readOnly"`
}

// IsValid checks if the filesystem honeytoken trap is valid.
// The file path must be absolute.
func (f *FilesystemHoneytoken) IsValid() error {
	// Check if the file path is absolute
	if !filepath.IsAbs(f.FilePath) {
		return fmt.Errorf("FilePath is not absolute: '%s'", f.FilePath)
	}

	return nil
}
