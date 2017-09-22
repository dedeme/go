// Copyright 05-Sep-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

var excluded = map[string]map[string]struct{}{
	// basic ---------------------------------------------------------------------
	"/deme/dmjs17/libs/basic/src/all.js": map[string]struct{}{
		"$": {}, "$$": {},
	},
	"/deme/dmjs17/libs/basic/src/Ui.js": map[string]struct{}{
		"location": {}, "search": {}, "onload": {},
		"host": {}, "onreadystatechange": {}, "readyState": {}, "responseText": {},
		"body": {}, "dataTransfer": {}, "dropEffect": {}, "files": {},
		"keyCode": {}, "selectionStart": {}, "selectionEnd": {}, "value": {},
		"png": {}, "frequency": {}, "destination": {}, "scrollLeft": {},
		"documentElement": {}, "clientX": {}, "scrollTop": {}, "clientY": {},
	},
	"/deme/dmjs17/libs/basic/src/Store.js": map[string]struct{}{
		"localStorage": {}, "sessionStorage": {},
	},
	"/deme/dmjs17/libs/basic/src/Domo.js": map[string]struct{}{
		"innerHTML": {}, "textContent": {}, "className": {}, "disabled": {},
		"checked": {}, "value": {}, "firstChild": {}, "nextSibling": {},
	},
	"/deme/dmjs17/libs/basic/src/Client.js": map[string]struct{}{
		"onreadystatechange": {}, "readyState": {}, "responseText": {}, "host": {},
	},
	"/deme/dmjs17/libs/basic/src/B64.js": map[string]struct{}{
		"fromCharCode": {},
	},
	// DmSudoku ------------------------------------------------------------------
	"/deme/dmjs17/app/client/DmSudoku/src/Main.js": map[string]struct{}{
		"newSudoku": {}, "onmessage": {}, "document": {}, "keyCode": {}, "data": {},
	},
	"/deme/dmjs17/app/client/DmSudoku/src/View.js": map[string]struct{}{
		"copySudoku": {}, "readSudoku": {}, "saveSudoku": {}, "changeLang": {},
		"downLevel": {}, "upLevel": {}, "changeDevice": {}, "clearSudoku": {},
		"helpSudoku": {}, "solveSudoku": {}, "copyAccept": {}, "copyCancel": {},
		"loadCancel": {}, "solveAccept": {}, "newSudoku": {},
	},
	"/deme/dmjs17/app/client/DmSudoku/src/sudokuMaker.js": map[string]struct{}{
		"data": {},
	},
	// Try -----------------------------------------------------------------------
	"/deme/dmjs17/app/client/Try/src/main.js": map[string]struct{}{
		"Test": {},
	},
	"/deme/dmjs17/app/client/Try/src/test.js": map[string]struct{}{
		"Test": {},
	},
	// wallpapers ----------------------------------------------------------------
	"/deme/dmjs17/app/client/wallpapers/src/Control.js": map[string]struct{}{
		"onload": {}, "result": {},
	},
	"/deme/dmjs17/app/client/wallpapers/src/view/Menu.js": map[string]struct{}{
		"dataTransfer": {}, "files": {},
	},
	"/deme/dmjs17/app/client/wallpapers/src/view/Viewer.js": map[string]struct{}{
		"width": {}, "height": {}, "onclick": {}, "top": {},
		"clientY": {}, "clientX": {}, "left": {}, "data": {},
	},
	// JsDoc ---------------------------------------------------------------------
	"/deme/dmjs17/app/server/JsDoc/src/Module/ModuleV.js": map[string]struct{}{
		"vendor": {}, "hash": {}, "className": {}, "href": {}, "sortf": {},
	},
	"/deme/dmjs17/app/server/JsDoc/src/Auth/AuthV.js": map[string]struct{}{
		"changeLanguage": {},
	},
}
