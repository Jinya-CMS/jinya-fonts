package database

import (
	"context"
	"github.com/DerKnerd/gorp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"jinya-fonts/config"
)

var dbMap *gorp.DbMap

func GetDbMap() *gorp.DbMap {
	return dbMap
}

func SetupDatabase() {
	if dbMap == nil {
		pool, err := pgxpool.New(context.Background(), config.LoadedConfiguration.PostgresUrl)
		if err != nil {
			panic(err)
		}

		conn := stdlib.OpenDBFromPool(pool)

		dialect := gorp.PostgresDialect{}

		dbMap = &gorp.DbMap{Db: conn, Dialect: dialect}

		dbMap.
			AddTableWithName(Webfont{}, "font")
		designer := dbMap.
			AddTableWithName(Designer{}, "designer")
		designer.
			SetKeys(false, "name", "font")

		dbMap.
			AddTableWithName(File{}, "file")
		dbMap.
			AddTableWithName(JinyaFontsSettings{}, "settings")

		err = dbMap.CreateTablesIfNotExists()
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(`
alter table designer
	drop constraint if exists designer_font_fkey;
`)
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(`
alter table designer
	add constraint designer_font_fkey foreign key (font) references font(name) on delete cascade;
`)
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(`
alter table file
	drop constraint if exists file_font_fkey;
`)
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(`
alter table file
	add constraint file_font_fkey foreign key (font) references font(name) on delete cascade;
`)
		if err != nil {
			panic(err)
		}
	}
}
