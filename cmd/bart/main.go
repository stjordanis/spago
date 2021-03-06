// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/nlpodyssey/spago/cmd/bart/app"
)

func main() {
	if err := app.NewBartApp().Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
