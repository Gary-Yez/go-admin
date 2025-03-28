package global

func Init() error {

	err := InitConfig()
	if err != nil {
		return err
	}
	err = InitDB()
	if err != nil {
		return err
	}
	return err
}
