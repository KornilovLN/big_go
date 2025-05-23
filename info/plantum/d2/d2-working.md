# Использование D2 для создания диаграмм
    D2 - это мощный инструмент для создания диаграмм с помощью декларативного языка.
    Основные шаги по использованию D2:

## Основы использования D2
* **Создать файл с расширением .d2:**
```bash
touch architecture.d2
```
* **Открыть файл в редакторе и добавить описание диаграммы:**

### Пример архитектуры Big Go
```
user: Пользователь {
  shape: person
}

services: Сервисы {
  generator: Generator {
    shape: rectangle
    style: {
      fill: "#CCFFCC"
    }
  }
  
  collector: Collector {
    shape: rectangle
    style: {
      fill: "#CCFFCC"
    }
  }
  
  user1: User1 {
    shape: rectangle
    style: {
      fill: "#CCFFCC"
    }
  }
  
  user2: User2 {
    shape: rectangle
    style: {
      fill: "#CCFFCC"
    }
  }
}

infrastructure: Инфраструктура {
  rabbitmq: RabbitMQ {
    shape: cylinder
    style: {
      fill: "#FFCCCC"
    }
  }
  
  postgres: PostgreSQL {
    shape: cylinder
    style: {
      fill: "#CCCCFF"
    }
  }
  
  redis: Redis {
    shape: cylinder
    style: {
      fill: "#FFFFCC"
    }
  }
}

### Связи
services.generator -> infrastructure.rabbitmq: Отправляет данные
infrastructure.rabbitmq -> services.collector: Передает данные
services.collector -> services.user1: HTTP
services.collector -> services.user2: HTTP
services.collector -> infrastructure.postgres: SQL
services.collector -> infrastructure.redis: Кэширует

user -> services.user1: Просматривает
user -> services.user2: Просматривает

services.user1 -> infrastructure.postgres: Читает данные
services.user2 -> infrastructure.postgres: Читает данные
infrastructure.rabbitmq -> infrastructure.redis: Сообщения
infrastructure.redis -> infrastructure.postgres: Передает данные
```

* **Сгенерировать диаграмму из файла:**
```bash
d2 architecture.d2 architecture.png
```

* **Результат:**
```bash
xdg-open architecture.png  # для Linux
```

* **Дополнительные возможности**
  * D2 поддерживает интерактивный режим, который автоматически обновляет диаграмму при изменении файла:
  ```bash
    d2 --watch architecture.d2 architecture.svg
  ```
* **Экспорт в разные форматы**
  * D2 поддерживает экспорт в различные форматы:
  ```bash 
  d2 architecture.d2 architecture.svg  # SVG формат
  d2 architecture.d2 architecture.pdf  # PDF формат
  d2 architecture.d2 architecture.png  # PNG формат
  ```
* **Настройка темы**
  * Вы можете использовать разные темы:
  ```bash
  d2 --theme 3 architecture.d2 architecture.png  # Использовать тему 3
  ```

## Создание сложных диаграмм
    D2 позволяет создавать сложные диаграммы с вложенными элементами, различными стилями и формами. 
    Вот пример более сложной диаграммы:

# Более сложная архитектура
```
system: "Big Go System" {
  user: Пользователь {
    shape: person
  }
  
  frontend: "Frontend" {
    user1: "User1 Dashboard" {
      shape: rectangle
      style: {
        fill: "#CCFFCC"
      }
    }
    
    user2: "User2 Dashboard" {
      shape: rectangle
      style: {
        fill: "#CCFFCC"
      }
    }
  }
  
  backend: "Backend Services" {
    generator: "Data Generator" {
      shape: rectangle
      style: {
        fill: "#FFCCCC"
      }
    }
    
    collector: "Data Collector" {
      shape: rectangle
      style: {
        fill: "#FFCCCC"
      }
    }
  }
  
  databases: "Databases & Messaging" {
    rabbitmq: "RabbitMQ" {
      shape: cylinder
      style: {
        fill: "#CCCCFF"
      }
    }
    
    postgres: "PostgreSQL" {
      shape: cylinder
      style: {
        fill: "#CCCCFF"
      }
    }
    
    redis: "Redis" {
      shape: cylinder
      style: {
        fill: "#CCCCFF"
      }
    }
  }
  
  # Связи
  backend.generator -> databases.rabbitmq: "Публикует данные"
  databases.rabbitmq -> backend.collector: "Потребляет данные"
  backend.collector -> frontend.user1: "HTTP API"
  backend.collector -> frontend.user2: "HTTP API"
  backend.collector -> databases.postgres: "Сохраняет"
  backend.collector -> databases.redis: "Кэширует"
  
  user -> frontend.user1: "Просматривает"
  user -> frontend.user2: "Просматривает"
  
  frontend.user1 -> databases.postgres: "Запрашивает данные"
  frontend.user2 -> databases.postgres: "Запрашивает данные"
  databases.rabbitmq -> databases.redis: "Уведомления"
  databases.redis -> databases.postgres: "Синхронизация"
}
```

## Документация
* **[D2 Documentation](https://github.com/terrastruct/d2)**