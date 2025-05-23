# По поводу конфигов
## Можно в коде использовать константы
```GO
const post_numbers = 4
const address_numbers = 10
const recipient_numbers = 2
```

## Но лучше применять конфиг-файлы
	Вынесение в конфигурационный файл: Это более гибкий подход, который позволяет изменять значения без перекомпиляции.

* **ZB: Если количество постов (post_numbers) может меняться или надо иметь возможность настраивать его без перекомпиляции, лучше вынести это значение в конфигурационный файл.**

* **Можно добавить это значение в существующий config_go.json или создать отдельный конфигурационный файл для генератора, например, config_generator.json:**
```json
{
  "post_number": 4,
  "generation_interval_min": 1,
  "generation_interval_max": 5
}
```
* **Затем в коде генератора вы можете загружать и использовать эти настройки:**
  * **(main.go)**
```GO
type GeneratorConfig struct {
    PostNumber          int `json:"post_number"`
    GenerationIntervalMin int `json:"generation_interval_min"`
    GenerationIntervalMax int `json:"generation_interval_max"`
}

func LoadGeneratorConfig(filename string) (*GeneratorConfig, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("could not open config file: %v", err)
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    config := &GeneratorConfig{}
    err = decoder.Decode(config)
    if err != nil {
        return nil, fmt.Errorf("could not decode config JSON: %v", err)
    }

    return config, nil
}

func main() {
    // Загрузка конфигурации генератора
    genConfig, err := LoadGeneratorConfig("config_generator.json")
    if err != nil {
        log.Fatalf("Ошибка загрузки конфигурации генератора: %v", err)
    }
    
    // Использование значения из конфигурации
    postNumber := genConfig.PostNumber
    
    // ...остальной код...
}
```

Этот подход дает вам гибкость в настройке параметров генератора без необходимости перекомпиляции кода.
```
