package main

func main() {
	flags := NewCliOptions()
	envs, err := NewEnvConfig()
	if err != nil {
		ErrorLog.Fatal(err)
	}
	app := NewApp(NewConfig(flags, envs))
	app.Run()
}
