package main

func main() {
	flags := NewCliOptions()
	envs, err := NewEnvConfig()
	if err != nil {
		ErrorLog.Fatal(err)
	}
	ErrorLog.Println(envs)
	ErrorLog.Println(NewConfig(flags, envs))
	app := NewApp(NewConfig(flags, envs))
	app.Run()
}
