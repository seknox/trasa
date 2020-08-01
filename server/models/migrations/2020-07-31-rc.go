package migrations

var migration = map[string]string{

	"Delete Browser ext table ": `
		drop table browser_ext;
	`,

	"Alter table browsers. Add extensions column ": `
		ALTER TABLE browsers ADD COLUMN extensions JSONB;
	`,
	"added trusted field in devices": `alter table devices ADD trusted bool;update devices set trusted=false where true;`,
}

// Add following to update migrations
// [+] org_id foreign key in browsers with delete cascade
//
