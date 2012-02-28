/*
Copyright 2011-2012 Paul Ruane.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package commands

import (
    "fmt"
    "path/filepath"
    "tmsu/common"
    "tmsu/database"
)

type RepairCommand struct{}

func (RepairCommand) Name() string {
	return "repair"
}

func (RepairCommand) Synopsis() string {
	return "Repair stale data caused by file moves and amendments"
}

func (RepairCommand) Description() string {
	return `tmsu repair

Updates the database with respect to changed and moved files.`
}

func (command RepairCommand) Exec(args []string) error {
    db, err := database.OpenDatabase()
    if err != nil {
        return err
    }
    defer db.Close()

    for _, path := range args {
        absPath, err := filepath.Abs(path)
        if err != nil {
            return err
        }

        files, err := db.FilesByDirectory(absPath)
        if err != nil {
            return err
        }

        for _, file := range files {
            fingerprint, err := common.Fingerprint(file.Path())
            if err != nil {
                //TODO detect missing files
                common.Warnf("Could not fingerprint '%v': %v", file.Path(), err)
            }

            if (file.Fingerprint != fingerprint) {
                fmt.Printf("File '%v' has changed.\n", file.Path())
            }
        }
    }

    //TODO check fingerprints of existing database entries -- update if necessary
    //TODO identify missing files
    //TODO find files with same fingerprints -- hook them up if only one
    return nil
}
