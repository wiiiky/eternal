package db

func Init(pgURL, mongoURL, mongoDBName string) error {
	if err := InitPG(pgURL); err != nil {
		return err
	} else if err := InitMongo(mongoURL, mongoDBName); err != nil {
		return err
	}
	return nil
}
