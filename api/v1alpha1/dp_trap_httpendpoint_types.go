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

import "errors"

// HttpEndpoint defines the configuration for an HTTP endpoint trap.
type HttpEndpoint struct {
	// TODO: Implement.
}

// IsValid checks if the HTTP endpoint trap is valid.
func (f *HttpEndpoint) IsValid() error {
	// TODO: Implement.
	return errors.New("HttpEndpoint validation not implemented yet")
}
