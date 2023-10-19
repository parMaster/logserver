package web

import (
	_ "embed"
)

//go:embed view.html
var View_html string

//go:embed chart_tpl.min.js
var Chart_tpl_min_js string
