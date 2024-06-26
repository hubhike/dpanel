package logic

type ImageTemplate struct {
}

type envItem struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	BaseImage string `json:"baseImage"`
}

type env struct {
	Name string    `json:"name"`
	Env  []envItem `json:"env"`
	Ext  []string  `json:"ext"`
}

func (self ImageTemplate) GetSupportEnv() map[string]env {
	supportEnv := make(map[string]env)

	supportEnv[LangPhp] = env{
		Name: LangPhp,
		Env: []envItem{
			{
				Name:      "php-72",
				BaseImage: "dpanel/base-image:php-72|7",
			},
			{
				Name:      "php-74",
				BaseImage: "dpanel/base-image:php-74|7",
			},
			{
				Name:      "php-81",
				BaseImage: "dpanel/base-image:php-81|81",
			},
		},
		Ext: []string{
			"intl", "pecl-apcu", "imap", "pecl-mongodb", "pdo_pgsql",
		},
	}

	supportEnv[LangJava] = env{
		Name: LangJava,
		Env: []envItem{
			{
				Name:      "jdk8",
				Version:   "8",
				BaseImage: "dpanel/base-image:java-8",
			},
			{
				Name:      "jdk11",
				Version:   "11",
				BaseImage: "dpanel/base-image:java-11",
			},
			{
				Name:      "jdk12",
				Version:   "12",
				BaseImage: "dpanel/base-image:java-12",
			},
		},
	}

	supportEnv[LangGolang] = env{
		Name: LangGolang,
		Env: []envItem{
			{
				Name:      "go1.21",
				Version:   "1.21",
				BaseImage: "dpanel/base-image:go-1.21|1.21",
			},
		},
	}

	supportEnv[LangNode] = env{
		Name: LangNode,
		Env: []envItem{
			{
				Name:      "node12",
				Version:   "12",
				BaseImage: "dpanel/base-image:node-12",
			},
			{
				Name:      "node14",
				Version:   "14",
				BaseImage: "dpanel/base-image:node-14",
			},
			{
				Name:      "node16",
				Version:   "16",
				BaseImage: "dpanel/base-image:node-16",
			},
			{
				Name:      "node18",
				Version:   "18",
				BaseImage: "dpanel/base-image:node-18",
			},
		},
	}

	supportEnv[LangHtml] = env{
		Name: LangHtml,
		Env: []envItem{
			{
				Name:      "common",
				Version:   "1.0.0",
				BaseImage: "dpanel/base-image:html-common",
			},
		},
	}
	return supportEnv
}
