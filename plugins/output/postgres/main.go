package postgre

_ "github.com/lib/pq"

type config struct {
	ConnectionString string `yaml:"connection_string"`
	Table            string `yaml:"table"`
	Truncate         bool   `yaml:"truncate"`
}

func (c *config) getQuery() string {
	return fmt.Sprintf("select * from %s", c.Table)
}

func (c *config) buildQuery(args ...interface{}) (string, error) {
	var fieldCount int
	if len(args) > 0 {
		fieldCount = args[0].(int)
	} else {
		return "", errors.New("Dynamic field resolution not supported yet")
	}

	fields := ""

	for i := 1; i <= fieldCount; i++ {
		fields = fields + ":" + strconv.Itoa(i)
		if i != fieldCount {
			fields = fields + ","
		}
	}
	return "insert into " + c.Table + " values(" + fields + ")", nil
}

// Description ...
func Description() string {
	return "Read data from oracle database"
}

// PluginType ...
func PluginType() api.PluginType {
	return api.InputPlugin
}


// Run plugin
func Run(wg *sync.WaitGroup, conf []byte, wire *api.Wire) error {

	cfg := config{Truncate: true}
	if err := api.ParseConfig(conf, &cfg); err != nil {
		return err
	}

	db, err := sql.Open("postgres", cfg.ConnectionString)
	if err != nil {
		logging.Error(err)
		return err
	}
	if err := sqlhelper.Truncate(db, cfg.Table); err != nil {
		logging.Error(err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		db.Close()
		return err
	}

	var query string

	isOpen := struct {
		MetadataChannel bool
		DataChannel     bool
		FeedChannel     bool
	}{MetadataChannel: true, DataChannel: true, FeedChannel: true}

	for {
		select {
		case md, ok := <-wire.Metadata:
			if !ok {
				isOpen.MetadataChannel = false
			} else if md.Type == api.FieldsMetadata {
				query, err = cfg.BuildQuery(len(md.Content.([]api.Field)))
				if err != nil {
					return err
				}
			}
		case data, ok := <-wire.Data:
			if !ok {
				isOpen.DataChannel = false
			} else {
				if query == "" {
					query, err = cfg.BuildQuery(len(data.Content[0]))
					if err != nil {
						return nil
					}
				}
				for _, d := range data.Content {
					_, err = tx.Exec(query, d...)
					if err != nil {
						logging.Error(err)
						tx.Rollback()
						return err
					}
					wire.Feed <- api.NewWriteCountFeed(1)
				}
			}
		}
		if !isOpen.DataChannel {
			break
		}
	}

	tx.Commit()
	db.Close()
	wg.Done()
	return nil
}