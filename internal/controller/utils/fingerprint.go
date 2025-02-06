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

package utils

import "fmt"

// TODO: Randomize on startup and sync with alerting system
const KoneyFingerprint = 1337

// EncodeFingerprintInEcho encodes a fingerprint in a call to `echo`, to be
// used, e.g. in a call such as `echo -e "foobar\c KONEY_FINGERPRINT_123"` after
// the `\c` escape sequence. Everything after the `\c` escape sequence is
// ignored by `echo` when `-e` is used, thus, the command will still work as
// expected. This is useful to mark `echo` calls from Koney, so that we won't
// alert on them later.
func EncodeFingerprintInEcho(code int) string {
	return fmt.Sprintf("KONEY_FINGERPRINT_%d", code)
}

// EncodeFingerprintInCat encodes a fingerprint in a call to `cat`, to be used,
// e.g. in a call such as `cat -u -uu -u -u -uu /foo/bar` where the `-u` flag is
// used to binary-encode the fingerprint (`-u` is 0, `-uu` is 1). The `-u` flag
// is ignored by `cat`, thus, the command will still work as expected. This is
// useful to mark `cat` calls from Koney, so that we won't alert on them later.
func EncodeFingerprintInCat(code int) string {
	binaryCode := fmt.Sprintf("%b", code)

	// iterate and replace 0 with -u and 1 with -uu
	cipher := ""
	for _, bit := range binaryCode {
		if bit == '0' {
			cipher += "-u "
		} else {
			cipher += "-uu "
		}
	}

	// remove the trailing space
	cipher = cipher[:len(cipher)-1]

	return cipher
}
