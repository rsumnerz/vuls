/* Vuls - Vulnerability Scanner
Copyright (C) 2016  Future Architect, Inc. Japan.

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
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	c "github.com/future-architect/vuls/config"
	"github.com/future-architect/vuls/report"
	"github.com/future-architect/vuls/util"
	"github.com/google/subcommands"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/context"
)

// TuiCmd is Subcommand of host discovery mode
type TuiCmd struct {
	lang        string
	debugSQL    bool
	jsonBaseDir string
}

// Name return subcommand name
func (*TuiCmd) Name() string { return "tui" }

// Synopsis return synopsis
func (*TuiCmd) Synopsis() string { return "Run Tui view to anayze vulnerabilites" }

// Usage return usage
func (*TuiCmd) Usage() string {
	return `tui:
	tui [-results-dir=/path/to/results]

`
}

// SetFlags set flag
func (p *TuiCmd) SetFlags(f *flag.FlagSet) {
	//  f.StringVar(&p.lang, "lang", "en", "[en|ja]")
	f.BoolVar(&p.debugSQL, "debug-sql", false, "debug SQL")
	f.BoolVar(&p.debug, "debug", false, "debug mode")

	defaultLogDir := util.GetDefaultLogDir()
	f.StringVar(&p.logDir, "log-dir", defaultLogDir, "/path/to/log")

	defaultDBPath := os.Getenv("PWD") + "/vuls.sqlite3"
	f.StringVar(&p.dbpath, "dbpath", defaultDBPath,
		fmt.Sprintf("/path/to/sqlite3 (default: %s)", defaultDBPath))
}

// Execute execute
func (p *TuiCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	c.Conf.Lang = "en"

	// Setup Logger
	c.Conf.Debug = p.debug
	c.Conf.DebugSQL = p.debugSQL
	c.Conf.DBPath = p.dbpath

	historyID := ""
	if 0 < len(f.Args()) {
		if _, err := strconv.Atoi(f.Args()[0]); err != nil {
			log.Errorf("First Argument have to be scan_histores record ID: %s", err)
			return subcommands.ExitFailure
		}
		historyID = f.Args()[0]
	} else {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Errorf("Failed to read stdin: %s", err)
			return subcommands.ExitFailure
		}
		fields := strings.Fields(string(bytes))
		if 0 < len(fields) {
			historyID = fields[0]
		}
	}
	return report.RunTui(historyID)
}
