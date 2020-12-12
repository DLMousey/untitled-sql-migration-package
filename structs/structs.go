package structs

type Projects struct {
	Projects []ProjectConfig `hcl:"project,block"`
}

type ProjectConfig struct {
	Name 	  	 string `hcl:"name,key"`
	Directory 	 string `hcl:"dir,label"`
	SqlDirectory string `hcl:"sqlDir,label"`
	Driver 		 string `hcl:"driver,label"`
	Host		 string `hcl:"host,label"`
	Username	 string `hcl:"username,label"`
	Password	 string `hcl:"password,label"`
	Database 	 string `hcl:"database,label"`
	Port		 int    `hcl:"port,label"`
}

type DatabaseState struct {
	Initialised bool
	Reachable 	bool
}
